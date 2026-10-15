package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
)

// BasicAuthenticator is the authenticator responsible for HTTP Basic authentication.
type BasicAuthenticator struct {
	Logger               log.Logger
	UserProvider         backend.UserBackend
	UserRoleAssigner     userroles.UserRoleAssigner
	MultiInstanceEnabled bool
}

// Authenticate implements the authenticator interface to authenticate requests via basic auth.
func (m BasicAuthenticator) Authenticate(r *http.Request) (*http.Request, error) {
	if isPublicPath(r.URL.Path) && isPublicWithShareToken(r) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
		return nil, ErrAuthenticationFailed
	}

	if m.MultiInstanceEnabled {
		// Basic auth credentials carry no OIDC claims, so there is no way to
		// verify tenant membership (OCIS_MULTI_INSTANCE_MEMBER_CLAIM/GUEST_CLAIM)
		// for them. Reject rather than let the request bypass the tenant
		// check in services/proxy/pkg/middleware/account_resolver.go's
		// resolveUserType.
		return nil, ErrAuthenticationFailed
	}

	login, password, ok := r.BasicAuth()
	if !ok {
		return nil, ErrAuthenticationFailed
	}

	user, _, err := m.UserProvider.Authenticate(r.Context(), login, password)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "basic").
			Str("path", r.URL.Path).
			Msg("failed to authenticate request")
		return nil, ErrAuthenticationFailed
	}

	user, err = m.UserRoleAssigner.ApplyUserRole(r.Context(), user)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "basic").
			Str("path", r.URL.Path).
			Msg("could not apply user role")
		return nil, ErrAuthenticationFailed
	}

	m.Logger.Debug().
		Str("authenticator", "basic").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	// the user is already authenticated, put it into the context directly
	// instead of faking oidc claims for the account resolver to re-resolve.
	return r.WithContext(revactx.ContextSetUser(r.Context(), user)), nil
}
