package admissionregistrationv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	MutatingWebhookConfigurationKind = admissionregistrationv1.SchemeGroupVersion.WithKind("MutatingWebhookConfiguration")
)

type MutatingWebhookConfiguration struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *admissionregistrationv1.MutatingWebhookConfiguration
}

func makeMutatingWebhookConfiguration(name string, obj *admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfiguration {
	mwc := &MutatingWebhookConfiguration{Name: name, Object: obj}
	mwc.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&mwc.Name, lifecycle.TypedObject{GVK: MutatingWebhookConfigurationKind, Object: mwc.Object})
	return mwc
}

func (mwc *MutatingWebhookConfiguration) Copy() *MutatingWebhookConfiguration {
	return makeMutatingWebhookConfiguration(mwc.Name, mwc.Object.DeepCopy())
}

func NewMutatingWebhookConfiguration(name string) *MutatingWebhookConfiguration {
	return makeMutatingWebhookConfiguration(name, &admissionregistrationv1.MutatingWebhookConfiguration{})
}

func NewMutatingWebhookConfigurationFromObject(obj *admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfiguration {
	return makeMutatingWebhookConfiguration(obj.GetName(), obj)
}

func NewMutatingWebhookConfigurationPatcher(upd, orig *MutatingWebhookConfiguration) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
