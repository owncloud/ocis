package service

import (
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
)

// GetActivitiesResponse is the response on GET activities requests
type GetActivitiesResponse struct {
	Activities []Activity `json:"value"`
}

// Activity represents an activity as it is returned to the client
type Activity struct {
	ID string `json:"id"`

	// TODO: Implement these
	Action    interface{} `json:"action"`
	DriveItem Resource    `json:"driveItem"`
	Actor     Actor       `json:"actor"`
	Times     Times       `json:"times"`

	Template Template `json:"template"`
}

// Resource represents an item such as a file or folder
type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Actor represents the user who performed the Action
type Actor struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// Times represents the timestamps of the Activity
type Times struct {
	RecordedTime time.Time `json:"recordedTime"`
}

// Template contains activity details
type Template struct {
	Message   string                 `json:"message"`
	Variables map[string]interface{} `json:"variables"`
}

// UploadReady converts a UploadReady events to an Activity
func UploadReady(eid string, e events.UploadReady) Activity {
	rid, _ := storagespace.FormatReference(e.FileRef)
	res := Resource{
		ID:   rid,
		Name: e.Filename,
	}
	return Activity{
		ID: eid,
		Template: Template{
			Message: "file created",
			Variables: map[string]interface{}{
				"resource": res,
			},
		},
	}
}
