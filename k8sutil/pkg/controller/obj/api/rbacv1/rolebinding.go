package rbacv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	RoleBindingKind = rbacv1.SchemeGroupVersion.WithKind("RoleBinding")
)

type RoleBinding struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *rbacv1.RoleBinding
}

func makeRoleBinding(key client.ObjectKey, obj *rbacv1.RoleBinding) *RoleBinding {
	rb := &RoleBinding{Key: key, Object: obj}
	rb.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&rb.Key, lifecycle.TypedObject{GVK: RoleBindingKind, Object: rb.Object})
	return rb
}

func (rb *RoleBinding) Copy() *RoleBinding {
	return makeRoleBinding(rb.Key, rb.Object.DeepCopy())
}

func NewRoleBinding(key client.ObjectKey) *RoleBinding {
	return makeRoleBinding(key, &rbacv1.RoleBinding{})
}

func NewRoleBindingFromObject(obj *rbacv1.RoleBinding) *RoleBinding {
	return makeRoleBinding(client.ObjectKeyFromObject(obj), obj)
}

func NewRoleBindingPatcher(upd, orig *RoleBinding) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
