package prometheus

import (
	"net/http"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
)

type Prometheus struct {
	namespace string
}

func (p *Prometheus) NewCounter(name string) (collectors.Counter, error) {
	c := prom.NewCounter(prom.CounterOpts{
		Namespace: p.namespace,
		Name:      name,
	})

	if err := prom.Register(c); err != nil {
		return nil, err
	}

	return &Counter{delegate: c}, nil
}

func (p *Prometheus) NewTimer(name string, opts collectors.TimerOptions) (collectors.Timer, error) {
	observer := prom.NewHistogram(prom.HistogramOpts{
		Namespace: p.namespace,
		Name:      name,
		Buckets:   opts.HistogramBoundaries,
	})

	if err := prom.Register(observer); err != nil {
		return nil, err
	}

	t := NewTimer(observer)

	return t, nil
}

func (p *Prometheus) NewHandler() http.Handler {
	return promhttp.Handler()
}

func New(namespace string) *Prometheus {
	return &Prometheus{
		namespace: namespace,
	}
}
