package fnlib_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/stretchr/testify/require"
)

func TestConditionals(t *testing.T) {
	equals := fnlib.EqualsDescriptor[ref.EmptyReferences]()
	notEquals := fnlib.NotEqualsDescriptor[ref.EmptyReferences]()

	cases := []struct {
		descriptor     fn.Descriptor[ref.EmptyReferences]
		args           []any
		expectedResult bool
		expectedError  error
	}{
		{
			descriptor:     equals,
			args:           []any{"foobar", "foobar"},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{10, 10},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{10.5, 10.5},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{[]any{"foo", "bar"}, []any{"foo", "bar"}},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{[]any{1, 2}, []any{1, 2}},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{[]any{1.1, 2.0}, []any{1.1, 2.0}},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{[]any{1.1, 2.0}, []any{1.1, 2.0}},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{map[string]any{"foo": "bar"}, map[string]any{"foo": "bar"}},
			expectedResult: true,
		},
		{
			descriptor:     equals,
			args:           []any{"true", true},
			expectedResult: false,
		},
		{
			descriptor:     equals,
			args:           []any{"10", 10},
			expectedResult: false,
		},
		{
			descriptor:     equals,
			args:           []any{10.5, 10},
			expectedResult: false,
		},
		{
			descriptor:     equals,
			args:           []any{1, 2},
			expectedResult: false,
		},
		{
			descriptor:    equals,
			args:          []any{1, 2, 3},
			expectedError: &fn.ArityError{Wanted: []int{2}, Got: 3},
		},
		{
			descriptor:     notEquals,
			args:           []any{"foobar", "barfoo"},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{10, 50},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{10.0, 50.5},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{[]any{"foo", "bar", "baz"}, []any{"foo", "bar"}},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{[]any{1, 2, 3}, []any{1, 2}},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{[]any{1.1, 2.0, 3.2}, []any{1.1, 2.0}},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{[]any{1.1, 2.0, 3.2}, []any{1.1, 2.0}},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{map[string]any{"foo": "bar", "baz": "biz"}, map[string]any{"foo": "bar"}},
			expectedResult: true,
		},
		{
			descriptor:     notEquals,
			args:           []any{true, true},
			expectedResult: false,
		},
		{
			descriptor:     notEquals,
			args:           []any{10, 10},
			expectedResult: false,
		},
		{
			descriptor:     notEquals,
			args:           []any{"foobar", "foobar"},
			expectedResult: false,
		},
		{
			descriptor:    notEquals,
			args:          []any{1, 2, 3},
			expectedError: &fn.ArityError{Wanted: []int{2}, Got: 3},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %v", i, c.args), func(t *testing.T) {
			invoker, err := fn.NewPositionalArgs(
				evaluate.DefaultEvaluator[ref.EmptyReferences](),
				c.args,
			).InvokerFor(c.descriptor)
			if c.expectedError != nil {
				require.EqualError(t, err, c.expectedError.Error())
			} else {
				require.NoError(t, err)

				r, err := invoker.Invoke(context.Background())
				require.NoError(t, err)

				require.True(t, r.OK())
				require.Equal(t, c.expectedResult, r.Value)
			}
		})
	}
}
