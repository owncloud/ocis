package eventSVC

import (
	"context"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

// Service defines the service handlers.
type Service struct {
	query  string
	log    log.Logger
	stream events.Stream
	engine engine.Engine
}

// New returns a service implementation for Service.
func New(stream events.Stream, logger log.Logger, engine engine.Engine, query string) (Service, error) {
	svc := Service{
		log:    logger,
		query:  query,
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
		switch ev := e.Event.(type) {
		case events.StartPostprocessingStep:
			if ev.StepToStart != "policies" {
				continue
			}

			env := engine.Environment{
				Stage: engine.StagePP,
				User:  *ev.ExecutingUser,
				Resource: engine.Resource{
					ID:   *ev.ResourceID,
					Name: ev.Filename,
					URL:  ev.URL,
					Size: ev.Filesize,
				},
			}

			result, err := s.engine.Evaluate(context.TODO(), s.query, env)
			if err != nil {
				s.log.Error().Err(err).Msg("unable evaluate policy")
			}

			outcome := events.PPOutcomeContinue
			if !result {
				outcome = events.PPOutcomeDelete
			}

			if err := events.Publish(s.stream, events.PostprocessingStepFinished{
				Outcome:       outcome,
				UploadID:      ev.UploadID,
				ExecutingUser: ev.ExecutingUser,
				Filename:      ev.Filename,
				FinishedStep:  ev.StepToStart,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
