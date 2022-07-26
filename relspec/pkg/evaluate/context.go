package evaluate

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type contextKey int

const (
	metadataContextKey contextKey = iota
)

func ContextWithMetadata[T ref.References[T]](ctx context.Context, md *Metadata[T]) context.Context {
	return context.WithValue(ctx, metadataContextKey, md)
}

func MetadataFromContext[T ref.References[T]](ctx context.Context) *Metadata[T] {
	md, ok := ctx.Value(metadataContextKey).(*Metadata[T])
	if !ok {
		md = EmptyMetadata[T]()
	}
	return md
}

func CoalesceEvaluatorsPreferringContext[T ref.References[T]](ctx context.Context, evs ...Evaluator[T]) Evaluator[T] {
	return CoalesceEvaluators(append([]Evaluator[T]{MetadataFromContext[T](ctx).Evaluator}, evs...)...)
}

func ContextWithEmptyMetadata[T ref.References[T]](ctx context.Context) (context.Context, *Metadata[T]) {
	md := EmptyMetadata[T]()
	md.MergeEvaluatorIfNotSet(MetadataFromContext[T](ctx))
	return ContextWithMetadata(ctx, md), md
}
