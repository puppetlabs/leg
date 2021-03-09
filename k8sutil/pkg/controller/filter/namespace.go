package filter

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type SingleNamespaceReconciler struct {
	delegate reconcile.Reconciler
	name     string
}

var _ reconcile.Reconciler = &SingleNamespaceReconciler{}
var _ inject.Injector = &SingleNamespaceReconciler{}

func (snr SingleNamespaceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (result reconcile.Result, err error) {
	// You can't be clever and use the built-in namespace restrictions or
	// predicates in controller-runtime to filter out the namespace before it
	// gets here. The caching applies to the same namespace filter, so the
	// namespaces used/created by this controller will appear to not exist!
	if snr.name != "" && req.Namespace != snr.name {
		return reconcile.Result{}, nil
	}

	return snr.delegate.Reconcile(ctx, req)
}

func (snr SingleNamespaceReconciler) InjectFunc(f inject.Func) error {
	return f(snr.delegate)
}

func NewSingleNamespaceReconciler(delegate reconcile.Reconciler, namespace string) *SingleNamespaceReconciler {
	return &SingleNamespaceReconciler{
		delegate: delegate,
		name:     namespace,
	}
}

func ChainSingleNamespaceReconciler(namespace string) Chainable {
	return ChainableFunc(func(delegate reconcile.Reconciler) reconcile.Reconciler {
		return NewSingleNamespaceReconciler(delegate, namespace)
	})
}
