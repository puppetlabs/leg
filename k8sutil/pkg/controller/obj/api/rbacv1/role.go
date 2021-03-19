package rbacv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	RoleKind = rbacv1.SchemeGroupVersion.WithKind("Role")
)

type Role struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *rbacv1.Role
}

func makeRole(key client.ObjectKey, obj *rbacv1.Role) *Role {
	r := &Role{Key: key, Object: obj}
	r.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&r.Key, lifecycle.TypedObject{GVK: RoleKind, Object: r.Object})
	return r
}

func (r *Role) Copy() *Role {
	return makeRole(r.Key, r.Object.DeepCopy())
}

func NewRole(key client.ObjectKey) *Role {
	return makeRole(key, &rbacv1.Role{})
}

func NewRoleFromObject(obj *rbacv1.Role) *Role {
	return makeRole(client.ObjectKeyFromObject(obj), obj)
}

func NewRolePatcher(upd, orig *Role) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
