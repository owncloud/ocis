package middleware

import (
	"fmt"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"net/http"
	"strings"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

func BasicAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)

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
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	m.logger.Warn().Msg("basic auth enabled, use only for testing or development")

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

func (m basicAuth) shouldServe(req *http.Request) bool {
	login, password, ok := req.BasicAuth()

	if ok && login == "public" && strings.HasPrefix(req.URL.Path, "/remote.php/dav/public-files/") {
		return true
	}

	if m.enabled && ok && login != "" && password != ""{
		return true
	}

	return false
}
