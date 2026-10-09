package middleware

import (
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
)

// AppAuthAuthenticator defines the app auth authenticator
type AppAuthAuthenticator struct {
	Logger               log.Logger
	RevaGatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	UserRoleAssigner     userroles.UserRoleAssigner
	MultiInstanceEnabled bool
}

// Authenticate implements the authenticator interface to authenticate requests via app auth.
func (m AppAuthAuthenticator) Authenticate(r *http.Request) (*http.Request, error) {
	if isPublicPath(r.URL.Path) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
		return nil, ErrAuthenticationFailed
	}

	if m.MultiInstanceEnabled {
		// App passwords carry no OIDC claims, so there is no way to verify
		// tenant membership (OCIS_MULTI_INSTANCE_MEMBER_CLAIM/GUEST_CLAIM)
		// for them. Reject rather than let the request bypass the tenant
		// check in services/proxy/pkg/middleware/account_resolver.go's
		// resolveUserType.
		return nil, ErrAuthenticationFailed
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, ErrAuthenticationFailed
	}
	next, err := m.RevaGatewaySelector.Next()
	if err != nil {
		return nil, ErrAuthenticationFailed
	}

	authenticateResponse, err := next.Authenticate(r.Context(), &gateway.AuthenticateRequest{
		Type:         "appauth",
		ClientId:     username,
		ClientSecret: password,
	})
	if err != nil {
		return nil, ErrAuthenticationFailed
	}
	if authenticateResponse.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
		// TODO: log???
		return nil, ErrAuthenticationFailed
	}

	user, err := m.UserRoleAssigner.ApplyUserRole(r.Context(), authenticateResponse.GetUser())
	if err != nil {
		m.Logger.Error().Err(err).Msg("could not apply user role")
		return nil, ErrAuthenticationFailed
	}

	// only mutate the request once we know we're returning success, so a
	// failed authentication attempt can't leave a stale access token behind
	// for a later authenticator in the chain to build on.
	r.Header.Set(revactx.TokenHeader, authenticateResponse.GetToken())
	// the user is already authenticated, put it into the context directly
	// instead of faking oidc claims for the account resolver to re-resolve.
	r = r.WithContext(revactx.ContextSetUser(r.Context(), user))

	return r, nil
}
