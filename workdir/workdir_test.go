package workdir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewNamespace(t *testing.T) {
	t.Parallel()

	testID := uuid.New().String()

	var cases = []struct {
		description string
		setup       func()
		dirType     dirType
		namespace   []string
		expected    string
		shouldError bool
	}{
		{
			description: "can create config dirs with XDG var set",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CONFIG_HOME", "/tmp/"))
			},
			dirType:   DirTypeConfig,
			namespace: []string{testID, "leg", "config-dir-test"},
			expected:  filepath.Join("/tmp", testID, "leg", "config-dir-test"),
		},
		{
			description: "can create cache dirs with XDG var set",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CACHE_HOME", "/tmp/"))
			},
			dirType:   DirTypeCache,
			namespace: []string{testID, "leg", "cache-dir-test"},
			expected:  filepath.Join("/tmp", testID, "leg", "cache-dir-test"),
		},
		{
			description: "can create data dirs with XDG var set",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_DATA_HOME", "/tmp/"))
			},
			dirType:   DirTypeData,
			namespace: []string{testID, "leg", "data-dir-test"},
			expected:  filepath.Join("/tmp", testID, "leg", "data-dir-test"),
		},
		{
			description: "can create config dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CONFIG_HOME", ""))
			},
			dirType:   DirTypeConfig,
			namespace: []string{testID, "leg", "config-dir-test"},
			expected:  filepath.Join(os.Getenv("HOME"), ".config", testID, "leg", "config-dir-test"),
		},
		{
			description: "can create cache dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CACHE_HOME", ""))
			},
			dirType:   DirTypeCache,
			namespace: []string{testID, "leg", "cache-dir-test"},
			expected:  filepath.Join(os.Getenv("HOME"), ".cache", testID, "leg", "cache-dir-test"),
		},
		{
			description: "can create data dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_DATA_HOME", ""))
			},
			dirType:   DirTypeData,
			namespace: []string{testID, "leg", "data-dir-test"},
			expected:  filepath.Join(os.Getenv("HOME"), ".local", "share", testID, "leg", "data-dir-test"),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if c.setup != nil {
				c.setup()
			}

			wd, err := NewNamespace(c.namespace).New(c.dirType, Options{})
			if c.shouldError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.expected, wd.Path)

			_, err = os.Stat(c.expected)
			require.NoError(t, err, "directory should exist")

			require.NoError(t, wd.Cleanup())

			_, err = os.Stat(c.expected)
			require.Error(t, err, "expected directory to be gone")
		})
	}
}
