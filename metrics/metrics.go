package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/puppetlabs/insights-instrumentation/errors"
	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
	"github.com/puppetlabs/insights-instrumentation/metrics/delegates"
	"github.com/puppetlabs/insights-instrumentation/metrics/internal/noop"
	logging "github.com/puppetlabs/insights-logging"
)

type errorBehavior int

const (
	// ErrorBehaviorPanic tells Must* functions to panic if they encounter an error
	ErrorBehaviorPanic errorBehavior = iota
	// ErrorBehaviorLog tells Must* functions to log the error and return a noop type
	ErrorBehaviorLog
)

// Options are configuration options used in creating a new Metrics
type Options struct {
	// DelegateType is the delegate to create for use in Metrics
	DelegateType  delegates.DelegateType
	ErrorBehavior errorBehavior
	Logger        logging.Logger
}

// Metrics provides a wrapper for a collector delegate to report metrics to.
type Metrics struct {
	Namespace     string
	counters      map[string]collectors.Counter
	timers        map[string]collectors.Timer
	delegate      delegates.Delegate
	errorBehavior errorBehavior
	logger        logging.Logger

	sync.Mutex
}

// RegisterTimer registeres a timer at name
func (m *Metrics) RegisterTimer(name string, opts collectors.TimerOptions) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.timers[name]; !ok {
		t, err := m.delegate.NewTimer(name, opts)
		if err != nil {
			return err
		}

		m.timers[name] = t
	}

	return nil
}

// MustRegisterTimer calls RegisterTimer and logs an error if one occurs
func (m *Metrics) MustRegisterTimer(name string, opts collectors.TimerOptions) {
	if err := m.RegisterTimer(name, opts); err != nil {
		m.handleError(err)
	}
}

// Timer returns a Timer metric at name
func (m *Metrics) Timer(name string) (collectors.Timer, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.timers[name]; !ok {
		return nil, fmt.Errorf("could not find %s timer", name)
	}

	return m.timers[name], nil
}

// MustTimer calls Timer and returns a Timer at the name, if an error occurs
// then a NoOpTimer is returned instead and the error is logged
func (m *Metrics) MustTimer(name string, labels ...collectors.Label) collectors.Timer {
	t, err := m.Timer(name)
	if err != nil {
		m.handleError(err)

		return noop.Timer{}
	}

	t, err = t.WithLabels(labels)
	if err != nil {
		m.handleError(err)

		return noop.Timer{}
	}

	return t
}

// OnTimer records the duration of a call to the user defined function fn
func (m *Metrics) OnTimer(t collectors.Timer, fn func()) {
	h := t.Start()
	fn()
	t.ObserveDuration(h)
}

func (m *Metrics) RegisterCounter(name string, opts collectors.CounterOptions) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.counters[name]; !ok {
		c, err := m.delegate.NewCounter(name, opts)
		if err != nil {
			return err
		}

		m.counters[name] = c
	}

	return nil
}

func (m *Metrics) MustRegisterCounter(name string, opts collectors.CounterOptions) {
	if err := m.RegisterCounter(name, opts); err != nil {
		m.handleError(err)
	}
}

// Counter returns a new Counter metric registered as name
func (m *Metrics) Counter(name string) (collectors.Counter, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.counters[name]; !ok {
		return nil, errors.NewMetricsNotFoundError(name, "counter")
	}

	return m.counters[name], nil
}

func (m *Metrics) MustCounter(name string, labels ...collectors.Label) collectors.Counter {
	c, err := m.Counter(name)
	if err != nil {
		m.handleError(err)

		return noop.Counter{}
	}

	c, err = c.WithLabels(labels)
	if err != nil {
		m.handleError(err)

		return noop.Counter{}
	}

	return c
}

// Handler returns the http handler from the delegate if there is one
func (m *Metrics) Handler() http.Handler {
	return m.delegate.NewHandler()
}

func (m *Metrics) handleError(err error) {
	if m.errorBehavior == ErrorBehaviorLog {
		m.logger.Error(err.Error())
	} else {
		panic(err)
	}
}

// NewNamespace returns a new Metrics object at namespace
func NewNamespace(namespace string, opts Options) (*Metrics, error) {
	delegate, err := delegates.New(namespace, opts.DelegateType)
	if err != nil {
		return nil, err
	}

	logger := log(context.Background())
	if opts.Logger != nil {
		logger = opts.Logger
	}

	return &Metrics{
		Namespace:     namespace,
		delegate:      delegate,
		counters:      make(map[string]collectors.Counter),
		timers:        make(map[string]collectors.Timer),
		errorBehavior: opts.ErrorBehavior,
		logger:        logger,
	}, nil
}

func NewLabel(name, value string) collectors.Label {
	return collectors.Label{Name: name, Value: value}
}
