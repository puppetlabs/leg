package backoff_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinMaxBound(t *testing.T) {
	tests := []struct {
		Name     string
		Factory  backoff.RuleFactory
		Initial  time.Duration
		Expected time.Duration
	}{
		{
			Name:     "Minimum below",
			Factory:  backoff.MinBound(5 * time.Second),
			Initial:  2 * time.Second,
			Expected: 5 * time.Second,
		},
		{
			Name:     "Minimum above",
			Factory:  backoff.MinBound(5 * time.Second),
			Initial:  10 * time.Second,
			Expected: 10 * time.Second,
		},
		{
			Name:     "Maximum below",
			Factory:  backoff.MaxBound(5 * time.Second),
			Initial:  2 * time.Second,
			Expected: 2 * time.Second,
		},
		{
			Name:     "Maximum above",
			Factory:  backoff.MaxBound(5 * time.Second),
			Initial:  10 * time.Second,
			Expected: 5 * time.Second,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r, err := test.Factory.New()
			require.NoError(t, err)

			generate, _, err := r.ApplyBefore()
			require.NoError(t, err)
			require.Equal(t, true, generate)

			final, err := r.ApplyAfter(test.Initial)
			require.NoError(t, err)
			assert.Equal(t, test.Expected, final)
		})
	}
}
