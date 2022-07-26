package evaluate

import (
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Annotation struct {
	Name        string
	Attributes  map[string]any
	Description string
}

type Step[T ref.References[T]] struct {
	*Result[T]
	Annotations []Annotation
}

type Metadata[T ref.References[T]] struct {
	Annotation Annotation
	References T
	Steps      []*Step[T]
	Evaluator  Evaluator[T]
}

func (m *Metadata[T]) SetAnnotation(annotation Annotation) {
	m.Annotation = annotation
}

func (m *Metadata[T]) MergeReferences(other *Metadata[T]) {
	m.References = m.References.Merge(other.References)
}

func (m *Metadata[T]) InsertSteps(other *Metadata[T]) {
	old := m.Steps
	m.Steps = make([]*Step[T], 0, len(other.Steps)+len(old))
	m.MergeSteps(other)
	m.Steps = append(m.Steps, old...)
}

func (m *Metadata[T]) addStep(step *Step[T]) {
	m.Steps = append(m.Steps, step)
	m.MergeReferences(step.Metadata)
}

func (m *Metadata[T]) AddStep(step *Result[T], annotations ...Annotation) {
	m.addStep(&Step[T]{
		Result:      step,
		Annotations: annotations,
	})
}

func (m *Metadata[T]) MergeSteps(other *Metadata[T]) {
	for _, step := range other.Steps {
		m.addStep(step)
	}
}

func (m *Metadata[T]) SetEvaluator(ev Evaluator[T]) {
	m.Evaluator = ev
}

func (m *Metadata[T]) MergeEvaluator(other *Metadata[T]) {
	if other.Evaluator != nil {
		m.SetEvaluator(other.Evaluator)
	}
}

func (m *Metadata[T]) MergeEvaluatorIfNotSet(other *Metadata[T]) {
	if m.Evaluator != nil {
		return
	}

	m.MergeEvaluator(other)
}

func (m *Metadata[T]) Reduce() {
	m.Steps = nil
}

func (m *Metadata[T]) Resolved() bool {
	return m.References.Resolved()
}

func (m *Metadata[T]) OK() bool {
	return m.References.OK()
}

func EmptyMetadata[T ref.References[T]]() *Metadata[T] {
	return &Metadata[T]{}
}

func NewMetadata[T ref.References[T]](refs T) *Metadata[T] {
	md := EmptyMetadata[T]()
	md.References = refs
	return md
}
