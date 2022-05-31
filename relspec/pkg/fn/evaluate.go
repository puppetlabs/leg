package fn

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func EvaluateInvocation[T ref.References[T], A Arg[T], V any, AS Args[T, A, V]](ctx context.Context, m Map[T], name string, args AS) (*evaluate.Result[T], error) {
	descriptor, err := m.Descriptor(name)
	if err != nil {
		return nil, &InvocationError{Name: name, Cause: err}
	}

	invoker, err := args.InvokerFor(descriptor)
	if err != nil {
		return nil, &InvocationError{Name: name, Cause: err}
	}

	r, err := invoker.Invoke(ctx)
	if err != nil {
		return nil, &InvocationError{Name: name, Cause: err}
	}

	return r, nil
}
