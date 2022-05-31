package fnlib_test

import (
	"context"
	"math/cmplx"
	"reflect"
	"testing"

	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/relspec"
	"github.com/stretchr/testify/require"
)

type test struct {
	Name                    string
	ObjectArg               any
	QArg                    any
	DefaultArg              any
	Expected                any
	ExpectedIncomplete      bool
	ExpectedPositionalError error
	ExpectedKeywordError    error
}

type tests []test

func (tts tests) RunAll(t *testing.T) {
	for _, tt := range tts {
		t.Run(tt.Name, tt.Run)
	}
}

func TestPath(t *testing.T) {
	tests{
		{
			Name: "path exists",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg:     "foo.bar",
			Expected: "baz",
		},
		{
			Name: "path does not exist",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg: "foo.quux",
			ExpectedPositionalError: &fn.PositionalArgError{
				Arg: 1,
				Cause: &evaluate.PathEvaluationError{
					Path: "foo",
					Cause: &evaluate.PathEvaluationError{
						Path:  "quux",
						Cause: &eval.UnknownKeyError{Key: "quux"},
					},
				},
			},
			ExpectedKeywordError: &fn.KeywordArgError{
				Arg: "object",
				Cause: &evaluate.PathEvaluationError{
					Path: "foo",
					Cause: &evaluate.PathEvaluationError{
						Path:  "quux",
						Cause: &eval.UnknownKeyError{Key: "quux"},
					},
				},
			},
		},
		{
			Name: "path exists with default",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg:       "foo.bar",
			DefaultArg: "quux",
			Expected:   "baz",
		},
		{
			Name: "path does not exist with default",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg:       "foo.quux",
			DefaultArg: 42,
			Expected:   42,
		},
		{
			Name: "query is not resolvable",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg:               map[string]any{"$type": "example"},
			ExpectedIncomplete: true,
		},
		{
			Name: "default is not resolvable",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			QArg:               "foo.quux",
			DefaultArg:         map[string]any{"$type": "example"},
			ExpectedIncomplete: true,
		},
		{
			Name:       "object is unsupported type (nil)",
			ObjectArg:  nil,
			QArg:       "foo.bar",
			DefaultArg: "bar",
			ExpectedPositionalError: &fn.PositionalArgError{
				Arg: 1,
				Cause: &evaluate.PathEvaluationError{
					Path:  "foo",
					Cause: &eval.UnsupportedValueTypeError{Value: nil, Field: "foo"},
				},
			},
			ExpectedKeywordError: &fn.KeywordArgError{
				Arg: "object",
				Cause: &evaluate.PathEvaluationError{
					Path:  "foo",
					Cause: &eval.UnsupportedValueTypeError{Value: nil, Field: "foo"},
				},
			},
		},
		{
			Name: "object is unsupported type (complex128)",
			ObjectArg: map[string]any{
				"foo": map[string]any{
					"bar": cmplx.Sqrt(12i),
				},
			},
			QArg:       "foo.bar",
			DefaultArg: "bar",
			ExpectedPositionalError: &fn.PositionalArgError{
				Arg: 1,
				Cause: &evaluate.PathEvaluationError{
					Path: "foo",
					Cause: &evaluate.PathEvaluationError{
						Path: "bar",
						Cause: &evaluate.UnsupportedValueError{
							Type: reflect.TypeOf(complex128(1)),
						},
					},
				},
			},
			ExpectedKeywordError: &fn.KeywordArgError{
				Arg: "object",
				Cause: &evaluate.PathEvaluationError{
					Path: "foo",
					Cause: &evaluate.PathEvaluationError{
						Path: "bar",
						Cause: &evaluate.UnsupportedValueError{
							Type: reflect.TypeOf(complex128(1)),
						},
					},
				},
			},
		},
	}.RunAll(t)
}

func (tt test) Run(t *testing.T) {
	desc := fnlib.PathDescriptor[testReferences]()

	t.Run(tt.Name, func(t *testing.T) {
		ctx := context.Background()

		ev := relspec.NewEvaluator[testReferences](
			relspec.WithMappingTypeResolvers[testReferences](map[string]relspec.MappingTypeResolver[testReferences]{
				"example": &testMappingTypeResolver{},
			}),
		)

		t.Run("positional", func(t *testing.T) {
			args := []any{tt.ObjectArg, tt.QArg}
			if tt.DefaultArg != nil {
				args = append(args, tt.DefaultArg)
			}

			invoker, err := fn.NewPositionalArgs(ev, args).InvokerFor(desc)
			require.NoError(t, err)

			r, err := invoker.Invoke(ctx)
			if tt.ExpectedPositionalError != nil {
				require.Equal(t, tt.ExpectedPositionalError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.ExpectedIncomplete, !r.Resolved())
				if !tt.ExpectedIncomplete {
					require.Equal(t, tt.Expected, r.Value)
				} else {
					expected := make([]any, len(args))
					for i, arg := range args {
						r, err := evaluate.EvaluateAll(ctx, ev, arg)
						require.NoError(t, err)
						expected[i] = r.Value
					}
					require.Equal(t, expected, r.Value)
				}
			}
		})

		t.Run("keyword", func(t *testing.T) {
			args := map[string]any{
				"object": tt.ObjectArg,
				"query":  tt.QArg,
			}
			if tt.DefaultArg != nil {
				args["default"] = tt.DefaultArg
			}

			invoker, err := fn.NewKeywordArgs(ev, args).InvokerFor(desc)
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			if tt.ExpectedKeywordError != nil {
				require.Equal(t, tt.ExpectedKeywordError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.ExpectedIncomplete, !r.Resolved())
				if !tt.ExpectedIncomplete {
					require.Equal(t, tt.Expected, r.Value)
				} else {
					expected := make(map[string]any)
					for k, arg := range args {
						r, err := evaluate.EvaluateAll(ctx, ev, arg)
						require.NoError(t, err)
						expected[k] = r.Value
					}
					require.Equal(t, expected, r.Value)
				}
			}
		})
	})
}
