package postprocessing

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
)

// Postprocessing handles postprocessing of a file
type Postprocessing struct {
	id  string
	url string
	m   map[events.Postprocessingstep]interface{}
	c   config.Postprocessing
}

// New returns a new postprocessing instance
func New(uploadID string, uploadURL string, c config.Postprocessing) *Postprocessing {
	return &Postprocessing{
		id:  uploadID,
		url: uploadURL,
		m:   make(map[events.Postprocessingstep]interface{}),
		c:   c,
	}
}

// Init is the first step of the postprocessing
func (pp *Postprocessing) Init(ev events.BytesReceived) interface{} {
	pp.m["init"] = ev
	return events.StartPostprocessingStep{
		UploadID:    pp.id,
		URL:         pp.url,
		StepToStart: events.PPStepAntivirus, // TODO: make order configurable
	}

}

// Virusscan is the virusscanning step of the postprocessing
func (pp *Postprocessing) Virusscan(ev events.VirusscanFinished) interface{} {
	pp.m["virusscan"] = ev

	var action string
	switch {
	case ev.Infected:
		action = "delete"
	case ev.Error != nil:
		action = "abort"
	default:
		action = "continue"
	}

	return events.PostprocessingFinished{
		UploadID: pp.id,
		Result:   pp.m,
		Action:   action,
	}

}
