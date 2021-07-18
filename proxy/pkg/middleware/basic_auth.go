package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
)

const publicFilesEndpoint = "/remote.php/dav/public-files/"

// BasicAuth provides a middleware to check if BasicAuth is provided
func BasicAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	if options.EnableBasicAuth {
		options.Logger.Warn().Msg("basic auth enabled, use only for testing or development")
	}

	h := basicAuth{
		logger:       logger,
		enabled:      options.EnableBasicAuth,
		userProvider: options.UserProvider,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				if h.isPublicLink(req) || !h.isBasicAuth(req) || h.isOIDCTokenAuth(req) {
					if !h.isPublicLink(req) {
						userAgentAuthenticateLockIn(w, req, options.CredentialsByUserAgent, "basic")
					}
					next.ServeHTTP(w, req)
					return
				}

				removeSuperfluousAuthenticate(w)
				login, password, _ := req.BasicAuth()
				user, err := h.userProvider.Authenticate(req.Context(), login, password)

				// touch is a user agent locking guard, when touched changes to true it indicates the User-Agent on the
				// request is configured to support only one challenge, it it remains untouched, there are no considera-
				// tions and we should write all available authentication challenges to the response.
				touch := false

				if err != nil {
					for k, v := range options.CredentialsByUserAgent {
						if strings.Contains(k, req.UserAgent()) {
							removeSuperfluousAuthenticate(w)
							w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(v), req.Host))
							touch = true
							break
						}
					}

					// if the request is not bound to any user agent, write all available challenges
					if !touch {
						writeSupportedAuthenticateHeader(w, req)
					}

					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				claims := &oidc.StandardClaims{
					OcisID:            user.Id.OpaqueId,
					Iss:               user.Id.Idp,
					PreferredUsername: user.Username,
					Email:             user.Mail,
				}

				next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), claims)))
			},
		)
	}
}

type basicAuth struct {
	logger       log.Logger
	enabled      bool
	userProvider backend.UserBackend
}

func (m basicAuth) isPublicLink(req *http.Request) bool {
	login, _, ok := req.BasicAuth()
	return ok && login == "public" && strings.HasPrefix(req.URL.Path, publicFilesEndpoint)
}

// The token auth endpoint uses basic auth for clients, see https://openid.net/specs/openid-connect-basic-1_0.html#TokenRequest
// > The Client MUST authenticate to the Token Endpoint using the HTTP Basic method, as described in 2.3.1 of OAuth 2.0.
func (m basicAuth) isOIDCTokenAuth(req *http.Request) bool {
	return req.URL.Path == "/konnect/v1/token"
}

func (m basicAuth) isBasicAuth(req *http.Request) bool {
	login, password, ok := req.BasicAuth()
	return m.enabled && ok && login != "" && password != ""
}
