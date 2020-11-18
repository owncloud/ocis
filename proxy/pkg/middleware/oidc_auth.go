package middleware

import (
	"context"
	"net/http"
	"strings"

	gOidc "github.com/coreos/go-oidc"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/cache"
	"golang.org/x/oauth2"
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

// OIDCAuth provides a middleware to check access secured by a static token.
func OIDCAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	tokenCache := cache.NewCache(
		cache.Size(options.UserinfoCacheSize),
		cache.TTL(options.UserinfoCacheTTL),
	)

	return func(next http.Handler) http.Handler {
		return &oidcAuth{
			next:         next,
			logger:       options.Logger,
			providerFunc: options.OIDCProviderFunc,
			httpClient:   options.HTTPClient,
			oidcIss:      options.OIDCIss,
			tokenCache:   &tokenCache,
		}
	}
}

type oidcAuth struct {
	next         http.Handler
	logger       log.Logger
	provider     OIDCProvider
	providerFunc func() (OIDCProvider, error)
	httpClient   *http.Client
	oidcIss      string
	tokenCache   *cache.Cache
}

func (m oidcAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	if m.getProvider() == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

	claims, status := m.getClaims(token, req)
	if status != 0 {
		w.WriteHeader(status)
		return
	}

	// inject claims to the request context for the account_uuid middleware.
	req = req.WithContext(oidc.NewContext(req.Context(), &claims))

	// store claims in context
	// uses the original context, not the one with probably reduced security
	m.next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), &claims)))
}

func (m oidcAuth) getClaims(token string, req *http.Request) (claims oidc.StandardClaims, status int) {
	hit := m.tokenCache.Get(token)
	if hit == nil {
		// TODO cache userinfo for access token if we can determine the expiry (which works in case it is a jwt based access token)
		oauth2Token := &oauth2.Token{
			AccessToken: token,
		}

		userInfo, err := m.provider.UserInfo(
			context.WithValue(req.Context(), oauth2.HTTPClient, m.httpClient),
			oauth2.StaticTokenSource(oauth2Token),
		)
		if err != nil {
			m.logger.Error().Err(err).Str("token", token).Msg("Failed to get userinfo")
			status = http.StatusUnauthorized
			return
		}

		if err := userInfo.Claims(&claims); err != nil {
			m.logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
			status = http.StatusInternalServerError
			return
		}

		m.logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Msg("unmarshalled userinfo")

		//TODO: This should be read from the token instead of config
		claims.Iss = m.oidcIss

		m.tokenCache.Set(token, claims)
		return
	}

	var ok = false
	if claims, ok = hit.V.(oidc.StandardClaims); !ok {
		status = http.StatusInternalServerError
		return
	}
	m.logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
	return
}

func (m oidcAuth) shouldServe(req *http.Request) bool {
	header := req.Header.Get("Authorization")

	if m.oidcIss == "" {
		return false
	}

	// todo: looks dirty, check later
	// TODO: make a PR to coreos/go-oidc for exposing userinfo endpoint on provider, see https://github.com/coreos/go-oidc/issues/248
	for _, ignoringPath := range []string{"/konnect/v1/userinfo"} {
		if req.URL.Path == ignoringPath {
			return false
		}
	}

	return strings.HasPrefix(header, "Bearer ")
}

func (m oidcAuth) getProvider() OIDCProvider {
	if m.provider == nil {
		// Lazily initialize a provider

		// provider needs to be cached as when it is created
		// it will fetch the keys from the issuer using the .well-known
		// endpoint
		provider, err := m.providerFunc()
		if err != nil {
			m.logger.Error().Err(err).Msg("could not initialize oidcAuth provider")
			return nil
		}

		m.provider = provider
	}
	return m.provider
}
