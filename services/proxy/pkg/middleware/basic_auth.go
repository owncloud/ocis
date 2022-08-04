package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
)

type BasicAuthenticator struct {
	Logger        log.Logger
	UserProvider  backend.UserBackend
	UserCS3Claim  string
	UserOIDCClaim string
}

func (m BasicAuthenticator) Authenticate(req *http.Request) (*http.Request, bool) {
	if isPublicPath(req.URL.Path) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
		return nil, false
	}

	login, password, ok := req.BasicAuth()
	if !ok {
		return nil, false
	}

	user, _, err := m.UserProvider.Authenticate(req.Context(), login, password)
	if err != nil {
		// TODO add log line
		return nil, false
	}

	// fake oidc claims
	claims := map[string]interface{}{
		oidc.Iss:               user.Id.Idp,
		oidc.PreferredUsername: user.Username,
		oidc.Email:             user.Mail,
		oidc.OwncloudUUID:      user.Id.OpaqueId,
	}

	if m.UserCS3Claim == "userid" {
		// set the custom user claim only if users will be looked up by the userid on the CS3api
		// OpaqueId contains the userid configured in STORAGE_LDAP_USER_SCHEMA_UID
		claims[m.UserOIDCClaim] = user.Id.OpaqueId

	}
	return req.WithContext(oidc.NewContext(req.Context(), claims)), true
}
