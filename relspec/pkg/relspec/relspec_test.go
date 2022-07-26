package relspec_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
	"github.com/puppetlabs/leg/gvalutil/pkg/template"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/fnlib"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/puppetlabs/leg/relspec/pkg/relspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testID struct {
	Name string
}

func (t testID) Less(other testID) bool {
	return t.Name < other.Name
}

var (
	errInvalid  = errors.New("invalid")
	errNotFound = errors.New("not found")
)

type testReferences = *ref.Log[testID]

type testEnv map[string]any

var (
	_ relspec.MappingTypeResolver[testReferences] = &testEnv{}
	_ eval.Indexable                              = &testEnv{}
	_ evaluate.Expandable[testReferences]         = &testEnv{}
)

func (te testEnv) ResolveMappingType(ctx context.Context, tm map[string]any) (*evaluate.Result[testReferences], error) {
	r := evaluate.ContextualizedResult(evaluate.NewMetadata(ref.NewLog[testID]()))
	r.SetEvaluator(evaluate.DefaultEvaluator[testReferences]())

	name, ok := tm["name"].(string)
	if !ok {
		return nil, errInvalid
	}

	if v, ok := te[name]; ok {
		r.References.Set(ref.OK(testID{Name: name}))
		r.SetValue(v)
	} else {
		r.References.Set(ref.Errored(testID{Name: name}, errNotFound))
		r.SetValue(tm)
	}

	return r, nil
}

func (te testEnv) Index(ctx context.Context, idx any) (any, error) {
	name, err := eval.StringValue(idx)
	if err != nil {
		return nil, err
	}

	r := evaluate.ContextualizedResult(evaluate.NewMetadata(ref.NewLog[testID]()))
	r.SetEvaluator(evaluate.DefaultEvaluator[testReferences]())

	if v, ok := te[name]; ok {
		r.References.Set(ref.OK(testID{Name: name}))
		r.SetValue(v)
	} else {
		r.References.Set(ref.Errored(testID{Name: name}, errNotFound))
	}

	return evaluate.StaticExpandable(r), nil
}

func (te testEnv) Expand(ctx context.Context, depth int) (*evaluate.Result[testReferences], error) {
	if depth == 0 {
		return evaluate.StaticResult[testReferences](te), nil
	}

	r := evaluate.NewResult(evaluate.NewMetadata(ref.NewLog[testID]()), map[string]any(te))
	r.SetEvaluator(evaluate.DefaultEvaluator[testReferences]())
	for name := range te {
		r.References.Set(ref.OK(testID{Name: name}))
	}
	return r, nil
}

func jsonInvocation(name string, args any) map[string]any {
	return map[string]any{fmt.Sprintf("$fn.%s", name): args}
}

func jsonEncoding(ty transfer.EncodingType, data any) map[string]any {
	return map[string]any{"$encoding": string(ty), "data": data}
}

func jsonEnv(name string) map[string]any {
	return map[string]any{"$type": "Env", "name": name}
}

type randomOrder []any

type test struct {
	Name               string
	Data               string
	Env                testEnv
	Depth              int
	QueryLanguage      query.Language[testReferences]
	Query              string
	Into               any
	ExpectedValue      any
	ExpectedReferences testReferences
	ExpectedError      error
}

