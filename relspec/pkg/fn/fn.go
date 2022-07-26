package fn

import (
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

// Descriptor is a type that describes how a function can be invoked by a caller
type Descriptor[T ref.References[T]] interface {
	// Description returns a string that describes what the function does
	Description() string
	// PositionalInvoker takes a slice of values that act like positional
	// arguments to the function. Enforcing order and length constraints is up
	// to the author of the function.
	PositionalInvoker(ev evaluate.Evaluator[T], args *PositionalArgs[T]) (Invoker[T], error)
	// KeywordInvoker takes its arguments as a map. This acts like labeled or named argments
	// to the function. Enforcing name and length constraints is up to the author
	// of the function.
	KeywordInvoker(ev evaluate.Evaluator[T], args *KeywordArgs[T]) (Invoker[T], error)
}

// DescriptorFuncs is an adapter that takes anonymous functions that handle
// methods defined in the Descriptor interface. This is a convenience type that
// allows simple wrapping of one-off functions.
type DescriptorFuncs[T ref.References[T]] struct {
	DescriptionFunc       func() string
	PositionalInvokerFunc func(ev evaluate.Evaluator[T], args *PositionalArgs[T]) (Invoker[T], error)
	KeywordInvokerFunc    func(ev evaluate.Evaluator[T], args *KeywordArgs[T]) (Invoker[T], error)
}

var _ Descriptor[ref.EmptyReferences] = DescriptorFuncs[ref.EmptyReferences]{}

func (df DescriptorFuncs[_]) Description() string {
	if df.DescriptionFunc == nil {
		return "<anonymous>"
	}

	return df.DescriptionFunc()
}

func (df DescriptorFuncs[T]) PositionalInvoker(ev evaluate.Evaluator[T], args *PositionalArgs[T]) (Invoker[T], error) {
	if df.PositionalInvokerFunc == nil {
		return nil, ErrPositionalArgsNotAccepted
	}

	return df.PositionalInvokerFunc(ev, args)
}

func (df DescriptorFuncs[T]) KeywordInvoker(ev evaluate.Evaluator[T], args *KeywordArgs[T]) (Invoker[T], error) {
	if df.KeywordInvokerFunc == nil {
		return nil, ErrKeywordArgsNotAccepted
	}

	return df.KeywordInvokerFunc(ev, args)
}

type Map[T ref.References[T]] interface {
	Descriptor(name string) (Descriptor[T], error)
}

type funcMap[T ref.References[T]] map[string]Descriptor[T]

func (fm funcMap[T]) Descriptor(name string) (Descriptor[T], error) {
	fd, found := fm[name]
	if !found {
		return nil, ErrFunctionNotFound
	}

	return fd, nil
}

func NewMap[T ref.References[T]](m map[string]Descriptor[T]) Map[T] {
	return funcMap[T](m)
}
