package prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	delegate prom.Counter
}

func (c *Counter) Add(n float64) {
	c.delegate.Add(n)
}

func (c *Counter) Inc() {
	c.delegate.Inc()
}
