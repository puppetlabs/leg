package fnlib

import (
	"context"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func ExistsDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string { return "Determines whether a value is set" },
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l != 1 {
				return nil, &fn.ArityError{Wanted: []int{1}, Got: l}
			}

			arg := args.Use(0)

			return fn.InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
				md := evaluate.EmptyMetadata[T]()

				r, err := arg.Expand(ctx, evaluate.DepthFull)
				if errmark.Matches(err, notExistsRule) {
					return evaluate.StaticResult[T](false), nil
				} else if err != nil {
					return nil, err
				}

				md.AddStep(r, arg.Annotation())

				if !r.OK() {
					return evaluate.NewResult(md, args.ToValue()), nil
				}

				return evaluate.NewResult(md, true), nil
			}), nil
		},
	}
}
