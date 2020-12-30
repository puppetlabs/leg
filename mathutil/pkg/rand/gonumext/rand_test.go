package gonumext_test

import (
	"testing"

	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	"github.com/puppetlabs/leg/mathutil/pkg/rand/gonumext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMersenneTwisterFactory(t *testing.T) {
	mt, err := gonumext.NewMersenneTwisterFactory(rand.NewConstantSeeder(5489)).New()
	require.NoError(t, err)

	rv, err := rand.Uint64(mt)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xd091bb5c22ae9ef6), rv)

	rv, err = rand.Uint64(mt)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xe7e1faeed5c31f79), rv)
}

func TestMersenneTwister64Factory(t *testing.T) {
	mt, err := gonumext.NewMersenneTwister64Factory(rand.NewConstantSeeder(5489)).New()
	require.NoError(t, err)

	rv, err := rand.Uint64(mt)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xc96d191cf6f6aea6), rv)

	rv, err = rand.Uint64(mt)
	require.NoError(t, err)
	assert.Equal(t, uint64(0x401f7ac78bc80f1c), rv)
}
