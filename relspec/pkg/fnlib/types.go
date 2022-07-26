package fnlib

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func ToStringDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Converts arbitrary scalar input data to a string" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 1 {
				return nil, &fn.ArityError{Wanted: []int{1}, Got: l}
			}

			arg := args.Use(0)

			return fn.EvaluatedInvoker[T, []any](args, []*fn.PositionalArg[T]{arg}, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (any, error) {
				v, err := toString(vs[arg])
				if err != nil {
					return nil, arg.WrapError(err)
				}

				return v, nil
			}), nil
		},
	}
}
