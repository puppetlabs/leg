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
		UserID:    act.UserID,
		EventName: act.Name,
		CreatedAt: act.OccuredAt.Unix(),
		Metadata:  convertActivityMetadataForIntercom(act.Metadata),
	})
}

func (d *Intercom) Close() error {
	return nil
}

func NewIntercom(appID, apiKey string) *Intercom {
	return &Intercom{
		client: intercom.NewClient(appID, apiKey),
	}
}

func convertActivityMetadataForIntercom(am activity.ActivityMetadata) map[string]interface{} {
	res := make(map[string]interface{}, len(am))

	for k, v := range am {
		res[k] = v
	}

	return res
}
