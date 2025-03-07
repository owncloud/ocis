package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/postprocessing"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/utils"
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
	// ErrFatal is returned when a fatal error occurs and we want to exit.
	ErrFatal = errors.New("fatal error")
	// ErrEvent is returned when something went wrong with a specific event.
	ErrEvent = errors.New("event error")
	// ErrNotFound is returned when a postprocessing is not found in the store.
	ErrNotFound = errors.New("postprocessing not found")
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
	wg := sync.WaitGroup{}

	for i := 0; i < pps.c.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range pps.events {
				if err := pps.processEvent(e); err != nil {
					switch {
					case errors.Is(err, ErrFatal):
						pps.log.Fatal().Err(err).Msg("fatal error - exiting")
					case errors.Is(err, ErrEvent):
						pps.log.Error().Err(err).Msg("continuing")
					default:
						pps.log.Fatal().Err(err).Msg("unknown error - exiting")
					}
				}
			}
		}()
	}
	wg.Wait()

	return nil
}

func (pps *PostprocessingService) processEvent(e events.Event) error {
	pps.log.Debug().Str("Type", e.Type).Str("ID", e.ID).Msg("processing event received")

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
			ID:                ev.UploadID,
			URL:               ev.URL,
			User:              ev.ExecutingUser,
			Filename:          ev.Filename,
			Filesize:          ev.Filesize,
			ResourceID:        ev.ResourceID,
			Steps:             pps.steps,
			InitiatorID:       e.InitiatorID,
			ImpersonatingUser: ev.ImpersonatingUser,
		}
		pps.log.Info().Str("UploadID", ev.UploadID).Msg("processing init")
		next = pp.Init(ev)
	case events.PostprocessingStepFinished:
		sublog := pps.log.Debug()
		if ev.Outcome == events.PPOutcomeRetry {
			sublog = pps.log.Error()
		}
		sublog.Str("UploadID", ev.UploadID).Str("step", string(ev.FinishedStep)).Str("outcome", string(ev.Outcome)).Msg("processing step finished")

		if ev.UploadID == "" {
			// no current upload - this was an on demand scan
			return nil
		}
		pp, err = pps.getPP(pps.store, ev.UploadID)
		if err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
			return fmt.Errorf("%w: cannot get upload", ErrEvent)
		}
		next = pp.NextStep(ev)

		if pp.Status.Outcome == events.PPOutcomeRetry {
			// schedule retry
			backoff := pp.BackoffDuration()
			pps.log.Info().Str("UploadID", ev.UploadID).Str("step", string(ev.FinishedStep)).Err(ev.Error).Msg("retrying step in " + backoff.String())
			go func() {
				time.Sleep(backoff)
				retryEvent := events.StartPostprocessingStep{
					UploadID:          pp.ID,
					URL:               pp.URL,
					ExecutingUser:     pp.User,
					Filename:          pp.Filename,
					Filesize:          pp.Filesize,
					ResourceID:        pp.ResourceID,
					StepToStart:       pp.Status.CurrentStep,
					ImpersonatingUser: pp.ImpersonatingUser,
				}
				err := events.Publish(ctx, pps.pub, retryEvent)
				if err != nil {
					pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot publish RestartPostprocessing event")
				}
			}()
		}

		if pp.Status.CurrentStep == events.PPStepFinished {
			pps.log.Info().Str("UploadID", e.ID).Msg("processing finished")
		}

	case events.StartPostprocessingStep:
		pps.log.Debug().Str("UploadID", ev.UploadID).Str("step", string(ev.StepToStart)).Msg("processing step started")
		if ev.StepToStart != events.PPStepDelay {
			return nil
		}
		pp, err = pps.getPP(pps.store, ev.UploadID)
		if err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
			return fmt.Errorf("%w: cannot get upload", ErrEvent)
		}
		pp.Delay(func(next interface{}) {
			if err := events.Publish(ctx, pps.pub, next); err != nil {
				pps.log.Error().Err(err).Msg("cannot publish event")
			}
		})
	case events.UploadReady:
		pps.log.Debug().Str("UploadID", e.ID).Str("filename", ev.Filename).Msg("processing UploadReady")
		if ev.Failed {
			// the upload failed - let's keep it around for a while - but mark it as finished
			pp, err = pps.getPP(pps.store, ev.UploadID)
			if err != nil {
				pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot get upload")
				return fmt.Errorf("%w: cannot get upload", ErrEvent)
			}
			pp.Finished = true
			return storePP(pps.store, pp)
		}

		// the storage provider thinks the upload is done - so no need to keep it any more
		if err := pps.store.Delete(ev.UploadID); err != nil {
			pps.log.Error().Str("uploadID", ev.UploadID).Err(err).Msg("cannot delete upload")
			return fmt.Errorf("%w: cannot delete upload", ErrEvent)
		}
	case events.ResumePostprocessing:
		pps.log.Info().Str("UploadID", ev.UploadID).Str("step", string(ev.Step)).Msg("processing resumed")
		return pps.handleResumePPEvent(ctx, ev)
	}

	if pp != nil {
		ctx = ctxpkg.ContextSetInitiator(ctx, pp.InitiatorID)

		if err := storePP(pps.store, pp); err != nil {
			pps.log.Error().Str("uploadID", pp.ID).Err(err).Msg("cannot store upload")
			return fmt.Errorf("%w: cannot store upload", ErrEvent)
		}
	}

	if next != nil {
		if err := events.Publish(ctx, pps.pub, next); err != nil {
			pps.log.Error().Err(err).Msg("unable to publish event")
			return fmt.Errorf("%w: unable to publish event", ErrFatal) // we can't publish -> we are screwed
		}
	}
	return nil
}

