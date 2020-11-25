package gcs

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	gcstorage "cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/puppetlabs/leg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	raw "google.golang.org/api/storage/v1"
)

func toUrl(s string) url.URL {
	u, err := url.ParseRequestURI(s)
	if nil != err {
		panic(err)
	}
	return *u
}

func TestRealGCS(t *testing.T) {
	bucketName := os.Getenv("GCS_BUCKET")
	if 0 == len(bucketName) || 0 == len(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")) {
		t.Skip("Define the GCS_BUCKET and GOOGLE_APPLICATION_CREDENTIALS environment variables to enable GCS tests")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	key, err := uuid.NewRandom()
	assert.NoError(t, err)

	content := []byte("TEST CONTENT")

	gcs, err := storage.NewBlobStore(
		toUrl("gs://" + bucketName))

	assert.NoError(t, err)
	err = gcs.Put(ctx, key.String(), func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	}, storage.PutOptions{
		ContentType: "application/testing",
	})
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = gcs.Get(ctx, key.String(), func(meta *storage.Meta, r io.Reader) error {
		assert.Equal(t, meta.ContentType, "application/testing")
		assert.Equal(t, meta.Offset, int64(8))
		assert.Equal(t, meta.Size, int64(12))
		_, err := io.Copy(&buf, r)
		return err
	}, storage.GetOptions{
		Offset: -4,
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("TENT"), buf.Bytes())

	buf.Reset()
	err = gcs.Get(ctx, key.String(), func(meta *storage.Meta, r io.Reader) error {
		assert.Equal(t, meta.ContentType, "application/testing")
		assert.Equal(t, meta.Offset, int64(5))
		assert.Equal(t, meta.Size, int64(12))
		_, err := io.Copy(&buf, r)
		return err
	}, storage.GetOptions{
		Offset: 5,
		Length: 3,
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("CON"), buf.Bytes())

	buf.Reset()
	err = gcs.Get(ctx, key.String(), func(meta *storage.Meta, r io.Reader) error {
		assert.Equal(t, meta.ContentType, "application/testing")
		assert.Equal(t, meta.Offset, int64(0))
		assert.Equal(t, meta.Size, int64(12))
		_, err := io.Copy(&buf, r)
		return err
	}, storage.GetOptions{
		Offset: -100,
	})
	assert.NoError(t, err)
	assert.Equal(t, content, buf.Bytes())

	assert.NoError(t, gcs.Delete(ctx, key.String(), storage.DeleteOptions{}))
}

func unexpectedEOFHandler(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("Unable to create hijacker")
		return
	}

	conn, bufrw, err := hj.Hijack()
	if err != nil {
		fmt.Println("Unable to hijack request")
		return
	}

	bufrw.Flush()
	conn.Close()
}

func withTestServer(t *testing.T, h http.Handler, fn func(gcs storage.BlobStore)) {
	ts := httptest.NewTLSServer(h)
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS: func(netw, addr string) (net.Conn, error) {
			return tls.Dial("tcp", ts.Listener.Addr().String(), tlsConfig)
		},
	}

	httpClient := &http.Client{Transport: transport}
	gcsClient, err := gcstorage.NewClient(context.Background(), option.WithHTTPClient(httpClient))
	require.NoError(t, err)

	gcs, err := newGCS(toUrl("gs://bucket=insights-dataflow-bulk-storage"), gcsClient)
	require.NoError(t, err)
	fn(gcs)

	transport.CloseIdleConnections()
	ts.Close()
}

func TestPutRetry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTestServer(t, http.HandlerFunc(unexpectedEOFHandler), func(gcs storage.BlobStore) {
		buf := bytes.NewBufferString("test string")
		err := gcs.Put(ctx, "test/key", func(w io.Writer) error {
			_, err := io.Copy(w, buf)
			return err
		}, storage.PutOptions{})

		require.Error(t, err)
	})
}

func TestGetRetry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTestServer(t, http.HandlerFunc(unexpectedEOFHandler), func(gcs storage.BlobStore) {
		buf := &bytes.Buffer{}
		err := gcs.Get(ctx, "test/key", func(meta *storage.Meta, r io.Reader) error {
			_, err := io.Copy(buf, r)
			return err
		}, storage.GetOptions{})

		require.Error(t, err)
	})
}

func TestDeleteRetry(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withTestServer(t, http.HandlerFunc(unexpectedEOFHandler), func(gcs storage.BlobStore) {
		err := gcs.Delete(ctx, "test/key", storage.DeleteOptions{})

		require.Error(t, err)
	})
}

func TestPut(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	payload := []byte("test string")

	h := func(w http.ResponseWriter, r *http.Request) {
		_, params, err := mime.ParseMediaType(r.Header["Content-Type"][0])
		if nil != err {
			t.Fatal(err)
		}
		rd := multipart.NewReader(r.Body, params["boundary"])
		rd.NextPart()
		part, err := rd.NextPart()
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, part)
		require.Equal(t, buf.Bytes(), payload)
		res := &raw.RewriteResponse{Done: true}
		bytes, err := res.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		w.Write(bytes)
	}

	withTestServer(t, http.HandlerFunc(h), func(gcs storage.BlobStore) {
		buf := bytes.NewBuffer(payload)

		require.NoError(t, gcs.Put(ctx, "test/key", func(w io.Writer) error {
			_, err := io.Copy(w, buf)
			return err
		}, storage.PutOptions{}))
	})
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	payload := []byte("{\"test\": \"payload\"}")

	h := func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}

	withTestServer(t, http.HandlerFunc(h), func(gcs storage.BlobStore) {
		buf := &bytes.Buffer{}
		err := gcs.Get(ctx, "test/key", func(meta *storage.Meta, r io.Reader) error {
			_, err := io.Copy(buf, r)
			return err
		}, storage.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, payload, buf.Bytes())
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	h := func(w http.ResponseWriter, r *http.Request) {
		res := &raw.RewriteResponse{Done: true}
		bytes, err := res.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		w.Write(bytes)
	}

	withTestServer(t, http.HandlerFunc(h), func(gcs storage.BlobStore) {
		require.NoError(t, gcs.Delete(ctx, "test/key", storage.DeleteOptions{}))
	})
}
