package scheduler_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/puppetlabs/leg/scheduler"
	"github.com/stretchr/testify/require"
)

func TestRestartDescriptorPropagatesErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	descriptor, _ := scheduler.NewRestartableDescriptor(scheduler.DescriptorFunc(func(ctx context.Context, pc chan<- scheduler.Process) error {
		return errors.New("boom")
	}))

	require.EqualError(t, descriptor.Run(ctx, make(chan<- scheduler.Process)), "boom")
}

func TestRestartDescriptorRestartableAfterFinishing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var i int32
	sig := make(chan int32)
	descriptor, restart := scheduler.NewRestartableDescriptor(scheduler.DescriptorFunc(func(ctx context.Context, pc chan<- scheduler.Process) error {
		select {
		case sig <- atomic.AddInt32(&i, 1):
		case <-ctx.Done():
			require.Fail(t, "context timed out waiting for test")
		}
		return nil
	}))

	lc := scheduler.NewScheduler(scheduler.ManySchedulableDescriptor(
		[]scheduler.Descriptor{descriptor},
		make(chan<- scheduler.Process),
	)).Start(scheduler.LifecycleStartOptions{})

	select {
	case ti := <-sig:
		require.Equal(t, int32(1), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for initial descriptor invocation")
	}

	restart()

	select {
	case ti := <-sig:
		require.Equal(t, int32(2), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for restarted descriptor invocation")
	}

	require.NoError(t, scheduler.CloseWaitContext(ctx, lc))
}

func TestRestartDescriptorRestartableWhileExecuting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var i int32
	sig := make(chan int32, 1)
	descriptor, restart := scheduler.NewRestartableDescriptor(scheduler.DescriptorFunc(func(ctx context.Context, pc chan<- scheduler.Process) error {
		sig <- atomic.AddInt32(&i, 1)
		<-ctx.Done()
		sig <- atomic.AddInt32(&i, 1)
		return nil
	}))

	lc := scheduler.NewScheduler(scheduler.ManySchedulableDescriptor(
		[]scheduler.Descriptor{descriptor},
		make(chan<- scheduler.Process),
	)).Start(scheduler.LifecycleStartOptions{})

	// Initial read before we try to restart the process.
	select {
	case ti := <-sig:
		require.Equal(t, int32(1), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for initial descriptor invocation before restart")
	}

	restart()

	// Second read immediately after restarting.
	select {
	case ti := <-sig:
		require.Equal(t, int32(2), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for initial descriptor invocation after restart")
	}

	// Third read for the restarted descriptor.
	select {
	case ti := <-sig:
		require.Equal(t, int32(3), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for restart descriptor invocation before close")
	}

	require.NoError(t, scheduler.CloseWaitContext(ctx, lc))

	// Fourth read after we close.
	select {
	case ti := <-sig:
		require.Equal(t, int32(4), ti)
	case <-ctx.Done():
		require.Fail(t, "context timed out waiting for restart descriptor invocation after close")
	}
}

func TestRestartDescriptorInMultipleSchedulers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var i int32
	sig := make(chan int32)
	descriptor, restart := scheduler.NewRestartableDescriptor(scheduler.DescriptorFunc(func(ctx context.Context, pc chan<- scheduler.Process) error {
		select {
		case sig <- atomic.AddInt32(&i, 1):
		case <-ctx.Done():
			require.Fail(t, "context timed out waiting for test")
		}
		return nil
	}))

	n := 3
	all := make([]scheduler.Lifecycle, n)
	for i := 0; i < n; i++ {
		all[i] = scheduler.NewScheduler(scheduler.ManySchedulableDescriptor(
			[]scheduler.Descriptor{descriptor},
			make(chan<- scheduler.Process),
		))
	}

	lc := scheduler.NewParent(all...).Start(scheduler.LifecycleStartOptions{})

	// Wait for descriptors to run once.
	for i := 0; i < n; i++ {
		select {
		case <-sig:
		case <-ctx.Done():
			require.Fail(t, "context timed out waiting for initial descriptor invocations")
		}
	}

	restart()

	// Wait for them to run again.
	for i := 0; i < n; i++ {
		select {
		case <-sig:
		case <-ctx.Done():
			require.Fail(t, "context timed out waiting for restarted descriptor invocations")
		}
	}

	require.NoError(t, scheduler.CloseWaitContext(ctx, lc))
	require.Equal(t, int32(n*2), i)
}
