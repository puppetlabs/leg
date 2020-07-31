package delegates

import (
	intercom "gopkg.in/intercom/intercom-go.v2"

	"github.com/puppetlabs/horsehead/v2/instrumentation/activities/activity"
)

// Intercom manages reporting activity information to Intercom.
type Intercom struct {
	client *intercom.Client
}

func (d *Intercom) Report(act activity.Activity) error {
	return d.client.Events.Save(&intercom.Event{
		ID:        act.ID,
		UserID:    act.UserID,
		EventName: act.Name,
		CreatedAt: act.OccuredAt.Unix(),
		Metadata:  convertActivityMetadataForIntercom(act.Metadata),
	})
}

func NewIntercom(accessToken string) *Intercom {
	// Intercom has changed around how authentication works with their API.
	// Historically, they've required an app ID and an API key. They have a new,
	// improved system that only requires an access token. To not break the
	// relevant interface, they've used this methodology for using an access
	// token only.
	client := intercom.NewClient(accessToken, "")

	return &Intercom{
		client: client,
	}
}

func convertActivityMetadataForIntercom(am activity.ActivityMetadata) map[string]interface{} {
	res := make(map[string]interface{}, len(am))

	for k, v := range am {
		res[k] = v
	}

	return res
}
