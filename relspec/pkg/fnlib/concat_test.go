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

func TestConcat(t *testing.T) {
	tests := []struct {
		Name     string
		Args     []any
		Expected any
	}{
		{
			Name:     "empty",
			Expected: "",
		},
		{
			Name:     "basic",
			Args:     []any{"Hello, ", "world!"},
			Expected: "Hello, world!",
		},
		{
			Name:     "type conversion",
			Args:     []any{"H", 3, "llo, world!"},
			Expected: "H3llo, world!",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			invoker, err := fn.NewPositionalArgs(
				evaluate.DefaultEvaluator[ref.EmptyReferences](),
				test.Args,
			).InvokerFor(fnlib.ConcatDescriptor[ref.EmptyReferences]())
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			require.NoError(t, err)
			require.True(t, r.OK())
			require.Equal(t, test.Expected, r.Value)
		})
	}
}
