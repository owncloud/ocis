package service

import (
	"fmt"
	"strings"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/postprocessing"
)

// PostprocessingService is an instance of the service handling postprocessing of files
type PostprocessingService struct {
	log    log.Logger
	events <-chan interface{}
	pub    events.Publisher
	steps  []events.Postprocessingstep
	c      config.Postprocessing
}

// NewPostprocessingService returns a new instance of a postprocessing service
func NewPostprocessingService(stream events.Stream, logger log.Logger, c config.Postprocessing) (*PostprocessingService, error) {
	evs, err := events.Consume(stream, "postprocessing",
		events.BytesReceived{},
		events.StartPostprocessingStep{},
		events.VirusscanFinished{},
		events.UploadReady{},
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
		switch ev := e.(type) {
		case events.BytesReceived:
			pp := postprocessing.New(ev.UploadID, ev.URL, ev.ExecutingUser, ev.Filename, ev.Filesize, ev.ResourceID, pps.steps, pps.c.Delayprocessing)
			current[ev.UploadID] = pp
			next = pp.Init(ev)
		case events.VirusscanFinished:
			pp := current[ev.UploadID]
			if pp == nil {
				// no current upload - this was an on demand scan
				continue
			}
			next = pp.Virusscan(ev)
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

	if c.Virusscan {
		if !contains(steps, events.PPStepAntivirus) {
			steps = append(steps, events.PPStepAntivirus)
			fmt.Printf("ATTENTION: POSTPROCESSING_VIRUSSCAN is deprecated. Use `POSTPROCESSING_STEPS=%v` in the future\n", join(steps))
		}
	}

	if c.Delayprocessing != 0 {
		if !contains(steps, events.PPStepDelay) {
			if len(steps) > 0 {
				fmt.Printf("Added delay step to the list of postprocessing steps. NOTE: Use envvar `POSTPROCESSING_STEPS=%v` to suppress this message and choose the order of postprocessing steps.\n", join(append(steps, events.PPStepDelay)))
			}

			steps = append(steps, events.PPStepDelay)
		}
	}

	return steps
}

func contains(all []events.Postprocessingstep, candidate events.Postprocessingstep) bool {
	for _, s := range all {
		if s == candidate {
			return true
		}
	}
	return false
}

func join(all []events.Postprocessingstep) string {
	var slice []string
	for _, s := range all {
		slice = append(slice, string(s))
	}
	return strings.Join(slice, ",")
}
