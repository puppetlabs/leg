package fs_test

import (
	"net/http"
	"testing"

	"github.com/puppetlabs/insights-stdlib/httputil/fs"
	"github.com/stretchr/testify/require"
)

func TestWrappedFileSystemsDoNotReportModificationTimes(t *testing.T) {
	dir := http.Dir(".")

	root, err := dir.Open("/")
	require.NoError(t, err)
	defer root.Close()

	files, err := root.Readdir(-1)
	require.NoError(t, err)
	require.True(t, len(files) > 0, "current directory bizarrely has no files")
	require.False(t, files[0].ModTime().IsZero(), "file has zero modification time on disk")

	wrapped := fs.FileSystemWithoutModTimes(dir)

	root, err = wrapped.Open("/")
	require.NoError(t, err)
	defer root.Close()

	files, err = root.Readdir(-1)
	require.NoError(t, err)
	require.True(t, files[0].ModTime().IsZero(), "wrapped file has modification time")
}
