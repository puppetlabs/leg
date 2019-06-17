package alerts

import "github.com/puppetlabs/horsehead/instrumentation/alerts/trackers"

type Delegate interface {
	NewCapturer() trackers.Capturer
}
