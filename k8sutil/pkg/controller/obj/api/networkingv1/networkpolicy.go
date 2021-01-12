package networkingv1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	NetworkPolicyKind = networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy")
)

type NetworkPolicy struct {
	Key    client.ObjectKey
	Object *networkingv1.NetworkPolicy
}

var _ lifecycle.Deleter = &NetworkPolicy{}
var _ lifecycle.LabelAnnotatableFrom = &NetworkPolicy{}
var _ lifecycle.Loader = &NetworkPolicy{}
var _ lifecycle.Ownable = &NetworkPolicy{}
var _ lifecycle.Persister = &NetworkPolicy{}

func (np *NetworkPolicy) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, np.Object, opts...)
}

func (np *NetworkPolicy) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&np.Object.ObjectMeta, from)
}

func (np *NetworkPolicy) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, np.Key, np.Object)
}

func (np *NetworkPolicy) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(np.Object, owner)
}

func (np *NetworkPolicy) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, np.Object, helper.WithObjectKey(np.Key)); err != nil {
		return err
	}

	np.Key = client.ObjectKeyFromObject(np.Object)
	return nil
}

func (np *NetworkPolicy) Copy() *NetworkPolicy {
	return &NetworkPolicy{
		Key:    np.Key,
		Object: np.Object.DeepCopy(),
	}
}

func (np *NetworkPolicy) AllowAll() {
	np.Object.Spec = networkingv1.NetworkPolicySpec{}
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
	return &NetworkPolicy{
		Key:    key,
		Object: &networkingv1.NetworkPolicy{},
	}
}

func NewNetworkPolicyFromObject(obj *networkingv1.NetworkPolicy) *NetworkPolicy {
	return &NetworkPolicy{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewNetworkPolicyPatcher(upd, orig *NetworkPolicy) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
