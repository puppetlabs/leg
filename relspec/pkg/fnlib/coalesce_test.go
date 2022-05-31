package fnlib_test

import (
	"context"
	"errors"
	"testing"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/stretchr/testify/require"
)

func TestCoalesce(t *testing.T) {
	tests := []struct {
		Name     string
		Args     []any
		Expected any
	}{
		{
			Name:     "null then values",
			Args:     []any{nil, 3, 4},
			Expected: 3,
		},
		{
			Name: "unresolvable then values",
			Args: []any{
				evaluate.StaticExpandable(
					evaluate.ContextualizedResult(evaluate.NewMetadata(
						ref.InitialLog(ref.Errored(testID{Name: "foo"}, errors.New("not found"))),
					)),
				),
				3,
				4,
			},
			Expected: 3,
		},
		{
			Name:     "values first",
			Args:     []any{1, nil},
			Expected: 1,
		},
		{
			Name:     "no arguments",
			Expected: nil,
		},
		{
			Name: "no unresolvable or non-null values",
			Args: []any{
				nil,
				evaluate.StaticExpandable(
					evaluate.ContextualizedResult(evaluate.NewMetadata(
						ref.InitialLog(ref.Errored(testID{Name: "foo"}, errors.New("not found"))),
					)),
				),
			},
			Expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			invoker, err := fn.NewPositionalArgs(
				evaluate.DefaultEvaluator[*ref.Log[testID]](),
				test.Args,
			).InvokerFor(fnlib.CoalesceDescriptor[*ref.Log[testID]]())
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			require.NoError(t, err)
			require.True(t, r.OK())
			require.Equal(t, test.Expected, r.Value)
		})
	}
}
