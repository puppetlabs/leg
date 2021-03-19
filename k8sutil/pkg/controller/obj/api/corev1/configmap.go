package corev1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ConfigMapKind = corev1.SchemeGroupVersion.WithKind("ConfigMap")
)

type ConfigMap struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.ConfigMap
}

func makeConfigMap(key client.ObjectKey, obj *corev1.ConfigMap) *ConfigMap {
	cm := &ConfigMap{Key: key, Object: obj}
	cm.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&cm.Key, lifecycle.TypedObject{GVK: ConfigMapKind, Object: cm.Object})
	return cm

}

func (cm *ConfigMap) Copy() *ConfigMap {
	return makeConfigMap(cm.Key, cm.Object.DeepCopy())
}

func NewConfigMap(key client.ObjectKey) *ConfigMap {
	return makeConfigMap(key, &corev1.ConfigMap{})
}

func NewConfigMapFromObject(obj *corev1.ConfigMap) *ConfigMap {
	return makeConfigMap(client.ObjectKeyFromObject(obj), obj)
}

func NewConfigMapPatcher(upd, orig *ConfigMap) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
