package fnlib

import (
	"context"
	"strings"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func ConcatDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Concatenates string arguments into a single string" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if args.Len() == 0 {
				return fn.StaticInvoker[T](""), nil
			}

			s := args.UseAll()

			return fn.EvaluatedInvoker[T, []any](args, s, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (m any, err error) {
				strs := make([]string, len(s))
				for i, arg := range s {
					vt, err := toString(vs[arg])
					if err != nil {
						return nil, arg.WrapError(err)
					}

					strs[i] = vt
				}

				return strings.Join(strs, ""), nil
			}), nil
		},
	}
}
