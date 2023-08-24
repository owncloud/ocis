package eventSVC

import (
	"context"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
	"go.opentelemetry.io/otel/trace"
)

// Service defines the service handlers.
type Service struct {
	ctx    context.Context
	query  string
	log    log.Logger
	stream events.Stream
	engine engine.Engine
	tp     trace.TracerProvider
}

// New returns a service implementation for Service.
func New(ctx context.Context, stream events.Stream, logger log.Logger, tp trace.TracerProvider, engine engine.Engine, query string) (Service, error) {
	svc := Service{
		ctx:    ctx,
		log:    logger,
		query:  query,
		tp:     tp,
		engine: engine,
		stream: stream,
	}

	return svc, nil
}

// Run to fulfil Runner interface
func (s Service) Run() error {
	ch, err := events.Consume(s.stream, "policies", events.StartPostprocessingStep{})
	if err != nil {
		return err
	}

	for e := range ch {
		err := s.processEvent(e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) processEvent(e events.Event) error {
	ctx := e.GetTraceContext(s.ctx)
	ctx, span := s.tp.Tracer("policies").Start(ctx, "processEvent")
	defer span.End()

	switch ev := e.Event.(type) {
	case events.StartPostprocessingStep:
		if ev.StepToStart != events.PPStepPolicies {
			return nil
		}

		outcome := events.PPOutcomeContinue

		if s.query != "" {
			env := engine.Environment{
				Stage: engine.StagePP,
				Resource: engine.Resource{
					Name: ev.Filename,
					URL:  ev.URL,
					Size: ev.Filesize,
				},
			}

			if ev.ExecutingUser != nil {
				env.User = *ev.ExecutingUser
			}

			if ev.ResourceID != nil {
				env.Resource.ID = *ev.ResourceID
			}

			result, err := s.engine.Evaluate(context.TODO(), s.query, env)
			if err != nil {
				s.log.Error().Err(err).Msg("unable evaluate policy")
			}

			if !result {
				outcome = events.PPOutcomeDelete
			}
		}

		if err := events.Publish(ctx, s.stream, events.PostprocessingStepFinished{
			Outcome:       outcome,
			UploadID:      ev.UploadID,
			ExecutingUser: ev.ExecutingUser,
			Filename:      ev.Filename,
			FinishedStep:  ev.StepToStart,
		}); err != nil {
			return err
		}
	}
	return nil
}
