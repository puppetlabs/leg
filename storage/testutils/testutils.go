package testutils

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/puppetlabs/horsehead/storage"
	"github.com/stretchr/testify/require"
)

func mustCreateTempDir(t *testing.T) string {
	tmp, err := ioutil.TempDir("", "storage-test")
	require.NoError(t, err)

	t.Logf("[blob storage] created temp directory: %s", tmp)

	return tmp
}

func NewTempFilesystemBlobStore(t *testing.T) (storage.BlobStore, func(), string) {
	tmp := mustCreateTempDir(t)

	u, err := url.Parse("file://" + tmp)
	require.NoError(t, err)
	fs, err := storage.NewBlobStore(*u)
	require.NoError(t, err)

	return fs, func() {
		require.NoError(t, os.RemoveAll(tmp))
		t.Logf("[blob storage] cleaned up temp directory %s", tmp)
	}, tmp
}
