package scheduler

import (
	"context"
	"time"

	errawr "github.com/puppetlabs/errawr-go"
)

type IntervalDescriptor struct {
	interval time.Duration
	process  Process
}

func (id *IntervalDescriptor) runOnce(ctx context.Context, pc chan<- Process) (bool, errawr.Error) {
	pc <- id.process

	t := time.NewTimer(id.interval)
	defer t.Stop()

	select {
	case <-ctx.Done():
		// XXX: Return error?
		return false, nil
	case <-t.C:
	}

	return true, nil
}

func (id *IntervalDescriptor) Run(ctx context.Context, pc chan<- Process) errawr.Error {
	for {
		if cont, err := id.runOnce(ctx, pc); err != nil {
			return err
		} else if !cont {
			break
		}
	}

	return nil
}

func NewIntervalDescriptor(interval time.Duration, process Process) *IntervalDescriptor {
	return &IntervalDescriptor{
		interval: interval,
		process:  process,
	}
}
