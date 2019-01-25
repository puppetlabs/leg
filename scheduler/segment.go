package scheduler

import (
	"context"

	"github.com/puppetlabs/insights-instrumentation/alerts"
	"github.com/puppetlabs/insights-instrumentation/alerts/trackers"
	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/insights-stdlib/request"
	"github.com/puppetlabs/insights-stdlib/scheduler/errors"
)

type StartedSegment struct {
	ctx      context.Context
	cancel   context.CancelFunc
	pc       chan Process
	capturer trackers.Capturer
}

func (ss *StartedSegment) run(process Process) {
	req := request.New()

	ctx := request.NewContext(ss.ctx, req)
	ctx = logging.NewContext(ctx, "request", req.Identifier)

	log(ctx).Debug("process running", "description", process.Description())

	err := ss.capturer.Try(ctx, func(ctx context.Context) {
		if err := process.Run(ctx); err != nil {
			log(ctx).Warn("process failed", "error", err)
			ss.capturer.Capture(err).AsWarning().Report(ctx)
		} else {
			log(ctx).Debug("process complete")
		}
	})
	if err != nil {
		log(ctx).Crit("process panic()!", "error", err)
	}
}

func (ss *StartedSegment) describe(desc Descriptor) {
	if err := desc.Run(ss.ctx, ss.pc); err != nil {
		log(ss.ctx).Warn("descriptor ended with error", "error", err)
	}
}

func (ss *StartedSegment) supervise() {
	for {
		select {
		case <-ss.ctx.Done():
			return
		case process := <-ss.pc:
			ss.run(process)
		}
	}
}

func (ss *StartedSegment) Close(ctx context.Context) errors.Error {
	ss.cancel()

	// TODO: Wait for processes to terminate?
	return nil
}

type Segment struct {
	concurrency int
	descriptors []Descriptor
	capturer    trackers.Capturer
}

func (s *Segment) WithCapturer(capturer trackers.Capturer) *Segment {
	s.capturer = capturer
	return s
}

func (s *Segment) Start() StartedLifecycle {
	ctx, cancel := context.WithCancel(context.Background())
	pc := make(chan Process)

	capturer := s.capturer
	if capturer == nil {
		capturer = alerts.NewAlerts(alerts.NoDelegate, alerts.Options{}).NewCapturer()
	}

	ctx = trackers.NewContextWithCapturer(ctx, capturer)

	ss := &StartedSegment{
		ctx:      ctx,
		cancel:   cancel,
		pc:       pc,
		capturer: capturer,
	}
	for _, desc := range s.descriptors {
		go ss.describe(desc)
	}
	for i := 0; i < s.concurrency; i++ {
		go ss.supervise()
	}

	return ss
}

func NewSegment(concurrency int, descriptors []Descriptor) *Segment {
	return &Segment{
		concurrency: concurrency,
		descriptors: descriptors,
	}
}
