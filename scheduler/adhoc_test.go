package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/errawr-go/v2/pkg/testutil"
	"github.com/puppetlabs/leg/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestAdhocQueue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ad, as := scheduler.NewAdhocDescriptor()
	lc := scheduler.
		NewSegment(1, []scheduler.Descriptor{ad}).
		WithErrorBehavior(scheduler.ErrorBehaviorDrop)

	p1 := as.Submit(scheduler.DescribeProcessFunc("p1", func(ctx context.Context) error {
		return nil
	}))
	p2 := as.Submit(scheduler.DescribeProcessFunc("p2", func(ctx context.Context) error {
		return nil
	}))

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	defer func() {
		assert.NoError(t, scheduler.CloseWaitContext(ctx, slc))
		assert.Empty(t, slc.Errs())
	}()

	select {
	case err := <-p1:
		assert.NoError(t, err, "p1 returned error")
	case <-ctx.Done():
		assert.Fail(t, "p1 context expired")
	}

	select {
	case err := <-p2:
		assert.NoError(t, err, "p2 returned error")
	case <-ctx.Done():
		assert.Fail(t, "p2 context expired")
	}
}

func TestAdhocErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ad, as := scheduler.NewAdhocDescriptor()
	lc := scheduler.
		NewSegment(1, []scheduler.Descriptor{ad}).
		WithErrorBehavior(scheduler.ErrorBehaviorDrop)

	p1 := as.Submit(scheduler.DescribeProcessFunc("p1", func(ctx context.Context) error {
		return testutil.NewStubError("p1")
	}))
	p2 := as.Submit(scheduler.DescribeProcessFunc("p2", func(ctx context.Context) error {
		panic(testutil.NewStubError("p2"))
	}))

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	defer func() {
		assert.NoError(t, scheduler.CloseWaitContext(ctx, slc))
		assert.Empty(t, slc.Errs())
	}()

	select {
	case err := <-p1:
		switch rerr := err.(type) {
		case errawr.Error:
			assert.Equal(t, "p1", rerr.Code())
		default:
			assert.Fail(t, "p1 did not return an error")
		}
	case <-ctx.Done():
		assert.Fail(t, "p1 context expired")
	}

	select {
	case err := <-p2:
		switch rerr := err.(type) {
		case errawr.Error:
			assert.Equal(t, "p2", rerr.Code())
		default:
			assert.Fail(t, "p2 did not return an error")
		}
	case <-ctx.Done():
		assert.Fail(t, "p2 context expired")
	}
}

func TestAdhocSubmissionMultipleLifecycles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ad, as := scheduler.NewAdhocDescriptor()

	slc1 := scheduler.
		NewSegment(1, []scheduler.Descriptor{ad}).
		WithErrorBehavior(scheduler.ErrorBehaviorDrop).
		Start(scheduler.LifecycleStartOptions{})
	defer func() {
		assert.NoError(t, scheduler.CloseWaitContext(ctx, slc1))
		assert.Empty(t, slc1.Errs())
	}()

	slc2 := scheduler.
		NewSegment(1, []scheduler.Descriptor{ad}).
		WithErrorBehavior(scheduler.ErrorBehaviorDrop).
		Start(scheduler.LifecycleStartOptions{})
	defer func() {
		assert.NoError(t, scheduler.CloseWaitContext(ctx, slc2))
		assert.Empty(t, slc2.Errs())
	}()

	var rc int32
	p1 := as.Submit(scheduler.DescribeProcessFunc("p1", func(ctx context.Context) error {
		// Should run exactly once, even if there are two active schedulers.
		atomic.AddInt32(&rc, 1)
		return nil
	}))

	assert.NoError(t, <-p1)
	assert.Equal(t, int32(1), rc)

	// Close one of the segments; we should still be able to run jobs on the
	// second one.
	assert.NoError(t, scheduler.CloseWaitContext(ctx, slc1))
	assert.Empty(t, slc1.Errs())

	p2 := as.Submit(scheduler.DescribeProcessFunc("p2", func(ctx context.Context) error {
		atomic.AddInt32(&rc, 1)
		return nil
	}))

	assert.NoError(t, <-p2)
	assert.Equal(t, int32(2), rc)

	// Close the second segment. Now a process should just go into the queue as
	// there is nothing to run it.
	assert.NoError(t, scheduler.CloseWaitContext(ctx, slc2))
	assert.Empty(t, slc2.Errs())

	as.Submit(scheduler.DescribeProcessFunc("p3", func(ctx context.Context) error {
		atomic.AddInt32(&rc, 1)
		return nil
	}))

	assert.Equal(t, 1, as.QueueLen())
	assert.Equal(t, int32(2), rc)
}
