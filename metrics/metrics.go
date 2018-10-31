package metrics

import (
	"net/http"
	"sync"

	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
	"github.com/puppetlabs/insights-instrumentation/metrics/delegates"
)

// Options are configuration options used in creating a new Metrics
type Options struct {
	// DelegateName is the delegate to create for use in Metrics
	DelegateName delegates.DelegateType
}

// Metrics provides a wrapper for a collector delegate to report metrics to.
type Metrics struct {
	Namespace string
	counters  map[string]collectors.Counter
	timers    map[string]collectors.Timer
	delegate  delegates.Delegate

	sync.Mutex
}

// Timer returns a new Timer metric registered as name
func (m *Metrics) Timer(name string, opts collectors.TimerOptions) (collectors.Timer, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.timers[name]; !ok {
		t, err := m.delegate.NewTimer(name, opts)
		if err != nil {
			return nil, err
		}

		m.timers[name] = t
	}

	return m.timers[name], nil
}

// OnTimer records the duration of a call to the user defined function fn
func (m *Metrics) OnTimer(name string, opts collectors.TimerOptions, fn func()) error {
	t, err := m.Timer(name, opts)
	if err != nil {
		return err
	}

	h := t.Start()
	fn()
	t.ObserveDuration(h)

	return nil
}

// Counter returns a new Counter metric registered as name
func (m *Metrics) Counter(name string) (collectors.Counter, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.counters[name]; !ok {
		c, err := m.delegate.NewCounter(name)
		if err != nil {
			return nil, err
		}

		m.counters[name] = c
	}

	return m.counters[name], nil
}

// OnCounter gets or created a Counter and passes it into fn for use by a user
// defined function
func (m *Metrics) OnCounter(name string, fn func(c collectors.Counter)) error {
	c, err := m.Counter(name)
	if err != nil {
		return err
	}

	fn(c)

	return nil
}

// Handler returns the http handler from the delegate if there is one
func (m *Metrics) Handler() http.Handler {
	return m.delegate.NewHandler()
}

// NewNamespace returns a new Metrics object at namespace
func NewNamespace(namespace string, kind delegates.DelegateType) (*Metrics, error) {
	delegate, err := delegates.New(namespace, kind)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Namespace: namespace,
		delegate:  delegate,
		counters:  make(map[string]collectors.Counter),
		timers:    make(map[string]collectors.Timer),
	}, nil
}
