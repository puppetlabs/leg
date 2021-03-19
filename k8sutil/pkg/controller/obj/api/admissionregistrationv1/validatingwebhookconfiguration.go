package admissionregistrationv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ValidatingWebhookConfigurationKind = admissionregistrationv1.SchemeGroupVersion.WithKind("ValidatingWebhookConfiguration")
)

type ValidatingWebhookConfiguration struct {
	*helper.ClusterScopedAPIObject

	Name   string
	Object *admissionregistrationv1.ValidatingWebhookConfiguration
}

func makeValidatingWebhookConfiguration(name string, obj *admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingWebhookConfiguration {
	vwc := &ValidatingWebhookConfiguration{Name: name, Object: obj}
	vwc.ClusterScopedAPIObject = helper.ForClusterScopedAPIObject(&vwc.Name, lifecycle.TypedObject{GVK: ValidatingWebhookConfigurationKind, Object: vwc.Object})
	return vwc
}

func (vwc *ValidatingWebhookConfiguration) Copy() *ValidatingWebhookConfiguration {
	return makeValidatingWebhookConfiguration(vwc.Name, vwc.Object.DeepCopy())
}

func NewValidatingWebhookConfiguration(name string) *ValidatingWebhookConfiguration {
	return makeValidatingWebhookConfiguration(name, &admissionregistrationv1.ValidatingWebhookConfiguration{})
}

func NewValidatingWebhookConfigurationFromObject(obj *admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingWebhookConfiguration {
	return makeValidatingWebhookConfiguration(obj.GetName(), obj)
}

func NewValidatingWebhookConfigurationPatcher(upd, orig *ValidatingWebhookConfiguration) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
