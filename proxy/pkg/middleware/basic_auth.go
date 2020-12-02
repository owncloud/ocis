package middleware

import (
	"fmt"
	"net/http"
	"strings"

	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
)

const publicFilesEndpoint = "/remote.php/dav/public-files/"

// BasicAuth provides a middleware to check if BasicAuth is provided
func BasicAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger
	oidcIss := options.OIDCIss

	if options.EnableBasicAuth {
		options.Logger.Warn().Msg("basic auth enabled, use only for testing or development")
	}

	h := basicAuth{
		logger:         logger,
		enabled:        options.EnableBasicAuth,
		accountsClient: options.AccountsClient,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				if h.isPublicLink(req) || !h.isBasicAuth(req) {
					// if we want to prevent duplicated Www-Authenticate headers coming from Reva consider using w.Header().Del("Www-Authenticate")
					// but this will require the proxy being aware of endpoints which authentication fallback to Reva.
					if !h.isPublicLink(req) {
						for i := 0; i < len(ProxyWwwAuthenticate); i++ {
							if strings.Contains(req.RequestURI, fmt.Sprintf("/%v/", ProxyWwwAuthenticate[i])) {
								for k, v := range options.CredentialsByUserAgent {
									if strings.Contains(k, req.UserAgent()) {
										w.Header().Del("Www-Authenticate")
										w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(v), req.Host))
										goto OUT
									}
								}
								w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", "Basic", req.Host))
							}
						}
					}
				OUT:
					next.ServeHTTP(w, req)
					return
				}

				w.Header().Del("Www-Authenticate")

				account, ok := h.getAccount(req)

				// touch is a user agent locking guard, when touched changes to true it indicates the User-Agent on the
				// request is configured to support only one challenge, it it remains untouched, there are no considera-
				// tions and we should write all available authentication challenges to the response.
				touch := false

				if !ok {
					// if the request is bound to a user agent the locked write Www-Authenticate for such user
					for k, v := range options.CredentialsByUserAgent {
						if strings.Contains(k, req.UserAgent()) {
							w.Header().Del("Www-Authenticate")
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
					OcisID: account.Id,
					Iss:    oidcIss,
				}

				next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), claims)))
			},
		)
	}
}

type basicAuth struct {
	logger         log.Logger
	enabled        bool
	accountsClient accounts.AccountsService
}

func (m basicAuth) isPublicLink(req *http.Request) bool {
	login, _, ok := req.BasicAuth()

	return ok && login == "public" && strings.HasPrefix(req.URL.Path, publicFilesEndpoint)
}

func (m basicAuth) isBasicAuth(req *http.Request) bool {
	login, password, ok := req.BasicAuth()

	return m.enabled && ok && login != "" && password != ""
}

func (m basicAuth) getAccount(req *http.Request) (*accounts.Account, bool) {
	login, password, _ := req.BasicAuth()

	account, status := getAccount(
		m.logger,
		m.accountsClient,
		fmt.Sprintf(
			"login eq '%s' and password eq '%s'",
			strings.ReplaceAll(login, "'", "''"),
			strings.ReplaceAll(password, "'", "''"),
		),
	)

	return account, status == 0
}
