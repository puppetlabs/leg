package sentry

import (
	raven "github.com/getsentry/raven-go"
	"github.com/puppetlabs/insights-instrumentation/alerts/trackers"
	"github.com/puppetlabs/insights-instrumentation/errors"
)

type Options struct {
	Environment string
	Release     string
}

type Sentry struct {
	client *raven.Client
}

func (s Sentry) NewCapturer() trackers.Capturer {
	return &Capturer{
		client: s.client,
	}
}

func NewSentry(dsn string, opts Options) (*Sentry, errors.Error) {
	client, err := raven.New(dsn)
	if err != nil {
		// XXX: FIXME
	}

	if opts.Environment != "" {
		client.SetEnvironment(opts.Environment)
	}

	if opts.Release != "" {
		client.SetRelease(opts.Release)
	}

	s := &Sentry{
		client: client,
	}
	return s, nil
}
