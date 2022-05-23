package eval

import "context"

type Indexable interface {
	Index(ctx context.Context, idx any) (any, error)
}
