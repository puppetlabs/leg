package fs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/horsehead/httputil/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMap(t *testing.T) {
	dir, err := ioutil.TempDir("", "idf-httputil-test-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "t1.txt"), []byte{}, 0644))
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "t2.dat"), []byte{}, 0644))
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "t3.txt"), []byte{}, 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "t4"), 0755))
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "t4", "a.txt"), []byte{}, 0644))

	fm := fs.FileMap(dir, []string{"t1.txt", "t2.dat", "t4/a.txt"})

	tests := []struct {
		Name     string
		Children []string
	}{
		{"/", []string{"t1.txt", "t2.dat", "t4"}},
		{"/t1.txt", nil},
		{"/t2.dat", nil},
		{"/t4", []string{"a.txt"}},
		{"/t4/a.txt", nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			f, err := fm.Open(test.Name)
			require.NoError(t, err)
			defer f.Close()

			s, err := f.Stat()
			require.NoError(t, err)

			if len(test.Children) > 0 {
				require.True(t, s.IsDir(), "file should be a directory but is not")

				contents, err := f.Readdir(-1)
				require.NoError(t, err)

				assert.Len(t, contents, len(test.Children))

				files := make([]string, len(contents))
				for i, fi := range contents {
					files[i] = fi.Name()
				}

				for _, child := range test.Children {
					assert.Contains(t, files, child)
				}
			} else {
				assert.False(t, s.IsDir(), "file should not be a directory but is")
			}
		})
	}
}
