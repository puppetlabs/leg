package ephemeralbloomfilter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEphemeralBackend(t *testing.T) {
	ebf := New()

	key1 := "test-key1"
	key2 := "test-key2"

	t.Run("key1 can be set and tested for existence", func(t *testing.T) {
		require.NoError(t, ebf.Set(key1))

		exists, err := ebf.Exists(key1)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("key2 should not exist yet", func(t *testing.T) {
		exists, err := ebf.Exists(key2)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("key2 can be set with CheckAndSet", func(t *testing.T) {
		exists, err := ebf.CheckAndSet(key2)
		require.NoError(t, err)
		require.False(t, exists)

		exists, err = ebf.CheckAndSet(key2)
		require.NoError(t, err)
		require.True(t, exists)
	})
}
