package noop

import "github.com/puppetlabs/insights-instrumentation/alerts/trackers"

type NoOp struct{}

func (NoOp) NewCapturer() trackers.Capturer {
	return &Capturer{}
}
