package testutils

import (
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/puppetlabs/horsehead/storage"
	"github.com/puppetlabs/horsehead/workdir"
	"github.com/stretchr/testify/require"
)

func NewTempFilesystemBlobStore(t *testing.T) (storage.BlobStore, func() error, string) {
	wd, err := workdir.NewNamespace([]string{"tmp-blob-store", uuid.New().String()}).New(workdir.DirTypeCache, workdir.Options{
		Mode: os.FileMode(0700),
	})
	require.NoError(t, err)

	u, err := url.Parse("file://" + wd.Path)
	require.NoError(t, err)
	fs, err := storage.NewBlobStore(*u)
	require.NoError(t, err)

	return fs, wd.Cleanup, wd.Path
}
