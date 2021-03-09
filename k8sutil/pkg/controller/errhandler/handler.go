package errhandler

import (
	"context"

	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ErrorHandler interface {
	OnError(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error)
}

type ErrorHandlerFunc func(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error)

var _ ErrorHandler = ErrorHandlerFunc(nil)

func (ehf ErrorHandlerFunc) OnError(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error) {
	return ehf(ctx, req, err)
}

type PanicHandler interface {
	OnPanic(ctx context.Context, req reconcile.Request, rv interface{}) (reconcile.Result, error)
}

type PanicHandlerFunc func(ctx context.Context, req reconcile.Request, rv interface{}) (reconcile.Result, error)

var _ PanicHandler = PanicHandlerFunc(nil)

func (phf PanicHandlerFunc) OnPanic(ctx context.Context, req reconcile.Request, rv interface{}) (reconcile.Result, error) {
	return phf(ctx, req, rv)
}

var (
	LoggingErrorHandler ErrorHandler = ErrorHandlerFunc(func(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error) {
		klog.ErrorDepth(2, err)
		return reconcile.Result{}, err
	})

	PropagatingErrorHandler ErrorHandler = ErrorHandlerFunc(func(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error) {
		return reconcile.Result{}, err
	})

	MaskingErrorHandler ErrorHandler = ErrorHandlerFunc(func(ctx context.Context, req reconcile.Request, err error) (reconcile.Result, error) {
		return reconcile.Result{}, nil
	})
)

type ErrorMatcher struct {
	Rule    errmark.Rule
	Handler ErrorHandler
}

type ErrorMatchersBuilder struct {
	first    []ErrorMatcher
	last     []ErrorMatcher
	fallback ErrorHandler
}

func (emb *ErrorMatchersBuilder) Prepend(rule errmark.Rule, hnd ErrorHandler) *ErrorMatchersBuilder {
	emb.first = append(emb.first, ErrorMatcher{
		Rule:    rule,
		Handler: hnd,
	})
	return emb
}

func (emb *ErrorMatchersBuilder) Append(rule errmark.Rule, hnd ErrorHandler) *ErrorMatchersBuilder {
	emb.last = append(emb.last, ErrorMatcher{
		Rule:    rule,
		Handler: hnd,
	})
	return emb
}

func (emb *ErrorMatchersBuilder) SetFallback(hnd ErrorHandler) *ErrorMatchersBuilder {
	emb.fallback = hnd
	return emb
}

func (emb *ErrorMatchersBuilder) Build() (matchers []ErrorMatcher) {
	// Prepended matchers.
	for i := len(emb.first) - 1; i >= 0; i-- {
		matchers = append(matchers, emb.first[i])
	}

	// Appended matchers.
	matchers = append(matchers, emb.last...)

	// Fallback matcher.
	if emb.fallback != nil {
		matchers = append(matchers, ErrorMatcher{
			Rule:    errmark.RuleAlways,
			Handler: emb.fallback,
		})
	}

	return
}

func NewErrorMatchersBuilder() *ErrorMatchersBuilder {
	return &ErrorMatchersBuilder{}
}

func NewDefaultErrorMatchersBuilder() *ErrorMatchersBuilder {
	return NewErrorMatchersBuilder().
		Append(
			errmark.RuleMarkedUser,
			MaskingErrorHandler,
		).
		Append(
			errmark.RuleAny(
				RuleIsConflict,
				RuleIsTimeout,
				errmark.RuleMarkedTransient,
			),
			PropagatingErrorHandler,
		).
		SetFallback(LoggingErrorHandler)
}

var DefaultErrorMatchers = NewDefaultErrorMatchersBuilder().Build()
