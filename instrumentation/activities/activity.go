package activities

import "time"
import "github.com/puppetlabs/leg/instrumentation/activities/activity"

// NewActivity is a helper function for quickly instantiating a new activity.
func NewActivity(userID, name string) activity.Activity {
	return activity.Activity{
		UserID:    userID,
		Name:      name,
		Metadata:  make(activity.ActivityMetadata),
		OccuredAt: time.Now(),
	}
}
