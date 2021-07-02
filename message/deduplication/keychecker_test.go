package deduplication

import (
	"testing"

	"github.com/puppetlabs/leg/message/deduplication/ephemeralbloomfilter"
	"github.com/stretchr/testify/require"
)

func TestKeySetterAndChecker(t *testing.T) {
	t.Run("setting a key should exist when checked", func(t *testing.T) {
		bf := NewBloomFilter(ephemeralbloomfilter.New())

		key := "test-key"

		err := bf.SetKeyAsSeen(key)
		require.NoError(t, err)

		exists, err := bf.KeyHasBeenSeen(key)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("checking a key that has not been set should return false", func(t *testing.T) {
		bf := NewBloomFilter(ephemeralbloomfilter.New())

		key := "missing-key"

		exists, err := bf.KeyHasBeenSeen(key)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("check and set should detect a key that has already been added before", func(t *testing.T) {
		bf := NewBloomFilter(ephemeralbloomfilter.New())

		key := "test-key"

		exists, err := bf.CheckAndSetKey(key)
		require.NoError(t, err)
		require.False(t, exists)

		exists, err = bf.CheckAndSetKey(key)
		require.NoError(t, err)
		require.True(t, exists)

		key2 := "test-another-key"

		exists, err = bf.CheckAndSetKey(key2)
		require.NoError(t, err)
		require.False(t, exists)
	})
}
