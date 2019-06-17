package scheduler

import (
	"context"
	"reflect"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
)

// RecoveryDescriptor wraps a given descriptor so that it restarts if the
// descriptor itself fails. This is useful for descriptors that work off of
// external information (APIs, events, etc.).
type RecoveryDescriptor struct {
	delegate Descriptor
}

func (rd *RecoveryDescriptor) runOnce(ctx context.Context, pc chan<- Process) (bool, errawr.Error) {
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

func (rd *RecoveryDescriptor) Run(ctx context.Context, pc chan<- Process) errawr.Error {
	for {
		if cont, err := rd.runOnce(ctx, pc); err != nil {
			return err
		} else if !cont {
			break
		}
	}

	return nil
}

func NewRecoveryDescriptor(delegate Descriptor) *RecoveryDescriptor {
	return &RecoveryDescriptor{delegate: delegate}
}
