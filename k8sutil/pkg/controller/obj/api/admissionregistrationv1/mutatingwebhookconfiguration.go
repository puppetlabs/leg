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
	MutatingWebhookConfigurationKind = admissionregistrationv1.SchemeGroupVersion.WithKind("MutatingWebhookConfiguration")
)

type MutatingWebhookConfiguration struct {
	Name   string
	Object *admissionregistrationv1.MutatingWebhookConfiguration
}

var _ lifecycle.Deleter = &MutatingWebhookConfiguration{}
var _ lifecycle.Persister = &MutatingWebhookConfiguration{}
var _ lifecycle.Loader = &MutatingWebhookConfiguration{}
var _ lifecycle.Ownable = &MutatingWebhookConfiguration{}
var _ lifecycle.LabelAnnotatableFrom = &MutatingWebhookConfiguration{}

func (mwc *MutatingWebhookConfiguration) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, mwc.Object, opts...)
}

func (mwc *MutatingWebhookConfiguration) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&mwc.Object.ObjectMeta, from)
}

func (mwc *MutatingWebhookConfiguration) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, client.ObjectKey{Name: mwc.Name}, mwc.Object)
}

func (mwc *MutatingWebhookConfiguration) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(mwc.Object, owner)
}

func (mwc *MutatingWebhookConfiguration) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, mwc.Object, helper.WithObjectKey(client.ObjectKey{Name: mwc.Name})); err != nil {
		return err
	}

	mwc.Name = mwc.Object.GetName()
	return nil
}

func (mwc *MutatingWebhookConfiguration) Copy() *MutatingWebhookConfiguration {
	return &MutatingWebhookConfiguration{
		Name:   mwc.Name,
		Object: mwc.Object.DeepCopy(),
	}
}

func NewMutatingWebhookConfiguration(name string) *MutatingWebhookConfiguration {
	return &MutatingWebhookConfiguration{
		Name:   name,
		Object: &admissionregistrationv1.MutatingWebhookConfiguration{},
	}
}

func NewMutatingWebhookConfigurationFromObject(obj *admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfiguration {
	return &MutatingWebhookConfiguration{
		Name:   obj.GetName(),
		Object: obj,
	}
}

func NewMutatingWebhookConfigurationPatcher(upd, orig *MutatingWebhookConfiguration) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(client.ObjectKey{Name: upd.Name}))
}
