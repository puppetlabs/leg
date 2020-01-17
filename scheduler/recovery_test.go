package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		errCount: 10,
		// make sure we never succeed
		successAfterCount: 15,
	}

	opts := RecoveryDescriptorOptions{
		MaxRetries: 10,
	}
	descriptor := NewRecoveryDescriptorWithOptions(mock, opts)

	pc := make(chan Process)

	defer cancel()

	require.Error(t, descriptor.Run(ctx, pc))
}

type mockRetryResetDescriptor struct {
	count           int
	successDuration time.Duration
	cancel          context.CancelFunc
}

func (d *mockRetryResetDescriptor) Run(ctx context.Context, pc chan<- Process) error {
	if d.count == 0 {
		<-time.After(d.successDuration)
		d.cancel()
		return nil
	}

	err := fmt.Errorf("err count %d", d.count)
	d.count--

	return err
}

func TestRecoverySchedulerRetryCountReset(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	successDuration := time.Second * 1

	mock := &mockRetryResetDescriptor{
		count:           5,
		cancel:          cancel,
		successDuration: successDuration,
	}

	opts := RecoveryDescriptorOptions{
		MaxRetries:        10,
		ResetRetriesAfter: successDuration - (time.Millisecond * 500),
	}
	descriptor := NewRecoveryDescriptorWithOptions(mock, opts)

	pc := make(chan Process)

	defer cancel()

	require.NoError(t, descriptor.Run(ctx, pc))
	require.Equal(t, 0, mock.count)
}
