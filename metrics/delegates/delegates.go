package delegates

import (
	"errors"
	"net/http"

	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
	"github.com/puppetlabs/insights-instrumentation/metrics/internal/prometheus"
)

// Delegate is an interface metrics collectors implement (i.e. prometheus)
type Delegate interface {
	NewCounter(name string) (collectors.Counter, error)
	NewTimer(name string, opts collectors.TimerOptions) (collectors.Timer, error)
	NewHandler() http.Handler
}

type DelegateType int

const (
	// PrometheusDelegate is a const that represents the prometheus backend
	PrometheusDelegate DelegateType = iota
)

func New(namespace string, t DelegateType) (Delegate, error) {
	switch t {
	case PrometheusDelegate:
		return prometheus.New(namespace), nil
	}

	return nil, errors.New("no delegate found")
}
