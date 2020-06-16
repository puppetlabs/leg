package filesystem_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/horsehead/v2/storage"
	"github.com/puppetlabs/horsehead/v2/storage/testutils"
	"github.com/stretchr/testify/require"
)

func withTempDir(t *testing.T, fn func(t *testing.T, fs storage.BlobStore, tmp string)) {
	fs, cleanup, tmp := testutils.NewTempFilesystemBlobStore(t)
	defer cleanup()

	fn(t, fs, tmp)
}

func TestNewFilesystemFailsIfPathDoesNotExist(t *testing.T) {
	u, err := url.Parse("file:///does/not/exist")
	require.NoError(t, err)
	_, err = storage.NewBlobStore(*u)

	require.Error(t, err)
}

func TestCanPutFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTempDir(t, func(t *testing.T, backend storage.BlobStore, tmp string) {
		var testPayload = `{"test": "payload"}`

		cases := []struct {
			key    string
			reader io.Reader
		}{
			{"root-key.json", bytes.NewBufferString(testPayload)},
			{"ci-results/account-1234/all.json", bytes.NewBufferString(testPayload)},
		}

		for _, c := range cases {
			require.NoError(t, backend.Put(ctx, c.key, func(w io.Writer) error {
				_, err := io.Copy(w, c.reader)
				return err
			}, storage.PutOptions{}))

			path := filepath.Join(tmp, "blob/"+c.key)

			_, err := os.Stat(path)
			require.NoError(t, err)

			b, err := ioutil.ReadFile(path)
			require.NoError(t, err)

			require.Equal(t, string(b), testPayload)
		}
	})
}

func TestRange(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTempDir(t, func(t *testing.T, backend storage.BlobStore, tmp string) {
		var testPayload = []byte("012345678Xabcdef\n012345678Yabcdef\n012345678Zabcdef\n")

		require.NoError(t, backend.Put(ctx, "3hex", func(w io.Writer) error {
			_, err := w.Write(testPayload)
			return err
		}, storage.PutOptions{
			ContentType: "hex",
		}))

		cases := []struct {
			offset         int64
			length         int64
			absoluteOffset int64
			expect         string
		}{
			{-17, 0, 2 * 17, "012345678Zabcdef\n"},
			{17, 17, 17, "012345678Yabcdef\n"},
			{0, 17, 0, "012345678Xabcdef\n"},
		}

		for _, c := range cases {
			buf := bytes.Buffer{}
			err := backend.Get(ctx, "3hex", func(meta *storage.Meta, r io.Reader) error {
				require.Equal(t, meta.ContentType, "hex")
				require.Equal(t, meta.Offset, c.absoluteOffset)
				require.Equal(t, meta.Size, int64(17*3))

				chunk := make([]byte, 5)
				n, err := r.Read(chunk)
				require.NoError(t, err)
				require.Equal(t, n, 5)
				buf.Write(chunk)

				n, err = r.Read(chunk)
				require.NoError(t, err)
				require.Equal(t, n, 5)
				buf.Write(chunk)

				n, err = r.Read(chunk)
				require.NoError(t, err)
				require.Equal(t, n, 5)
				buf.Write(chunk)

				n, err = r.Read(chunk)
				require.NoError(t, err)
				require.Equal(t, n, 2)
				buf.Write(chunk[0:n])

				n, err = r.Read(chunk)
				require.Equal(t, err, io.EOF)
				return err
			}, storage.GetOptions{
				Offset: c.offset,
				Length: c.length,
			})
			require.NoError(t, err)
			require.Equal(t, c.expect, buf.String())
		}

		buf := bytes.Buffer{}
		err := backend.Get(ctx, "3hex", func(meta *storage.Meta, r io.Reader) error {
			require.Equal(t, meta.ContentType, "hex")
			require.Equal(t, meta.Offset, int64(0))
			require.Equal(t, meta.Size, int64(17*3))

			_, err := io.Copy(&buf, r)
			return err
		}, storage.GetOptions{
			// Get everything:
			Offset: -100,
		})
		require.NoError(t, err)
		require.Equal(t, testPayload, buf.Bytes())
	})
}

func TestCanGetFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTempDir(t, func(t *testing.T, backend storage.BlobStore, tmp string) {
		var testPayload = `{"test": "payload"}`

		cases := []struct {
			key    string
			reader *bytes.Buffer
		}{
			{"root-key.json", bytes.NewBufferString(testPayload)},
			{"ci-results/account-1234/all.json", bytes.NewBufferString(testPayload)},
		}

		for _, c := range cases {
			var err error

			expected := &bytes.Buffer{}
			tee := io.TeeReader(c.reader, expected)

			require.NoError(t, backend.Put(ctx, c.key, func(w io.Writer) error {
				_, err := io.Copy(w, tee)
				return err
			}, storage.PutOptions{}))

			buf := bytes.Buffer{}
			err = backend.Get(ctx, c.key, func(_ *storage.Meta, r io.Reader) error {
				_, err := buf.ReadFrom(r)
				return err
			}, storage.GetOptions{})
			require.NoError(t, err)

			require.Equal(t, expected.String(), buf.String())
		}
	})
}

func TestCanUpdateFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTempDir(t, func(t *testing.T, backend storage.BlobStore, tmp string) {
		var testPayload = `{"test": "payload"}`
		var updatedPayload = `{"test": "updated"}`

		cases := []struct {
			key    string
			reader *bytes.Buffer
		}{
			{"root-key.json", bytes.NewBufferString(testPayload)},
			{"ci-results/account-1234/all.json", bytes.NewBufferString(testPayload)},
		}

		for _, c := range cases {
			{
				var err error

				expected := &bytes.Buffer{}
				tee := io.TeeReader(c.reader, expected)

				require.NoError(t, backend.Put(ctx, c.key, func(w io.Writer) error {
					_, err := io.Copy(w, tee)
					return err
				}, storage.PutOptions{}))

				buf := bytes.Buffer{}
				err = backend.Get(ctx, c.key, func(_ *storage.Meta, r io.Reader) error {
					_, err := buf.ReadFrom(r)
					return err
				}, storage.GetOptions{})
				require.NoError(t, err)

				require.Equal(t, expected.String(), buf.String())
			}

			updatedReader := bytes.NewBufferString(updatedPayload)

			{
				var err error

				expected := &bytes.Buffer{}
				tee := io.TeeReader(updatedReader, expected)

				require.NoError(t, backend.Put(ctx, c.key, func(w io.Writer) error {
					_, err := io.Copy(w, tee)
					return err
				}, storage.PutOptions{}))

				buf := bytes.Buffer{}
				err = backend.Get(ctx, c.key, func(_ *storage.Meta, r io.Reader) error {
					_, err := buf.ReadFrom(r)
					return err
				}, storage.GetOptions{})
				require.NoError(t, err)

				require.Equal(t, expected.String(), buf.String())
			}
		}
	})
}

func TestCanDeleteFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTempDir(t, func(t *testing.T, backend storage.BlobStore, tmp string) {
		var testPayload = `{"test": "payload"}`

		cases := []struct {
			key    string
			reader *bytes.Buffer
		}{
			{"root-key.json", bytes.NewBufferString(testPayload)},
			{"ci-results/account-1234/all.json", bytes.NewBufferString(testPayload)},
		}

		for _, c := range cases {
			var err error

			expected := &bytes.Buffer{}
			tee := io.TeeReader(c.reader, expected)

			require.NoError(t, backend.Put(ctx, c.key, func(w io.Writer) error {
				_, err := io.Copy(w, tee)
				return err
			}, storage.PutOptions{}))

			buf := bytes.Buffer{}
			err = backend.Get(ctx, c.key, func(_ *storage.Meta, r io.Reader) error {
				_, err := buf.ReadFrom(r)
				return err
			}, storage.GetOptions{})
			require.NoError(t, err)

			require.Equal(t, expected.String(), buf.String())

			require.NoError(t, backend.Delete(ctx, c.key, storage.DeleteOptions{}))

			_, err = os.Stat(filepath.Join(tmp, "blob/"+c.key))
			require.Error(t, err)

			err = backend.Get(ctx, c.key, func(_ *storage.Meta, _ io.Reader) error {
				return nil
			}, storage.GetOptions{})
			require.Error(t, err)
		}
	})
}
