package lifecycle

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deleter is the type of an entity that can be deleted from a cluster.
type Deleter interface {
	// Delete removes this entity from the cluster.
	Delete(ctx context.Context, cl client.Client) (bool, error)
}
