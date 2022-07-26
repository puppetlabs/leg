package relspec

import (
	"context"

	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type templateEnvironment[T ref.References[T]] map[string]evaluate.Expandable[T]

var (
	_ eval.Indexable                           = templateEnvironment[ref.EmptyReferences](nil)
	_ evaluate.Expandable[ref.EmptyReferences] = templateEnvironment[ref.EmptyReferences](nil)
)

func (te templateEnvironment[_]) Index(ctx context.Context, idx any) (any, error) {
	k, err := eval.StringValue(idx)
	if err != nil {
		return nil, err
	}

	r, ok := te[k]
	if !ok {
		return nil, &eval.UnknownKeyError{Key: k}
	}

	return r, nil
}

func (te templateEnvironment[T]) Expand(ctx context.Context, depth int) (*evaluate.Result[T], error) {
	if depth == 0 {
		return evaluate.StaticResult[T](te), nil
	}

	vm := make(map[string]any, len(te))
	r := evaluate.StaticResult[T](vm)

	for key, v := range te {
		vr, err := v.Expand(ctx, depth-1)
		if err != nil {
			return nil, err
		}

		if vr == nil {
			continue
		}

		r.AddStep(vr)
		vm[key] = vr.Value
	}

	return r, nil
}

func evaluateTemplate[T ref.References[T]](o *Options[T]) func(ctx context.Context, s string, depth int, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
	ev := evaluate.DefaultEvaluator[T]()
	pl := pathlang.New[T](
		pathlang.WithFunctionMap[T]{Map: o.FunctionMap},
		pathlang.WithEagerEvaluation[T](o.EagerEvaluation),
	)
	env := templateEnvironment[T](o.TemplateEnvironment)

	return func(ctx context.Context, s string, depth int, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
		r, err := query.EvaluateQuery(ctx, ev, pl.Template, env, s)
		if err != nil {
			return nil, err
		}

		if !r.OK() {
			r.SetValue(s)
		}

		if !o.KeepSteps {
			r.Reduce()
		}

		return r, nil
	}
}