func (tt test) Run(t *testing.T) {
	ctx := context.Background()

	var tree any
	require.NoError(t, json.Unmarshal([]byte(tt.Data), &tree))

	ev := relspec.NewEvaluator[testReferences](
		relspec.WithMappingTypeResolvers[testReferences](map[string]relspec.MappingTypeResolver[testReferences]{
			"Env": tt.Env,
		}),
		relspec.WithTemplateEnvironment[testReferences](map[string]evaluate.Expandable[testReferences]{
			"env": tt.Env,
		}),
	)

	check := func(t *testing.T, err error) {
		if tt.ExpectedError != nil {
			require.Equal(t, tt.ExpectedError, err)
		} else {
			require.NoError(t, err)
		}
	}

	var (
		r   *evaluate.Result[testReferences]
		err error
	)
	if tt.Query != "" {
		lang := tt.QueryLanguage
		if lang == nil {
			lang = pathlang.New[testReferences](
				pathlang.WithFunctionMap[testReferences]{Map: fnlib.Library[testReferences]()},
			).Expression
		}

		r, err = query.EvaluateQuery(ctx, ev, lang, tree, tt.Query)
		check(t, err)
	} else if tt.Into != nil {
		md, err := evaluate.EvaluateInto(ctx, ev, tree, tt.Into)
		check(t, err)

		r = evaluate.NewResult(md, tt.Into)
	} else {
		depth := tt.Depth
		if depth == 0 {
			depth = evaluate.DepthFull
		}

		r, err = ev.Evaluate(ctx, tree, depth)
		check(t, err)
	}

	require.Equal(t, tt.ExpectedError == nil, r != nil)
	if r == nil {
		return
	}

	expected := tt.ExpectedValue
	if ro, ok := expected.(randomOrder); ok {
		expected = []any(ro)

		// Requests sorting before continuing.
		if actual, ok := r.Value.([]any); ok {
			sort.Slice(actual, func(i, j int) bool {
				return fmt.Sprintf("%T %v", actual[i], actual[i]) < fmt.Sprintf("%T %v", actual[j], actual[j])
			})
		}
	}

	assert.Equal(t, expected, r.Value)

	refs := tt.ExpectedReferences
	if refs == nil {
		refs = ref.NewLog[testID]()
	}

	assert.Equal(t, refs, r.References)
}

type tests []test

func (tts tests) RunAll(t *testing.T) {
	for _, tt := range tts {
		t.Run(tt.Name, tt.Run)
	}
}

