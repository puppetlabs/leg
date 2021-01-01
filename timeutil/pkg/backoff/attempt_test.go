package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxAttempts(t *testing.T) {
	r, err := backoff.MaxAttempts(3).New()
	require.NoError(t, err)

	// Run through three attempts.
	for i := 1; i <= 3; i++ {
		generate, _, err := r.ApplyBefore()
		assert.NoError(t, err)
		assert.True(t, generate)

		final, err := r.ApplyAfter(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, 5*time.Second, final)
	}

	// Now do a final attempt, which should produce an error.
	generate, _, err := r.ApplyBefore()
	assert.Equal(t, &backoff.MaxAttemptsReachedError{N: 3}, err)
	assert.False(t, generate)
}
