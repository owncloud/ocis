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
	if options.EnableBasicAuth {
		options.Logger.Warn().Msg("basic auth enabled, use only for testing or development")
	}

	return func(next http.Handler) http.Handler {
		return &basicAuth{
			next:           next,
			logger:         options.Logger,
			enabled:        options.EnableBasicAuth,
			accountsClient: options.AccountsClient,
			oidcIss:        options.OIDCIss,
		}
	}
}

type basicAuth struct {
	next           http.Handler
	logger         log.Logger
	enabled        bool
	accountsClient accounts.AccountsService
	oidcIss        string
}

func (m basicAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m.isPublicLink(req) || !m.isBasicAuth(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	login, password, _ := req.BasicAuth()

	account, status := getAccount(m.logger, m.accountsClient, fmt.Sprintf("login eq '%s' and password eq '%s'", strings.ReplaceAll(login, "'", "''"), strings.ReplaceAll(password, "'", "''")))

	if status != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &oidc.StandardClaims{
		OcisID: account.Id,
		Iss:    m.oidcIss,
	}

	m.next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), claims)))
}

func (m basicAuth) isPublicLink(req *http.Request) bool {
	login, _, ok := req.BasicAuth()

	return ok && login == "public" && strings.HasPrefix(req.URL.Path, publicFilesEndpoint)
}

func (m basicAuth) isBasicAuth(req *http.Request) bool {
	login, password, ok := req.BasicAuth()

	return m.enabled && ok && login != "" && password != ""
}
