package pathlang

import (
	"context"
	"unicode"

	"github.com/PaesslerAG/gval"
	"github.com/generikvault/gvalstrings"
	"github.com/puppetlabs/leg/gvalutil/pkg/template"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Language[T ref.References[T]] struct {
	o *Options[T]
}

func (l *Language[T]) Expression(ev evaluate.Evaluator[T]) gval.Language {
	return gval.NewLanguage(
		base,
		ident(ev, l.o),
		operators(ev, l.o),
		gval.VariableSelector(variableSelector(ev)),
		gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
			p.SetIsIdentRuneFunc(func(ch rune, i int) bool {
				switch {
				case ch == '_' || unicode.IsLetter(ch):
					return true
				case i == 0:
					return false
				case ch == '-' || unicode.IsDigit(ch):
					return true
				default:
					return false
				}
			})

			eval, err := p.ParseExpression(ctx)
			if err != nil {
				return nil, err
			}

			return func(ctx context.Context, parameter any) (any, error) {
				md := evaluate.MetadataFromContext[T](ctx)

				ev := evaluate.CoalesceEvaluators(md.Evaluator, ev)

				r, err := query.EvaluateGval(ctx, ev, eval, parameter)
				if err != nil {
					return nil, err
				}

				md.AddStep(r, ExpressionAnnotation)
				md.SetEvaluator(evaluate.CoalesceEvaluators(r.Evaluator, ev))
				return r.Value, nil
			}, nil
		}),
	)
}

func (l *Language[T]) Template(ev evaluate.Evaluator[T]) gval.Language {
	return template.Language(
		template.WithJoiner{
			Joiner: template.NewStringJoiner(template.WithEmptyStringsEliminated(true)),
		},
		template.WithDelimitedLanguage{
			DelimitedLanguage: &template.DelimitedLanguage{
				Start:    "${",
				End:      "}",
				Language: l.Expression(ev),
			},
		},
	)
}

func New[T ref.References[T]](opts ...Option[T]) *Language[T] {
	o := &Options[T]{
		FunctionMap: fn.NewMap(map[string]fn.Descriptor[T]{}),
	}
	o.ApplyOptions(opts)

	if o.EagerEvaluation {
		o.FunctionMap = fn.NewEagerEvaluationMap(o.FunctionMap)
	}

	return &Language[T]{o: o}
}

func DefaultExpression[T ref.References[T]]() gval.Language {
	return New[T]().Expression(evaluate.DefaultEvaluator[T]())
}

func DefaultTemplate[T ref.References[T]]() gval.Language {
	return New[T]().Template(evaluate.DefaultEvaluator[T]())
}

var base = gval.NewLanguage(
	gval.Base(),
	gval.JSON(),
	gvalstrings.SingleQuoted(),
	gval.Constant("null", nil),
	gval.PrefixExtension('$', func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		return identVar(ctx, p, nil)
	}),
	gval.PostfixOperator("|>", func(c context.Context, p *gval.Parser, pre gval.Evaluable) (gval.Evaluable, error) {
		post, err := p.ParseExpression(c)
		if err != nil {
			return nil, err
		}
		return func(c context.Context, v any) (any, error) {
			v, err := pre(c, v)
			if err != nil {
				return nil, err
			}
			return post(c, v)
		}, nil
	}),
)
