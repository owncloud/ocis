package grpcSVC

import (
	"context"
	"encoding/json"

	"github.com/owncloud/ocis/v2/protogen/gen/ocis/services/authz/v0"
	"github.com/owncloud/ocis/v2/services/authz/pkg/authz"
	"google.golang.org/protobuf/encoding/protojson"
)

// Service defines the service handlers.
type Service struct {
	authorizers []authz.Authorizer
}

// New returns a service implementation for Service.
func New(authorizers []authz.Authorizer) (Service, error) {
	svc := Service{
		authorizers: authorizers,
	}

	return svc, nil
}

func (s Service) Allowed(ctx context.Context, request *v0.AllowedRequest, response *v0.AllowedResponse) error {
	rData, err := protojson.Marshal(request)
	if err != nil {
		return err
	}

	env := authz.Environment{}
	if json.Unmarshal(rData, &env); err != nil {
		return err
	}

	allowed, err := authz.Authorized(ctx, env, s.authorizers...)
	response.Allowed = allowed

	return err
}
