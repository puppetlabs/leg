package pathlang

import (
	"context"
	"fmt"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func variableSelector[T ref.References[T]](ev evaluate.Evaluator[T]) func(path gval.Evaluables) gval.Evaluable {
	return func(path gval.Evaluables) gval.Evaluable {
		return func(ctx context.Context, iv any) (rv any, err error) {
			var parents []any
			defer func() {
				if err != nil {
					for i := len(parents) - 1; i >= 0; i-- {
						err = &evaluate.PathEvaluationError{
							Path:  fmt.Sprintf("%v", parents[i]),
							Cause: err,
						}
					}
				}
			}()

			// We're going to add traversal directly on top of the existing
			// metadata, so it completely retains any context as it evaluates
			// its own steps.
			md := evaluate.MetadataFromContext[T](ctx)

			// We're going to choose the base evaluator to use for evaluating
			// both keys and values. For values, we'll successively update the
			// evaluator each time we nest further into the data. However, key
			// evaluation is always relative to the base. Consider the contrived
			// example:
			//
			//   data[funcA(data.foo)][funcB(data.bar)]
			//
			// Clearly, for non-constant keys, each time the key is referenced
			// it is relative to the same root being indexed into.

			// This is the base evaluator for this traversal operation.
			ev := evaluate.CoalesceEvaluators(md.Evaluator, ev)

			// Performs one iteration of traversal in isolation.
			traverse := func(cr *evaluate.Result[T], key gval.Evaluable) (*evaluate.Result[T], error) {
				r := evaluate.EmptyResult[T]()

				kr, err := query.EvaluateGval(ctx, ev, key, iv)
				if err != nil {
					return nil, err
				}

				r.AddStep(kr, query.TraversalKeyAnnotation)

				if !r.OK() {
					return r, nil
				}

				parents = append(parents, kr.Value)

				switch vt := cr.Value.(type) {
				case eval.Indexable:
					v, err := vt.Index(ctx, kr.Value)
					if err != nil {
						return nil, err
					}

					r.SetValue(v)
				default:
					nr, err := cr.Evaluator.Evaluate(ctx, vt, 1)
					if err != nil {
						return nil, err
					}

					r.AddStep(nr, query.TraversalDataAnnotation)

					if !r.OK() {
						return r, nil
					}

					v, err := eval.Select(ctx, nr.Value, kr.Value)
					if err != nil {
						return nil, err
					}

					r.MergeEvaluator(nr.Metadata)
					r.SetValue(v)
				}

				return r, nil
			}

			// This is the incremental result, updated each iteration through
			// the traversal loop.
			cr := evaluate.StaticResult[T](iv)
			cr.SetEvaluator(ev)

			for _, key := range path {
				r, err := traverse(cr, key)
				if err != nil {
					return nil, err
				}

				md.AddStep(r, query.TraversalAnnotation)

				if !r.OK() {
					return nil, nil
				}

				r.MergeEvaluatorIfNotSet(cr.Metadata)
				cr = r
			}

			r, err := cr.Evaluator.Evaluate(ctx, cr.Value, 1)
			if err != nil {
				return nil, err
			}

			md.AddStep(r, query.TraversalAnnotation)

			if !r.OK() {
				return nil, nil
			}

			md.SetEvaluator(evaluate.CoalesceEvaluators(r.Evaluator, cr.Evaluator))
			return r.Value, nil
		}
	}
}
