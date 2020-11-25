package scheduler_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/puppetlabs/leg/scheduler"
)

func ExampleScheduler() {
	c := int32(5)

	worker := scheduler.SchedulableFunc(func(ctx context.Context, er scheduler.ErrorReporter) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if atomic.AddInt32(&c, -1) < 0 {
				return
			}

			fmt.Println("Mmm... pie.")
		}
	})

	l := scheduler.NewScheduler(scheduler.OneSchedulable(worker))

	sl := l.Start(scheduler.LifecycleStartOptions{})
	<-sl.Done()

	// Output: Mmm... pie.
	// Mmm... pie.
	// Mmm... pie.
	// Mmm... pie.
	// Mmm... pie.
}
