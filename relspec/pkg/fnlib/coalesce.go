package fnlib

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func CoalesceDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string {
			return "Finds and returns the first resolvable non-null argument, returning null otherwise"
		},
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			s := args.UseAll()

			return fn.InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
				for _, arg := range s {
					r, err := arg.Expand(ctx, evaluate.DepthFull)
					if err != nil {
						return nil, err
					}

					if r.OK() && r.Value != nil {
						md := evaluate.EmptyMetadata[T]()
						md.AddStep(r, arg.Annotation())

						return evaluate.NewResult(md, r.Value), nil
					}
				}

				return evaluate.EmptyResult[T](), nil
			}), nil
		},
	}
}
