package query

import (
	"context"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type evalExpandable[T ref.References[T]] struct {
	eval      gval.Evaluable
	parameter any
}

var _ evaluate.Expandable[ref.EmptyReferences] = &evalExpandable[ref.EmptyReferences]{}

func (ee *evalExpandable[T]) Expand(ctx context.Context, depth int) (*evaluate.Result[T], error) {
	ctx, md := evaluate.ContextWithEmptyMetadata[T](ctx)

	v, err := ee.eval(ctx, ee.parameter)
	if err != nil {
		return nil, err
	}

	return evaluate.NewResult(md, v), nil
}

func EvalExpandable[T ref.References[T]](eval gval.Evaluable, parameter any) evaluate.Expandable[T] {
	return &evalExpandable[T]{
		eval:      eval,
		parameter: parameter,
	}
}
