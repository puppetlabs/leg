package query

import (
	"context"

	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func variableVisitor[T ref.References[T]](ev evaluate.Evaluator[T]) jsonpath.VariableVisitor {
	return jsonpath.VariableVisitorFuncs{
		VisitChildFunc: func(ctx context.Context, parameter any, key any, next func(ctx context.Context, pv jsonpath.PathValue) error) error {
			md := evaluate.MetadataFromContext[T](ctx)

			ev := evaluate.CoalesceEvaluators(md.Evaluator, ev)

			cr := evaluate.EmptyResult[T]()
			cr.AddStep(evaluate.StaticResult[T](key), TraversalKeyAnnotation)

			// We need to evaluate the base parameter before indexing in. This
			// is because the base parameter could be itself a $type, $encoding,
			// etc.
			nr, err := ev.Evaluate(ctx, parameter, 1)
			if err != nil {
				return err
			}

			cr.AddStep(nr, TraversalDataAnnotation)
			nr.SetEvaluator(evaluate.CoalesceEvaluators(nr.Evaluator, ev))

			err = jsonpath.DefaultVariableVisitor().VisitChild(ctx, nr.Value, key, func(ctx context.Context, pv jsonpath.PathValue) error {
				// Expand just this value without recursing.
				r, err := nr.Evaluator.Evaluate(ctx, pv.Value, 1)
				if err != nil {
					return err
				}

				cr.AddStep(r, TraversalDataAnnotation)
				cr.SetValue(r.Value)

				if !cr.OK() {
					return nil
				}

				pv.Value = cr.Value
				return next(ctx, pv)
			})
			if err != nil {
				return err
			}

			md.AddStep(cr, TraversalAnnotation)

			return nil
		},
	}
}
