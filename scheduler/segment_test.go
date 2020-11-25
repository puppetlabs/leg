package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/errawr-go/v2/pkg/testutil"
	"github.com/puppetlabs/leg/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestSegmentProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var ok bool

	lc := scheduler.NewSegment(1, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) error {
			ok = true
			return nil
		})),
	})

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	defer slc.Close()

	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.True(t, ok, "process did not run")
}

func TestSegmentClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(1, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})),
	})

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.CloseWaitContext(ctx, slc))
	assert.Empty(t, slc.Errs())
}

func TestErrorBehaviorCollect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) error {
			return testutil.NewStubError("boom 1")
		})),
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) error {
			return testutil.NewStubError("boom 2")
		})),
	}).WithErrorBehavior(scheduler.ErrorBehaviorCollect)

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.Len(t, slc.Errs(), 2)
}

func TestErrorBehaviorTerminate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewIntervalDescriptor(10*time.Second, scheduler.ProcessFunc(func(ctx context.Context) error {
			return testutil.NewStubError("boom")
		})),
		scheduler.NewIntervalDescriptor(10*time.Second, scheduler.ProcessFunc(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
				t.Fatal("terminating error behavior did not terminate segment")
			}

			return nil
		})),
	}).WithErrorBehavior(scheduler.ErrorBehaviorTerminate)

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.NotEmpty(t, slc.Errs())
}

func TestErrorBehaviorPanic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewIntervalDescriptor(10*time.Second, scheduler.ProcessFunc(func(ctx context.Context) error {
			panic("no")
		})),
		scheduler.NewIntervalDescriptor(10*time.Second, scheduler.ProcessFunc(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
				t.Fatal("terminating error behavior did not terminate segment")
			}

			return nil
		})),
	}).WithErrorBehavior(scheduler.ErrorBehaviorTerminate)

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.NotEmpty(t, slc.Errs())
}

func TestErrorBehaviorDrop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) error {
			return testutil.NewStubError("boom")
		})),
	}).WithErrorBehavior(scheduler.ErrorBehaviorDrop)

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.Empty(t, slc.Errs())
}

func TestProcessErrorBehaviorTerminate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// It is conceivable that someone would (possibly erroneously, but whatever)
	// configure a segment to collect descriptor errors but terminate on process
	// errors.
	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewIntervalDescriptor(2*time.Second, scheduler.ProcessFunc(func(ctx context.Context) error {
			return testutil.NewStubError("boom")
		})),
	}).
		WithDescriptorErrorBehavior(scheduler.ErrorBehaviorCollect).
		WithProcessErrorBehavior(scheduler.ErrorBehaviorTerminate)

	slc := lc.Start(scheduler.LifecycleStartOptions{})
	assert.NoError(t, scheduler.WaitContext(ctx, slc))
	assert.Len(t, slc.Errs(), 1)
}
