package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRules(t *testing.T) {
	// We pick an early-out rule that should give us generate == false.
	r1, err := backoff.NonSliding.New()
	require.NoError(t, err)

	// Note that we apply the min-bound rule after non-sliding (i.e., order
	// matters).
	r2, err := backoff.MinBound(5 * time.Second).New()
	require.NoError(t, err)

	rs := backoff.Rules{r1, r2}

	generate, next, err := rs.ApplyBefore()
	require.NoError(t, err)
	assert.Equal(t, false, generate)

	final, err := rs.ApplyAfter(next)
	require.NoError(t, err)
	assert.Equal(t, 5*time.Second, final)
}
