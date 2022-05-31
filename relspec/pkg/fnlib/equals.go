package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func EqualsDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Checks if the left side equals the right side" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Got: l}
			}

			lhs := args.Use(0)
			rhs := args.Use(1)

			return fn.EvaluatedInvoker[T, []any](args, []*fn.PositionalArg[T]{lhs, rhs}, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (m any, err error) {
				return reflect.DeepEqual(vs[lhs], vs[rhs]), nil
			}), nil
		},
	}
}

func NotEqualsDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Checks if the left side does not equal the right side" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Got: l}
			}

			lhs := args.Use(0)
			rhs := args.Use(1)

			return fn.EvaluatedInvoker[T, []any](args, []*fn.PositionalArg[T]{lhs, rhs}, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (m any, err error) {
				return !reflect.DeepEqual(vs[lhs], vs[rhs]), nil
			}), nil
		},
	}
}
