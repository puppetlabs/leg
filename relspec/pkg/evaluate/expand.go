package evaluate

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

type Expandable[T ref.References[T]] interface {
	Expand(ctx context.Context, depth int) (*Result[T], error)
}

type staticExpandable[T ref.References[T]] struct {
	r *Result[T]
}

var _ Expandable[ref.EmptyReferences] = &staticExpandable[ref.EmptyReferences]{}

func (se *staticExpandable[T]) Expand(ctx context.Context, depth int) (*Result[T], error) {
	return se.r, nil
}

func StaticExpandable[T ref.References[T]](r *Result[T]) Expandable[T] {
	return &staticExpandable[T]{r: r}
}

type memoizedExpandable[T ref.References[T]] struct {
	delegate Expandable[T]
	cache    map[int]*Result[T]
}

func (me *memoizedExpandable[T]) Expand(ctx context.Context, depth int) (r *Result[T], err error) {
	if IsFullDepth(depth) {
		depth = DepthFull
	}

	r, ok := me.cache[depth]
	if !ok {
		r, err = me.delegate.Expand(ctx, depth)
		if err != nil {
			return
		}

		me.cache[depth] = r
	}

	return
}

func MemoizedExpandable[T ref.References[T]](delegate Expandable[T]) Expandable[T] {
	return &memoizedExpandable[T]{
		delegate: delegate,
		cache:    make(map[int]*Result[T]),
	}
}
