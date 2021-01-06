package corev1

import (
	"context"
	"errors"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrPodTerminated = errors.New("pod terminated")
	ErrPodWaiting    = errors.New("pod waiting to start")
)

var (
	PodKind = corev1.SchemeGroupVersion.WithKind("Pod")
)

type Pod struct {
	Key    client.ObjectKey
	Object *corev1.Pod
}

var _ lifecycle.Deleter = &Pod{}
var _ lifecycle.LabelAnnotatableFrom = &Pod{}
var _ lifecycle.Loader = &Pod{}
var _ lifecycle.Ownable = &Pod{}
var _ lifecycle.Persister = &Pod{}

func (p *Pod) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, p.Object, opts...)
}

func (p *Pod) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&p.Object.ObjectMeta, from)
}

func (p *Pod) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, p.Key, p.Object)
}

func (p *Pod) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, p.Object, owner)
}

func (p *Pod) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, p.Object, helper.WithObjectKey(p.Key)); err != nil {
		return err
	}

	p.Key = client.ObjectKeyFromObject(p.Object)
	return nil
}

func (p *Pod) Copy() *Pod {
	return &Pod{
		Key:    p.Key,
		Object: p.Object.DeepCopy(),
	}
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
	return &Pod{
		Key:    key,
		Object: &corev1.Pod{},
	}
}

func NewPodFromObject(obj *corev1.Pod) *Pod {
	return &Pod{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
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
