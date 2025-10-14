package grpcSVC

import (
	"context"

	"github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

// Service defines the service handlers.
type Service struct {
	engine engine.Engine
}

// New returns a service implementation for Service.
func New(engine engine.Engine) (Service, error) {
	svc := Service{
		engine: engine,
	}

	return svc, nil
}

// Evaluate exposes the engine policy evaluation.
func (s Service) Evaluate(ctx context.Context, request *v0.EvaluateRequest, response *v0.EvaluateResponse) error {
	env, err := engine.NewEnvironmentFromPB(request.Environment)
	if err != nil {
		return err
	}

	result, err := s.engine.Evaluate(ctx, request.Query, env)
	response.Result = result

	return err
}
