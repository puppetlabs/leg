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
	RoleKind = rbacv1.SchemeGroupVersion.WithKind("Role")
)

type Role struct {
	Key    client.ObjectKey
	Object *rbacv1.Role
}

var _ lifecycle.Deleter = &Role{}
var _ lifecycle.LabelAnnotatableFrom = &Role{}
var _ lifecycle.Loader = &Role{}
var _ lifecycle.Ownable = &Role{}
var _ lifecycle.Persister = &Role{}

func (r *Role) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, r.Object, opts...)
}

func (r *Role) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&r.Object.ObjectMeta, from)
}

func (r *Role) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, r.Key, r.Object)
}

func (r *Role) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, r.Object, owner)
}

func (r *Role) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, r.Object, helper.WithObjectKey(r.Key)); err != nil {
		return err
	}

	r.Key = client.ObjectKeyFromObject(r.Object)
	return nil
}

func (r *Role) Copy() *Role {
	return &Role{
		Key:    r.Key,
		Object: r.Object.DeepCopy(),
	}
}

func NewRole(key client.ObjectKey) *Role {
	return &Role{
		Key:    key,
		Object: &rbacv1.Role{},
	}
}

func NewRoleFromObject(obj *rbacv1.Role) *Role {
	return &Role{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewRolePatcher(upd, orig *Role) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
