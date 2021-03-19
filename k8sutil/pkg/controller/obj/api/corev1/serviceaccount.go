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
	ErrServiceAccountMissingDefaultTokenSecret = errors.New("service account has no default token secret")
)

var (
	ServiceAccountKind = corev1.SchemeGroupVersion.WithKind("ServiceAccount")
)

type ServiceAccount struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *corev1.ServiceAccount
}

func makeServiceAccount(key client.ObjectKey, obj *corev1.ServiceAccount) *ServiceAccount {
	sa := &ServiceAccount{Key: key, Object: obj}
	sa.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&sa.Key, lifecycle.TypedObject{GVK: ServiceAccountKind, Object: sa.Object})
	return sa
}

func (sa *ServiceAccount) Copy() *ServiceAccount {
	return makeServiceAccount(sa.Key, sa.Object.DeepCopy())
}

func NewServiceAccount(key client.ObjectKey) *ServiceAccount {
	return makeServiceAccount(key, &corev1.ServiceAccount{})
}

func NewServiceAccountFromObject(obj *corev1.ServiceAccount) *ServiceAccount {
	return makeServiceAccount(client.ObjectKeyFromObject(obj), obj)
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
