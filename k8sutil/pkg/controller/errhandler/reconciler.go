package errhandler

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/filter"
	"k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type Reconciler struct {
	delegate      reconcile.Reconciler
	errorMatchers []ErrorMatcher
	panicHandler  PanicHandler
}

var _ reconcile.Reconciler = &Reconciler{}
var _ inject.Injector = &Reconciler{}

func (r Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (result reconcile.Result, err error) {
	defer runtime.HandleCrash()

	var finished bool
	defer func() {
		switch {
		case !finished:
			if r.panicHandler != nil {
				result, err = r.panicHandler.OnPanic(ctx, req, recover())
			}
		case err != nil:
			for _, matcher := range r.errorMatchers {
				if matcher.Rule.Matches(err) {
					result, err = matcher.Handler.OnError(ctx, req, err)
					return
				}
			}
		}
	}()

	result, err = r.delegate.Reconcile(ctx, req)
	finished = true
	return
}

func (r *Reconciler) InjectFunc(f inject.Func) error {
	return f(r.delegate)
}

type ReconcilerOptions struct {
	ErrorMatchers []ErrorMatcher
	PanicHandler  PanicHandler
}

type ReconcilerOption interface {
	ApplyToReconcilerOptions(target *ReconcilerOptions)
}

func (o *ReconcilerOptions) ApplyOptions(opts []ReconcilerOption) {
	for _, opt := range opts {
		opt.ApplyToReconcilerOptions(o)
	}
}

type ReconcilerOptionFunc func(target *ReconcilerOptions)

var _ ReconcilerOption = ReconcilerOptionFunc(nil)

func (rof ReconcilerOptionFunc) ApplyToReconcilerOptions(target *ReconcilerOptions) {
	rof(target)
}

type WithErrorMatchers []ErrorMatcher

var _ ReconcilerOption = WithErrorMatchers(nil)

func (wem WithErrorMatchers) ApplyToReconcilerOptions(target *ReconcilerOptions) {
	target.ErrorMatchers = wem
}

func WithPanicHandler(ph PanicHandler) ReconcilerOption {
	return ReconcilerOptionFunc(func(target *ReconcilerOptions) {
		target.PanicHandler = ph
	})
}

func NewReconciler(delegate reconcile.Reconciler, opts ...ReconcilerOption) *Reconciler {
	o := &ReconcilerOptions{
		ErrorMatchers: DefaultErrorMatchers,
	}
	o.ApplyOptions(opts)

	return &Reconciler{
		delegate:      delegate,
		errorMatchers: o.ErrorMatchers,
		panicHandler:  o.PanicHandler,
	}
}

func ChainReconciler(opts ...ReconcilerOption) filter.Chainable {
	return filter.ChainableFunc(func(delegate reconcile.Reconciler) reconcile.Reconciler {
		return NewReconciler(delegate, opts...)
	})
}
