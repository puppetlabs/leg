package fnlib_test

import (
	"context"
	"testing"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	desc := fnlib.MergeDescriptor[ref.EmptyReferences]()

	tt := []struct {
		Name          string
		Req           func() (fn.Invoker[ref.EmptyReferences], error)
		Expected      any
		ExpectedError string
	}{
		{
			Name: "positional",
			Req: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewPositionalArgs(
					evaluate.DefaultEvaluator[ref.EmptyReferences](),
					[]any{
						map[string]any{"a": "b", "c": map[string]any{"d": "e"}},
						map[string]any{"a": "overwritten", "c": map[string]any{"f": "added"}},
					},
				).InvokerFor(desc)
			},
			Expected: map[string]any{
				"a": "overwritten",
				"c": map[string]any{"d": "e", "f": "added"},
			},
		},
		{
			Name: "keyword deep",
			Req: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(
					evaluate.DefaultEvaluator[ref.EmptyReferences](),
					map[string]any{
						"mode": "deep",
						"objects": []any{
							map[string]any{"a": "b", "c": map[string]any{"d": "e"}},
							map[string]any{"a": "overwritten", "c": map[string]any{"f": "added"}},
						},
					},
				).InvokerFor(desc)
			},
			Expected: map[string]any{
				"a": "overwritten",
				"c": map[string]any{"d": "e", "f": "added"},
			},
		},
		{
			Name: "keyword shallow",
			Req: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(
					evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
						"mode": "shallow",
						"objects": []any{
							map[string]any{"a": "b", "c": map[string]any{"d": "e"}},
							map[string]any{"a": "overwritten", "c": map[string]any{"f": "overwritten"}},
						},
					},
				).InvokerFor(desc)
			},
			Expected: map[string]any{
				"a": "overwritten",
				"c": map[string]any{"f": "overwritten"},
			},
		},
		{
			Name: "invalid mode",
			Req: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(
					evaluate.DefaultEvaluator[ref.EmptyReferences](),
					map[string]any{
						"mode":    "secret",
						"objects": []any{},
					},
				).InvokerFor(desc)
			},
			ExpectedError: `arg "mode": unexpected value "secret", wanted one of "deep" or "shallow"`,
		},
		{
			Name: "merge candidate is not a map",
			Req: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewPositionalArgs(
					evaluate.DefaultEvaluator[ref.EmptyReferences](),
					[]any{
						map[string]any{"a": "b", "c": map[string]any{"d": "e"}},
						"hi",
					},
				).InvokerFor(desc)
			},
			ExpectedError: `arg 2: unexpected type string (wanted map[string]interface {})`,
		},
	}
	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			invoker, err := test.Req()
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			if test.ExpectedError != "" {
				require.EqualError(t, err, test.ExpectedError)
			} else {
				require.NoError(t, err)

				require.True(t, r.OK())
				require.Equal(t, test.Expected, r.Value)
			}
		})
	}
}
