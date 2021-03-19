package appsv1

import (
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrDeploymentNotAvailable = errors.New("deployment not available")
)

var (
	DeploymentKind = appsv1.SchemeGroupVersion.WithKind("Deployment")
)

type Deployment struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *appsv1.Deployment
}

func makeDeployment(key client.ObjectKey, obj *appsv1.Deployment) *Deployment {
	d := &Deployment{Key: key, Object: obj}
	d.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&d.Key, lifecycle.TypedObject{GVK: DeploymentKind, Object: d.Object})
	return d
}

func (d *Deployment) Copy() *Deployment {
	return makeDeployment(d.Key, d.Object.DeepCopy())
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
	return makeDeployment(key, &appsv1.Deployment{})
}

func NewDeploymentFromObject(obj *appsv1.Deployment) *Deployment {
	return makeDeployment(client.ObjectKeyFromObject(obj), obj)
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
