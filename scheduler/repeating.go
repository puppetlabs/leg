package scheduler

import (
	"context"
	"fmt"

	errawr "github.com/puppetlabs/errawr-go"
)

type repeatingProcess struct {
	ch       chan struct{}
	delegate Process
}

func (rp *repeatingProcess) Description() string {
	return fmt.Sprintf("(Repeating) %s", rp.delegate.Description())
}

func (rp *repeatingProcess) Run(ctx context.Context) errawr.Error {
	defer close(rp.ch)

	return rp.delegate.Run(ctx)
}

type RepeatingDescriptor struct {
	process Process
}

func (rd *RepeatingDescriptor) runOnce(ctx context.Context, pc chan<- Process) bool {
	ch := make(chan struct{})

	pc <- &repeatingProcess{ch: ch, delegate: rd.process}

	select {
	case <-ctx.Done():
		return false
	case <-ch:
	}

	return true
}

func (rd *RepeatingDescriptor) Run(ctx context.Context, pc chan<- Process) errawr.Error {
	for {
		if !rd.runOnce(ctx, pc) {
			break
		}
	}

	return nil
}

func NewRepeatingDescriptor(process Process) *RepeatingDescriptor {
	return &RepeatingDescriptor{
		process: process,
	}
}
