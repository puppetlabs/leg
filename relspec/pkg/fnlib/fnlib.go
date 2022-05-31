package fnlib

import (
	"github.com/puppetlabs/leg/relspec/pkg/fn"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

// Library creates an fn.Map of all the core functions supported
// by the platform.
func Library[T ref.References[T]]() fn.Map[T] {
	return fn.NewMap(map[string]fn.Descriptor[T]{
		"append":          AppendDescriptor[T](),
		"coalesce":        CoalesceDescriptor[T](),
		"concat":          ConcatDescriptor[T](),
		"convertMarkdown": ConvertMarkdownDescriptor[T](),
		"equals":          EqualsDescriptor[T](),
		"exists":          ExistsDescriptor[T](),
		"jsonMarshal":     JSONMarshalDescriptor[T](),
		"jsonUnmarshal":   JSONUnmarshalDescriptor[T](),
		"merge":           MergeDescriptor[T](),
		"notEquals":       NotEqualsDescriptor[T](),
		"now":             NowDescriptor[T](),
		"path":            PathDescriptor[T](),
		"toString":        ToStringDescriptor[T](),
	})
}
