package authorization

import (
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	pgMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/authz/v0"
	pgService "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/authz/v0"
)

// Authz uses rpc calls to authorize a request
type Authz struct {
	authzClient pgService.AuthzProviderService
}

// NewAuthz returns a ready to use Authz Authorizer.
func NewAuthz() (Authz, error) {
	opaAuthorizer := Authz{
		authzClient: pgService.NewAuthzProviderService("com.owncloud.api.authz", grpc.DefaultClient()),
	}

	return opaAuthorizer, nil
}

// Authorize implements the Authorizer interface to authorize requests via grpc call.
func (a Authz) Authorize(r *http.Request) (bool, error) {
	req := &pgService.AllowedRequest{
		Name:   r.URL.Path,
		Method: r.Method,
		Stage:  pgMessage.Stage_STAGE_HTTP,
	}

	if user, ok := revactx.ContextGetUser(r.Context()); ok {
		req.User = &pgMessage.User{
			Username:    user.Username,
			Mail:        user.Mail,
			DisplayName: user.DisplayName,
			Groups:      user.Groups,
		}
	}

	rsp, err := a.authzClient.Allowed(r.Context(), req)
	if err != nil {
		return false, err
	}

	return rsp.Allowed, nil
}
