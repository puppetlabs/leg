package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonSliding(t *testing.T) {
	r, err := backoff.NonSliding.New()
	require.NoError(t, err)

	generate, next, err := r.ApplyBefore()
	require.NoError(t, err)

	// First application should give us generate == false. Then it should be
	// true.
	assert.Equal(t, generate, false)
	assert.Equal(t, time.Duration(0), next)

	final, err := r.ApplyAfter(next)
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), final)

	for i := 0; i < 10; i++ {
		generate, _, err := r.ApplyBefore()
		require.NoError(t, err)
		assert.Equal(t, generate, true)

		final, err := r.ApplyAfter(5 * time.Second)
		require.NoError(t, err)
		assert.Equal(t, 5*time.Second, final)
	}
}
