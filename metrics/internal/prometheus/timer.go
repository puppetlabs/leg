package prometheus

import (
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/puppetlabs/insights-instrumentation/errors"
	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
)

type Timer struct {
	vector   prom.ObserverVec
	delegate prom.Observer
	timers   map[*collectors.TimerHandle]*prom.Timer

	sync.Mutex
}

func (t *Timer) WithLabels(labels []collectors.Label) (collectors.Timer, error) {
	delegate, err := t.vector.GetMetricWith(convertLabels(labels))
	if err != nil {
		return nil, errors.NewMetricsUnknownError("prometheus").WithCause(err)
	}

	return &Timer{
		vector:   t.vector,
		delegate: delegate,
		timers:   make(map[*collectors.TimerHandle]*prom.Timer),
	}, nil
}

func (t *Timer) Start() *collectors.TimerHandle {
	t.Lock()
	defer t.Unlock()

	h := &collectors.TimerHandle{}
	promt := prom.NewTimer(prom.ObserverFunc(func(v float64) {
		t.delegate.Observe(v)
	}))

	t.timers[h] = promt

	return h
}

func (t *Timer) ObserveDuration(h *collectors.TimerHandle) {
	if promt, ok := t.timers[h]; ok {
		promt.ObserveDuration()
	}
}

func NewTimer(vector prom.ObserverVec) *Timer {
	return &Timer{
		vector: vector,
		timers: make(map[*collectors.TimerHandle]*prom.Timer),
	}
}