func TestEvaluate(t *testing.T) {
	tests{
		{
			Name:          "literal",
			Data:          `{"foo": "bar"}`,
			ExpectedValue: map[string]any{"foo": "bar"},
		},
		{
			Name: "unresolvable data",
			Data: `{"foo": {"$type": "Env", "name": "bar"}}`,
			ExpectedValue: map[string]any{
				"foo": jsonEnv("bar"),
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "bar"}, errNotFound)),
		},
		{
			Name: "unresolvable data in template",
			Data: `{"foo": "${env.bar}"}`,
			ExpectedValue: map[string]any{
				"foo": "${env.bar}",
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "bar"}, errNotFound)),
		},
		{
			Name: "invalid invocation",
			Data: `{"foo": {"$fn.foo": "bar"}}`,
			ExpectedError: &evaluate.PathEvaluationError{
				Path:  "foo",
				Cause: &fn.InvocationError{Name: "foo", Cause: fn.ErrFunctionNotFound},
			},
		},
		{
			Name: "invalid invocation in template",
			Data: `{"foo": "${foo('bar')}"}`,
			ExpectedError: &evaluate.PathEvaluationError{
				Path: "foo",
				Cause: &template.EvaluationError{
					Start: "${",
					Cause: &fn.InvocationError{Name: "foo", Cause: fn.ErrFunctionNotFound},
				},
			},
		},
		{
			Name: "many unresolvable",
			Data: `{
				"a": {"$type": "Env", "name": "foo"},
				"b": {"$type": "Env", "name": "bar"},
				"c": "hello"
			}`,
			ExpectedValue: map[string]any{
				"a": jsonEnv("foo"),
				"b": jsonEnv("bar"),
				"c": "hello",
			},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "foo"}, errNotFound),
				ref.Errored(testID{Name: "bar"}, errNotFound),
			),
		},
		{
			Name: "many unresolvable in template",
			Data: `{
				"a": "${env.foo}",
				"b": "${env.bar}",
				"c": "hello"
			}`,
			ExpectedValue: map[string]any{
				"a": "${env.foo}",
				"b": "${env.bar}",
				"c": "hello",
			},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "foo"}, errNotFound),
				ref.Errored(testID{Name: "bar"}, errNotFound),
			),
		},
		{
			Name: "unresolvable at depth",
			Data: `{
				"foo": [
					{"a": {"$type": "Env", "name": "foo"}},
					{"$type": "Env", "name": "bar"}
				],
				"bar": {"$type": "Env", "name": "frob"}
			}`,
			Depth: 3,
			ExpectedValue: map[string]any{
				"foo": []any{
					map[string]any{"a": jsonEnv("foo")},
					jsonEnv("bar"),
				},
				"bar": jsonEnv("frob"),
			},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "bar"}, errNotFound),
				ref.Errored(testID{Name: "frob"}, errNotFound),
			),
		},
		{
			Name: "unresolvable at depth in template",
			Data: `{
				"foo": [
					{"a": "${env.foo}"},
					"${env.bar}"
				],
				"bar": "${env.frob}"
			}`,
			Depth: 3,
			ExpectedValue: map[string]any{
				"foo": []any{
					map[string]any{"a": "${env.foo}"},
					"${env.bar}",
				},
				"bar": "${env.frob}",
			},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "bar"}, errNotFound),
				ref.Errored(testID{Name: "frob"}, errNotFound),
			),
		},
		{
			Name: "resolvable",
			Data: `{
				"a": {"$type": "Env", "name": "foo"},
				"b": {"$type": "Env", "name": "quux"},
				"c": {"$fn.concat": ["foo", "bar"]},
				"d": "hello"
			}`,
			Env: testEnv{
				"foo":  "v3ry s3kr3t!",
				"quux": []any{1, 2, 3},
			},
			ExpectedValue: map[string]any{
				"a": "v3ry s3kr3t!",
				"b": []any{1, 2, 3},
				"c": "foobar",
				"d": "hello",
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "foo"}),
				ref.OK(testID{Name: "quux"}),
			),
		},
		{
			Name: "resolvable in template",
			Data: `{
				"a": "${env.foo}",
				"b": "${env.quux}",
				"c": "${concat('foo', 'bar')}",
				"d": "hello"
			}`,
			Env: testEnv{
				"foo":  "v3ry s3kr3t!",
				"quux": []any{1, 2, 3},
			},
			ExpectedValue: map[string]any{
				"a": "v3ry s3kr3t!",
				"b": []any{1, 2, 3},
				"c": "foobar",
				"d": "hello",
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "foo"}),
				ref.OK(testID{Name: "quux"}),
			),
		},
		{
			Name: "resolvable expansion in template",
			Data: `{
				"foo": "${$}"
			}`,
			Env: testEnv{
				"foo":  "v3ry s3kr3t!",
				"quux": []any{1, 2, 3},
			},
			ExpectedValue: map[string]any{
				"foo": map[string]any{
					"env": map[string]any{
						"foo":  "v3ry s3kr3t!",
						"quux": []any{1, 2, 3},
					},
				},
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "foo"}),
				ref.OK(testID{Name: "quux"}),
			),
		},
		{
			Name: "nested resolvable",
			Data: `{
				"aws": {
					"accessKeyID": {"$type": "Env", "name": "accessKeyID"},
					"secretAccessKey": {"$type": "Env", "name": "secretAccessKey"}
				},
				"instanceID": {"$type": "Env", "name": "instanceID"}
			}`,
			Env: testEnv{
				"accessKeyID":     "AKIANOAHISCOOL",
				"secretAccessKey": "abcdefs3cr37s",
				"instanceID":      "i-abcdef123456",
			},
			ExpectedValue: map[string]any{
				"aws": map[string]any{
					"accessKeyID":     "AKIANOAHISCOOL",
					"secretAccessKey": "abcdefs3cr37s",
				},
				"instanceID": "i-abcdef123456",
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "accessKeyID"}),
				ref.OK(testID{Name: "secretAccessKey"}),
				ref.OK(testID{Name: "instanceID"}),
			),
		},
		{
			Name: "nested resolvable in template",
			Data: `{
				"aws": {
					"accessKeyID": "${env.accessKeyID}",
					"secretAccessKey": "${env.secretAccessKey}"
				},
				"instanceID": "${env.instanceID}"
			}`,
			Env: testEnv{
				"accessKeyID":     "AKIANOAHISCOOL",
				"secretAccessKey": "abcdefs3cr37s",
				"instanceID":      "i-abcdef123456",
			},
			ExpectedValue: map[string]any{
				"aws": map[string]any{
					"accessKeyID":     "AKIANOAHISCOOL",
					"secretAccessKey": "abcdefs3cr37s",
				},
				"instanceID": "i-abcdef123456",
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "accessKeyID"}),
				ref.OK(testID{Name: "secretAccessKey"}),
				ref.OK(testID{Name: "instanceID"}),
			),
		},
		{
			Name: "resolvable data traversal",
			Data: `{
				"accessKeyID": "${env.aws.accessKeyID}"
			}`,
			Env: testEnv{
				"aws": map[string]any{"accessKeyID": "foo", "secretAccessKey": "bar"},
			},
			ExpectedValue: map[string]any{
				"accessKeyID": "foo",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "aws"})),
		},
		{
			Name: "resolvable data in invocation argument",
			Data: `{
				"aws": {"$fn.jsonUnmarshal": {"$type": "Env", "name": "aws"}}
			}`,
			Env: testEnv{
				"aws": `{"accessKeyID": "foo", "secretAccessKey": "bar"}`,
			},
			ExpectedValue: map[string]any{
				"aws": map[string]any{
					"accessKeyID":     "foo",
					"secretAccessKey": "bar",
				},
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "aws"})),
		},
		{
			Name: "resolvable data in invocation argument in partial template",
			Data: `{
				"aws": {"$fn.jsonUnmarshal": "${env.aws}"}
			}`,
			Env: testEnv{
				"aws": `{"accessKeyID": "foo", "secretAccessKey": "bar"}`,
			},
			ExpectedValue: map[string]any{
				"aws": map[string]any{
					"accessKeyID":     "foo",
					"secretAccessKey": "bar",
				},
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "aws"})),
		},
		{
			Name: "resolvable data in invocation argument in template",
			Data: `{
				"aws": "${jsonUnmarshal(env.aws)}"
			}`,
			Env: testEnv{
				"aws": `{"accessKeyID": "foo", "secretAccessKey": "bar"}`,
			},
			ExpectedValue: map[string]any{
				"aws": map[string]any{
					"accessKeyID":     "foo",
					"secretAccessKey": "bar",
				},
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "aws"})),
		},
		{
			Name: "unresolvable data in invocation argument",
			Data: `{
				"aws": {"$fn.jsonUnmarshal": {"$type": "Env", "name": "aws"}}
			}`,
			ExpectedValue: map[string]any{
				"aws": jsonInvocation("jsonUnmarshal", []any{jsonEnv("aws")}),
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "aws"}, errNotFound)),
		},
		{
			Name: "unresolvable data in invocation argument in partial template",
			Data: `{
				"aws": {"$fn.jsonUnmarshal": "${env.aws}"}
			}`,
			ExpectedValue: map[string]any{
				"aws": jsonInvocation("jsonUnmarshal", []any{"${env.aws}"}),
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "aws"}, errNotFound)),
		},
		{
			Name: "unresolvable data in invocation argument in template",
			Data: `{
				"aws": "${jsonUnmarshal(env.aws)}"
			}`,
			ExpectedValue: map[string]any{
				"aws": "${jsonUnmarshal(env.aws)}",
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "aws"}, errNotFound)),
		},
		{
			Name: "partially resolvable invocation",
			Data: `{
				"foo": {
					"$fn.concat": [
						{"$type": "Env", "name": "first"},
						{"$type": "Env", "name": "second"}
					]
				}
			}`,
			Env: testEnv{
				"first": "bar",
			},
			ExpectedValue: map[string]any{
				"foo": jsonInvocation("concat", []any{
					"bar",
					jsonEnv("second"),
				}),
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "first"}),
				ref.Errored(testID{Name: "second"}, errNotFound),
			),
		},
		{
			Name: "partially resolvable invocation in partial template",
			Data: `{
				"foo": {
					"$fn.concat": [
						"${env.first}",
						"${env.second}"
					]
				}
			}`,
			Env: testEnv{
				"first": "bar",
			},
			ExpectedValue: map[string]any{
				"foo": jsonInvocation("concat", []any{
					"bar",
					"${env.second}",
				}),
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "first"}),
				ref.Errored(testID{Name: "second"}, errNotFound),
			),
		},
		{
			Name: "partially resolvable invocation in template",
			Data: `{
				"foo": "${concat(env.first, env.second)}"
			}`,
			Env: testEnv{
				"first": "bar",
			},
			ExpectedValue: map[string]any{
				"foo": "${concat(env.first, env.second)}",
			},
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "first"}),
				ref.Errored(testID{Name: "second"}, errNotFound),
			),
		},
		{
			Name: "successful invocation of fn.convertMarkdown to Jira syntax",
			Data: `{
				"foo": {
					"$fn.convertMarkdown": [
						"jira",` +
				"\"--- `code` ---\"" + `
					]
				}
			}`,
			ExpectedValue: map[string]any{
				"foo": "\n----\n{code}code{code}\n----\n",
			},
		},
		{
			Name: "successful invocation of fn.convertMarkdown to Jira syntax in template",
			Data: fmt.Sprintf(`{"foo": %q}`, "${convertMarkdown('jira', '--- `code` ---')}"),
			ExpectedValue: map[string]any{
				"foo": "\n----\n{code}code{code}\n----\n",
			},
		},
		{
			Name: "encoded string",
			Data: `{
				"foo": {
					"$encoding": "base64",
					"data": "SGVsbG8sIJCiikU="
				}
			}`,
			ExpectedValue: map[string]any{
				"foo": "Hello, \x90\xA2\x8A\x45",
			},
		},
		{
			Name: "encoded string from data",
			Data: `{
				"foo": {
					"$encoding": "base64",
					"data": {"$type": "Env", "name": "bar"}
				}
			}`,
			Env: testEnv{
				"bar": "SGVsbG8sIJCiikU=",
			},
			ExpectedValue: map[string]any{
				"foo": "Hello, \x90\xA2\x8A\x45",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name: "encoded string from data in template",
			Data: `{
				"foo": {
					"$encoding": "base64",
					"data": "${env.bar}"
				}
			}`,
			Env: testEnv{
				"bar": "SGVsbG8sIJCiikU=",
			},
			ExpectedValue: map[string]any{
				"foo": "Hello, \x90\xA2\x8A\x45",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name: "encoded string from unresolvable data",
			Data: `{
				"foo": {
					"$encoding": "base64",
					"data": {"$type": "Env", "name": "bar"}
				}
			}`,
			ExpectedValue: map[string]any{
				"foo": jsonEncoding(transfer.Base64EncodingType, jsonEnv("bar")),
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "bar"}, errNotFound)),
		},
		{
			Name: "encoded string from unresolvable data in template",
			Data: `{
				"foo": {
					"$encoding": "base64",
					"data": "${env.bar}"
				}
			}`,
			ExpectedValue: map[string]any{
				"foo": jsonEncoding(transfer.Base64EncodingType, "${env.bar}"),
			},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "bar"}, errNotFound)),
		},
		{
			Name:          "invocation with array arguments",
			Data:          `{"$fn.concat": ["bar", "baz"]}`,
			ExpectedValue: "barbaz",
		},
		{
			Name:          "invocation with array arguments in template",
			Data:          `"${concat('bar', 'baz')}"`,
			ExpectedValue: "barbaz",
		},
		{
			Name:          "invocation with object arguments",
			Data:          `{"$fn.path": {"object": {"foo": ["bar"]}, "query": "foo[0]"}}`,
			ExpectedValue: "bar",
		},
		{
			Name:          "invocation with object arguments in template",
			Data:          `"${path(object: {'foo': ['bar']}, query: 'foo[0]')}"`,
			ExpectedValue: "bar",
		},
		{
			Name: "bad invocation",
			Data: `{"$fn.append": [1, 2, 3]}`,
			ExpectedError: &fn.InvocationError{
				Name: "append",
				Cause: &fn.PositionalArgError{
					Arg: 1,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf([]any(nil))},
						Got:    reflect.TypeOf(float64(0)),
					},
				},
			},
		},
		{
			Name: "bad invocation in template",
			Data: `"${append(1, 2, 3)}"`,
			ExpectedError: &template.EvaluationError{
				Start: "${",
				Cause: &fn.InvocationError{
					Name: "append",
					Cause: &fn.PositionalArgError{
						Arg: 1,
						Cause: &fn.UnexpectedTypeError{
							Wanted: []reflect.Type{reflect.TypeOf([]any(nil))},
							Got:    reflect.TypeOf(float64(0)),
						},
					},
				},
			},
		},
		{
			Name: "resolvable template dereferencing",
			Data: `"${env['regions.' + env.region]}"`,
			Env: testEnv{
				"region":            "us-east-1",
				"regions.us-east-1": "EAST",
				"regions.us-west-1": "WEST",
			},
			ExpectedValue: "EAST",
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "region"}),
				ref.OK(testID{Name: "regions.us-east-1"}),
			),
		},
		{
			Name: "unresolvable template dereferencing",
			Data: `"${env['regions.' + env.region]}"`,
			Env: testEnv{
				"region":            "us-east-2",
				"regions.us-east-1": "EAST",
				"regions.us-west-1": "WEST",
			},
			ExpectedValue: `${env['regions.' + env.region]}`,
			ExpectedReferences: ref.InitialLog(
				ref.OK(testID{Name: "region"}),
				ref.Errored(testID{Name: "regions.us-east-2"}, errNotFound),
			),
		},
		{
			Name:               "nested unresolvable template dereferencing",
			Data:               `"${env['regions.' + env.region]}"`,
			ExpectedValue:      `${env['regions.' + env.region]}`,
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "region"}, errNotFound)),
		},
		{
			Name: "data expansion in template",
			Data: `{
				"foo": "${env}"
			}`,
			Env: testEnv{
				"region": "us-east-1",
			},
			ExpectedValue: map[string]any{
				"foo": map[string]any{
					"region": "us-east-1",
				},
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "region"})),
		},
		{
			Name: "template interpolation",
			Data: `{
				"foo": "Hello, ${env.who}!"
			}`,
			Env: testEnv{
				"who": "friend",
			},
			ExpectedValue: map[string]any{
				"foo": "Hello, friend!",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "who"})),
		},
		{
			Name: "template interpolation with mapping type",
			Data: `{
				"foo": "Some secret people:\n${env}"
			}`,
			Env: testEnv{
				"who": "friend",
			},
			ExpectedValue: map[string]any{
				"foo": `Some secret people:
{
	"who": "friend"
}`,
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "who"})),
		},
	}.RunAll(t)
}

