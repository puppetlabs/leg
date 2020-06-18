package activities

import (
	"context"
	"sync"
	"time"

	"github.com/puppetlabs/horsehead/v2/instrumentation/activities/activity"
)

var (
	DefaultReportTimeout = 2 * time.Second
)

type Reporter struct {
	// ReportTimeout is the amount of time clients will wait when attempting to
	// report activities to activity services.
	ReportTimeout time.Duration

	shutdown  chan struct{}
	open      bool
	mut       sync.RWMutex
	ch        chan activity.Activity
	delegates []Delegate
}

func (r *Reporter) doReport(act activity.Activity) error {
	var lastErr error

	for _, d := range r.delegates {
		// we don't immediately abandon ship in case only one of the delegates has
		// an issue, thus we cant still get activity data to some services.
		//
		// not really sure what to do if multiple errors occur with delegates and
		// it's probably not really worth worrying about. one possibility would be
		// to log the relevant errors but meh...
		if err := d.Report(act); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (r *Reporter) doReportLoop() {
	for {
		// TODO: The reporter can be shut down without delivering all of it's
		// pending activities. We should find a way around that.
		act := <-r.ch

		if err := r.doReport(act); err != nil {
			log(context.Background()).Error("failed to report activity", "error", err, "user_id", act.UserID)
		}

		select {
		case <-r.shutdown:
			break
		default:
			// noop
		}
	}
}

func (r *Reporter) Report(act activity.Activity) error {
	r.mut.RLock()
	defer r.mut.RUnlock()

	if !r.open {
		return ErrReporterShutdown
	}

	select {
	case r.ch <- act:
		return nil
	case <-time.After(r.ReportTimeout):
		return ErrReportTimeout
	}
}

// Close terminates delivery of activity information to external services.
func (r *Reporter) Close() error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.shutdown <- struct{}{}
	r.open = false

	close(r.shutdown)
	close(r.ch)

	for _, d := range r.delegates {
		d.Close()
	}

	return nil
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
		ReportTimeout: DefaultReportTimeout,
		ch:            make(chan activity.Activity, 8),
		shutdown:      make(chan struct{}),
		open:          true,
	}

	go r.doReportLoop()

	return r
}
