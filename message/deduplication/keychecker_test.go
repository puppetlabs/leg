package deduplication

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeySetterAndChecker(t *testing.T) {
	bf := NewBloomFilter(NewEphemeralBloomFilterBackend())

	err := bf.SetKeyAsSeen("this is a test key")
	require.NoError(t, err)

	exists, err := bf.KeyHasBeenSeen("this is a test key")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = bf.KeyHasBeenSeen("this key isn't in there")
	require.NoError(t, err)
	require.False(t, exists)
}
