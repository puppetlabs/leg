package rbacv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ClusterRoleKind = rbacv1.SchemeGroupVersion.WithKind("ClusterRole")
)

type ClusterRole struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *rbacv1.ClusterRole
}

func makeClusterRole(name string, obj *rbacv1.ClusterRole) *ClusterRole {
	r := &ClusterRole{Name: name, Object: obj}
	r.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(
		&name,
		lifecycle.TypedObject{GVK: ClusterRoleKind, Object: r.Object},
	)

	return r
}

func (r *ClusterRole) Copy() *ClusterRole {
	return makeClusterRole(r.Name, r.Object.DeepCopy())
}

func NewClusterRole(name string) *ClusterRole {
	return makeClusterRole(name, &rbacv1.ClusterRole{})
}

func NewClusterRoleFromObject(obj *rbacv1.ClusterRole) *ClusterRole {
	return makeClusterRole(obj.GetName(), obj)
}

func NewClusterRolePatcher(upd, orig *ClusterRole) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
