package pathlang_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/puppetlabs/leg/timeutil/pkg/clockctx"
	"github.com/stretchr/testify/require"
	k8sclock "k8s.io/apimachinery/pkg/util/clock"
)

type testID struct {
	Name string
}

func (ti testID) Less(other testID) bool {
	return ti.Name < other.Name
}

var (
	errNotFound = errors.New("not found")
)

type testReferences = *ref.Log[testID]

func TestExpressions(t *testing.T) {
	ctx := context.Background()

	now := time.Date(2021, time.May, 31, 0, 0, 0, 0, time.UTC)
	fc := k8sclock.NewFakeClock(now)
	ctx = clockctx.WithClock(ctx, k8sext.NewClock(fc))

	input := map[string]any{
		"a": []any{1, 2, []any{3, 4}},
		"b": 10,
		"c": map[string]any{
			"x": 1,
			"y": 2,
			"x y": map[string]any{
				"z": 3,
			},
			"false": 4,
		},
		"d": evaluate.StaticExpandable(
			evaluate.ContextualizedResult(evaluate.NewMetadata(
				ref.InitialLog(ref.Errored(testID{Name: "foo"}, errNotFound)),
			)),
		),
		"e": evaluate.StaticExpandable(
			evaluate.NewResult(evaluate.NewMetadata(ref.InitialLog(ref.OK(testID{Name: "bar"}))), true),
		),
	}

	tests := []struct {
		Name               string
		EagerEvaluation    bool
		Expression         string
		Expected           any
		ExpectedReferences testReferences
		ExpectedError      string
	}{
		{
			Name:       "map path with bracket syntax",
			Expression: "c['x y']['z']",
			Expected:   3,
		},
		{
			Name:       "map path with dot syntax",
			Expression: "c.x",
			Expected:   1,
		},
		{
			Name:       "map path with dot syntax and constant string",
			Expression: "c.'x y'",
			Expected: map[string]any{
				"z": 3,
			},
		},
		{
			Name:       "map path with dot syntax and constant string traversal",
			Expression: "c.'x y'.z",
			Expected:   3,
		},
		{
			Name:       "map path with reserved identifier",
			Expression: "c.false",
			Expected:   4,
		},
		{
			Name:       "map path with operators",
			Expression: "c.('x ' + 'y').z",
			Expected:   3,
		},
		{
			Name:       "array path with bracket syntax",
			Expression: "a[0]",
			Expected:   1,
		},
		{
			Name:       "array path with dot syntax",
			Expression: "a.0",
			Expected:   1,
		},
		{
			Name:       "array path with dot syntax traversal with dot syntax",
			Expression: "a.2.1",
			Expected:   4,
		},
		{
			Name:       "array path with dot syntax traversal with bracket syntax",
			Expression: "a.2[1]",
			Expected:   4,
		},
		{
			Name:          "invalid map path",
			Expression:    "a.x",
			ExpectedError: `path "a.x": unexpected string index "x" for slice, must be convertible to int: strconv.ParseInt: parsing "x": invalid syntax`,
		},
		{
			Name:       "root",
			Expression: "$",
			Expected: map[string]any{
				"a": []any{1, 2, []any{3, 4}},
				"b": 10,
				"c": map[string]any{
					"x": 1,
					"y": 2,
					"x y": map[string]any{
						"z": 3,
					},
					"false": 4,
				},
				"d": nil,
				"e": true,
			},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "foo"}, errNotFound),
				ref.OK(testID{Name: "bar"}),
			),
		},
		{
			Name:       "map path from root",
			Expression: "$.c.x",
			Expected:   1,
		},
		{
			Name:          "invalid character after root",
			Expression:    "$$",
			ExpectedError: "unexpected \"$\" while scanning",
		},
		{
			Name:       "map path from root with dot syntax and constant string",
			Expression: "$.c.'x y'.z",
			Expected:   3,
		},
		{
			Name:       "map path from root with bracket syntax",
			Expression: "$['c'].x",
			Expected:   1,
		},
		{
			Name:          "invalid path",
			Expression:    "$.",
			ExpectedError: "unexpected",
		},
		{
			Name:       "simple arithmetic",
			Expression: "c.x + c.y",
			Expected:   float64(3),
		},
		{
			Name:       "parentheses",
			Expression: "(c.x + c.y) * a[1]",
			Expected:   float64(6),
		},
		{
			Name:       "map creation",
			Expression: `{'foo': a[0] + a[1], 'bar': a[1] + a[2][0]}`,
			Expected:   map[string]any{"foo": float64(3), "bar": float64(5)},
		},
		{
			Name:       "array creation",
			Expression: `["foo", a[0], c.y]`,
			Expected:   []any{"foo", 1, 2},
		},
		{
			Name:       "function call with no arguments",
			Expression: "now()",
			Expected:   now,
		},
		{
			Name:       "function call with one positional argument",
			Expression: `jsonUnmarshal('{"foo": "bar"}')`,
			Expected:   map[string]any{"foo": "bar"},
		},
		{
			Name:       "function call with more than one positional argument",
			Expression: "concat('a', b, 'c')",
			Expected:   "a10c",
		},
		{
			Name:       "function call with one keyword argument",
			Expression: "merge(objects: [{'a': 5}, {'b': b}])",
			Expected:   map[string]any{"a": float64(5), "b": 10},
		},
		{
			Name:       "function call with multiple keyword arguments",
			Expression: "path(object: $, query: 'c.y')",
			Expected:   2,
		},
		{
			Name:       "call stack",
			Expression: "path(merge({'a': 5}, {'b': b}), 'b')",
			Expected:   10,
		},
		{
			Name:       "pipe",
			Expression: `jsonUnmarshal('{"x": {"y": "z"}}') |> x.y`,
			Expected:   "z",
		},
		{
			Name:       "pipe scope",
			Expression: `b * (jsonUnmarshal('{"x": 20}') |> x) + b`,
			Expected:   float64(210),
		},
		{
			Name:       "pipe dot equivalence",
			Expression: `jsonUnmarshal('{"x": {"y": "z"}}') |> x |> y`,
			Expected:   "z",
		},
		{
			Name:       "resolvable coalesce evaluation",
			Expression: `coalesce(c.x, d)`,
			Expected:   1,
		},
		{
			Name:       "unresolvable coalesce evaluation",
			Expression: `coalesce(null, d)`,
			Expected:   nil,
		},
		{
			Name:            "resolvable coalesce with eager evaluation",
			EagerEvaluation: true,
			Expression:      `coalesce(c.x, d)`,
			Expected:        1,
			ExpectedReferences: ref.InitialLog(
				ref.Unused(ref.Errored(testID{Name: "foo"}, errNotFound)),
			),
		},
		{
			Name:       "resolvable path query",
			Expression: `path($, 'c.x')`,
			Expected:   1,
		},
		{
			Name:               "unresolvable path query",
			Expression:         `path($, 'd.g')`,
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "foo"}, errNotFound)),
		},
		{
			Name:               "unresolvable exists evaluation",
			Expression:         `exists(d.g)`,
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "foo"}, errNotFound)),
		},
		{
			Name:       "resolvable exists evaluation for nonexistent key",
			Expression: `exists(c.nope)`,
			Expected:   false,
		},
		{
			Name:       "resolvable exists evaluation for existent key",
			Expression: `exists(c.x)`,
			Expected:   true,
		},
		{
			Name:               "binary operator short-circuiting without eager evaluation",
			Expression:         `e || d`,
			Expected:           true,
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name:            "binary operator short-circuiting with eager evaluation",
			EagerEvaluation: true,
			Expression:      `e || d`,
			Expected:        true,
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "bar"}),
				ref.Unused(ref.Errored(testID{Name: "foo"}, errNotFound)),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			pl := pathlang.New[testReferences](
				pathlang.WithFunctionMap[testReferences]{Map: fnlib.Library[testReferences]()},
				pathlang.WithEagerEvaluation[testReferences](test.EagerEvaluation),
			)

			t.Run("expression", func(t *testing.T) {
				ctx, md := evaluate.ContextWithEmptyMetadata[*ref.Log[testID]](ctx)

				r, err := pl.Expression(evaluate.DefaultEvaluator[testReferences]()).EvaluateWithContext(ctx, test.Expression, input)
				if test.ExpectedError != "" {
					require.NotNil(t, err)
					require.Contains(t, err.Error(), test.ExpectedError)
				} else {
					require.NoError(t, err)
					require.Equal(t, test.Expected, r)

					refs := test.ExpectedReferences
					if refs == nil {
						refs = ref.NewLog[testID]()
					}

					require.Equal(t, refs, md.References)
				}
			})

			t.Run("template", func(t *testing.T) {
				ctx, md := evaluate.ContextWithEmptyMetadata[*ref.Log[testID]](ctx)

				r, err := pl.Template(evaluate.DefaultEvaluator[testReferences]()).EvaluateWithContext(ctx, `${`+test.Expression+`}`, input)
				if test.ExpectedError != "" {
					require.NotNil(t, err)
					require.Contains(t, err.Error(), test.ExpectedError)
				} else {
					require.NoError(t, err)
					require.Equal(t, test.Expected, r)

					refs := test.ExpectedReferences
					if refs == nil {
						refs = ref.NewLog[testID]()
					}

					require.Equal(t, refs, md.References)
				}
			})
		})
	}
}

