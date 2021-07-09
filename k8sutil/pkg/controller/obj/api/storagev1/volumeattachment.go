package storagev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	storagev1 "k8s.io/api/storage/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	VolumeAttachmentKind = storagev1.SchemeGroupVersion.WithKind("VolumeAttachment")
)

type VolumeAttachment struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *storagev1.VolumeAttachment
}

func makeVolumeAttachment(name string, obj *storagev1.VolumeAttachment) *VolumeAttachment {
	va := &VolumeAttachment{Name: name, Object: obj}
	va.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&va.Name, lifecycle.TypedObject{GVK: VolumeAttachmentKind, Object: va.Object})
	return va
}

func (va *VolumeAttachment) Copy() *VolumeAttachment {
	return makeVolumeAttachment(va.Name, va.Object.DeepCopy())
}

func NewVolumeAttachment(name string) *VolumeAttachment {
	return makeVolumeAttachment(name, &storagev1.VolumeAttachment{})
}

func NewVolumeAttachmentFromObject(obj *storagev1.VolumeAttachment) *VolumeAttachment {
	return makeVolumeAttachment(obj.GetName(), obj)
}

func NewVolumeAttachmentPatcher(upd, orig *VolumeAttachment) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
