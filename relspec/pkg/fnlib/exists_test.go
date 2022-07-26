package fnlib_test

import (
	"context"
	"testing"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	tests := []struct {
		Name          string
		Expr          string
		Param         any
		Expected      bool
		ExpectedError string
	}{
		{
			Name:     "basic",
			Expr:     "foo",
			Param:    map[string]any{"foo": "bar"},
			Expected: true,
		},
		{
			Name:     "missing key",
			Expr:     "foo.bar",
			Param:    map[string]any{"foo": map[string]any{"baz": "quux"}},
			Expected: false,
		},
		{
			Name:     "missing key at depth",
			Expr:     "foo.bar.baz",
			Param:    map[string]any{"foo": map[string]any{"baz": "quux"}},
			Expected: false,
		},
		{
			Name:     "index out of bounds",
			Expr:     "foo[1].bar",
			Param:    map[string]any{"foo": []any{map[string]any{"baz": "quux"}}},
			Expected: false,
		},
		{
			Name:          "invalid type",
			Expr:          "foo.bar",
			Param:         map[string]any{"foo": "bar"},
			ExpectedError: "unsupported value type string",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx, md := evaluate.ContextWithEmptyMetadata[ref.EmptyReferences](context.Background())

			eval, err := pathlang.DefaultExpression[ref.EmptyReferences]().NewEvaluable(test.Expr)
			require.NoError(t, err)

			invoker, err := fn.NewPositionalArgs(
				evaluate.DefaultEvaluator[ref.EmptyReferences](),
				[]any{
					query.EvalExpandable[ref.EmptyReferences](eval, test.Param),
				},
			).InvokerFor(fnlib.ExistsDescriptor[ref.EmptyReferences]())
			require.NoError(t, err)

			r, err := invoker.Invoke(ctx)
			if test.ExpectedError != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), test.ExpectedError)
			} else {
				require.NoError(t, err)
				require.True(t, md.OK())
				require.Equal(t, test.Expected, r.Value)
			}
		})
	}
}
