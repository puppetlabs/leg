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

func TestJSONUnmarshal(t *testing.T) {
	invoker, err := fn.NewPositionalArgs(
		evaluate.DefaultEvaluator[ref.EmptyReferences](),
		[]any{`{"foo": "bar"}`},
	).InvokerFor(fnlib.JSONUnmarshalDescriptor[ref.EmptyReferences]())
	require.NoError(t, err)

	r, err := invoker.Invoke(context.Background())
	require.NoError(t, err)
	require.True(t, r.OK())
	require.Equal(t, map[string]any{"foo": "bar"}, r.Value)
}
