package ownerext

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type EnqueueRequestForAnnotatedDependent struct {
	Manager   *Manager
	OwnerType runtime.Object
	gvk       schema.GroupVersionKind
}

var _ handler.EventHandler = &EnqueueRequestForAnnotatedDependent{}
var _ inject.Scheme = &EnqueueRequestForAnnotatedDependent{}

func (e *EnqueueRequestForAnnotatedDependent) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

func (e *EnqueueRequestForAnnotatedDependent) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.ObjectOld, q)
	e.add(evt.ObjectNew, q)
}

func (e *EnqueueRequestForAnnotatedDependent) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

func (e *EnqueueRequestForAnnotatedDependent) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

func (e *EnqueueRequestForAnnotatedDependent) add(target client.Object, q workqueue.RateLimitingInterface) {
	dep, found, err := e.Manager.GetDependencyOf(target)
	if err != nil {
		klog.V(4).Infof("enqueue: unable to unmarshal JSON from dependency: %+v", err)
		return
	} else if !found {
		klog.V(4).Infof("enqueue: no annotation on %s/%s", target.GetNamespace(), target.GetName())
		return
	}

	depGroupVersion, _ := schema.ParseGroupVersion(dep.APIVersion)

	if e.gvk.Kind != dep.Kind || e.gvk.GroupVersion() != depGroupVersion {
		klog.V(4).Infof("enqueue: dependency points at GVK %s, Kind=%s, but we want %s", depGroupVersion, dep.Kind, e.gvk)
		return
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: dep.Namespace,
			Name:      dep.Name,
		},
	}
	q.Add(req)
	klog.V(4).Infof("enqueue: successful enqueue of %s %s", e.gvk, req.NamespacedName)
}

func (e *EnqueueRequestForAnnotatedDependent) InjectScheme(s *runtime.Scheme) error {
	kinds, _, err := s.ObjectKinds(e.OwnerType)
	if err != nil {
		return err
	}

	e.gvk = kinds[0]
	return nil
}
