package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/errawr-go/v2/pkg/testutil"
	"github.com/puppetlabs/horsehead/scheduler"
	"github.com/puppetlabs/horsehead/scheduler/errors"
	"github.com/stretchr/testify/assert"
)

func TestSegmentProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ok bool

	lc := scheduler.NewSegment(1, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			ok = true
			return nil
		})),
	})

	slc := lc.Start()
	defer slc.Close(ctx)

	assert.NoError(t, slc.Wait(ctx))
	assert.True(t, ok, "process did not run")
}

func TestSegmentClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(1, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			<-ctx.Done()
			return nil
		})),
	})

	slc := lc.Start()
	assert.NoError(t, slc.Close(ctx))
}

func TestSegmentErrorBehaviorCollect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			return testutil.NewStubError("boom 1")
		})),
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			return testutil.NewStubError("boom 2")
		})),
	}).WithErrorBehavior(scheduler.SegmentErrorBehaviorCollect)

	slc := lc.Start()
	err := slc.Wait(ctx)
	assert.True(t, errors.IsLifecycleExecutionError(err))
	assert.Len(t, err.Causes(), 2)
}

func TestSegmentErrorBehaviorTerminate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			return testutil.NewStubError("boom")
		})),
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
				t.Fatal("terminating error behavior did not terminate segment")
			}

			return nil
		})),
	}).WithErrorBehavior(scheduler.SegmentErrorBehaviorTerminate)

	slc := lc.Start()
	assert.True(t, errors.IsLifecycleExecutionError(slc.Wait(ctx)))
}

func TestSegmentErrorBehaviorPanic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			panic("no")
		})),
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
				t.Fatal("terminating error behavior did not terminate segment")
			}

			return nil
		})),
	}).WithErrorBehavior(scheduler.SegmentErrorBehaviorTerminate)

	slc := lc.Start()
	assert.True(t, errors.IsLifecycleExecutionError(slc.Wait(ctx)))
}

func TestSegmentErrorBehaviorDrop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lc := scheduler.NewSegment(2, []scheduler.Descriptor{
		scheduler.NewImmediateDescriptor(scheduler.ProcessFunc(func(ctx context.Context) errawr.Error {
			return testutil.NewStubError("boom")
		})),
	}).WithErrorBehavior(scheduler.SegmentErrorBehaviorDrop)

	slc := lc.Start()
	assert.NoError(t, slc.Wait(ctx))
}
