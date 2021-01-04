package corev1

import (
	"context"
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrEndpointsNotBound = errors.New("endpoints not bound")
)

var (
	EndpointsKind = corev1.SchemeGroupVersion.WithKind("Endpoints")
)

type Endpoints struct {
	Service *Service
	Object  *corev1.Endpoints
}

var _ lifecycle.Loader = &Endpoints{}

func (e *Endpoints) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, e.Service.Key, e.Object)
}

func (e *Endpoints) Bound() bool {
	for _, subset := range e.Object.Subsets {
		if len(subset.Addresses) > 0 {
			return true
		}
	}

	return false
}

func NewEndpoints(svc *Service) *Endpoints {
	return &Endpoints{
		Service: svc,
		Object:  &corev1.Endpoints{},
	}
}

func NewEndpointsBoundPoller(eps *Endpoints) lifecycle.RetryLoader {
	return lifecycle.NewRetryLoader(eps, func(ok bool, err error) (bool, error) {
		if !ok || err != nil {
			return ok, err
		}

		if !eps.Bound() {
			return false, ErrEndpointsNotBound
		}

		return true, nil
	})
}
