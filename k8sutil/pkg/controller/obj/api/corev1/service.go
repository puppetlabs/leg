package corev1

import (
	"context"
	"fmt"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ServiceKind = corev1.SchemeGroupVersion.WithKind("Service")
)

type Service struct {
	Key    client.ObjectKey
	Object *corev1.Service
}

var _ lifecycle.Deleter = &Service{}
var _ lifecycle.LabelAnnotatableFrom = &Service{}
var _ lifecycle.Loader = &Service{}
var _ lifecycle.Ownable = &Service{}
var _ lifecycle.Persister = &Service{}

func (s *Service) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, s.Object, opts...)
}

func (s *Service) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&s.Object.ObjectMeta, from)
}

func (s *Service) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, s.Key, s.Object)
}

func (s *Service) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(s.Object, owner)
}

func (s *Service) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, s.Object, helper.WithObjectKey(s.Key)); err != nil {
		return err
	}

	s.Key = client.ObjectKeyFromObject(s.Object)
	return nil
}

func (s *Service) Copy() *Service {
	return &Service{
		Key:    s.Key,
		Object: s.Object.DeepCopy(),
	}
}

func (s *Service) DNSName() string {
	return fmt.Sprintf("%s.%s.svc", s.Key.Name, s.Key.Namespace)
}

func NewService(key client.ObjectKey) *Service {
	return &Service{
		Key:    key,
		Object: &corev1.Service{},
	}
}

func NewServiceFromObject(obj *corev1.Service) *Service {
	return &Service{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewServicePatcher(upd, orig *Service) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
