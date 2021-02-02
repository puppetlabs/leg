package admissionregistrationv1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ValidatingWebhookConfigurationKind = admissionregistrationv1.SchemeGroupVersion.WithKind("ValidatingWebhookConfiguration")
)

type ValidatingWebhookConfiguration struct {
	Name   string
	Object *admissionregistrationv1.ValidatingWebhookConfiguration
}

var _ lifecycle.Deleter = &ValidatingWebhookConfiguration{}
var _ lifecycle.Persister = &ValidatingWebhookConfiguration{}
var _ lifecycle.Loader = &ValidatingWebhookConfiguration{}
var _ lifecycle.Ownable = &ValidatingWebhookConfiguration{}
var _ lifecycle.LabelAnnotatableFrom = &ValidatingWebhookConfiguration{}

func (vwc *ValidatingWebhookConfiguration) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, vwc.Object, opts...)
}

func (vwc *ValidatingWebhookConfiguration) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&vwc.Object.ObjectMeta, from)
}

func (vwc *ValidatingWebhookConfiguration) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, client.ObjectKey{Name: vwc.Name}, vwc.Object)
}

func (vwc *ValidatingWebhookConfiguration) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(vwc.Object, owner)
}

func (vwc *ValidatingWebhookConfiguration) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, vwc.Object, helper.WithObjectKey(client.ObjectKey{Name: vwc.Name})); err != nil {
		return err
	}

	vwc.Name = vwc.Object.GetName()
	return nil
}

func (vwc *ValidatingWebhookConfiguration) Copy() *ValidatingWebhookConfiguration {
	return &ValidatingWebhookConfiguration{
		Name:   vwc.Name,
		Object: vwc.Object.DeepCopy(),
	}
}

func NewValidatingWebhookConfiguration(name string) *ValidatingWebhookConfiguration {
	return &ValidatingWebhookConfiguration{
		Name:   name,
		Object: &admissionregistrationv1.ValidatingWebhookConfiguration{},
	}
}

func NewValidatingWebhookConfigurationFromObject(obj *admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingWebhookConfiguration {
	return &ValidatingWebhookConfiguration{
		Name:   obj.GetName(),
		Object: obj,
	}
}

func NewValidatingWebhookConfigurationPatcher(upd, orig *ValidatingWebhookConfiguration) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
