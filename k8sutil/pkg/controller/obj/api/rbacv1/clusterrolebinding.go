package rbacv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ClusterRoleBindingKind = rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
)

type ClusterRoleBinding struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *rbacv1.ClusterRoleBinding
}

func makeClusterRoleBinding(name string, obj *rbacv1.ClusterRoleBinding) *ClusterRoleBinding {
	rb := &ClusterRoleBinding{Name: name, Object: obj}
	rb.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(
		&name,
		lifecycle.TypedObject{GVK: ClusterRoleBindingKind, Object: rb.Object},
	)

	return rb
}

func (rb *ClusterRoleBinding) Copy() *ClusterRoleBinding {
	return makeClusterRoleBinding(rb.Name, rb.Object.DeepCopy())
}

func NewClusterRoleBinding(name string) *ClusterRoleBinding {
	return makeClusterRoleBinding(name, &rbacv1.ClusterRoleBinding{})
}

func NewClusterRoleBindingFromObject(obj *rbacv1.ClusterRoleBinding) *ClusterRoleBinding {
	return makeClusterRoleBinding(obj.GetName(), obj)
}

func NewClusterRoleBindingPatcher(upd, orig *ClusterRoleBinding) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
