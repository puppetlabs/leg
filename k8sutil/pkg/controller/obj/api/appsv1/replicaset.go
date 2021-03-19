package appsv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ReplicaSetKind = appsv1.SchemeGroupVersion.WithKind("ReplicaSet")
)

type ReplicaSet struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *appsv1.ReplicaSet
}

func makeReplicaSet(key client.ObjectKey, obj *appsv1.ReplicaSet) *ReplicaSet {
	rs := &ReplicaSet{Key: key, Object: obj}
	rs.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&rs.Key, lifecycle.TypedObject{GVK: ReplicaSetKind, Object: rs.Object})
	return rs
}

func (rs *ReplicaSet) Copy() *ReplicaSet {
	return makeReplicaSet(rs.Key, rs.Object.DeepCopy())
}

func NewReplicaSet(key client.ObjectKey) *ReplicaSet {
	return makeReplicaSet(key, &appsv1.ReplicaSet{})
}

func NewReplicaSetFromObject(obj *appsv1.ReplicaSet) *ReplicaSet {
	return makeReplicaSet(client.ObjectKeyFromObject(obj), obj)
}

func NewReplicaSetPatcher(upd, orig *ReplicaSet) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
