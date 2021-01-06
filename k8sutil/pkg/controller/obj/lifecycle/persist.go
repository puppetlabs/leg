package lifecycle

import (
	"context"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Persister is the type of an entity that can be saved to a cluster.
type Persister interface {
	// Persist saves this entity using the given client.
	Persist(ctx context.Context, cl client.Client) error
}

// IgnoreNilPersister is an adapter for a persistable entity that makes sure the
// entity has a value.
type IgnoreNilPersister struct {
	Persister
}

// Persist saves this entity using the given client, or does nothing if the
// underlying entity is nil or is an interface with a nil value.
func (inp IgnoreNilPersister) Persist(ctx context.Context, cl client.Client) error {
	if inp.Persister == nil || reflect.ValueOf(inp.Persister).IsNil() {
		return nil
	}

	return inp.Persister.Persist(ctx, cl)
}
