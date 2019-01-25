package scheduler

import (
	"context"

	errawr "github.com/puppetlabs/errawr-go"
)

type Process interface {
	Description() string
	Run(ctx context.Context) errawr.Error
}

type Descriptor interface {
	Run(ctx context.Context, pc chan<- Process) errawr.Error
}
