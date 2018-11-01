package noop

import "github.com/puppetlabs/insights-instrumentation/metrics/collectors"

type Timer struct{}

func (n Timer) WithLabels([]collectors.Label) (collectors.Timer, error) { return n, nil }
func (n Timer) Start() *collectors.TimerHandle                          { return &collectors.TimerHandle{} }
func (n Timer) ObserveDuration(*collectors.TimerHandle)                 {}
