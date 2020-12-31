package backoff_test

import (
	"math"
	"testing"
	"time"

	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJitter(t *testing.T) {
	rf := rand.NewTestFactory(rand.ZeroSeeder)

	r, err := backoff.FullJitter(backoff.JitterWithRandFactory(rf)).New()
	require.NoError(t, err)

	prev := time.Duration(-1)
	for i := 0; i < 10_000; i++ {
		final, err := r.ApplyAfter(time.Duration(math.MaxInt64))
		require.NoError(t, err)

		// Using our test RNG each following jitter value should be
		// monotonically increasing.
		assert.Greater(t, int64(final), int64(prev), "iteration #%d", i)
		prev = final
	}
}

func TestJitterPercentage(t *testing.T) {
	tests := []struct {
		Name            string
		Factory         backoff.RuleFactory
		Initial         time.Duration
		ExpectedAtLeast time.Duration
	}{
		{
			Name:            "Equal",
			Factory:         backoff.EqualJitter(backoff.JitterWithRandFactory(randFactory)),
			Initial:         10 * time.Second,
			ExpectedAtLeast: 5 * time.Second,
		},
		{
			Name:            "Tenth",
			Factory:         backoff.Jitter(0.1, backoff.JitterWithRandFactory(randFactory)),
			Initial:         10 * time.Second,
			ExpectedAtLeast: 9 * time.Second,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r, err := test.Factory.New()
			require.NoError(t, err)

			for i := 0; i < 10_000; i++ {
				generate, _, err := r.ApplyBefore()
				require.NoError(t, err)
				require.Equal(t, true, generate)

				final, err := r.ApplyAfter(test.Initial)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, int64(final), int64(test.ExpectedAtLeast))
				assert.Less(t, int64(final), int64(test.Initial))
			}
		})
	}
}
