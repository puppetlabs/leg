package corev1

import (
	"fmt"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ServiceKind = corev1.SchemeGroupVersion.WithKind("Service")
)

type Service struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.Service
}

func makeService(key client.ObjectKey, obj *corev1.Service) *Service {
	s := &Service{Key: key, Object: obj}
	s.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&s.Key, lifecycle.TypedObject{GVK: ServiceKind, Object: s.Object})
	return s
}

func (s *Service) Copy() *Service {
	return makeService(s.Key, s.Object.DeepCopy())
}

func (s *Service) DNSName() string {
	return fmt.Sprintf("%s.%s.svc", s.Key.Name, s.Key.Namespace)
}

func NewService(key client.ObjectKey) *Service {
	return makeService(key, &corev1.Service{})
}

func NewServiceFromObject(obj *corev1.Service) *Service {
	return makeService(client.ObjectKeyFromObject(obj), obj)
}

func NewServicePatcher(upd, orig *Service) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
