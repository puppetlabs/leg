package scheduler_test

import (
	"context"
	"testing"
	"time"

	errawr "github.com/puppetlabs/errawr-go"
	"github.com/puppetlabs/insights-stdlib/scheduler"
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
