package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	NodeKind = corev1.SchemeGroupVersion.WithKind("Node")
)

type Node struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *corev1.Node
}

func makeNode(name string, obj *corev1.Node) *Node {
	n := &Node{Name: name, Object: obj}
	n.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&n.Name, lifecycle.TypedObject{GVK: NodeKind, Object: n.Object})
	return n
}

func (n *Node) Copy() *Node {
	return makeNode(n.Name, n.Object.DeepCopy())
}

func NewNode(name string) *Node {
	return makeNode(name, &corev1.Node{})
}

func NewNodeFromObject(obj *corev1.Node) *Node {
	return makeNode(obj.GetName(), obj)
}

func NewNodePatcher(upd, orig *Node) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
