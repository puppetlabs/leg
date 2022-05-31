package fnlib

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func JSONMarshalDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Marshals arbitrary data into a JSON-encoded string" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 1 {
				return nil, &fn.ArityError{Wanted: []int{1}, Got: l}
			}

			arg := args.Use(0)

			return fn.EvaluatedInvoker[T, []any](args, []*fn.PositionalArg[T]{arg}, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (any, error) {
				b, err := json.Marshal(vs[arg])
				if err != nil {
					return nil, arg.WrapError(err)
				}

				return string(b), nil
			}), nil
		},
	}
}

func JSONUnmarshalDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Unmarshals a JSON-encoded string into the specification" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 1 {
				return nil, &fn.ArityError{Wanted: []int{1}, Got: l}
			}

			arg := args.Use(0)

			return fn.EvaluatedInvoker[T, []any](args, []*fn.PositionalArg[T]{arg}, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (m any, err error) {
				var b []byte

				switch v := vs[arg].(type) {
				case string:
					b = []byte(v)
				default:
					return nil, arg.WrapError(&fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(arg),
					})
				}

				err = json.Unmarshal(b, &m)
				return
			}), nil
		},
	}
}
