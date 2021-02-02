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
	LimitRangeKind = corev1.SchemeGroupVersion.WithKind("LimitRange")
)

type LimitRange struct {
	Key    client.ObjectKey
	Object *corev1.LimitRange
}

var _ lifecycle.Deleter = &LimitRange{}
var _ lifecycle.LabelAnnotatableFrom = &LimitRange{}
var _ lifecycle.Loader = &LimitRange{}
var _ lifecycle.Ownable = &LimitRange{}
var _ lifecycle.Persister = &LimitRange{}

func (lr *LimitRange) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, lr.Object, opts...)
}

func (lr *LimitRange) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&lr.Object.ObjectMeta, from)
}

func (lr *LimitRange) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, lr.Key, lr.Object)
}

func (lr *LimitRange) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(lr.Object, owner)
}

func (lr *LimitRange) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, lr.Object, helper.WithObjectKey(lr.Key)); err != nil {
		return err
	}

	lr.Key = client.ObjectKeyFromObject(lr.Object)
	return nil
}

func (lr *LimitRange) Copy() *LimitRange {
	return &LimitRange{
		Key:    lr.Key,
		Object: lr.Object.DeepCopy(),
	}
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
	return &LimitRange{
		Key:    key,
		Object: &corev1.LimitRange{},
	}
}

func NewLimitRangeFromObject(obj *corev1.LimitRange) *LimitRange {
	return &LimitRange{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewLimitRangePatcher(upd, orig *LimitRange) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
