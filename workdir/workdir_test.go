package workdir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testID := uuid.New().String()

	var cases = []struct {
		description string
		setup       func()
		provided    string
		dirType     dirType
		namespace   []string
		expected    string
		shouldError bool
	}{
		{
			description: "make sure we can create config dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CONFIG_HOME", "/tmp/"))
			},
			dirType:   DirTypeConfig,
			namespace: []string{testID, "horsehead", "config-dir-test"},
			expected:  filepath.Join("/tmp", testID, "horsehead", "config-dir-test"),
		},
		{
			description: "make sure we can create cache dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_CACHE_HOME", "/tmp/"))
			},
			dirType:   DirTypeCache,
			namespace: []string{testID, "horsehead", "cache-dir-test"},
			expected:  filepath.Join("/tmp", testID, "horsehead", "cache-dir-test"),
		},
		{
			description: "make sure we can create data dirs",
			setup: func() {
				require.NoError(t, os.Setenv("XDG_DATA_HOME", "/tmp/"))
			},
			dirType:   DirTypeData,
			namespace: []string{testID, "horsehead", "data-dir-test"},
			expected:  filepath.Join("/tmp", testID, "horsehead", "data-dir-test"),
		},
		{
			description: "make sure we can create provided dirs",
			provided:    filepath.Join("/tmp", testID, "horsehead", "provided-dir-test"),
			expected:    filepath.Join("/tmp", testID, "horsehead", "provided-dir-test"),
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if c.setup != nil {
				c.setup()
			}

			wd, err := New(c.provided, c.dirType, c.namespace, Options{})
			if c.shouldError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.expected, wd.Path)

			_, err = os.Stat(c.expected)
			require.NoError(t, err)

			require.NoError(t, wd.Cleanup())
		})
	}
}
