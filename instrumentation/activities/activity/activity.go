package activity

import "time"

type ActivityMetadata map[string]string

type Activity struct {
	// ID is a unique identifier for this activity. Activity data consumers will
	// use this identifier (when provided) to deduplicated activity data.
	ID string

	// Name is the name of the activity to report. e.g. "workflow-run"
	Name string

	// UserID is the identifier of the user that this activity belongs to.
	UserID    string
	Metadata  ActivityMetadata
	OccuredAt time.Time
}
