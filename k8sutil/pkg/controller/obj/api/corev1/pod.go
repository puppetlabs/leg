package corev1

import (
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrPodTerminated = errors.New("pod terminated")
	ErrPodRunning    = errors.New("pod running")
	ErrPodWaiting    = errors.New("pod waiting to start")
)

var (
	PodKind = corev1.SchemeGroupVersion.WithKind("Pod")
)

type Pod struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.Pod
}

func makePod(key client.ObjectKey, obj *corev1.Pod) *Pod {
	p := &Pod{Key: key, Object: obj}
	p.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&p.Key, lifecycle.TypedObject{GVK: PodKind, Object: p.Object})
	return p
}

func (p *Pod) Copy() *Pod {
	return makePod(p.Key, p.Object.DeepCopy())
}

func (p *Pod) Phase() corev1.PodPhase {
	return p.Object.Status.Phase
}

func (p *Pod) Terminated() bool {
	return p.Phase() == corev1.PodFailed || p.Phase() == corev1.PodSucceeded
}

func (p *Pod) Running() bool {
	return p.Phase() == corev1.PodRunning
}

func NewPod(key client.ObjectKey) *Pod {
	return makePod(key, &corev1.Pod{})
}

func NewPodFromObject(obj *corev1.Pod) *Pod {
	return makePod(client.ObjectKeyFromObject(obj), obj)
}

func NewPodPatcher(upd, orig *Pod) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}

func NewPodRunningPoller(pod *Pod) lifecycle.RetryLoader {
	return lifecycle.NewRetryLoader(pod, func(ok bool, err error) (bool, error) {
		if !ok || err != nil {
			return ok, err
		}

		switch {
		case pod.Running():
			return true, nil
		case pod.Terminated():
			return true, ErrPodTerminated
		default:
			return false, ErrPodWaiting
		}
	})
}

func NewPodTerminatedPoller(pod *Pod) lifecycle.RetryLoader {
	return lifecycle.NewRetryLoader(pod, func(ok bool, err error) (bool, error) {
		if !ok || err != nil {
			return ok, err
		}

		switch {
		case pod.Terminated():
			return true, nil
		case pod.Running():
			return false, ErrPodRunning
		default:
			return false, ErrPodWaiting
		}
	})
}
