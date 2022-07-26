package relspec

import (
	"context"
	"reflect"
	"strings"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type MappingTypeResolver[T ref.References[T]] interface {
	ResolveMappingType(ctx context.Context, tm map[string]any) (*evaluate.Result[T], error)
}

func evaluateType[T ref.References[T]](ctx context.Context, resolvers map[string]MappingTypeResolver[T], tm map[string]any) (*evaluate.Result[T], error) {
	ty, ok := tm["$type"].(string)
	if !ok {
		return evaluate.StaticResult[T](tm), nil
	}

	resolver, found := resolvers[ty]
	if !found {
		return evaluate.StaticResult[T](tm), nil
	}

	r, err := resolver.ResolveMappingType(ctx, tm)
	if err != nil {
		return nil, &InvalidTypeError{Type: ty, Cause: err}
	}

	return r, nil
}

func evaluateEncoding[T ref.References[T]](ctx context.Context, em map[string]any, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
	ty, ok := em["$encoding"].(string)
	if !ok {
		return evaluate.StaticResult[T](em), nil
	}

	r, err := evaluate.EvaluateAll(ctx, next, em["data"])
	if err != nil {
		return nil, &InvalidEncodingError{Type: ty, Cause: err}
	}

	if !r.OK() {
		r.SetValue(map[string]any{
			"$encoding": ty,
			"data":      r.Value,
		})
	} else {
		data, ok := r.Value.(string)
		if !ok {
			return nil, &InvalidEncodingError{
				Type: ty,
				Cause: &fn.UnexpectedTypeError{
					Wanted: []reflect.Type{reflect.TypeOf("")},
					Got:    reflect.TypeOf(r.Value),
				},
			}
		}

		decoded, err := transfer.JSON{
			EncodingType: transfer.EncodingType(ty),
			Data:         data,
		}.Decode()
		if err != nil {
			return nil, &InvalidEncodingError{Type: ty, Cause: err}
		}

		r.SetValue(string(decoded))
	}

	return r, nil
}

func evaluateInvocationPositional[T ref.References[T]](ctx context.Context, fns fn.Map[T], name string, args []any, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
	r, err := fn.EvaluateInvocation[T, *fn.PositionalArg[T], []any](ctx, fns, name, fn.NewPositionalArgs(next, args))
	if err != nil {
		return nil, err
	}

	r.SetAnnotation(fn.PositionalInvocationAnnotation(name))

	return r, nil
}

func evaluateInvocationKeyword[T ref.References[T]](ctx context.Context, fns fn.Map[T], name string, args map[string]any, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
	r, err := fn.EvaluateInvocation[T, *fn.KeywordArg[T], map[string]any](ctx, fns, name, fn.NewKeywordArgs(next, args))
	if err != nil {
		return nil, err
	}

	r.SetAnnotation(fn.KeywordInvocationAnnotation(name))

	return r, nil
}

func evaluateInvocation[T ref.References[T]](ctx context.Context, fns fn.Map[T], im map[string]any, next evaluate.Evaluator[T]) (r *evaluate.Result[T], err error) {
	var key string
	var value any
	for key, value = range im {
	}

	name := strings.TrimPrefix(key, "$fn.")

	// Evaluate one level to determine whether we should do a positional or
	// keyword invocation.
	a, err := next.Evaluate(ctx, value, 1)
	if err != nil {
		return nil, &fn.InvocationError{Name: name, Cause: err}
	} else if !a.OK() {
		// The top level couldn't be resolved, so we'll pass it in unmodified as
		// a single-argument parameter.
		r, err = evaluateInvocationPositional(ctx, fns, name, []any{value}, next)
	} else {
		next = evaluate.CoalesceEvaluators(a.Evaluator, next)

		switch ra := a.Value.(type) {
		case []any:
			r, err = evaluateInvocationPositional(ctx, fns, name, ra, next)
		case map[string]any:
			r, err = evaluateInvocationKeyword(ctx, fns, name, ra, next)
		default:
			r, err = evaluateInvocationPositional(ctx, fns, name, []any{ra}, next)
		}
	}
	if err != nil {
		return
	}

	r.InsertSteps(a.Metadata)
	r.MergeReferences(a.Metadata)

	if !r.OK() {
		r.SetValue(map[string]any{key: r.Value})
	}

	return
}

func evaluateMap[T ref.References[T]](o *Options[T]) func(ctx context.Context, m map[string]any, depth int, next evaluate.Evaluator[T]) (*evaluate.Result[T], error) {
	fns := o.FunctionMap
	if o.EagerEvaluation {
		fns = fn.NewEagerEvaluationMap(fns)
	}

	return func(ctx context.Context, m map[string]any, depth int, next evaluate.Evaluator[T]) (r *evaluate.Result[T], err error) {
		if !o.KeepSteps {
			defer func() {
				if r != nil {
					r.Reduce()
				}
			}()
		}

		if _, ok := m["$type"]; ok {
			return evaluateType(ctx, o.MappingTypeResolvers, m)
		} else if _, ok := m["$encoding"]; ok {
			return evaluateEncoding(ctx, m, next)
		} else if len(m) == 1 {
			var first string
			for first = range m {
			}

			if strings.HasPrefix(first, "$fn.") {
				return evaluateInvocation(ctx, fns, m, next)
			}
		}

		return evaluate.DefaultVisitor[T]().VisitMap(ctx, m, depth, next)
	}
}
