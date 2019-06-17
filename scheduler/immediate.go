package scheduler

import (
	"context"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
)

type ImmediateDescriptor struct {
	process Process
}

func (id *ImmediateDescriptor) Run(ctx context.Context, pc chan<- Process) errawr.Error {
	pc <- id.process
	return nil
}

func NewImmediateDescriptor(process Process) *ImmediateDescriptor {
	return &ImmediateDescriptor{
		process: process,
	}
}
