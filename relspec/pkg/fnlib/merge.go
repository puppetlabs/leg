package fnlib

import (
	"context"
	"fmt"
	"reflect"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func merge(dst, src map[string]any, deep bool) {
	for k, v := range src {
		if deep {
			if dm, ok := dst[k].(map[string]any); ok {
				if sm, ok := v.(map[string]any); ok {
					merge(dm, sm, deep)
					continue
				}
			}
		}

		dst[k] = v
	}
}

func MergeDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string {
			return `Merges a series of objects, with each object overwriting prior entries.

Merges are performed deeply by default. Use the keyword form and set mode: shallow to perform a shallow merge.`
		},
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if args.Len() == 0 {
				return fn.StaticInvoker[T](map[string]any{}), nil
			}

			s := args.UseAll()

			return fn.EvaluatedInvoker[T, []any](args, s, func(ctx context.Context, vs map[*fn.PositionalArg[T]]any) (any, error) {
				r := make(map[string]any)

				for _, arg := range s {
					obj, ok := vs[arg].(map[string]any)
					if !ok {
						return nil, arg.WrapError(&fn.UnexpectedTypeError{
							Wanted: []reflect.Type{reflect.TypeOf(map[string]any(nil))},
							Got:    reflect.TypeOf(vs[arg]),
						})
					}

					merge(r, obj, true)
				}

				return r, nil
			}), nil
		},
		KeywordInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.KeywordArgs[T]) (fn.Invoker[T], error) {
			if !args.Has("objects") {
				return nil, &fn.KeywordArgError{Arg: "objects", Cause: fn.ErrArgNotFound}
			}

			objects := args.Use("objects")
			mode := args.Use("mode")

			return fn.EvaluatedInvoker[T, map[string]any](args, []*fn.KeywordArg[T]{objects, mode}, func(ctx context.Context, vs map[*fn.KeywordArg[T]]any) (any, error) {
				deep := true
				if fn.ArgPresent[T](mode) {
					switch vs[mode] {
					case "deep":
					case "shallow":
						deep = false
					default:
						return nil, mode.WrapError(fmt.Errorf(`unexpected value %q, wanted one of "deep" or "shallow"`, vs[mode]))
					}
				}

				os, ok := vs[objects].([]any)
				if !ok {
					return nil, objects.WrapError(&fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf([]any(nil))},
						Got:    reflect.TypeOf(vs[objects]),
					})
				}

				r := make(map[string]any)

				for i, o := range os {
					obj, ok := o.(map[string]any)
					if !ok {
						return nil, objects.WrapError(fmt.Errorf("array index %d: %w", i, &fn.UnexpectedTypeError{
							Wanted: []reflect.Type{reflect.TypeOf(map[string]any(nil))},
							Got:    reflect.TypeOf(o),
						}))
					}

					merge(r, obj, deep)
				}

				return r, nil
			}), nil
		},
	}
}
