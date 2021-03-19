package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	PersistentVolumeKind = corev1.SchemeGroupVersion.WithKind("PersistentVolume")
)

type PersistentVolume struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *corev1.PersistentVolume
}

func makePersistentVolume(name string, obj *corev1.PersistentVolume) *PersistentVolume {
	pv := &PersistentVolume{Name: name, Object: obj}
	pv.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&pv.Name, lifecycle.TypedObject{GVK: PersistentVolumeKind, Object: pv.Object})
	return pv
}

func (pv *PersistentVolume) Copy() *PersistentVolume {
	return makePersistentVolume(pv.Name, pv.Object.DeepCopy())
}

func NewPersistentVolume(name string) *PersistentVolume {
	return makePersistentVolume(name, &corev1.PersistentVolume{})
}

func NewPersistentVolumeFromObject(obj *corev1.PersistentVolume) *PersistentVolume {
	return makePersistentVolume(obj.GetName(), obj)
}

func NewPersistentVolumePatcher(upd, orig *PersistentVolume) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
