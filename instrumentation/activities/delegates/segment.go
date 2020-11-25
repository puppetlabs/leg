package delegates

import (
	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/puppetlabs/leg/instrumentation/activities/activity"
)

// Segment reports activity data to Segment, which in turn reports it to other
// activity providers.
type Segment struct {
	client analytics.Client
}

func (d *Segment) Report(act activity.Activity) error {
	track := analytics.Track{
		MessageId:  act.ID,
		UserId:     act.UserID,
		Event:      act.Name,
		Timestamp:  act.OccuredAt,
		Properties: convertActivityMetadataForSegment(act.Metadata),
	}

	return d.client.Enqueue(track)
}

func convertActivityMetadataForSegment(am activity.ActivityMetadata) analytics.Properties {
	p := analytics.NewProperties()

	for k, v := range am {
		p.Set(k, v)
	}

	return p
}

func NewSegment(writeKey string) *Segment {
	client := analytics.New(writeKey)
	return &Segment{client}
}
