package scheduler

import (
	"context"

	"github.com/puppetlabs/insights-stdlib/scheduler/errors"
)

type StartedLifecycle interface {
	Close(ctx context.Context) errors.Error
}

type Lifecycle interface {
	Start() StartedLifecycle
}