func (pps *PostprocessingService) getPP(sto store.Store, uploadID string) (*postprocessing.Postprocessing, error) {
	recs, err := sto.Read(uploadID)
	if err != nil {
		if err == store.ErrNotFound {
			pps.log.Error().Str("uploadID", uploadID).Err(err).Msg("reading from store: upload not found in the store")
			return nil, ErrNotFound
		}
		pps.log.Error().Str("uploadID", uploadID).Err(err).Msg("reading from store: upload store read error")
		return nil, err
	}

	if len(recs) == 0 {
		pps.log.Error().Str("uploadID", uploadID).Err(err).Msg("reading from store: empty upload records in the store")
		return nil, ErrNotFound
	}

	if len(recs) > 1 {
		pps.log.Error().Str("uploadID", uploadID).Int("records", len(recs)).Err(err).Msg("reading from store: expected only one result")
		return nil, fmt.Errorf("expected only one result for '%s', got %d", uploadID, len(recs))
	}

	pp := postprocessing.New(pps.c)
	err = json.Unmarshal(recs[0].Value, pp)
	if err != nil {
		pps.log.Error().Str("uploadID", uploadID).Err(err).Msg("reading from store: unmarshaling error")
		return nil, err
	}

	return pp, nil
}

func getSteps(c config.Postprocessing) []events.Postprocessingstep {
	// NOTE: improved version only allows configuring order of postprocessing steps
	// But we aim for a system where postprocessing steps can be configured per space, ideally by the spaceadmin itself
	// We need to iterate over configuring PP service when we see fit
	steps := make([]events.Postprocessingstep, 0, len(c.Steps))
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

func (pps *PostprocessingService) handleResumePPEvent(ctx context.Context, ev events.ResumePostprocessing) error {
	ids := []string{ev.UploadID}
	if ev.Step != "" {
		ids = pps.findUploadsByStep(ev.Step)
	}

	for _, id := range ids {
		if err := pps.resumePP(ctx, id); err != nil {
			pps.log.Error().Str("uploadID", id).Err(err).Msg("cannot resume upload")
		}
	}
	return nil
}

func (pps *PostprocessingService) resumePP(ctx context.Context, uploadID string) error {
	pp, err := pps.getPP(pps.store, uploadID)
	if err != nil {
		if err == ErrNotFound {
			if err := events.Publish(ctx, pps.pub, events.RestartPostprocessing{
				UploadID:  uploadID,
				Timestamp: utils.TSNow(),
			}); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("cannot get upload: %w", err)
	}

	if pp.Finished {
		// dont retry finished uploads
		return nil
	}

	return events.Publish(ctx, pps.pub, pp.CurrentStep())
}

func (pps *PostprocessingService) findUploadsByStep(step events.Postprocessingstep) []string {
	var ids []string

	keys, err := pps.store.List()
	if err != nil {
		pps.log.Error().Err(err).Msg("cannot list uploads")
	}

	for _, k := range keys {
		rec, err := pps.store.Read(k)
		if err != nil {
			pps.log.Error().Err(err).Msg("cannot read upload")
			continue
		}

		if len(rec) != 1 {
			pps.log.Error().Err(err).Msg("expected only one result")
			continue
		}

		pp := &postprocessing.Postprocessing{}
		err = json.Unmarshal(rec[0].Value, pp)
		if err != nil {
			pps.log.Error().Err(err).Msg("cannot unmarshal upload")
			continue
		}

		if pp.Status.CurrentStep == step {
			ids = append(ids, pp.ID)
		}
	}

	return ids
}
