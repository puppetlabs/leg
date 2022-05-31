package relspec

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func NewEvaluator[T ref.References[T]](opts ...Option[T]) evaluate.Evaluator[T] {
	o := &Options[T]{
		FunctionMap: fnlib.Library[T](),
	}
	o.ApplyOptions(opts)

	fns := &evaluate.VisitorFuncs[T]{
		VisitMapFunc:    evaluateMap(o),
		VisitStringFunc: evaluateTemplate(o),
	}
	if !o.KeepSteps {
		fns.VisitSliceFunc = func(ctx context.Context, s []any, depth int, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
			r, err := evaluate.DefaultVisitor[T]().VisitSlice(ctx, s, depth, next)
			if err != nil {
				return nil, err
			}

			r.Reduce()

			return r, nil
		}
	}

	return evaluate.NewEvaluator[T](fns)
}
