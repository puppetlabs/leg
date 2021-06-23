package retry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/puppetlabs/leg/timeutil/pkg/clock"
	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

func TestWait(t *testing.T) {
	// Global test context just in case we really mess something up.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		Name           string
		BackoffFactory *backoff.Factory
		Attempt        func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error)
		Step           func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func()
	}{
		{
			Name: "Always succeed",
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					return true, nil
				})
				require.NoError(t, err)
			},
			Step: func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func() {
				return func() { c.Step(1 * time.Millisecond) }
			},
		},
		{
			Name: "Succeeds after 3 attempts",
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				i := 0
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					i++
					return i >= 3, nil
				})
				require.NoError(t, err)
				assert.Equal(t, 3, i)
			},
			Step: func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func() {
				return func() { c.Step(1 * time.Minute) }
			},
		},
		{
			Name: "Returns error after 3 attempts",
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				i := 0
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					i++
					return i >= 3, fmt.Errorf("boom %d", i)
				})
				require.EqualError(t, err, "boom 3")
				assert.Equal(t, 3, i)
			},
			Step: func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func() {
				return func() { c.Step(1 * time.Minute) }
			},
		},
		{
			Name: "Context cancellation without error",
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					return false, nil
				})
				require.EqualError(t, err, "context canceled")
			},
			Step: func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func() {
				i := 0
				return func() {
					c.Step(1 * time.Minute)
					if i > 2 {
						cancel()
					}
					i++
				}
			},
		},
		{
			Name: "Context cancellation with error",
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				i := 0
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					i++
					return i >= 3, fmt.Errorf("boom %d", i)
				})
				require.EqualError(t, err, "boom 2")
			},
			Step: func(ctx context.Context, cancel context.CancelFunc, c *testclock.FakeClock) func() {
				i := 0
				return func() {
					if i >= 2 {
						cancel()
					}
					i++
					c.Step(1 * time.Minute)
				}
			},
		},
		{
			Name: "Limited backoff without error",
			BackoffFactory: backoff.Build(
				backoff.Immediate,
				backoff.MaxAttempts(3),
			),
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					return false, nil
				})
				require.EqualError(t, err, "maximum attempts of 3 reached")
			},
		},
		{
			Name: "Limited backoff with error",
			BackoffFactory: backoff.Build(
				backoff.Immediate,
				backoff.MaxAttempts(3),
			),
			Attempt: func(t *testing.T, ctx context.Context, fn func(context.Context, retry.WorkFunc) error) {
				i := 0
				err := fn(ctx, func(ctx context.Context) (bool, error) {
					i++
					return false, fmt.Errorf("boom %d", i)
				})
				require.EqualError(t, err, "boom 3")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			c := testclock.NewFakeClock(time.Now())

			var opts []retry.WaitOption
			if test.Step != nil {
				step := test.Step(ctx, cancel, c)
				opts = append(opts, retry.WithClock(
					clock.NewTimerCallbackClock(
						k8sext.NewClock(c),
						func(d time.Duration) { step() },
					),
				))
			}
			if test.BackoffFactory != nil {
				opts = append(opts, retry.WithBackoffFactory(test.BackoffFactory))
			}

			test.Attempt(t, ctx, func(ctx context.Context, work retry.WorkFunc) error {
				return retry.Wait(ctx, work, opts...)
			})
			assert.False(t, c.HasWaiters())
		})
	}
}

func TestWaitAsync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fc := testclock.NewFakeClock(time.Now())
	defer func() {
		assert.False(t, fc.HasWaiters())
	}()

	tcc := clock.NewTimerCallbackClock(
		k8sext.NewClock(fc),
		func(d time.Duration) {
			fc.Step(d)
		},
	)

	wt := make(chan int)

	i := 0
	ch := retry.WaitAsync(ctx, func(ctx context.Context) (bool, error) {
		i++
		wt <- i
		return i >= 3, fmt.Errorf("boom %d", i)
	}, retry.WithClock(tcc))

	for j := 1; j <= 3; j++ {
		// Channel should be empty.
		select {
		case err := <-ch:
			require.Fail(t, "asynchronous retry returned early", "attempt #%d, error %+v", j, err)
		default:
		}

		// Wait for more work (i.e., we've gone through one cycle of waiting on
		// the internal timer in Wait()).
		select {
		case ci := <-wt:
			assert.Equal(t, j, ci)
		case <-ctx.Done():
			require.Fail(t, "asynchronous retry did not wake up", "attempt #%d", j)
		}
	}

	select {
	case err := <-ch:
		require.EqualError(t, err, "boom 3")
	case <-ctx.Done():
		require.Fail(t, "asynchronous retry did not provide result on channel")
	}
}
