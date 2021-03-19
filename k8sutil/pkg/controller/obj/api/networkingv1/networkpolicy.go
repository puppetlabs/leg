package networkingv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	NetworkPolicyKind = networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy")
)

type NetworkPolicy struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *networkingv1.NetworkPolicy
}

func makeNetworkPolicy(key client.ObjectKey, obj *networkingv1.NetworkPolicy) *NetworkPolicy {
	np := &NetworkPolicy{Key: key, Object: obj}
	np.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&np.Key, lifecycle.TypedObject{GVK: NetworkPolicyKind, Object: np.Object})
	return np
}

func (np *NetworkPolicy) Copy() *NetworkPolicy {
	return makeNetworkPolicy(np.Key, np.Object.DeepCopy())
}

func (np *NetworkPolicy) AllowAll() {
	np.Object.Spec = networkingv1.NetworkPolicySpec{
		Ingress: []networkingv1.NetworkPolicyIngressRule{{}},
		Egress:  []networkingv1.NetworkPolicyEgressRule{{}},
		PolicyTypes: []networkingv1.PolicyType{
			networkingv1.PolicyTypeIngress,
			networkingv1.PolicyTypeEgress,
		},
	}
}

func (np *NetworkPolicy) DenyAll() {
	np.Object.Spec = networkingv1.NetworkPolicySpec{
		PolicyTypes: []networkingv1.PolicyType{
			networkingv1.PolicyTypeIngress,
			networkingv1.PolicyTypeEgress,
		},
	}
}

func NewNetworkPolicy(key client.ObjectKey) *NetworkPolicy {
	return makeNetworkPolicy(key, &networkingv1.NetworkPolicy{})
}

func NewNetworkPolicyFromObject(obj *networkingv1.NetworkPolicy) *NetworkPolicy {
	return makeNetworkPolicy(client.ObjectKeyFromObject(obj), obj)
}

func NewNetworkPolicyPatcher(upd, orig *NetworkPolicy) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
