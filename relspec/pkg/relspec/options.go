package relspec

import (
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Options[T ref.References[T]] struct {
	FunctionMap          fn.Map[T]
	MappingTypeResolvers map[string]MappingTypeResolver[T]
	TemplateEnvironment  map[string]evaluate.Expandable[T]
	KeepSteps            bool
	EagerEvaluation      bool
}

type Option[T ref.References[T]] interface {
	ApplyToOptions(target *Options[T])
}

func (o *Options[T]) ApplyOptions(opts []Option[T]) {
	for _, opt := range opts {
		opt.ApplyToOptions(o)
	}
}

type WithFunctionMap[T ref.References[T]] struct{ fn.Map[T] }

var _ Option[ref.EmptyReferences] = WithFunctionMap[ref.EmptyReferences]{}

func (wfm WithFunctionMap[T]) ApplyToOptions(target *Options[T]) {
	target.FunctionMap = wfm.Map
}

type WithMappingTypeResolvers[T ref.References[T]] map[string]MappingTypeResolver[T]

var _ Option[ref.EmptyReferences] = WithMappingTypeResolvers[ref.EmptyReferences](nil)

func (wmtr WithMappingTypeResolvers[T]) ApplyToOptions(target *Options[T]) {
	target.MappingTypeResolvers = map[string]MappingTypeResolver[T](wmtr)
}

type WithTemplateEnvironment[T ref.References[T]] map[string]evaluate.Expandable[T]

var _ Option[ref.EmptyReferences] = WithTemplateEnvironment[ref.EmptyReferences](nil)

func (wte WithTemplateEnvironment[T]) ApplyToOptions(target *Options[T]) {
	target.TemplateEnvironment = map[string]evaluate.Expandable[T](wte)
}

type WithKeepSteps[T ref.References[T]] bool

var _ Option[ref.EmptyReferences] = WithKeepSteps[ref.EmptyReferences](false)

func (wks WithKeepSteps[T]) ApplyToOptions(target *Options[T]) {
	target.KeepSteps = bool(wks)
}

type WithEagerEvaluation[T ref.References[T]] bool

var _ Option[ref.EmptyReferences] = WithEagerEvaluation[ref.EmptyReferences](false)

func (wee WithEagerEvaluation[T]) ApplyToOptions(target *Options[T]) {
	target.EagerEvaluation = bool(wee)
}
