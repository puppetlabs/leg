package noop

import "github.com/puppetlabs/horsehead/instrumentation/alerts/trackers"

type NoOp struct{}

func (NoOp) NewCapturer() trackers.Capturer {
	return &Capturer{}
}
