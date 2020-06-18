package activity

import "time"

type ActivityMetadata map[string]string

type Activity struct {
	UserID    string
	Name      string
	Metadata  ActivityMetadata
	OccuredAt time.Time
}
