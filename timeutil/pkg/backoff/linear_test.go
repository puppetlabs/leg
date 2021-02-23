package backoff_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinear(t *testing.T) {
	ctx := context.Background()

	g, err := backoff.Linear(5 * time.Second).New()
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		wait, err := g.Next(ctx)
		require.NoError(t, err)
		assert.Equal(t, 5*time.Duration(i+1)*time.Second, wait, "iteration #%d", i)
	}
}

func TestLinearWraparound(t *testing.T) {
	ctx := context.Background()

	g, err := backoff.Linear(1 << 61).New()
	require.NoError(t, err)

	expected := []time.Duration{
		1 << 61,
		1 << 62,
		1<<62 | 1<<61,
		math.MaxInt64,
		math.MaxInt64,
		math.MaxInt64,
	}
	for i, step := range expected {
		wait, err := g.Next(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(step), int64(wait), "step #%d", i)
	}
}
