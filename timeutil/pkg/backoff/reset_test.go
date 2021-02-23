package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestResetAfter(t *testing.T) {
	ctx := context.Background()

	type expected struct {
		Duration  time.Duration
		Error     error
		StepAfter time.Duration
	}

	tests := []struct {
		Name     string
		Factory  func(c clock.PassiveClock) backoff.GeneratorFactory
		Expected []expected
	}{
		{
			Name: "Maximum attempts",
			Factory: func(c clock.PassiveClock) backoff.GeneratorFactory {
				return backoff.ResetAfter(
					backoff.Build(
						backoff.Linear(5*time.Second),
						backoff.MaxAttempts(3),
					),
					30*time.Second,
					backoff.ResetAfterWithClock(c),
				)
			},
			Expected: []expected{
				{
					Duration:  5 * time.Second,
					StepAfter: 10 * time.Second,
				},
				{
					Duration:  10 * time.Second,
					StepAfter: 10 * time.Second,
				},
				{
					Duration:  15 * time.Second,
					StepAfter: 35 * time.Second,
				},
				{
					Duration:  5 * time.Second,
					StepAfter: 10 * time.Second,
				},
				{
					Duration:  10 * time.Second,
					StepAfter: 10 * time.Second,
				},
				{
					Duration:  15 * time.Second,
					StepAfter: 15 * time.Second,
				},
				{
					Error: &backoff.MaxAttemptsReachedError{N: 3},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fc := clock.NewFakeClock(time.Now())

			g, err := test.Factory(fc).New()
			require.NoError(t, err)

			for i, step := range test.Expected {
				next, err := g.Next(ctx)
				assert.Equal(t, step.Error, err, "step #%d", i)
				assert.Equal(t, step.Duration, next, "step #%d", i)

				fc.Step(step.StepAfter)
			}
		})
	}
}
