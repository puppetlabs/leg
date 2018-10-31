package prometheus

import (
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
)

type Timer struct {
	observer prom.Observer
	timers   map[*collectors.TimerHandle]*prom.Timer

	sync.Mutex
}

func (t *Timer) Start() *collectors.TimerHandle {
	t.Lock()
	defer t.Unlock()

	h := &collectors.TimerHandle{}
	promt := prom.NewTimer(t.observer)

	t.timers[h] = promt

	return h
}

func (t *Timer) ObserveDuration(h *collectors.TimerHandle) {
	if promt, ok := t.timers[h]; ok {
		promt.ObserveDuration()
	}
}

func NewTimer(observer prom.Observer) *Timer {
	return &Timer{
		observer: observer,
		timers:   make(map[*collectors.TimerHandle]*prom.Timer),
	}
}
