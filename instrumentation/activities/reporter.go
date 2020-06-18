package activities

import (
	"context"

	"github.com/puppetlabs/horsehead/v2/instrumentation/activities/activity"
	"github.com/puppetlabs/horsehead/v2/scheduler"
)

type Reporter struct {
	ch        chan activity.Activity
	delegates []Delegate
}

func (r *Reporter) doReport(act activity.Activity) {
	for _, d := range r.delegates {
		// we don't immediately abandon ship in case only one of the delegates has
		// an issue, thus we cant still get activity data to some services.
		if err := d.Report(act); err != nil {
			log(context.Background()).Error("failed to report activity", "error", err, "user_id", act.UserID)
		}
	}
}

func (r *Reporter) Report(ctx context.Context, act activity.Activity) error {
	select {
	case r.ch <- act:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Reporter) Run(ctx context.Context, pc chan<- scheduler.Process) error {
	for {
		select {
		case <-ctx.Done():
			// we want to drain the remaining activities to report before shutting
			// down.
			for {
				select {
				case act := <-r.ch:
					r.doReport(act)
				default:
					break
				}
			}

			return nil
		case act := <-r.ch:
			r.doReport(act)
		}
	}
}

// Scheduler instantiates a scheduler that can be used for scheduling this work
// directly.
func (r *Reporter) Scheduler() scheduler.Lifecycle {
	return scheduler.NewSegment(1, []scheduler.Descriptor{r})
}

// AddDelegate adds a new delegate to the set of delegates to report activities
// to.
func (r *Reporter) AddDelegate(d Delegate) {
	r.delegates = append(r.delegates, d)
}

// NewReporter creates a new reporter for asynchronously reporting activities
// to services. The process for pumping activities to services is started
// automatically as part of this.
func NewReporter() *Reporter {
	r := &Reporter{
		ch: make(chan activity.Activity, 8),
	}

	return r
}
