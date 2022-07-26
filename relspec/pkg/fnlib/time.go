package fnlib

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/puppetlabs/leg/timeutil/pkg/clockctx"
)

func NowDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Retrieves the current system time in UTC" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 0 {
				return nil, &fn.ArityError{Wanted: []int{0}, Got: l}
			}

			return fn.InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
				return evaluate.StaticResult[T](clockctx.Clock(ctx).Now()), nil
			}), nil
		},
	}
}
