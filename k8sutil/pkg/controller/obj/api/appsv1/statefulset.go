package appsv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	StatefulSetKind = appsv1.SchemeGroupVersion.WithKind("StatefulSet")
)

type StatefulSet struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *appsv1.StatefulSet
}

func makeStatefulSet(key client.ObjectKey, obj *appsv1.StatefulSet) *StatefulSet {
	rs := &StatefulSet{Key: key, Object: obj}
	rs.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&rs.Key, lifecycle.TypedObject{GVK: StatefulSetKind, Object: rs.Object})
	return rs
}

func (ss *StatefulSet) Copy() *StatefulSet {
	return makeStatefulSet(ss.Key, ss.Object.DeepCopy())
}

func NewStatefulSet(key client.ObjectKey) *StatefulSet {
	return makeStatefulSet(key, &appsv1.StatefulSet{})
}

func NewStatefulSetFromObject(obj *appsv1.StatefulSet) *StatefulSet {
	return makeStatefulSet(client.ObjectKeyFromObject(obj), obj)
}

func NewStatefulSetPatcher(upd, orig *StatefulSet) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
