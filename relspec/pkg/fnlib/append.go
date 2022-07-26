package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func AppendDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Adds new items to a given array, returning a new array" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l < 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Variadic: true, Got: l}
			}

			s := args.UseAll()

			return fn.EvaluatedInvoker[T, []any](args, s, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (any, error) {
				base, ok := vs[s[0]].([]any)
				if !ok {
					return nil, s[0].WrapError(&fn.UnexpectedTypeError{
						Wanted: []reflect.Type{
							reflect.TypeOf([]any(nil)),
						},
						Got: reflect.TypeOf(vs[s[0]]),
					})
				}

				appended := append([]any{}, base...)
				for _, arg := range s[1:] {
					appended = append(appended, vs[arg])
				}
				return appended, nil
			}), nil
		},
	}
}
