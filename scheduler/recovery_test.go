package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/puppetlabs/leg/timeutil/pkg/clock"
	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

type mockErrorDescriptor struct {
	errCount          int
	successAfterCount int
}

func (d *mockErrorDescriptor) Run(ctx context.Context, pc chan<- Process) error {
	if d.errCount >= 0 && d.successAfterCount != 0 {
		err := fmt.Errorf("err count %d", d.errCount)
		d.errCount--
		d.successAfterCount--

		return err
	}

	return nil
}

func TestRecoverySchedulerStops(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	mock := &mockErrorDescriptor{
		errCount: 5,
		// make sure we never succeed
		successAfterCount: 15,
	}

	descriptor := NewRecoveryDescriptor(mock, RecoveryDescriptorWithBackoffFactory(
		backoff.Build(
			backoff.Immediate,
			backoff.MaxRetries(5),
		),
	))

	pc := make(chan Process)

	defer cancel()

	err := descriptor.Run(ctx, pc)
	assert.Equal(t, &backoff.MaxAttemptsReachedError{N: 6}, err)
	assert.Equal(t, -1, mock.errCount)
}

type mockRetryResetDescriptor struct {
	count           int
	fc              *testclock.FakeClock
	successDuration time.Duration
}

func (d *mockRetryResetDescriptor) Run(ctx context.Context, pc chan<- Process) error {
	if d.count%5 == 0 {
		d.fc.Step(d.successDuration)
	}

	if d.count >= 0 {
		err := fmt.Errorf("err count %d", d.count)
		d.count--
		return err
	}

	return nil
}

func TestRecoverySchedulerRetryCountReset(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fc := testclock.NewFakeClock(time.Now())
	dc := clock.NewTimerCallbackClock(
		k8sext.NewClock(fc),
		func(d time.Duration) { fc.Step(d) },
	)

	successDuration := 5 * time.Second

	mock := &mockRetryResetDescriptor{
		count:           14,
		fc:              fc,
		successDuration: successDuration,
	}

	descriptor := NewRecoveryDescriptor(
		mock,
		RecoveryDescriptorWithBackoffFactory(
			backoff.Build(
				backoff.ResetAfter(
					backoff.Build(
						backoff.Constant(250*time.Millisecond),
						backoff.MaxRetries(10),
					),
					successDuration-(500*time.Millisecond),
					backoff.ResetAfterWithClock(dc),
				),
			),
		),
		RecoveryDescriptorWithClock(dc),
	)

	pc := make(chan Process)

	defer cancel()

	require.NoError(t, descriptor.Run(ctx, pc))
	require.Equal(t, -1, mock.count)
}
