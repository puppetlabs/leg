package backoff_test

import (
	"math"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExponential(t *testing.T) {
	g, err := backoff.Exponential(10*time.Second, 60.0).New()
	require.NoError(t, err)

	expected := []time.Duration{
		10 * time.Second,
		10 * time.Minute,
		10 * time.Hour,
		600 * time.Hour,
		36_000 * time.Hour,
		2_160_000 * time.Hour,
		math.MaxInt64,
		math.MaxInt64,
		math.MaxInt64,
	}
	for i, step := range expected {
		wait, err := g.Next()
		require.NoError(t, err)
		assert.Equal(t, step, wait, "step #%d", i)
	}
}

func TestDecorrelatedExponential(t *testing.T) {
	tests := []struct {
		Name     string
		Factory  backoff.GeneratorFactory
		Rule     backoff.RuleFactory
		Expected []time.Duration
	}{
		{
			Name:    "Basic",
			Factory: backoff.DecorrelatedExponential(5*time.Second, 2, backoff.DecorrelatedExponentialWithRandFactory(randFactory)),
			Expected: []time.Duration{
				0x17fd09d9b,
				0x2c5c962d7,
				0x4a8d71d9c,
				0x3d7ceeb3d,
				0x1dacb2137,
			},
		},
		{
			Name:    "Capped",
			Factory: backoff.DecorrelatedExponential(5*time.Second, 2, backoff.DecorrelatedExponentialWithRandFactory(randFactory)),
			Rule:    backoff.MaxBound(15 * time.Second),
			Expected: []time.Duration{
				0x17fd09d9b,
				0x2c5c962d7,
				0x37e11d600,
				0x31384da01,
				0x1b13a8ddb,
			},
		},
		{
			Name:    "Non-sliding",
			Factory: backoff.DecorrelatedExponential(5*time.Second, 2, backoff.DecorrelatedExponentialWithRandFactory(randFactory)),
			Rule:    backoff.NonSliding,
			Expected: []time.Duration{
				0,
				0x17fd09d9b,
				0x2c5c962d7,
				0x4a8d71d9c,
				0x3d7ceeb3d,
				0x1dacb2137,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			g, err := test.Factory.New()
			require.NoError(t, err)

			if test.Rule != nil {
				r, err := test.Rule.New()
				require.NoError(t, err)

				g.(backoff.RuleInjector).InjectRule(r)
			}

			for i, step := range test.Expected {
				wait, err := g.Next()
				require.NoError(t, err)
				assert.Equal(t, step, wait, "step #%d", i)
			}
		})
	}
}
