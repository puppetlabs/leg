package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	PersistentVolumeClaimKind = corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")
)

type PersistentVolumeClaim struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.PersistentVolumeClaim
}

func makePersistentVolumeClaim(key client.ObjectKey, obj *corev1.PersistentVolumeClaim) *PersistentVolumeClaim {
	pvc := &PersistentVolumeClaim{Key: key, Object: obj}
	pvc.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&pvc.Key, lifecycle.TypedObject{GVK: PersistentVolumeClaimKind, Object: pvc.Object})
	return pvc
}

func (pvc *PersistentVolumeClaim) Copy() *PersistentVolumeClaim {
	return makePersistentVolumeClaim(pvc.Key, pvc.Object.DeepCopy())
}

func NewPersistentVolumeClaim(key client.ObjectKey) *PersistentVolumeClaim {
	return makePersistentVolumeClaim(key, &corev1.PersistentVolumeClaim{})
}

func NewPersistentVolumeClaimFromObject(obj *corev1.PersistentVolumeClaim) *PersistentVolumeClaim {
	return makePersistentVolumeClaim(client.ObjectKeyFromObject(obj), obj)
}

func NewPersistentVolumeClaimPatcher(upd, orig *PersistentVolumeClaim) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