func TestEvaluateQuery(t *testing.T) {
	tests{
		{
			Name:          "literal",
			Data:          `{"foo": "bar"}`,
			Query:         `foo`,
			ExpectedValue: "bar",
		},
		{
			Name:  "nonexistent key",
			Data:  `{"foo": [{"bar": "baz"}]}`,
			Query: `foo[0].quux`,
			ExpectedError: &evaluate.PathEvaluationError{
				Path: "foo",
				Cause: &evaluate.PathEvaluationError{
					Path: "0",
					Cause: &evaluate.PathEvaluationError{
						Path:  "quux",
						Cause: &eval.UnknownKeyError{Key: "quux"},
					},
				},
			},
		},
		{
			Name:  "nonexistent index",
			Data:  `{"foo": [{"bar": "baz"}]}`,
			Query: `foo[1].quux`,
			ExpectedError: &evaluate.PathEvaluationError{
				Path: "foo",
				Cause: &evaluate.PathEvaluationError{
					Path:  "1",
					Cause: &eval.IndexOutOfBoundsError{Index: 1},
				},
			},
		},
		{
			Name: "traverses data",
			Data: `{
				"foo": {"$type": "Env", "name": "bar"}
			}`,
			Query: "foo.bar.baz",
			Env: testEnv{
				"bar": map[string]any{
					"bar": map[string]any{"baz": "quux"},
				},
			},
			ExpectedValue:      "quux",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name: "JSONPath traverses data",
			Data: `{
				"foo": {"$type": "Env", "name": "bar"}
			}`,
			QueryLanguage: query.JSONPathLanguage[testReferences],
			Query:         "$.foo.bar.baz",
			Env: testEnv{
				"bar": map[string]any{
					"bar": map[string]any{"baz": "quux"},
				},
			},
			ExpectedValue:      "quux",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name: "JSONPath template traverses data",
			Data: `{
				"foo": {"$type": "Env", "name": "bar"}
			}`,
			QueryLanguage: query.JSONPathTemplateLanguage[testReferences],
			Query:         "{.foo.bar.baz}",
			Env: testEnv{
				"bar": map[string]any{
					"bar": map[string]any{"baz": "quux"},
				},
			},
			ExpectedValue:      "quux",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "bar"})),
		},
		{
			Name: "unresolvable",
			Data: `{
				"foo": {"$type": "Env", "name": "bar"}
			}`,
			Query:              "foo.bar.baz",
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "bar"}, errNotFound)),
		},
		{
			Name: "JSONPath unresolvable",
			Data: `{
				"a": {"name": "aa", "value": {"$type": "Env", "name": "foo"}},
				"b": {"name": "bb", "value": {"$type": "Env", "name": "bar"}},
				"c": {"name": "cc", "value": "gggggg"}
			}`,
			QueryLanguage: query.JSONPathLanguage[testReferences],
			Query:         "$.*.value",
			ExpectedValue: []any{"gggggg"},
			ExpectedReferences: ref.InitialLog(
				ref.Errored(testID{Name: "foo"}, errNotFound),
				ref.Errored(testID{Name: "bar"}, errNotFound),
			),
		},
		{
			Name: "unresolvable not evaluated because not in path",
			Data: `{
				"a": {"$type": "Env", "name": "bar"},
				"b": {"c": {"$type": "Env", "name": "foo"}}
			}`,
			Query: "b.c",
			Env: testEnv{
				"foo": "very secret",
			},
			ExpectedValue:      "very secret",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
		{
			Name: "JSONPath object unresolvable not evaluated because not in path",
			Data: `{
				"a": {"name": "aa", "value": {"$type": "Env", "name": "bar"}},
				"b": {"name": "bb", "value": {"$type": "Env", "name": "foo"}}
			}`,
			QueryLanguage: query.JSONPathLanguage[testReferences],
			Query:         "$.*.name",
			ExpectedValue: randomOrder{"aa", "bb"},
		},
		{
			Name: "JSONPath array unresolvable not evaluated because not in path",
			Data: `[
				{"name": "aa", "value": {"$type": "Env", "name": "bar"}},
				{"name": "bb", "value": {"$type": "Env", "name": "foo"}}
			]`,
			QueryLanguage: query.JSONPathLanguage[testReferences],
			Query:         "$.*.name",
			ExpectedValue: randomOrder{"aa", "bb"},
		},
		{
			Name: "type resolver returns an unsupported type",
			Data: `{
				"a": {"$type": "Env", "name": "foo"}
			}`,
			Query: "a.inner",
			Env: testEnv{
				"foo": map[string]string{"inner": "test"},
			},
			ExpectedError: &evaluate.PathEvaluationError{
				Path: "a",
				Cause: &evaluate.PathEvaluationError{
					Path: "inner",
					Cause: &evaluate.UnsupportedValueError{
						Type: reflect.TypeOf(map[string]string(nil)),
					},
				},
			},
		},
		{
			Name: "type resolver returns an unsupported type in JSONPath",
			Data: `{
				"a": {"$type": "Env", "name": "foo"},
				"b": {"$type": "Env", "name": "bar"}
			}`,
			QueryLanguage: query.JSONPathLanguage[testReferences],
			Query:         "$.a.inner",
			Env: testEnv{
				"foo": map[string]string{"inner": "test"},
				"bar": map[string]any{"inner": "test"},
			},
			ExpectedError: &evaluate.UnsupportedValueError{
				Type: reflect.TypeOf(map[string]string(nil)),
			},
		},
		{
			Name: "type resolver returns an unsupported type in JSONPath template",
			Data: `{
				"a": {"$type": "Env", "name": "foo"}
			}`,
			QueryLanguage: query.JSONPathTemplateLanguage[testReferences],
			Query:         "{.a.inner}",
			Env: testEnv{
				"foo": map[string]string{"inner": "test"},
			},
			ExpectedError: &template.EvaluationError{
				Start: "{",
				Cause: &evaluate.UnsupportedValueError{
					Type: reflect.TypeOf(map[string]string(nil)),
				},
			},
		},
		{
			Name: "JSONPath template traverses object",
			Data: `{
				"args": {
					"a": "undo",
					"b": {"$fn.concat": ["deployment.v1.apps/", {"$type": "Env", "name": "deployment"}]}
				}
			}`,
			QueryLanguage: query.JSONPathTemplateLanguage[testReferences],
			Query:         "{.args}",
			Env: testEnv{
				"deployment": "my-test-deployment",
			},
			ExpectedValue:      "map[a:undo b:deployment.v1.apps/my-test-deployment]",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "deployment"})),
		},
		{
			Name: "JSONPath template traverses array",
			Data: `{
				"args": [
					"undo",
					{"$fn.concat": ["deployment.v1.apps/", {"$type": "Env", "name": "deployment"}]}
				]
			}`,
			QueryLanguage: query.JSONPathTemplateLanguage[testReferences],
			Query:         "{.args}",
			Env: testEnv{
				"deployment": "my-test-deployment",
			},
			ExpectedValue:      "undo deployment.v1.apps/my-test-deployment",
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "deployment"})),
		},
		{
			Name: "JSONPath template traverses array with unresolvables",
			Data: `{
				"args": [
					"undo",
					{"$fn.concat": ["deployment.v1.apps/", {"$type": "Env", "name": "deployment"}]}
				]
			}`,
			QueryLanguage:      query.JSONPathTemplateLanguage[testReferences],
			Query:              "{.args}",
			ExpectedValue:      "undo map[$fn.concat:[deployment.v1.apps/ map[$type:Env name:deployment]]]",
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "deployment"}, errNotFound)),
		},
		{
			Name:  "query has an error under a path",
			Data:  `{"foo": {"bar": ["baz", "quux"]}}`,
			Query: "foo.bar[0].nope",
			ExpectedError: &evaluate.PathEvaluationError{
				Path: "foo",
				Cause: &evaluate.PathEvaluationError{
					Path: "bar",
					Cause: &evaluate.PathEvaluationError{
						Path: "0",
						Cause: &evaluate.PathEvaluationError{
							Path: "nope",
							Cause: &eval.UnsupportedValueTypeError{
								Value: "baz",
								Field: "nope",
							},
						},
					},
				},
			},
		},
		{
			Name: "attempted escape from data",
			Data: `{
				"basic": {"$type": "Env", "name": "foo"},
				"ref": {"$fn.path": {"object": {"$type": "Env", "name": "foo"}, "query": "ref"}},
				"expr": {"$fn.path": {"object": {"$type": "Env", "name": "foo"}, "query": "expr"}}
			}`,
			Env: testEnv{
				"foo": map[string]any{
					"ref":  jsonEnv("bar"),
					"expr": "${env.quux}",
				},
			},
			ExpectedValue: map[string]any{
				"basic": map[string]any{
					"ref":  jsonEnv("bar"),
					"expr": "${env.quux}",
				},
				"ref":  jsonEnv("bar"),
				"expr": "${env.quux}",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
		{
			Name: "attempted escape from data using expressions",
			Data: `{
				"basic": "${env.foo}",
				"ref": "${path(env.foo, 'ref')}",
				"expr": "${path(env.foo, 'expr')}"
			}`,
			Env: testEnv{
				"foo": map[string]any{
					"ref":  jsonEnv("bar"),
					"expr": "${env.quux}",
				},
			},
			ExpectedValue: map[string]any{
				"basic": map[string]any{
					"ref":  jsonEnv("bar"),
					"expr": "${env.quux}",
				},
				"ref":  jsonEnv("bar"),
				"expr": "${env.quux}",
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
		{
			Name: "attempted self-referential escape from data",
			Data: `{
				"expand": {"$fn.path": {"object": {"$type": "Env", "name": "foo"}, "query": "ref"}}
			}`,
			Env: testEnv{
				"foo": map[string]any{
					"ref": jsonEnv("foo"),
				},
			},
			ExpectedValue: map[string]any{
				"expand": jsonEnv("foo"),
			},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
	}.RunAll(t)
}

func TestEvaluateIntoBasic(t *testing.T) {
	type foo struct {
		Bar string `spec:"bar"`
	}

	type root struct {
		Foo foo `spec:"foo"`
	}

	tests{
		{
			Name:          "basic",
			Data:          `{"foo": {"bar": "baz"}}`,
			Into:          &root{},
			ExpectedValue: &root{Foo: foo{Bar: "baz"}},
		},
		{
			Name: "resolvable",
			Data: `{"foo": {"bar": {"$type": "Env", "name": "foo"}}}`,
			Env: testEnv{
				"foo": "v3ry s3kr3t!",
			},
			Into:               &root{},
			ExpectedValue:      &root{Foo: foo{Bar: "v3ry s3kr3t!"}},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
		{
			Name:               "unresolvable",
			Data:               `{"foo": {"bar": {"$type": "Env", "name": "foo"}}}`,
			Into:               &root{Foo: foo{Bar: "masked"}},
			ExpectedValue:      &root{},
			ExpectedReferences: ref.InitialLog(ref.Errored(testID{Name: "foo"}, errNotFound)),
		},
		{
			Name: "map",
			Data: `{"foo": {"bar": {"$type": "Env", "name": "foo"}}}`,
			Env: testEnv{
				"foo": "v3ry s3kr3t!",
			},
			Into:               &map[string]any{},
			ExpectedValue:      &map[string]any{"foo": map[string]any{"bar": "v3ry s3kr3t!"}},
			ExpectedReferences: ref.InitialLog(ref.OK(testID{Name: "foo"})),
		},
	}.RunAll(t)
}
