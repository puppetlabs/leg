package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	LimitRangeKind = corev1.SchemeGroupVersion.WithKind("LimitRange")
)

type LimitRange struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.LimitRange
}

func makeLimitRange(key client.ObjectKey, obj *corev1.LimitRange) *LimitRange {
	lr := &LimitRange{Key: key, Object: obj}
	lr.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&lr.Key, lifecycle.TypedObject{GVK: LimitRangeKind, Object: lr.Object})
	return lr
}

func (lr *LimitRange) Copy() *LimitRange {
	return makeLimitRange(lr.Key, lr.Object.DeepCopy())
}

func (lr *LimitRange) MergeItem(item corev1.LimitRangeItem) {
	for i := range lr.Object.Spec.Limits {
		target := &lr.Object.Spec.Limits[i]

		if target.Type != item.Type {
			continue
		}

		if item.Max != nil {
			target.Max = item.Max
		}
		if item.Min != nil {
			target.Min = item.Min
		}
		if item.Default != nil {
			target.Default = item.Default
		}
		if item.DefaultRequest != nil {
			target.DefaultRequest = item.DefaultRequest
		}
		if item.MaxLimitRequestRatio != nil {
			target.MaxLimitRequestRatio = item.MaxLimitRequestRatio
		}

		return
	}

	lr.Object.Spec.Limits = append(lr.Object.Spec.Limits, item)
}

func (lr *LimitRange) SetContainerMax(lim corev1.ResourceList) {
	lr.MergeItem(corev1.LimitRangeItem{
		Type: corev1.LimitTypeContainer,
		Max:  lim,
	})
}

func (lr *LimitRange) SetContainerDefault(lim corev1.ResourceList) {
	lr.MergeItem(corev1.LimitRangeItem{
		Type:    corev1.LimitTypeContainer,
		Default: lim,
	})
}

func (lr *LimitRange) SetContainerDefaultRequest(lim corev1.ResourceList) {
	lr.MergeItem(corev1.LimitRangeItem{
		Type:           corev1.LimitTypeContainer,
		DefaultRequest: lim,
	})
}

func NewLimitRange(key client.ObjectKey) *LimitRange {
	return makeLimitRange(key, &corev1.LimitRange{})
}

func NewLimitRangeFromObject(obj *corev1.LimitRange) *LimitRange {
	return makeLimitRange(client.ObjectKeyFromObject(obj), obj)
}

func NewLimitRangePatcher(upd, orig *LimitRange) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
