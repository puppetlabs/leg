package scheduler

import (
	"context"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
)

type Process interface {
	Description() string
	Run(ctx context.Context) errawr.Error
}

type ProcessFunc func(ctx context.Context) errawr.Error

func (ProcessFunc) Description() string {
	return "<anonymous>"
}

func (p ProcessFunc) Run(ctx context.Context) errawr.Error {
	return p(ctx)
}

type Descriptor interface {
	Run(ctx context.Context, pc chan<- Process) errawr.Error
}
