package scheduler

import (
	"context"

	errawr "github.com/puppetlabs/errawr-go"
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
