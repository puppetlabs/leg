package endtoend

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/norm"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NamespaceOptions struct {
	NamePrefix string
	Labels     map[string]string
}

type NamespaceOption interface {
	ApplyToNamespaceOptions(target *NamespaceOptions)
}

var _ NamespaceOption = &NamespaceOptions{}

func (o *NamespaceOptions) ApplyToNamespaceOptions(target *NamespaceOptions) {
	*target = *o
}

func (o *NamespaceOptions) ApplyOptions(opts []NamespaceOption) {
	for _, opt := range opts {
		opt.ApplyToNamespaceOptions(o)
	}
}

type NamespaceWithNamePrefix string

var _ NamespaceOption = NamespaceWithNamePrefix("")

func (nwnp NamespaceWithNamePrefix) ApplyToNamespaceOptions(target *NamespaceOptions) {
	target.NamePrefix = string(nwnp)
}

type NamespaceWithLabels map[string]string

var _ NamespaceOption = NamespaceWithLabels(nil)

func (nwl NamespaceWithLabels) ApplyToNamespaceOptions(target *NamespaceOptions) {
	target.Labels = nwl
}

type NamespaceFactory interface {
	New(ctx context.Context, e *Environment) (*corev1.Namespace, error)
}

type NamespaceFactoryFunc func(ctx context.Context, e *Environment) (*corev1.Namespace, error)

var _ NamespaceFactory = NamespaceFactoryFunc(nil)

func (nff NamespaceFactoryFunc) New(ctx context.Context, e *Environment) (*corev1.Namespace, error) {
	return nff(ctx, e)
}

func NewGenerateNameNamespaceFactory(opts ...NamespaceOption) NamespaceFactory {
	o := &NamespaceOptions{
		NamePrefix: "test-",
	}
	o.ApplyOptions(opts)

	return NamespaceFactoryFunc(func(ctx context.Context, e *Environment) (*corev1.Namespace, error) {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: norm.MetaGenerateName(o.NamePrefix),
				Labels:       o.Labels,
			},
		}
		if err := e.ControllerClient.Create(ctx, ns); err != nil {
			return nil, err
		}

		// Wait for default service account to be populated.
		if err := retry.Wait(ctx, func(ctx context.Context) (bool, error) {
			err := e.ControllerClient.Get(ctx, client.ObjectKey{Namespace: ns.GetName(), Name: "default"}, &corev1.ServiceAccount{})
			return !errors.IsNotFound(err), err
		}); err != nil {
			return nil, err
		}

		return ns, nil
	})
}

func NewTestNamespaceFactory(t *testing.T, opts ...NamespaceOption) NamespaceFactory {
	o := &NamespaceOptions{}
	o.ApplyOptions(opts)

	return NewGenerateNameNamespaceFactory(
		o,
		NamespaceWithNamePrefix(o.NamePrefix+t.Name()),
	)
}

func WithNamespace(ctx context.Context, e *Environment, nf NamespaceFactory, fn func(ns *corev1.Namespace)) (err error) {
	ns, err := nf.New(ctx, e)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = e.ControllerClient.Delete(ctx, ns)
	}()

	fn(ns)
	return
}
