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
	PersistentVolumeKind = corev1.SchemeGroupVersion.WithKind("PersistentVolume")
)

type PersistentVolume struct {
	Name   string
	Object *corev1.PersistentVolume
}

var _ lifecycle.Deleter = &PersistentVolume{}
var _ lifecycle.Persister = &PersistentVolume{}
var _ lifecycle.Loader = &PersistentVolume{}
var _ lifecycle.Ownable = &PersistentVolume{}
var _ lifecycle.LabelAnnotatableFrom = &PersistentVolume{}

func (pv *PersistentVolume) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, pv.Object, opts...)
}

func (pv *PersistentVolume) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&pv.Object.ObjectMeta, from)
}

func (pv *PersistentVolume) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, client.ObjectKey{Name: pv.Name}, pv.Object)
}

func (pv *PersistentVolume) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(pv.Object, owner)
}

func (pv *PersistentVolume) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, pv.Object, helper.WithObjectKey(client.ObjectKey{Name: pv.Name})); err != nil {
		return err
	}

	pv.Name = pv.Object.GetName()
	return nil
}

func (pv *PersistentVolume) Copy() *PersistentVolume {
	return &PersistentVolume{
		Name:   pv.Name,
		Object: pv.Object.DeepCopy(),
	}
}

func NewPersistentVolume(name string) *PersistentVolume {
	return &PersistentVolume{
		Name:   name,
		Object: &corev1.PersistentVolume{},
	}
}

func NewPersistentVolumeFromObject(obj *corev1.PersistentVolume) *PersistentVolume {
	return &PersistentVolume{
		Name:   obj.GetName(),
		Object: obj,
	}
}

func NewPersistentVolumePatcher(upd, orig *PersistentVolume) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
