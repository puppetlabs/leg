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
	PersistentVolumeClaimKind = corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")
)

type PersistentVolumeClaim struct {
	Key    client.ObjectKey
	Object *corev1.PersistentVolumeClaim
}

var _ lifecycle.Deleter = &PersistentVolumeClaim{}
var _ lifecycle.LabelAnnotatableFrom = &PersistentVolumeClaim{}
var _ lifecycle.Loader = &PersistentVolumeClaim{}
var _ lifecycle.Ownable = &PersistentVolumeClaim{}
var _ lifecycle.Persister = &PersistentVolumeClaim{}

func (pvc *PersistentVolumeClaim) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, pvc.Object, opts...)
}

func (pvc *PersistentVolumeClaim) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&pvc.Object.ObjectMeta, from)
}

func (pvc *PersistentVolumeClaim) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, pvc.Key, pvc.Object)
}

func (pvc *PersistentVolumeClaim) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, pvc.Object, owner)
}

func (pvc *PersistentVolumeClaim) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, pvc.Object, helper.WithObjectKey(pvc.Key)); err != nil {
		return err
	}

	pvc.Key = client.ObjectKeyFromObject(pvc.Object)
	return nil
}

func (pvc *PersistentVolumeClaim) Copy() *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		Key:    pvc.Key,
		Object: pvc.Object.DeepCopy(),
	}
}

func NewPersistentVolumeClaim(key client.ObjectKey) *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		Key:    key,
		Object: &corev1.PersistentVolumeClaim{},
	}
}

func NewPersistentVolumeClaimFromObject(obj *corev1.PersistentVolumeClaim) *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewPersistentVolumeClaimPatcher(upd, orig *PersistentVolumeClaim) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
