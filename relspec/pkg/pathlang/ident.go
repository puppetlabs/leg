package pathlang

import (
	"context"
	"text/scanner"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func evalIdentFuncPositional[T ref.References[T]](ev evaluate.Evaluator[T], o *Options[T], name string, args []gval.Evaluable) gval.Evaluable {
	return func(ctx context.Context, parameter any) (any, error) {
		l := make([]any, len(args))
		for i, eval := range args {
			l[i] = evaluate.MemoizedExpandable(query.EvalExpandable[T](eval, parameter))
		}

		r, err := fn.EvaluateInvocation[T, *fn.PositionalArg[T], []any](ctx, o.FunctionMap, name, fn.NewPositionalArgs(ev, l))
		if err != nil {
			return nil, err
		}

		r.SetAnnotation(fn.PositionalInvocationAnnotation(name))

		if !r.OK() {
			r.SetValue(nil)
		}

		return evaluate.StaticExpandable(r), nil
	}
}

func identFuncPositional[T ref.References[T]](ctx context.Context, p *gval.Parser, ev evaluate.Evaluator[T], o *Options[T], name string, args []gval.Evaluable) (gval.Evaluable, error) {
	arg, err := p.ParseExpression(ctx)
	if err != nil {
		return nil, err
	}

	args = append(args, arg)

	switch p.Scan() {
	case ')':
		return evalIdentFuncPositional(ev, o, name, args), nil
	case ',':
		// Another arg.
		return identFuncPositional(ctx, p, ev, o, name, args)
	default:
		return nil, p.Expected("function argument", ',', ')')
	}
}

func evalIdentFuncKeyword[T ref.References[T]](ev evaluate.Evaluator[T], o *Options[T], name string, args map[string]gval.Evaluable) gval.Evaluable {
	return func(ctx context.Context, parameter any) (any, error) {
		m := make(map[string]any, len(args))
		for k, eval := range args {
			m[k] = evaluate.MemoizedExpandable(query.EvalExpandable[T](eval, parameter))
		}

		r, err := fn.EvaluateInvocation[T, *fn.KeywordArg[T], map[string]any](ctx, o.FunctionMap, name, fn.NewKeywordArgs(ev, m))
		if err != nil {
			return nil, err
		}

		r.SetAnnotation(fn.PositionalInvocationAnnotation(name))

		if !r.OK() {
			r.SetValue(nil)
		}

		return evaluate.StaticExpandable(r), nil
	}
}

func identFuncKeyword[T ref.References[T]](ctx context.Context, p *gval.Parser, ev evaluate.Evaluator[T], o *Options[T], name string, args map[string]gval.Evaluable) (gval.Evaluable, error) {
	if p.Scan() != scanner.Ident {
		return nil, p.Expected("function argument name", scanner.Ident)
	}

	key := p.TokenText()

	if p.Scan() != ':' {
		return nil, p.Expected("function argument name", ':')
	}

	arg, err := p.ParseExpression(ctx)
	if err != nil {
		return nil, err
	}

	if args == nil {
		args = make(map[string]gval.Evaluable)
	}

	args[key] = arg

	switch p.Scan() {
	case ')':
		return evalIdentFuncKeyword(ev, o, name, args), nil
	case ',':
		// Another arg.
		return identFuncKeyword(ctx, p, ev, o, name, args)
	default:
		return nil, p.Expected("function argument", ',', ')')
	}
}

func identFunc[T ref.References[T]](ctx context.Context, p *gval.Parser, ev evaluate.Evaluator[T], o *Options[T], name string) (gval.Evaluable, error) {
	switch p.Scan() {
	case ')':
		// Positional invocation with no arguments.
		return evalIdentFuncPositional(ev, o, name, nil), nil
	case scanner.Ident:
		candidate := p.TokenText()
		switch p.Peek() {
		case ':':
			p.Scan() // == ':'

			arg, err := p.ParseExpression(ctx)
			if err != nil {
				return nil, err
			}

			switch p.Scan() {
			case ',':
				// Keyword invocation with several arguments.
				return identFuncKeyword(ctx, p, ev, o, name, map[string]gval.Evaluable{candidate: arg})
			case ')':
				// Keyword invocation with a single argument.
				return evalIdentFuncKeyword(ev, o, name, map[string]gval.Evaluable{candidate: arg}), nil
			default:
				return nil, p.Expected("function argument", ',', ')')
			}
		default:
			// Assume positional invocation with data argument.
			p.Camouflage("function argument", ':', ',', ')')
			arg, err := p.ParseExpression(ctx)
			if err != nil {
				return nil, err
			}

			switch p.Scan() {
			case ',':
				return identFuncPositional(ctx, p, ev, o, name, []gval.Evaluable{arg})
			case ')':
				return evalIdentFuncPositional(ev, o, name, []gval.Evaluable{arg}), nil
			default:
				return nil, p.Expected("function argument", ',', ')')
			}
		}
	default:
		// Must be positional invocation.
		p.Camouflage("function argument", ':', ',', ')')
		return identFuncPositional(ctx, p, ev, o, name, nil)
	}
}

func identVar(ctx context.Context, p *gval.Parser, vars []gval.Evaluable) (gval.Evaluable, error) {
	switch scanIgnoreFloats(p) {
	case '.':
		switch scanIgnoreFloats(p) {
		case scanner.Ident:
			return identVar(ctx, p, append(vars, p.Const(p.TokenText())))
		default:
			p.Camouflage("field")
			eval, err := p.ParseNextExpression(ctx)
			if err != nil {
				return nil, err
			}

			return identVar(ctx, p, append(vars, eval))
		}
	case '[':
		key, err := p.ParseExpression(ctx)
		if err != nil {
			return nil, err
		}

		switch p.Scan() {
		case ']':
			return identVar(ctx, p, append(vars, key))
		default:
			return nil, p.Expected("key", ']')
		}
	default:
		p.Camouflage("variable", '.', '[')
		return p.Var(vars...), nil
	}
}

func ident[T ref.References[T]](ev evaluate.Evaluator[T], o *Options[T]) gval.Language {
	return gval.NewLanguage(
		gval.PrefixMetaPrefix(scanner.Ident, func(ctx context.Context, p *gval.Parser) (call string, alternative func() (gval.Evaluable, error), err error) {
			call = p.TokenText()
			alternative = func() (gval.Evaluable, error) {
				switch scanIgnoreFloats(p) {
				case '(':
					return identFunc(ctx, p, ev, o, call)
				default:
					p.Camouflage("identifier", '.', '[', '(')
					return identVar(ctx, p, []gval.Evaluable{p.Const(call)})
				}
			}
			return
		}),
	)
}

func scanIgnoreFloats(p *gval.Parser) rune {
	p.SetMode((scanner.GoTokens | scanner.ScanInts) &^ scanner.ScanFloats)
	defer p.SetMode(scanner.GoTokens)
	return p.Scan()
}
