package scheduler

import (
	"context"

	"github.com/puppetlabs/insights-stdlib/scheduler/errors"
)

type StartedParent struct {
	delegates []StartedLifecycle
}

func (sp *StartedParent) Close(ctx context.Context) errors.Error {
	err := errors.NewLifecycleCloseError()

	for _, d := range sp.delegates {
		if cerr := d.Close(ctx); cerr != nil {
			err = err.WithCause(cerr)
		}
	}

	switch len(err.Causes()) {
	case 0:
		return nil
	case 1:
		return err.Causes()[0]
	default:
		return err
	}
}

type Parent struct {
	delegates []Lifecycle
}

func (p *Parent) Start() StartedLifecycle {
	sd := make([]StartedLifecycle, len(p.delegates))
	for i, d := range p.delegates {
		sd[i] = d.Start()
	}

	return &StartedParent{delegates: sd}
}

func NewParent(delegates ...Lifecycle) *Parent {
	return &Parent{
		delegates: delegates,
	}
}
