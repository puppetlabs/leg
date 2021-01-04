package lifecycle

import (
	"context"

	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Finalizable is the type of an entity that supports finalizing.
type Finalizable interface {
	// Finalizing returns true if this entity is currently in a finalizing
	// state.
	//
	// Usually, this is indicated by the presence of a deletion timestamp in
	// Kubernetes object metadata.
	Finalizing() bool

	// AddFinalizer sets the given finalizer in object metadata. It returns true
	// if the finalizer was added, and false if it already exists.
	AddFinalizer(ctx context.Context, name string) bool

	// RemoveFinalizer removes the given finalizer from object metadata. It
	// returns true if the finalizer was removed, and false if the finalizer did
	// not exist.
	RemoveFinalizer(ctx context.Context, name string) bool
}

// FinalizablePersister is a combined interface for a finalizable entity that
// can be saved to the cluster.
type FinalizablePersister interface {
	Finalizable
	Persister
}

// Finalize provides a lifecycle for finalization for an entity that implements
// FinalizablePersister.
//
// If the entity is not in a finalizing state, a finalizer is added with the
// given name. The entity is automatically persisted.
//
// If the entity is in a finalizing state, the given callback function is run.
// If it succeeds, the finalizer is removed and the entity is persisted.
//
// This function returns true if the finalization callback successfully ran and
// the entity was updated. It returns false otherwise.
func Finalize(ctx context.Context, cl client.Client, name string, obj FinalizablePersister, run func() error) (bool, error) {
	if obj.Finalizing() {
		klog.Infof("running finalizer %s", name)

		if err := run(); err != nil {
			return false, err
		}

		obj.RemoveFinalizer(ctx, name)

		err := obj.Persist(ctx, cl)
		return err == nil, err
	} else if obj.AddFinalizer(ctx, name) {
		return false, obj.Persist(ctx, cl)
	}

	return false, nil
}
