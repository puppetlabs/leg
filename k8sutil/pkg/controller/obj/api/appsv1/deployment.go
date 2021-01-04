package appsv1

import (
	"context"
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrDeploymentNotAvailable = errors.New("deployment not available")
)

var (
	DeploymentKind = appsv1.SchemeGroupVersion.WithKind("Deployment")
)

type Deployment struct {
	Key    client.ObjectKey
	Object *appsv1.Deployment
}

var _ lifecycle.Deleter = &Deployment{}
var _ lifecycle.LabelAnnotatableFrom = &Deployment{}
var _ lifecycle.Loader = &Deployment{}
var _ lifecycle.Ownable = &Deployment{}
var _ lifecycle.Persister = &Deployment{}

func (d *Deployment) Delete(ctx context.Context, cl client.Client) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, d.Object)
}

func (d *Deployment) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&d.Object.ObjectMeta, from)
}

func (d *Deployment) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, d.Key, d.Object)
}

func (d *Deployment) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, d.Object, owner)
}

func (d *Deployment) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, d.Object, helper.WithObjectKey(d.Key)); err != nil {
		return err
	}

	d.Key = client.ObjectKeyFromObject(d.Object)
	return nil
}

func (d *Deployment) Copy() *Deployment {
	return &Deployment{
		Key:    d.Key,
		Object: d.Object.DeepCopy(),
	}
}

func (d *Deployment) Condition(typ appsv1.DeploymentConditionType) (appsv1.DeploymentCondition, bool) {
	for _, cond := range d.Object.Status.Conditions {
		if cond.Type == typ {
			return cond, true
		}
	}
	return appsv1.DeploymentCondition{Type: typ}, false
}

func (d *Deployment) AvailableCondition() (appsv1.DeploymentCondition, bool) {
	return d.Condition(appsv1.DeploymentAvailable)
}

func (d *Deployment) ProgressingCondition() (appsv1.DeploymentCondition, bool) {
	return d.Condition(appsv1.DeploymentProgressing)
}

func NewDeployment(key client.ObjectKey) *Deployment {
	return &Deployment{
		Key:    key,
		Object: &appsv1.Deployment{},
	}
}

func NewDeploymentFromObject(obj *appsv1.Deployment) *Deployment {
	return &Deployment{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewDeploymentPatcher(upd, orig *Deployment) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}

func NewDeploymentAvailablePoller(deployment *Deployment) lifecycle.RetryLoader {
	return lifecycle.NewRetryLoader(deployment, func(ok bool, err error) (bool, error) {
		if !ok || err != nil {
			return ok, err
		}

		cond, ok := deployment.AvailableCondition()
		if !ok || cond.Status != corev1.ConditionTrue {
			return false, ErrDeploymentNotAvailable
		}

		return true, nil
	})
}
