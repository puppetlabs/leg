package helper

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NamespaceScopedAPIObject allows an arbitrary Kubernetes object to be
// represented using the lifecycle constructs or embedded into a specific type.
type NamespaceScopedAPIObject struct {
	key *client.ObjectKey
	obj lifecycle.TypedObject
}

var (
	_ lifecycle.Deleter              = &NamespaceScopedAPIObject{}
	_ lifecycle.Finalizable          = &NamespaceScopedAPIObject{}
	_ lifecycle.LabelAnnotatableFrom = &NamespaceScopedAPIObject{}
	_ lifecycle.Ownable              = &NamespaceScopedAPIObject{}
	_ lifecycle.Owner                = &NamespaceScopedAPIObject{}
	_ lifecycle.Persister            = &NamespaceScopedAPIObject{}
)

func (ao *NamespaceScopedAPIObject) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return DeleteIgnoreNotFound(ctx, cl, ao.obj.Object, opts...)
}

func (ao *NamespaceScopedAPIObject) Finalizing() bool {
	return !ao.obj.Object.GetDeletionTimestamp().IsZero()
}

func (ao *NamespaceScopedAPIObject) AddFinalizer(ctx context.Context, name string) bool {
	return AddFinalizer(ao.obj.Object, name)
}

func (ao *NamespaceScopedAPIObject) RemoveFinalizer(ctx context.Context, name string) bool {
	return RemoveFinalizer(ao.obj.Object, name)
}

func (ao *NamespaceScopedAPIObject) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	CopyLabelsAndAnnotations(ao.obj.Object, from)
}

func (ao *NamespaceScopedAPIObject) Load(ctx context.Context, cl client.Client) (bool, error) {
	return GetIgnoreNotFound(ctx, cl, *ao.key, ao.obj.Object)
}

func (ao *NamespaceScopedAPIObject) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return Own(ao.obj.Object, owner)
}

func (ao *NamespaceScopedAPIObject) Own(ctx context.Context, other lifecycle.Ownable) error {
	return other.Owned(ctx, ao.obj)
}

func (ao *NamespaceScopedAPIObject) Persist(ctx context.Context, cl client.Client) error {
	if err := CreateOrUpdate(ctx, cl, ao.obj.Object, WithObjectKey(*ao.key)); err != nil {
		return err
	}

	*ao.key = client.ObjectKeyFromObject(ao.obj.Object)
	return nil
}

// ForNamespaceScopedAPIObject creates a lifecycle-compatible representation of
// a Kubernetes object that will automatically reflect changes back to the given
// arguments.
func ForNamespaceScopedAPIObject(key *client.ObjectKey, obj lifecycle.TypedObject) *NamespaceScopedAPIObject {
	obj.Object.SetNamespace(key.Namespace)
	obj.Object.SetName(key.Name)

	return &NamespaceScopedAPIObject{
		key: key,
		obj: obj,
	}
}

// ClusterScopedAPIObject allows an arbitrary Kubernetes object to be
// represented using the lifecycle constructs or embedded into a specific type.
type ClusterScopedAPIObject struct {
	name *string
	obj  lifecycle.TypedObject
}

var (
	_ lifecycle.Deleter              = &ClusterScopedAPIObject{}
	_ lifecycle.Finalizable          = &ClusterScopedAPIObject{}
	_ lifecycle.LabelAnnotatableFrom = &ClusterScopedAPIObject{}
	_ lifecycle.Ownable              = &ClusterScopedAPIObject{}
	_ lifecycle.Owner                = &ClusterScopedAPIObject{}
	_ lifecycle.Persister            = &ClusterScopedAPIObject{}
)

func (ao *ClusterScopedAPIObject) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return DeleteIgnoreNotFound(ctx, cl, ao.obj.Object, opts...)
}

func (ao *ClusterScopedAPIObject) Finalizing() bool {
	return !ao.obj.Object.GetDeletionTimestamp().IsZero()
}

func (ao *ClusterScopedAPIObject) AddFinalizer(ctx context.Context, name string) bool {
	return AddFinalizer(ao.obj.Object, name)
}

func (ao *ClusterScopedAPIObject) RemoveFinalizer(ctx context.Context, name string) bool {
	return RemoveFinalizer(ao.obj.Object, name)
}

func (ao *ClusterScopedAPIObject) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	CopyLabelsAndAnnotations(ao.obj.Object, from)
}

func (ao *ClusterScopedAPIObject) Load(ctx context.Context, cl client.Client) (bool, error) {
	return GetIgnoreNotFound(ctx, cl, client.ObjectKey{Name: *ao.name}, ao.obj.Object)
}

func (ao *ClusterScopedAPIObject) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return Own(ao.obj.Object, owner)
}

func (ao *ClusterScopedAPIObject) Own(ctx context.Context, other lifecycle.Ownable) error {
	return other.Owned(ctx, ao.obj)
}

func (ao *ClusterScopedAPIObject) Persist(ctx context.Context, cl client.Client) error {
	if err := CreateOrUpdate(ctx, cl, ao.obj.Object, WithObjectKey(client.ObjectKey{Name: *ao.name})); err != nil {
		return err
	}

	*ao.name = ao.obj.Object.GetName()
	return nil
}

// ForClusterScopedAPIObject creates a lifecycle-compatible representation of a
// Kubernetes object that will automatically reflect changes back to the given
// arguments.
func ForClusterScopedAPIObject(name *string, obj lifecycle.TypedObject) *ClusterScopedAPIObject {
	obj.Object.SetName(*name)

	return &ClusterScopedAPIObject{
		name: name,
		obj:  obj,
	}
}
