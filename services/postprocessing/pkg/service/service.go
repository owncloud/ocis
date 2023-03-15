package service

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/postprocessing"
)

// PostprocessingService is an instance of the service handling postprocessing of files
type PostprocessingService struct {
	log    log.Logger
	events <-chan events.Event
	pub    events.Publisher
	steps  []events.Postprocessingstep
	c      config.Postprocessing
}

// NewPostprocessingService returns a new instance of a postprocessing service
func NewPostprocessingService(stream events.Stream, logger log.Logger, c config.Postprocessing) (*PostprocessingService, error) {
	evs, err := events.Consume(stream, "postprocessing",
		events.BytesReceived{},
		events.StartPostprocessingStep{},
		events.UploadReady{},
		events.PostprocessingStepFinished{},
	)
	if err != nil {
		return nil, err
	}

	return &PostprocessingService{
		log:    logger,
		events: evs,
		pub:    stream,
		steps:  getSteps(c),
		c:      c,
	}, nil
}

// Run to fulfil Runner interface
func (pps *PostprocessingService) Run() error {
	current := make(map[string]*postprocessing.Postprocessing)
	for e := range pps.events {
		var next interface{}
		switch ev := e.Event.(type) {
		case events.BytesReceived:
			pp := postprocessing.New(ev.UploadID, ev.URL, ev.ExecutingUser, ev.Filename, ev.Filesize, ev.ResourceID, pps.steps, pps.c.Delayprocessing)
			current[ev.UploadID] = pp
			next = pp.Init(ev)
		case events.PostprocessingStepFinished:
			pp := current[ev.UploadID]
			if pp == nil {
				// no current upload - this was an on demand scan
				continue
			}
			next = pp.NextStep(ev)
		case events.StartPostprocessingStep:
			if ev.StepToStart != events.PPStepDelay {
				continue
			}
			pp := current[ev.UploadID]
			next = pp.Delay(ev)
		case events.UploadReady:
			// the storage provider thinks the upload is done - so no need to keep it any more
			delete(current, ev.UploadID)
		}

		if next != nil {
			if err := events.Publish(pps.pub, next); err != nil {
				pps.log.Error().Err(err).Msg("unable to publish event")
				return err // we can't publish -> we are screwed
			}
		}

	}
	return nil
}

func getSteps(c config.Postprocessing) []events.Postprocessingstep {
	// NOTE: improved version only allows configuring order of postprocessing steps
	// But we aim for a system where postprocessing steps can be configured per space, ideally by the spaceadmin itself
	// We need to iterate over configuring PP service when we see fit
	var steps []events.Postprocessingstep
	for _, s := range c.Steps {
		steps = append(steps, events.Postprocessingstep(s))
	}

	return steps
}
