package pathlang

import (
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Options[T ref.References[T]] struct {
	FunctionMap     fn.Map[T]
	EagerEvaluation bool
}

type Option[T ref.References[T]] interface {
	ApplyToOptions(target *Options[T])
}

var _ Option[ref.EmptyReferences] = &Options[ref.EmptyReferences]{}

func (o *Options[T]) ApplyToOptions(target *Options[T]) {
	if o.FunctionMap != nil {
		target.FunctionMap = o.FunctionMap
	}
}

func (o *Options[T]) ApplyOptions(opts []Option[T]) {
	for _, opt := range opts {
		opt.ApplyToOptions(o)
	}
}

type WithEagerEvaluation[T ref.References[T]] bool

var _ Option[ref.EmptyReferences] = WithEagerEvaluation[ref.EmptyReferences](false)

func (wee WithEagerEvaluation[T]) ApplyToOptions(target *Options[T]) {
	target.EagerEvaluation = bool(wee)
}

type WithFunctionMap[T ref.References[T]] struct{ fn.Map[T] }

var _ Option[ref.EmptyReferences] = WithFunctionMap[ref.EmptyReferences]{}

func (wfm WithFunctionMap[T]) ApplyToOptions(target *Options[T]) {
	target.FunctionMap = wfm.Map
}
