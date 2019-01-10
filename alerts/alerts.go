package alerts

import (
	"github.com/puppetlabs/insights-instrumentation/alerts/internal/noop"
	"github.com/puppetlabs/insights-instrumentation/alerts/internal/sentry"
	"github.com/puppetlabs/insights-instrumentation/alerts/trackers"
	"github.com/puppetlabs/insights-instrumentation/errors"
)

type Options struct {
	Environment string
	Version     string
}

type DelegateFunc func(opts Options) (Delegate, errors.Error)

func NoDelegate(opts Options) (Delegate, errors.Error) {
	return &noop.NoOp{}, nil
}

func DelegateToSentry(dsn string) DelegateFunc {
	return func(opts Options) (Delegate, errors.Error) {
		return sentry.NewSentry(dsn, sentry.Options{
			Environment: opts.Environment,
			Release:     opts.Version,
		})
	}
}

type Alerts struct {
	delegate Delegate
}

func (a *Alerts) NewCapturer() trackers.Capturer {
	return a.delegate.NewCapturer()
}

func NewAlerts(fn DelegateFunc, opts Options) (*Alerts, errors.Error) {
	delegate, err := fn(opts)
	if err != nil {
		return nil, err
	}

	a := &Alerts{
		delegate: delegate,
	}
	return a, nil
}
