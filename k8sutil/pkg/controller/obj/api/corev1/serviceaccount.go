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
	ErrServiceAccountMissingDefaultTokenSecret = errors.New("service account has no default token secret")
)

var (
	ServiceAccountKind = corev1.SchemeGroupVersion.WithKind("ServiceAccount")
)

type ServiceAccount struct {
	Key    client.ObjectKey
	Object *corev1.ServiceAccount
}

var _ lifecycle.Deleter = &ServiceAccount{}
var _ lifecycle.LabelAnnotatableFrom = &ServiceAccount{}
var _ lifecycle.Loader = &ServiceAccount{}
var _ lifecycle.Ownable = &ServiceAccount{}
var _ lifecycle.Persister = &ServiceAccount{}

func (sa *ServiceAccount) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, sa.Object, opts...)
}

func (sa *ServiceAccount) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&sa.Object.ObjectMeta, from)
}

func (sa *ServiceAccount) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, sa.Key, sa.Object)
}

func (sa *ServiceAccount) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(sa.Object, owner)
}

func (sa *ServiceAccount) Persist(ctx context.Context, cl client.Client) error {
	if err := helper.CreateOrUpdate(ctx, cl, sa.Object, helper.WithObjectKey(sa.Key)); err != nil {
		return err
	}

	sa.Key = client.ObjectKeyFromObject(sa.Object)
	return nil
}

func (sa *ServiceAccount) Copy() *ServiceAccount {
	return &ServiceAccount{
		Key:    sa.Key,
		Object: sa.Object.DeepCopy(),
	}
}

func NewServiceAccount(key client.ObjectKey) *ServiceAccount {
	return &ServiceAccount{
		Key:    key,
		Object: &corev1.ServiceAccount{},
	}
}

func NewServiceAccountFromObject(obj *corev1.ServiceAccount) *ServiceAccount {
	return &ServiceAccount{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewServiceAccountPatcher(upd, orig *ServiceAccount) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}

type ServiceAccountTokenSecrets struct {
	ServiceAccount *ServiceAccount

	DefaultTokenSecret     *ServiceAccountTokenSecret
	AdditionalTokenSecrets []*ServiceAccountTokenSecret
}

var _ lifecycle.Loader = &ServiceAccountTokenSecrets{}

func (sats *ServiceAccountTokenSecrets) Load(ctx context.Context, cl client.Client) (bool, error) {
	if len(sats.ServiceAccount.Object.Secrets) == 0 || sats.ServiceAccount.Object.Secrets[0].Name == "" {
		return false, ErrServiceAccountMissingDefaultTokenSecret
	}

	sats.DefaultTokenSecret = NewServiceAccountTokenSecret(client.ObjectKey{
		Namespace: sats.ServiceAccount.Key.Namespace,
		Name:      sats.ServiceAccount.Object.Secrets[0].Name,
	})
	loaders := lifecycle.Loaders{sats.DefaultTokenSecret}

	sats.AdditionalTokenSecrets = nil
	for _, secret := range sats.ServiceAccount.Object.Secrets[1:] {
		add := NewServiceAccountTokenSecret(client.ObjectKey{
			Namespace: sats.ServiceAccount.Key.Namespace,
			Name:      secret.Name,
		})

		sats.AdditionalTokenSecrets = append(sats.AdditionalTokenSecrets, add)
		loaders = append(loaders, sats)
	}

	return loaders.Load(ctx, cl)
}

func NewServiceAccountTokenSecrets(sa *ServiceAccount) *ServiceAccountTokenSecrets {
	return &ServiceAccountTokenSecrets{
		ServiceAccount: sa,
	}
}

func NewServiceAccountTokenSecretsDefaultPresentPoller(sats *ServiceAccountTokenSecrets) lifecycle.RetryLoader {
	return lifecycle.NewRetryLoader(
		lifecycle.LoaderFunc(func(ctx context.Context, cl client.Client) (bool, error) {
			ok, err := sats.ServiceAccount.Load(ctx, cl)
			if err != nil || !ok {
				return ok, err
			}

			return sats.Load(ctx, cl)
		}),
		func(ok bool, err error) (bool, error) {
			return ok || err != ErrServiceAccountMissingDefaultTokenSecret, err
		},
	)
}
