package activities

import "github.com/puppetlabs/horsehead/v2/instrumentation/activities/delegates"
import "github.com/puppetlabs/horsehead/v2/instrumentation/activities/activity"

// Delegate represents a manager for reporting activities to a given
// service.
type Delegate interface {
	// Report sends an activity to the relevant backend.
	Report(activity.Activity) error

	// Close tells the underlying delegate we're done with it.
	Close() error
}

// NewIntercomDelegate instantiates a delegate for reporting
// activities to Intercom.
func NewIntercomDelegate(appID, apiKey string) Delegate {
	return delegates.NewIntercom(appID, apiKey)
}
