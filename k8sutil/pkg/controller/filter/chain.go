package filter

import (
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Chainable interface {
	ChainTo(reconcile.Reconciler) reconcile.Reconciler
}

type ChainableFunc func(reconcile.Reconciler) reconcile.Reconciler

var _ Chainable = ChainableFunc(nil)

func (cf ChainableFunc) ChainTo(r reconcile.Reconciler) reconcile.Reconciler {
	return cf(r)
}

func ChainR(last reconcile.Reconciler, links ...Chainable) reconcile.Reconciler {
	for i := len(links) - 1; i >= 0; i-- {
		last = links[i].ChainTo(last)
	}

	return last
}
