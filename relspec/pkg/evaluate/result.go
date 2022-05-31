package evaluate

import (
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Result[T ref.References[T]] struct {
	*Metadata[T]
	Value any
}

func (r *Result[T]) SetValue(v any) {
	r.Value = v
}

func NewResult[T ref.References[T]](m *Metadata[T], v any) *Result[T] {
	return &Result[T]{
		Metadata: m,
		Value:    v,
	}
}

func ContextualizedResult[T ref.References[T]](m *Metadata[T]) *Result[T] {
	return NewResult(m, nil)
}

func EmptyResult[T ref.References[T]]() *Result[T] {
	return ContextualizedResult(EmptyMetadata[T]())
}

func StaticResult[T ref.References[T]](v any) *Result[T] {
	return NewResult(EmptyMetadata[T](), v)
}
