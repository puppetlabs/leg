package template

import (
	"context"
	"fmt"
	"strings"
	"text/scanner"
	"unicode/utf8"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/graph"
	"github.com/puppetlabs/leg/gvalutil/pkg/langctx"
	"github.com/puppetlabs/leg/gvalutil/pkg/langext"
)

var StopFlag = langctx.NewBoolFlag()

type Options struct {
	Joiner                     Joiner
	DelimitedLanguageFactories []DelimitedLanguageFactory
}

type Option interface {
	ApplyToOptions(target *Options)
}

func (o *Options) ApplyOptions(opts []Option) {
	for _, opt := range opts {
		opt.ApplyToOptions(o)
	}
}

func Language(opts ...Option) gval.Language {
	o := &Options{
		Joiner: DefaultStringJoiner,
	}
	o.ApplyOptions(opts)

	return langext.Late(func(outer gval.Language) gval.Language {
		g := graph.NewSimpleDirectedGraph()
		g.AddVertex("")

		for _, lf := range o.DelimitedLanguageFactories {
			lang := lf.DelimitedLanguage(outer)

			// No delimiter means the only valid language is the expression
			// language itself.
			if lang.Start == "" {
				return lang.Language
			}

			for i, w := 0, 0; i < len(lang.Start); i += w {
				_, w = utf8.DecodeRuneInString(lang.Start[i:])
				g.AddVertex(lang.Start[:i+w])
				_ = g.Connect(lang.Start[:i], lang.Start[:i+w])
			}

			g.AddVertex(lang)
			_ = g.Connect(lang.Start, lang)
		}

		var (
			langs []gval.Language
			inits string
		)

		edges, _ := g.OutgoingEdgesOf("")
		_ = edges.ForEach(func(edge graph.Edge) error {
			next, _ := graph.OppositeVertexOf(g, edge, "")

			if nt, ok := next.(string); ok {
				r, _ := utf8.DecodeLastRuneInString(nt)
				lang := delimLang(g, nt)

				langs = append(langs, gval.PrefixExtension(r, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
					ctx, stop := StopFlag.With(ctx)

					a, err := p.ParseSublanguage(ctx, lang)
					if err != nil {
						return nil, err
					}

					if *stop {
						return a, nil
					}

					b, err := p.ParseExpression(ctx)
					if err != nil {
						return nil, err
					}

					return func(ctx context.Context, parameter any) (any, error) {
						return o.Joiner.Join(ctx, a, b, parameter)
					}, nil
				}))
				inits += string(r)
			}

			return nil
		})

		return gval.NewLanguage(
			gval.NewLanguage(langs...),
			gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				p.SetWhitespace()
				p.SetMode(scanner.ScanIdents)
				p.SetIsIdentRuneFunc(func(ch rune, i int) bool { return ch > 0 && !strings.ContainsRune(inits, ch) })

				return p.ParseExpression(ctx)
			}),
			gval.PrefixExtension(scanner.Ident, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				a := p.Const(p.TokenText())

				b, err := p.ParseExpression(ctx)
				if err != nil {
					return nil, err
				}

				return func(ctx context.Context, parameter any) (any, error) {
					return o.Joiner.Join(ctx, a, b, parameter)
				}, nil
			}),
			gval.PrefixExtension(scanner.EOF, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				return p.Const(""), nil
			}),
		)
	})
}

func delimLang(g graph.DirectedGraph, prefix string) gval.Language {
	langs := make(map[rune]gval.Language)

	edges, _ := g.OutgoingEdgesOf(prefix)
	for _, edge := range edges.AsSlice() {
		next, _ := graph.OppositeVertexOf(g, edge, prefix)
		switch nt := next.(type) {
		case *DelimitedLanguage:
			// We can just return as-is here as we'll always pick the shortest
			// prefix for template strings.
			return gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				p.SetWhitespace()
				p.SetMode(0)

				eval, err := p.ParseSublanguage(ctx, nt.Language)
				if err != nil {
					return nil, err
				}

				for _, expected := range nt.End {
					actual := p.Scan()
					if expected != actual {
						return nil, p.Expected(fmt.Sprintf("template %q end", nt.Start), expected)
					}
				}

				return func(ctx context.Context, parameter any) (any, error) {
					v, err := eval(ctx, parameter)
					if err != nil {
						return nil, &EvaluationError{
							Start: nt.Start,
							Cause: err,
						}
					}

					return v, nil
				}, nil
			})
		case string:
			r, _ := utf8.DecodeLastRuneInString(nt)
			langs[r] = delimLang(g, nt)
		}
	}

	return gval.NewLanguage(
		gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
			p.SetWhitespace()
			p.SetMode(0)

			r := p.Scan()
			if r < 0 {
				return p.Const(prefix), nil
			} else if lang, ok := langs[r]; ok {
				return p.ParseSublanguage(ctx, lang)
			}
			return p.Const(prefix + string(r)), nil
		}),
	)
}
