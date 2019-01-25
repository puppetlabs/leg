package scheduler

import (
	"context"

	"github.com/puppetlabs/insights-stdlib/scheduler/errors"
)

// StartedLifecycle represents a fully configured, operating scheduler.
type StartedLifecycle interface {
	// Wait returns when all the descriptors in the lifecycle have completed and
	// all the processes created by those descriptors have completed.
	//
	// If the context ends before the lifecycle completes, this method returns
	// lis_lifecycle_timeout_error.
	Wait(ctx context.Context) errors.Error

	// Close terminates descriptors, dropping any processes emitted by those
	// descriptors, and asks any running processes to terminate.
	//
	// If the context ends before the lifecycle completes, the lifecycle gives
	// up on terminating those processes and returns lis_lifecycle_close_error.
	Close(ctx context.Context) errors.Error
}

// Lifecycle represents a partially or fully configured scheduler instance.
// Starting a lifecycle will dispatch the descriptors attached to the given
// lifecycle.
type Lifecycle interface {
	Start() StartedLifecycle
}
