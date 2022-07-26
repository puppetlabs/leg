package fn

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type EagerEvaluationDescriptor[T ref.References[T]] struct {
	delegate Descriptor[T]
}

var _ Descriptor[ref.EmptyReferences] = &EagerEvaluationDescriptor[ref.EmptyReferences]{}

func (eed *EagerEvaluationDescriptor[_]) Description() string {
	return eed.delegate.Description()
}

func (eed *EagerEvaluationDescriptor[T]) wrapInvoker(ev evaluate.Evaluator[T], args any, inner Invoker[T]) Invoker[T] {
	return InvokerFunc[T](func(ctx context.Context) (*evaluate.Result[T], error) {
		er, err := evaluate.CoalesceEvaluatorsPreferringContext(ctx, ev).Evaluate(ctx, args, evaluate.DepthFull)
		if err != nil {
			return nil, err
		}

		r, err := inner.Invoke(ctx)
		if err != nil {
			return nil, err
		}

		er.References.SetUsed(false)
		r.MergeReferences(er.Metadata)

		return r, nil
	})
}

func (eed *EagerEvaluationDescriptor[T]) PositionalInvoker(ev evaluate.Evaluator[T], args *PositionalArgs[T]) (Invoker[T], error) {
	invoker, err := eed.delegate.PositionalInvoker(ev, args)
	if err != nil {
		return nil, err
	}

	return eed.wrapInvoker(ev, args.ToValue(), invoker), nil
}

func (eed *EagerEvaluationDescriptor[T]) KeywordInvoker(ev evaluate.Evaluator[T], args *KeywordArgs[T]) (Invoker[T], error) {
	invoker, err := eed.delegate.KeywordInvoker(ev, args)
	if err != nil {
		return nil, err
	}

	return eed.wrapInvoker(ev, args.ToValue(), invoker), nil
}

func NewEagerEvaluationDescriptor[T ref.References[T]](delegate Descriptor[T]) *EagerEvaluationDescriptor[T] {
	return &EagerEvaluationDescriptor[T]{
		delegate: delegate,
	}
}

type EagerEvaluationMap[T ref.References[T]] struct {
	delegate Map[T]
}

var _ Map[ref.EmptyReferences] = &EagerEvaluationMap[ref.EmptyReferences]{}

func (eem *EagerEvaluationMap[T]) Descriptor(name string) (Descriptor[T], error) {
	desc, err := eem.delegate.Descriptor(name)
	if err != nil {
		return nil, err
	}

	return NewEagerEvaluationDescriptor(desc), nil
}

func NewEagerEvaluationMap[T ref.References[T]](delegate Map[T]) *EagerEvaluationMap[T] {
	return &EagerEvaluationMap[T]{
		delegate: delegate,
	}
}
