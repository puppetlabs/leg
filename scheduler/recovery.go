package scheduler

import (
	"context"
	"reflect"
)

// RecoveryDescriptor wraps a given descriptor so that it restarts if the
// descriptor itself fails. This is useful for descriptors that work off of
// external information (APIs, events, etc.).
type RecoveryDescriptor struct {
	delegate Descriptor
}

var _ Descriptor = &RecoveryDescriptor{}

func (rd *RecoveryDescriptor) runOnce(ctx context.Context, pc chan<- Process) (bool, error) {
	err := rd.delegate.Run(ctx, pc)

	select {
	case <-ctx.Done():
		return false, err
	default:
	}

	if err != nil {
		log(ctx).Warn("restarting failing descriptor", "descriptor", reflect.TypeOf(rd.delegate).String(), "error", err)
	}

	return true, nil
}

// Run delegates work to another descriptor, catching any errors are restarting
// the descriptor immediately if an error occurs. It never returns an error. It
// only terminates when the context is done.
func (rd *RecoveryDescriptor) Run(ctx context.Context, pc chan<- Process) error {
	for {
		if cont, err := rd.runOnce(ctx, pc); err != nil {
			return err
		} else if !cont {
			break
		}
	}

	return nil
}

// NewRecoveryDescriptor creates a new recovering descriptor wrapping the given
// delegate descriptor.
func NewRecoveryDescriptor(delegate Descriptor) *RecoveryDescriptor {
	return &RecoveryDescriptor{delegate: delegate}
}
