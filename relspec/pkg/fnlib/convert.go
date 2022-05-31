package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/leg/relspec/pkg/convert"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func convertMarkdownInvoker[T ref.References[T], V any, A fn.ComparableArg[T], AS fn.Args[T, A, V]](args AS, to, content A) fn.Invoker[T] {
	return fn.EvaluatedInvoker[T, V](args, []A{to, content}, func(ctx context.Context, vs map[A]any) (m any, err error) {
		toV, found := vs[to].(string)
		if !found {
			return nil, to.WrapError(&fn.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf("")},
				Got:    reflect.TypeOf(vs[to]),
			})
		}

		switch md := vs[content].(type) {
		case string:
			r, err := convert.ConvertMarkdown(convert.ConvertType(toV), []byte(md))
			if err != nil {
				return nil, err
			}
			return string(r), nil
		default:
			return nil, content.WrapError(&fn.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf("")},
				Got:    reflect.TypeOf(md),
			})
		}
	})
}

func ConvertMarkdownDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Converts a string in markdown format to another applicable syntax" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Got: l}
			}

			to := args.Use(0)
			content := args.Use(1)

			return convertMarkdownInvoker[T, []any](args, to, content), nil
		},
		KeywordInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.KeywordArgs[T]) (fn.Invoker[T], error) {
			for _, arg := range []string{"to", "content"} {
				if !args.Has(arg) {
					return nil, &fn.KeywordArgError{Arg: arg, Cause: fn.ErrArgNotFound}
				}
			}

			to := args.Use("to")
			content := args.Use("content")

			return convertMarkdownInvoker[T, map[string]any](args, to, content), nil
		},
	}
}
