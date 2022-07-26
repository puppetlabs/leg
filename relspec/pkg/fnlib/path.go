package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func path[T ref.References[T], V any, A fn.Arg[T], AS fn.Args[T, A, V]](ctx context.Context, ev evaluate.Evaluator[T], upd AS, obj, q, def A) (*evaluate.Result[T], error) {
	r := evaluate.EmptyResult[T]()

	qr, err := q.Expand(ctx, evaluate.DepthFull)
	if err != nil {
		return nil, err
	}

	upd.Update(q, qr.Value)
	r.AddStep(qr, q.Annotation())

	if !r.OK() {
		return r, nil
	}

	qs, ok := qr.Value.(string)
	if !ok {
		return nil, q.WrapError(&fn.UnexpectedTypeError{
			Wanted: []reflect.Type{reflect.TypeOf("")},
			Got:    reflect.TypeOf(qr.Value),
		})
	}

	or, err := query.EvaluateQuery(ctx, ev, pathlang.New[T]().Expression, fn.ArgValue[T](obj), qs)
	if err != nil {
		if !fn.ArgPresent[T](def) || !errmark.Matches(err, notExistsRule) {
			return nil, obj.WrapError(err)
		}
	} else {
		r.AddStep(or, obj.Annotation())

		if or.OK() {
			r.SetValue(or.Value)
		}

		return r, nil
	}

	dr, err := def.Expand(ctx, evaluate.DepthFull)
	if err != nil {
		return nil, err
	}

	upd.Update(def, dr.Value)
	r.AddStep(dr, def.Annotation())

	if r.OK() {
		r.SetValue(dr.Value)
	}

	return r, nil
}

func pathInvoker[T ref.References[T], V any, A fn.ComparableArg[T], AS fn.Args[T, A, V]](ev evaluate.Evaluator[T], args AS, obj, q, def A) fn.Invoker[T] {
	return fn.InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
		r, err := path[T, V](ctx, ev, args, obj, q, def)
		if err != nil {
			return nil, err
		}

		if !r.OK() {
			r.SetValue(args.ToValue())
		}

		return r, nil
	})
}

func PathDescriptor[T ref.References[T]]() fn.Descriptor[T] {
	return fn.DescriptorFuncs[T]{
		DescriptionFunc: func() string {
			return "Looks up a value at a given path in an object, optionally returning a default value if the path does not exist"
		},
		PositionalInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.PositionalArgs[T]) (fn.Invoker[T], error) {
			if l := args.Len(); l < 2 || l > 3 {
				return nil, &fn.ArityError{Wanted: []int{2, 3}, Got: l}
			}

			obj, q, def := args.Use(0), args.Use(1), args.Use(2)

			return pathInvoker[T, []any](ev, args, obj, q, def), nil
		},
		KeywordInvokerFunc: func(ev evaluate.Evaluator[T], args *fn.KeywordArgs[T]) (fn.Invoker[T], error) {
			for _, arg := range []string{"object", "query"} {
				if !args.Has(arg) {
					return nil, &fn.KeywordArgError{Arg: arg, Cause: fn.ErrArgNotFound}
				}
			}

			obj, q, def := args.Use("object"), args.Use("query"), args.Use("default")

			return pathInvoker[T, map[string]any](ev, args, obj, q, def), nil
		},
	}
}
