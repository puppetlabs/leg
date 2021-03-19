package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	NamespaceKind = corev1.SchemeGroupVersion.WithKind("Namespace")
)

type Namespace struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *corev1.Namespace
}

func makeNamespace(name string, obj *corev1.Namespace) *Namespace {
	n := &Namespace{Name: name, Object: obj}
	n.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&n.Name, lifecycle.TypedObject{GVK: NamespaceKind, Object: n.Object})
	return n

}

func (n *Namespace) Copy() *Namespace {
	return makeNamespace(n.Name, n.Object.DeepCopy())
}

func NewNamespace(name string) *Namespace {
	return makeNamespace(name, &corev1.Namespace{})
}

func NewNamespaceFromObject(obj *corev1.Namespace) *Namespace {
	return makeNamespace(obj.GetName(), obj)
}

func NewNamespacePatcher(upd, orig *Namespace) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
