package corev1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	NamespaceKind = corev1.SchemeGroupVersion.WithKind("Namespace")
)

type Namespace struct {
	Name   string
	Object *corev1.Namespace
}

var _ lifecycle.Deleter = &Namespace{}
var _ lifecycle.LabelAnnotatableFrom = &Namespace{}
var _ lifecycle.Loader = &Namespace{}
var _ lifecycle.Persister = &Namespace{}

func (n *Namespace) Delete(ctx context.Context, cl client.Client) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, n.Object)
}

func (n *Namespace) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&n.Object.ObjectMeta, from)
}

func (n *Namespace) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, client.ObjectKey{Name: n.Name}, n.Object)
}

func (n *Namespace) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, n.Object, helper.WithObjectKey(client.ObjectKey{Name: n.Name})); err != nil {
		return err
	}

	n.Name = n.Object.GetName()
	return nil
}

func (n *Namespace) Copy() *Namespace {
	return &Namespace{
		Name:   n.Name,
		Object: n.Object.DeepCopy(),
	}
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		Name:   name,
		Object: &corev1.Namespace{},
	}
}

func NewNamespaceFromObject(obj *corev1.Namespace) *Namespace {
	return &Namespace{
		Name:   obj.GetName(),
		Object: obj,
	}
}

func NewNamespacePatcher(upd, orig *Namespace) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
