package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/puppetlabs/horsehead/instrumentation/alerts"
	"github.com/puppetlabs/horsehead/instrumentation/alerts/trackers"
	"github.com/puppetlabs/horsehead/logging"
	"github.com/puppetlabs/horsehead/request"
	"github.com/puppetlabs/horsehead/scheduler/errors"
)

type StartedSegment struct {
	ctx    context.Context
	cancel context.CancelFunc
	pc     chan Process

	// Error handling and reporting.
	errorHandler segmentErrorHandler
	capturer     trackers.Capturer

	// Used to execute a close of the Process channel exactly once.
	done sync.Once

	// These channels are used to track whether descriptor and supervisor
	// goroutines have exited.
	descriptors []chan struct{}
	supervisors []chan struct{}
}

func (ss *StartedSegment) run(process Process) {
	req := request.New()

	ctx := request.NewContext(ss.ctx, req)
	ctx = logging.NewContext(ctx, "request", req.Identifier)

	log(ctx).Debug("process running", "description", process.Description())

	err := ss.capturer.Try(ctx, func(ctx context.Context) {
		if err := process.Run(ctx); err != nil {
			log(ctx).Warn("process failed", "error", err)

			ss.errorHandler.ForProcess(req, process, err)
			ss.capturer.Capture(err).AsWarning().Report(ctx)
		} else {
			log(ctx).Debug("process complete")
		}
	})
	if err != nil {
		log(ctx).Crit("process panic()!", "error", err)

		switch et := err.(type) {
		case errors.Error:
			ss.errorHandler.ForProcess(req, process, et)
		case error:
			ss.errorHandler.ForProcess(req, process, errors.NewProcessPanicError().WithCause(et))
		default:
			ss.errorHandler.ForProcess(req, process, errors.NewProcessPanicError().WithCause(fmt.Errorf("panic: %+v", et)))
		}
	}
}

func (ss *StartedSegment) describe(i int, desc Descriptor) {
	defer func() {
		close(ss.descriptors[i])
	}()

	if err := desc.Run(ss.ctx, ss.pc); err != nil {
		log(ss.ctx).Warn("descriptor ended with error", "error", err)

		ss.errorHandler.ForDescriptor(i, desc, err)
	}
}

func (ss *StartedSegment) supervise(i int) {
	defer func() {
		close(ss.supervisors[i])
	}()

	for {
		process, ok := <-ss.pc
		if !ok {
			return
		}

		ss.run(process)
	}
}

func (ss *StartedSegment) waitErrorHandler() {
	// If we have a potentially terminating error handler, we will force
	// ourselves to close if it returns.
	select {
	case <-ss.ctx.Done():
	case <-ss.errorHandler.Ch():
		ss.cancel()
	}
}

func (ss *StartedSegment) Wait(ctx context.Context) errors.Error {
	for _, ch := range ss.descriptors {
		select {
		case <-ctx.Done():
			return errors.NewLifecycleTimeoutError().WithCause(ctx.Err())
		case <-ch:
		}
	}

	// At this point, all descriptors have exited, so we can close the process
	// channel. Then we will wait for supervisors to clean up and exit.
	ss.done.Do(func() { close(ss.pc) })

	for _, ch := range ss.supervisors {
		select {
		case <-ctx.Done():
			return errors.NewLifecycleTimeoutError().WithCause(ctx.Err())
		case <-ch:
		}
	}

	// Now we will have collected all errors from running processes. If our
	// error handler recorded any errors, let's return them now.
	if errs := ss.errorHandler.Errors(); len(errs) > 0 {
		rerr := errors.NewLifecycleExecutionError()

		for _, err := range errs {
			rerr = rerr.WithCause(err)
		}

		return rerr
	}

	return nil
}

func (ss *StartedSegment) Close(ctx context.Context) errors.Error {
	// The important bit: cancel the context used by everything under this
	// segment.
	ss.cancel()

	// The supervise code above is simply waiting for the contex to cancel. The
	// descriptors need to respond to the context and exit on their own.
	if err := ss.Wait(ctx); err != nil {
		return errors.NewLifecycleCloseError().WithCause(err)
	}

	return nil
}

type Segment struct {
	concurrency   int
	descriptors   []Descriptor
	errorBehavior SegmentErrorBehavior
	capturer      trackers.Capturer
}

func (s *Segment) WithErrorBehavior(errorBehavior SegmentErrorBehavior) *Segment {
	s.errorBehavior = errorBehavior
	return s
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
		ctx:    ctx,
		cancel: cancel,
		pc:     pc,

		errorHandler: s.errorBehavior.handler(),
		capturer:     capturer,

		descriptors: make([]chan struct{}, len(s.descriptors)),
		supervisors: make([]chan struct{}, s.concurrency),
	}

	for i, desc := range s.descriptors {
		ss.descriptors[i] = make(chan struct{})

		go ss.describe(i, desc)
	}

	for i := 0; i < s.concurrency; i++ {
		ss.supervisors[i] = make(chan struct{})

		go ss.supervise(i)
	}

	go ss.waitErrorHandler()

	return ss
}

func NewSegment(concurrency int, descriptors []Descriptor) *Segment {
	return &Segment{
		concurrency:   concurrency,
		descriptors:   descriptors,
		errorBehavior: SegmentErrorBehaviorCollect,
	}
}
