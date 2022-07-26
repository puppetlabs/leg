package query

import (
	"context"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/gvalutil/pkg/template"
	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
	jsonpathtemplate "github.com/puppetlabs/leg/jsonutil/pkg/jsonpath/template"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Language[T ref.References[T]] func(ev evaluate.Evaluator[T]) gval.Language

func JSONPathLanguage[T ref.References[T]](ev evaluate.Evaluator[T]) gval.Language {
	return gval.NewLanguage(
		jsonpathtemplate.ExpressionLanguage(),
		gval.VariableSelector(jsonpath.VariableSelector(variableVisitor(ev))),
	)
}

var _ Language[ref.EmptyReferences] = JSONPathLanguage[ref.EmptyReferences]

func JSONPathTemplateLanguage[T ref.References[T]](ev evaluate.Evaluator[T]) gval.Language {
	return jsonpathtemplate.TemplateLanguage(
		jsonpathtemplate.WithExpressionLanguageVariableVisitor(variableVisitor(ev)),
		jsonpathtemplate.WithFormatter(template.StringFormatterFunc(func(ctx context.Context, v any) (string, error) {
			md := evaluate.MetadataFromContext[T](ctx)

			rv, err := evaluate.EvaluateAll(ctx, evaluate.CoalesceEvaluators(md.Evaluator, ev), v)
			if err != nil {
				return "", err
			}

			md.AddStep(rv)

			if rv.OK() {
				v = rv.Value
			}

			return jsonpathtemplate.DefaultFormatter.FormatString(ctx, v)
		})),
	)
}

var _ Language[ref.EmptyReferences] = JSONPathTemplateLanguage[ref.EmptyReferences]

func EvaluateGval[T ref.References[T]](ctx context.Context, ev evaluate.Evaluator[T], eval gval.Evaluable, parameter any) (*evaluate.Result[T], error) {
	// Push new metadata onto stack for tracking purposes.
	ctx, md := evaluate.ContextWithEmptyMetadata[T](ctx)
	md.SetAnnotation(QueryAnnotation)

	// Explicitly set this evaluator into the context for this eval call to make
	// sure it propagates.
	md.SetEvaluator(ev)

	v, err := eval(ctx, parameter)
	if err != nil {
		return nil, err
	}

	if !md.OK() {
		return evaluate.NewResult(md, v), nil
	}

	r, err := evaluate.EvaluateAll(ctx, md.Evaluator, v)
	if err != nil {
		return nil, err
	}

	md.AddStep(r, QueryResultAnnotation)
	md.MergeEvaluator(r.Metadata)

	return evaluate.NewResult(md, r.Value), nil
}

func EvaluateQuery[T ref.References[T]](ctx context.Context, ev evaluate.Evaluator[T], lang Language[T], tree any, query string) (*evaluate.Result[T], error) {
	path, err := lang(ev).NewEvaluable(query)
	if err != nil {
		return nil, err
	}

	return EvaluateGval(ctx, ev, path, tree)
}
