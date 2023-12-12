package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/postprocessing"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
)

// PostprocessingService is an instance of the service handling postprocessing of files
type PostprocessingService struct {
	ctx    context.Context
	log    log.Logger
	events <-chan events.Event
	pub    events.Publisher
	steps  []events.Postprocessingstep
	store  store.Store
	c      config.Postprocessing
	tp     trace.TracerProvider
}

var (
	// errFatal is returned when a fatal error occurs and we want to exit.
	errFatal = errors.New("fatal error")
	// ErrEvent is returned when something went wrong with a specific event.
	errEvent = errors.New("event error")
)

// NewPostprocessingService returns a new instance of a postprocessing service
func NewPostprocessingService(ctx context.Context, stream events.Stream, logger log.Logger, sto store.Store, tp trace.TracerProvider, c config.Postprocessing) (*PostprocessingService, error) {
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
		ctx:    ctx,
		log:    logger,
		events: evs,
		pub:    stream,
		steps:  getSteps(c),
		store:  sto,
		c:      c,
		tp:     tp,
	}, nil
}

// Run to fulfil Runner interface
func (pps *PostprocessingService) Run() error {
	for e := range pps.events {
		err := pps.processEvent(e)
		if err != nil {
			switch {
			case errors.Is(err, errFatal):
				return err
			case errors.Is(err, errEvent):
				continue
			default:
				pps.log.Fatal().Err(err).Msg("unknown error - exiting")
			}
		}
	}
	return nil
}

func (pps *PostprocessingService) processEvent(e events.Event) error {
	var (
		next interface{}
		pp   *postprocessing.Postprocessing
		err  error
	)

	ctx := e.GetTraceContext(pps.ctx)
	ctx, span := pps.tp.Tracer("postprocessing").Start(ctx, "processEvent")
	defer span.End()

	switch ev := e.Event.(type) {
	case events.BytesReceived:
		pp = &postprocessing.Postprocessing{
			ID:         ev.UploadID,
			URL:        ev.URL,
			User:       ev.ExecutingUser,
			Filename:   ev.Filename,
			Filesize:   ev.Filesize,
			ResourceID: ev.ResourceID,
			Steps:      pps.steps,
		}
		next = pp.Init(ev)
	case events.PostprocessingStepFinished:
		if ev.UploadID == "" {
			// no current upload - this was an on demand scan
			return nil
		}
		pp, err = pps.getPP(pps.store, ev.UploadID)
		if err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
			return fmt.Errorf("%w: cannot get upload", errEvent)
		}
		next = pp.NextStep(ev)

		switch pp.Status.Outcome {
		case events.PPOutcomeRetry:
			// schedule retry
			backoff := pp.BackoffDuration()
			go func() {
				time.Sleep(backoff)
				retryEvent := events.StartPostprocessingStep{
					UploadID:      pp.ID,
					URL:           pp.URL,
					ExecutingUser: pp.User,
					Filename:      pp.Filename,
					Filesize:      pp.Filesize,
					ResourceID:    pp.ResourceID,
					StepToStart:   pp.Status.CurrentStep,
				}
				err := events.Publish(ctx, pps.pub, retryEvent)
				if err != nil {
					pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot publish RestartPostprocessing event")
				}
			}()
		}
	case events.StartPostprocessingStep:
		if ev.StepToStart != events.PPStepDelay {
			return nil
		}
		pp, err = pps.getPP(pps.store, ev.UploadID)
		if err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
			return fmt.Errorf("%w: cannot get upload", errEvent)
		}
		next = pp.Delay(ev)
	case events.UploadReady:
		// the storage provider thinks the upload is done - so no need to keep it any more
		if err := pps.store.Delete(ev.UploadID); err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot delete upload")
			return fmt.Errorf("%w: cannot delete upload", errEvent)
		}
	case events.ResumePostprocessing:
		pp, err = pps.getPP(pps.store, ev.UploadID)
		if err != nil {
			if err == store.ErrNotFound {
				if err := events.Publish(ctx, pps.pub, events.RestartPostprocessing{
					UploadID:  ev.UploadID,
					Timestamp: ev.Timestamp,
				}); err != nil {
					pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot publish RestartPostprocessing event")
				}
				return fmt.Errorf("%w: cannot publish RestartPostprocessing event", errEvent)
			}
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
			return fmt.Errorf("%w: cannot get upload", errEvent)
		}
		next = pp.CurrentStep()
	}

	if pp != nil {
		if err := storePP(pps.store, pp); err != nil {
			pps.log.Error().Str("uploadID", pp.ID).Err(err).Msg("cannot store upload")
			return fmt.Errorf("%w: cannot store upload", errEvent)
		}
	}
	if next != nil {
		if err := events.Publish(ctx, pps.pub, next); err != nil {
			pps.log.Error().Err(err).Msg("unable to publish event")
			return fmt.Errorf("%w: unable to publish event", errFatal) // we can't publish -> we are screwed
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

func (pps *PostprocessingService) getPP(sto store.Store, uploadID string) (*postprocessing.Postprocessing, error) {
	recs, err := sto.Read(uploadID)
	if err != nil {
		return nil, err
	}

	if len(recs) != 1 {
		return nil, fmt.Errorf("expected only one result for '%s', got %d", uploadID, len(recs))
	}

	pp := postprocessing.New(pps.c)
	err = json.Unmarshal(recs[0].Value, pp)
	if err != nil {
		return nil, err
	}

	return pp, nil
}
