package middleware

import (
	"crypto/sha1"
	"fmt"
	"github.com/owncloud/ocis/proxy/pkg/cache"
	"net/http"
	"strings"
	"sync"
	"time"

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
	accountsCache := cache.NewCache(cache.Size(5000))
	if options.EnableBasicAuth {
		options.Logger.Warn().Msg("basic auth enabled, use only for testing or development")
	}

	h := basicAuth{
		logger:         logger,
		enabled:        options.EnableBasicAuth,
		accountsClient: options.AccountsClient,
		accountsCache:  &accountsCache,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				if h.isPublicLink(req) || !h.isBasicAuth(req) {
					next.ServeHTTP(w, req)
					return
				}

				account, ok := h.getAccount(req)

				if !ok {
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
	accountsCache  *cache.Cache
	m              sync.Mutex
}

func (m *basicAuth) isPublicLink(req *http.Request) bool {
	login, _, ok := req.BasicAuth()

	return ok && login == "public" && strings.HasPrefix(req.URL.Path, publicFilesEndpoint)
}

func (m *basicAuth) isBasicAuth(req *http.Request) bool {
	login, password, ok := req.BasicAuth()

	return m.enabled && ok && login != "" && password != ""
}

func (m *basicAuth) getAccount(req *http.Request) (*accounts.Account, bool) {
	var ok bool
	var hit *cache.Entry

	login, password, _ := req.BasicAuth()

	h := sha1.New()
	h.Write([]byte(login))

	lookup := fmt.Sprintf("%x", h.Sum([]byte(password)))

	m.m.Lock()
	defer m.m.Unlock()

	if hit = m.accountsCache.Get(lookup); hit == nil {
		account, status := getAccount(
			m.logger,
			m.accountsClient,
			fmt.Sprintf(
				"login eq '%s' and password eq '%s'",
				strings.ReplaceAll(login, "'", "''"),
				strings.ReplaceAll(password, "'", "''"),
			),
		)

		if ok = status == 0; ok {
			m.accountsCache.Set(lookup, account, time.Now().Add(10*time.Minute))
		}

		return account, ok
	}

	account, ok := hit.V.(*accounts.Account)

	return account, ok
}
