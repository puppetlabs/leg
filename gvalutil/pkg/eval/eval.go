package eval

import "context"

type Indexable interface {
	Index(ctx context.Context, idx interface{}) (interface{}, error)
}
