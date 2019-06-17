package alerts

import "github.com/puppetlabs/insights-instrumentation/alerts/trackers"

type Delegate interface {
	NewCapturer() trackers.Capturer
}
