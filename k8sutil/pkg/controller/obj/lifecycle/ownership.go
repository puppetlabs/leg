package lifecycle

import (
	"context"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Ownable is the type of an entity that can be a dependent of another object
// using the Kubernetes ownerReferences mechanism.
type Ownable interface {
	// Owned sets the ownerReferences controller for this entity to the given
	// object.
	Owned(ctx context.Context, owner TypedObject) error
}

// IgnoreNilOwnable is an adapter for an ownable entity that makes sure the
// entity has a value.
type IgnoreNilOwnable struct {
	Ownable
}

// Owned sets the ownerReferences controller for this entity to the given
// object, or does nothing if the underlying entity is nil or an interface with
// a nil value.
func (ino IgnoreNilOwnable) Owned(ctx context.Context, owner TypedObject) error {
	if ino.Ownable == nil || reflect.ValueOf(ino.Ownable).IsNil() {
		return nil
	}

	return ino.Ownable.Owned(ctx, owner)
}

// OwnablePersister is a combined interface for an object that can be owned and
// can be saved to the cluster.
type OwnablePersister interface {
	Ownable
	Persister
}

// IgnoreNilOwnablePersister combines IgnoreNilOwnable and IgnoreNilPersister.
type IgnoreNilOwnablePersister struct {
	OwnablePersister
}

// Owned sets the ownerReferences controller for this entity to the given
// object, or does nothing if the underlying entity is nil or an interface with
// a nil value.
func (inop IgnoreNilOwnablePersister) Owned(ctx context.Context, owner TypedObject) error {
	return IgnoreNilOwnable{inop.OwnablePersister}.Owned(ctx, owner)
}

// Persist saves this entity using the given client, or does nothing if the
// underlying entity is nil or is an interface with a nil value.
func (inop IgnoreNilOwnablePersister) Persist(ctx context.Context, cl client.Client) error {
	return IgnoreNilPersister{inop.OwnablePersister}.Persist(ctx, cl)
}

// OwnablePersisters allows a collection of ownable persisters to be used as a
// single entity.
type OwnablePersisters []OwnablePersister

var _ OwnablePersister = OwnablePersisters(nil)

// Owned sets the ownerReferences controller for each entity in this collection
// to the given object.
func (ops OwnablePersisters) Owned(ctx context.Context, owner TypedObject) error {
	for _, op := range ops {
		if err := op.Owned(ctx, owner); err != nil {
			return err
		}
	}

	return nil
}

// Persist saves each of the entities in this collection to the cluster using
// the given client.
func (ops OwnablePersisters) Persist(ctx context.Context, cl client.Client) error {
	for _, op := range ops {
		if err := op.Persist(ctx, cl); err != nil {
			return err
		}
	}

	return nil
}

// Owner is the type of an entity that can own other objects using the
// Kubernetes ownerReferences mechanism.
type Owner interface {
	// Own sets the owner of the specified dependent entity to this entity.
	Own(ctx context.Context, other Ownable) error
}

// OwnerPersister is a combined interface for an object that can own other
// objects and can be saved to the cluster.
type OwnerPersister interface {
	Owner
	Persister
}

// OwnershipPersister provides a lifecycle for managing the ownership of an
// entity or set of entities.
//
// It first persists the owner so that its unique object ID is known. Then it
// sets the ownerReferences for the intended dependent. Finally, it persists all
// the dependents.
type OwnershipPersister struct {
	Owner     OwnerPersister
	Dependent OwnablePersister
}

var _ Persister = OwnershipPersister{}

// Persist saves the owner and dependent, guaranteeing an ownerReference
// relationship between the dependent and the owner.
func (op OwnershipPersister) Persist(ctx context.Context, cl client.Client) error {
	if err := op.Owner.Persist(ctx, cl); err != nil {
		return err
	}

	if err := op.Owner.Own(ctx, op.Dependent); err != nil {
		return err
	}

	if err := op.Dependent.Persist(ctx, cl); err != nil {
		return err
	}

	return nil
}
