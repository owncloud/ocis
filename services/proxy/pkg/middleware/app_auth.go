package middleware

import (
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// AppAuthAuthenticator defines the app auth authenticator
type AppAuthAuthenticator struct {
	Logger              log.Logger
	RevaGatewaySelector *pool.Selector[gateway.GatewayAPIClient]
}

// Authenticate implements the authenticator interface to authenticate requests via app auth.
func (m AppAuthAuthenticator) Authenticate(r *http.Request) (*http.Request, bool) {
	if isPublicPath(r.URL.Path) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
		return nil, false
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, false
	}
	next, err := m.RevaGatewaySelector.Next()
	if err != nil {
		return nil, false
	}
	authenticateResponse, err := next.Authenticate(r.Context(), &gateway.AuthenticateRequest{
		Type:         "appauth",
		ClientId:     username,
		ClientSecret: password,
	})
	if err != nil {
		return nil, false
	}
	r.Header.Add(_headerRevaAccessToken, authenticateResponse.GetToken())

	return r, true
}
