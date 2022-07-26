package fn

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Invoker[T ref.References[T]] interface {
	Invoke(ctx context.Context) (*evaluate.Result[T], error)
}

type InvokerFunc[T ref.References[T]] func(ctx context.Context) (*evaluate.Result[T], error)

var _ Invoker[ref.EmptyReferences] = InvokerFunc[ref.EmptyReferences](nil)

func (fn InvokerFunc[T]) Invoke(ctx context.Context) (*evaluate.Result[T], error) {
	return fn(ctx)
}

func StaticInvoker[T ref.References[T]](value any) Invoker[T] {
	return InvokerFunc[T](func(_ context.Context) (*evaluate.Result[T], error) { return evaluate.StaticResult[T](value), nil })
}

func EvaluatedInvoker[T ref.References[T], V any, A ComparableArg[T], AS Args[T, A, V]](upd AS, args []A, fn func(ctx context.Context, args map[A]any) (any, error)) Invoker[T] {
	return InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
		md := evaluate.EmptyMetadata[T]()

		m := make(map[A]any, len(args))
		for _, arg := range args {
			if !ArgPresent[T](arg) {
				continue
			}

			r, err := arg.Expand(ctx, evaluate.DepthFull)
			if err != nil {
				return nil, err
			}

			upd.Update(arg, r.Value)
			md.AddStep(r, arg.Annotation())

			m[arg] = r.Value
		}

		if !md.OK() {
			return evaluate.NewResult(md, upd.ToValue()), nil
		}

		rv, err := fn(ctx, m)
		if err != nil {
			return nil, err
		}

		return evaluate.NewResult(md, rv), nil
	})
}
