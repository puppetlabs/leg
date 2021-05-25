package template_test

import (
	"context"
	"fmt"
	"testing"
	"text/scanner"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/gvalutil/pkg/template"
	"github.com/stretchr/testify/require"
)

func TestLanguage(t *testing.T) {
	tests := []struct {
		Name     string
		Language gval.Language
		Template string
		Input    interface{}
		Expected interface{}
	}{
		{
			Name: "simple",
			Language: template.Language(
				template.WithDelimitedLanguage{
					DelimitedLanguage: &template.DelimitedLanguage{
						Start:    "${",
						End:      "}",
						Language: gval.Full(),
					},
				},
			),
			Template: "hello ${foo}, friend",
			Input:    map[string]interface{}{"foo": "there"},
			Expected: "hello there, friend",
		},
		{
			Name: "multiple delimiters",
			Language: template.Language(
				template.WithDelimitedLanguage{
					DelimitedLanguage: &template.DelimitedLanguage{
						Start:    "${",
						End:      "}",
						Language: gval.Arithmetic(),
					},
				},
				template.WithDelimitedLanguage{
					DelimitedLanguage: &template.DelimitedLanguage{
						Start:    "#{",
						End:      "}",
						Language: gval.PropositionalLogic(),
					},
				},
			),
			Template: "math ${1 + 1}, logic #{true && false}",
			Expected: "math 2, logic false",
		},
		{
			Name: "nesting",
			Language: template.Language(
				template.WithDelimitedLanguageFactory{
					DelimitedLanguageFactory: template.DelimitedLanguageFactoryFunc(func(tpl gval.Language) *template.DelimitedLanguage {
						return &template.DelimitedLanguage{
							Start: "{",
							Language: gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
								if p.Scan() != scanner.Ident {
									return nil, p.Expected("template identifier")
								}

								switch p.TokenText() {
								case "dup":
									if p.Scan() != '}' {
										return nil, p.Expected("template identifier", '}')
									}

									eval, err := p.ParseSublanguage(ctx, tpl)
									if err != nil {
										return nil, err
									}

									if p.Scan() != scanner.Ident || p.TokenText() != "end" {
										return nil, p.Expected("template identifier")
									}

									if p.Scan() != '}' {
										return nil, p.Expected("template identifier", '}')
									}

									return func(ctx context.Context, parameter interface{}) (interface{}, error) {
										v, err := eval(ctx, parameter)
										if err != nil {
											return nil, err
										}

										return fmt.Sprintf("%v%v", v, v), nil
									}, nil
								case "end":
									template.StopFlag.MustSet(ctx)

									p.Camouflage("template identifier")
									return p.Const(""), nil
								default:
									return nil, p.Expected("template identifier")
								}
							}),
						}
					}),
				},
			),
			Template: "try not to {dup}dupl{dup}i{end}cate{end} words",
			Expected: "try not to dupliicatedupliicate words",
		},
		{
			Name: "custom joiner",
			Language: template.Language(
				template.WithJoiner{
					Joiner: template.NewStringJoiner(template.WithEmptyStringsEliminated(true)),
				},
				template.WithDelimitedLanguage{
					DelimitedLanguage: &template.DelimitedLanguage{
						Start:    "${",
						End:      "}",
						Language: gval.Full(),
					},
				},
			),
			Template: "${a + b}",
			Input:    map[string]interface{}{"a": 1, "b": 2},
			Expected: float64(3),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			eval, err := test.Language.NewEvaluable(test.Template)
			require.NoError(t, err)

			actual, err := eval(context.Background(), test.Input)
			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})
	}
}