func TestIdent(t *testing.T) {
	input := map[string]any{
		"a-b": 42,
		"x": map[string]any{
			"y-z": "foo",
		},
	}

	tests := []struct {
		Name          string
		Expression    string
		Expected      any
		ExpectedError string
	}{
		{
			Name:       "dash at beginning",
			Expression: "-a-b",
			Expected:   float64(-42),
		},
		{
			Name:       "dash in middle",
			Expression: "a-b",
			Expected:   42,
		},
		{
			Name:       "dash in middle with subtraction",
			Expression: "a-b - a-b",
			Expected:   float64(0),
		},
		{
			Name:       "dash in traversal",
			Expression: "x.y-z",
			Expected:   "foo",
		},
		{
			Name:          "nonexistent",
			Expression:    "x-y-z",
			ExpectedError: `path "x-y-z": unknown key x-y-z`,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("expression", func(t *testing.T) {
				r, err := pathlang.DefaultExpression[ref.EmptyReferences]().Evaluate(test.Expression, input)
				if test.ExpectedError != "" {
					require.NotNil(t, err)
					require.Contains(t, err.Error(), test.ExpectedError)
				} else {
					require.NoError(t, err)
					require.Equal(t, test.Expected, r)
				}
			})

			t.Run("template", func(t *testing.T) {
				r, err := pathlang.DefaultTemplate[ref.EmptyReferences]().Evaluate(`${`+test.Expression+`}`, input)
				if test.ExpectedError != "" {
					require.NotNil(t, err)
					require.Contains(t, err.Error(), test.ExpectedError)
				} else {
					require.NoError(t, err)
					require.Equal(t, test.Expected, r)
				}
			})
		})
	}
}
