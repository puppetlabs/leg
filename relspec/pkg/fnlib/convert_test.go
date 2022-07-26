package fnlib_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/puppetlabs/leg/relspec/pkg/convert"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/stretchr/testify/require"
)

func TestConvertMarkdown(t *testing.T) {
	desc := fnlib.ConvertMarkdownDescriptor[ref.EmptyReferences]()

	tcs := []struct {
		Name        string
		ConvertType convert.ConvertType
		Markdown    string
		Expected    string
	}{
		{
			Name:        "Sample monitor event",
			ConvertType: convert.ConvertTypeJira,
			Markdown:    "%%% @contact [![imageTitle](imageUrl)](imageRedirect) `data{context} > threshold` Detailed description. - - - [[linkTitle1](link1)] · [[linkTitle2](link2)] %%%",
			Expected:    "@contact \n\n[!imageUrl!|imageRedirect] {code}data{context} > threshold{code} Detailed description.\n----\n[[linkTitle1|link1]] · [[linkTitle2|link2]]",
		},
	}

	for _, test := range tcs {
		t.Run(test.Name, func(t *testing.T) {
			invoker, err := fn.NewPositionalArgs(
				evaluate.DefaultEvaluator[ref.EmptyReferences](),
				[]any{
					test.ConvertType.String(),
					test.Markdown,
				},
			).InvokerFor(desc)
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			require.NoError(t, err)

			require.True(t, r.OK())
			require.Equal(t, test.Expected, r.Value)

			invoker, err = fn.NewKeywordArgs(
				evaluate.DefaultEvaluator[ref.EmptyReferences](),
				map[string]any{
					"to":      test.ConvertType.String(),
					"content": test.Markdown,
				},
			).InvokerFor(desc)
			require.NoError(t, err)

			r, err = invoker.Invoke(context.Background())
			require.NoError(t, err)

			require.True(t, r.OK())
			require.Equal(t, test.Expected, r.Value)
		})
	}
}

func TestConvertMarkdownFunction(t *testing.T) {
	desc := fnlib.ConvertMarkdownDescriptor[ref.EmptyReferences]()

	tcs := []struct {
		Name                 string
		Invoker              func() (fn.Invoker[ref.EmptyReferences], error)
		ExpectedInvokeError  error
		ExpectedInvokerError error
	}{
		{
			Name: "keyword invoker with unsupported convert type",
			Invoker: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
					"to":      "foo",
					"content": "bar",
				}).InvokerFor(desc)
			},
			ExpectedInvokerError: convert.ErrConvertTypeNotSupported,
		},
		{
			Name: "keyword invoker with invalid to keyword type",
			Invoker: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
					"to":      false,
					"content": "bar",
				}).InvokerFor(desc)
			},
			ExpectedInvokerError: &fn.KeywordArgError{
				Arg: "to",
				Cause: &fn.UnexpectedTypeError{
					Wanted: []reflect.Type{
						reflect.TypeOf(""),
					},
					Got: reflect.TypeOf(false),
				},
			},
		},
		{
			Name: "keyword invoker with invalid content keyword type",
			Invoker: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
					"to":      "jira",
					"content": false,
				}).InvokerFor(desc)
			},
			ExpectedInvokerError: &fn.KeywordArgError{
				Arg: "content",
				Cause: &fn.UnexpectedTypeError{
					Wanted: []reflect.Type{
						reflect.TypeOf(""),
					},
					Got: reflect.TypeOf(false),
				},
			},
		},
		{
			Name: "keyword invoker with missing to keyword",
			Invoker: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
					"content": "bar",
				}).InvokerFor(desc)
			},
			ExpectedInvokeError: &fn.KeywordArgError{
				Arg:   "to",
				Cause: fn.ErrArgNotFound,
			},
		},
		{
			Name: "keyword invoker with missing content",
			Invoker: func() (fn.Invoker[ref.EmptyReferences], error) {
				return fn.NewKeywordArgs(evaluate.DefaultEvaluator[ref.EmptyReferences](), map[string]any{
					"to": "jira",
				}).InvokerFor(desc)
			},
			ExpectedInvokeError: &fn.KeywordArgError{
				Arg:   "content",
				Cause: fn.ErrArgNotFound,
			},
		},
	}

	for _, test := range tcs {
		t.Run(test.Name, func(t *testing.T) {
			invoker, err := test.Invoker()
			if test.ExpectedInvokeError != nil {
				require.Equal(t, test.ExpectedInvokeError, err)
			} else {
				require.NoError(t, err)

				_, err = invoker.Invoke(context.Background())
				if test.ExpectedInvokerError != nil {
					require.Equal(t, test.ExpectedInvokerError, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}
