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
	ConfigMapKind = corev1.SchemeGroupVersion.WithKind("ConfigMap")
)

type ConfigMap struct {
	Key    client.ObjectKey
	Object *corev1.ConfigMap
}

var _ lifecycle.Deleter = &ConfigMap{}
var _ lifecycle.LabelAnnotatableFrom = &ConfigMap{}
var _ lifecycle.Loader = &ConfigMap{}
var _ lifecycle.Ownable = &ConfigMap{}
var _ lifecycle.Owner = &ConfigMap{}
var _ lifecycle.Persister = &ConfigMap{}

func (cm *ConfigMap) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, cm.Object, opts...)
}

func (cm *ConfigMap) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&cm.Object.ObjectMeta, from)
}

func (cm *ConfigMap) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, cm.Key, cm.Object)
}

func (cm *ConfigMap) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(cm.Object, owner)
}

func (cm *ConfigMap) Own(ctx context.Context, other lifecycle.Ownable) error {
	return other.Owned(ctx, lifecycle.TypedObject{GVK: ConfigMapKind, Object: cm.Object})
}

func (cm *ConfigMap) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, cm.Object, helper.WithObjectKey(cm.Key)); err != nil {
		return err
	}

	cm.Key = client.ObjectKeyFromObject(cm.Object)
	return nil
}

func (cm *ConfigMap) Copy() *ConfigMap {
	return &ConfigMap{
		Key:    cm.Key,
		Object: cm.Object.DeepCopy(),
	}
}

func NewConfigMap(key client.ObjectKey) *ConfigMap {
	return &ConfigMap{
		Key:    key,
		Object: &corev1.ConfigMap{},
	}
}

func NewConfigMapFromObject(obj *corev1.ConfigMap) *ConfigMap {
	return &ConfigMap{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewConfigMapPatcher(upd, orig *ConfigMap) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
