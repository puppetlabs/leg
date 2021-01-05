package rbacv1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	RoleBindingKind = rbacv1.SchemeGroupVersion.WithKind("RoleBinding")
)

type RoleBinding struct {
	Key    client.ObjectKey
	Object *rbacv1.RoleBinding
}

var _ lifecycle.Deleter = &RoleBinding{}
var _ lifecycle.LabelAnnotatableFrom = &RoleBinding{}
var _ lifecycle.Loader = &RoleBinding{}
var _ lifecycle.Ownable = &RoleBinding{}
var _ lifecycle.Persister = &RoleBinding{}

func (rb *RoleBinding) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, rb.Object, opts...)
}

func (rb *RoleBinding) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&rb.Object.ObjectMeta, from)
}

func (rb *RoleBinding) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, rb.Key, rb.Object)
}

func (rb *RoleBinding) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, rb.Object, owner)
}

func (rb *RoleBinding) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, rb.Object, helper.WithObjectKey(rb.Key)); err != nil {
		return err
	}

	rb.Key = client.ObjectKeyFromObject(rb.Object)
	return nil
}

func (rb *RoleBinding) Copy() *RoleBinding {
	return &RoleBinding{
		Key:    rb.Key,
		Object: rb.Object.DeepCopy(),
	}
}

func NewRoleBinding(key client.ObjectKey) *RoleBinding {
	return &RoleBinding{
		Key:    key,
		Object: &rbacv1.RoleBinding{},
	}
}

func NewRoleBindingFromObject(obj *rbacv1.RoleBinding) *RoleBinding {
	return &RoleBinding{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewRoleBindingPatcher(upd, orig *RoleBinding) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
