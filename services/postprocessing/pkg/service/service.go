package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/postprocessing"
	"go-micro.dev/v4/store"
)

// PostprocessingService is an instance of the service handling postprocessing of files
type PostprocessingService struct {
	log    log.Logger
	events <-chan events.Event
	pub    events.Publisher
	steps  []events.Postprocessingstep
	store  store.Store
	c      config.Postprocessing
}

// NewPostprocessingService returns a new instance of a postprocessing service
func NewPostprocessingService(stream events.Stream, logger log.Logger, sto store.Store, c config.Postprocessing) (*PostprocessingService, error) {
	evs, err := events.Consume(stream, "postprocessing",
		events.BytesReceived{},
		events.StartPostprocessingStep{},
		events.UploadReady{},
		events.PostprocessingStepFinished{},
		events.ResumePostprocessing{},
	)
	if err != nil {
		return nil, err
	}

	return &PostprocessingService{
		log:    logger,
		events: evs,
		pub:    stream,
		steps:  getSteps(c),
		store:  sto,
		c:      c,
	}, nil
}

// Run to fulfil Runner interface
func (pps *PostprocessingService) Run() error {
	ctx := context.Background()
	for e := range pps.events {
		var (
			next interface{}
			pp   *postprocessing.Postprocessing
			err  error
		)

		ctx = e.GetTraceContext(ctx)

		switch ev := e.Event.(type) {
		case events.BytesReceived:
			pp = postprocessing.New(ev.UploadID, ev.URL, ev.ExecutingUser, ev.Filename, ev.Filesize, ev.ResourceID, pps.steps, pps.c.Delayprocessing)
			next = pp.Init(ev)
		case events.PostprocessingStepFinished:
			if ev.UploadID == "" {
				// no current upload - this was an on demand scan
				continue
			}
			pp, err = getPP(pps.store, ev.UploadID)
			if err != nil {
				pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
				continue
			}
			next = pp.NextStep(ev)
		case events.StartPostprocessingStep:
			if ev.StepToStart != events.PPStepDelay {
				continue
			}
			pp, err = getPP(pps.store, ev.UploadID)
			if err != nil {
				pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
				continue
			}
			next = pp.Delay(ev)
		case events.UploadReady:
			// the storage provider thinks the upload is done - so no need to keep it any more
			if err := pps.store.Delete(ev.UploadID); err != nil {
				pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot delete upload")
				continue
			}
		case events.ResumePostprocessing:
			pp, err = getPP(pps.store, ev.UploadID)
			if err != nil {
				if err == store.ErrNotFound {
					if err := events.Publish(ctx, pps.pub, events.RestartPostprocessing{
						UploadID:  ev.UploadID,
						Timestamp: ev.Timestamp,
					}); err != nil {
						pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot publish RestartPostprocessing event")
					}
					continue
				}
				pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
				continue
			}
			next = pp.CurrentStep()
		}

		if pp != nil {
			if err := storePP(pps.store, pp); err != nil {
				pps.log.Error().Str("uploadID", pp.ID).Err(err).Msg("cannot store upload")
				continue // TODO: should we really continue here?
			}
		}
		if next != nil {
			if err := events.Publish(ctx, pps.pub, next); err != nil {
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

func storePP(sto store.Store, pp *postprocessing.Postprocessing) error {
	b, err := json.Marshal(pp)
	if err != nil {
		return err
	}

	return sto.Write(&store.Record{
		Key:   pp.ID,
		Value: b,
	})
}

func getPP(sto store.Store, uploadID string) (*postprocessing.Postprocessing, error) {
	recs, err := sto.Read(uploadID)
	if err != nil {
		return nil, err
	}

	if len(recs) != 1 {
		return nil, fmt.Errorf("expected only one result for '%s', got %d", uploadID, len(recs))
	}

	var pp postprocessing.Postprocessing
	return &pp, json.Unmarshal(recs[0].Value, &pp)
}
