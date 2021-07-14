package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	gOidc "github.com/coreos/go-oidc"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"golang.org/x/oauth2"
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

// OIDCAuth provides a middleware to check access secured by a static token.
func OIDCAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	tokenCache := sync.NewCache(options.UserinfoCacheSize)

	h := oidcAuth{
		logger:        options.Logger,
		providerFunc:  options.OIDCProviderFunc,
		httpClient:    options.HTTPClient,
		oidcIss:       options.OIDCIss,
		tokenCache:    &tokenCache,
		tokenCacheTTL: options.UserinfoCacheTTL,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// there is no bearer token on the request,
			if !h.shouldServe(req) {
				// oidc supported but token not present, add header and handover to the next middleware.
				userAgentAuthenticateLockIn(w, req, options.CredentialsByUserAgent, "bearer")
				next.ServeHTTP(w, req)
				return
			}

			if h.getProvider() == nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

			claims, status := h.getClaims(token, req)
			if status != 0 {
				w.WriteHeader(status)
				return
			}

			// inject claims to the request context for the account_uuid middleware.
			req = req.WithContext(oidc.NewContext(req.Context(), &claims))

			// store claims in context
			// uses the original context, not the one with probably reduced security
			next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), &claims)))
		})
	}
}

type oidcAuth struct {
	logger        log.Logger
	provider      OIDCProvider
	providerFunc  func() (OIDCProvider, error)
	httpClient    *http.Client
	oidcIss       string
	tokenCache    *sync.Cache
	tokenCacheTTL time.Duration
}

func (m oidcAuth) getClaims(token string, req *http.Request) (claims oidc.StandardClaims, status int) {
	hit := m.tokenCache.Load(token)
	if hit == nil {
		// TODO cache userinfo for access token if we can determine the expiry (which works in case it is a jwt based access token)
		oauth2Token := &oauth2.Token{
			AccessToken: token,
		}

		userInfo, err := m.getProvider().UserInfo(
			context.WithValue(req.Context(), oauth2.HTTPClient, m.httpClient),
			oauth2.StaticTokenSource(oauth2Token),
		)
		if err != nil {
			m.logger.Error().Err(err).Msg("Failed to get userinfo")
			status = http.StatusUnauthorized
			return
		}

		// TODO allow extracting arbitrary claims ... or require idp to send a specific claim
		if err := userInfo.Claims(&claims); err != nil {
			m.logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
			status = http.StatusInternalServerError
			return
		}

		//TODO: This should be read from the token instead of config
		claims.Iss = m.oidcIss

		expiration := m.extractExpiration(token)
		m.tokenCache.Store(token, claims, expiration)

		m.logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Time("expiration", expiration.UTC()).Msg("unmarshalled and cached userinfo")
		return
	}

	var ok bool
	if claims, ok = hit.V.(oidc.StandardClaims); !ok {
		status = http.StatusInternalServerError
		return
	}
	m.logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
	return
}

// extractExpiration tries to parse and extract the expiration from the provided token. It might not even be a jwt.
// defaults to the configured fallback TTL.
// TODO: use introspection endpoint if available in the oidc configuration. Still needs a fallback to configured TTL.
func (m oidcAuth) extractExpiration(token string) time.Time {
	defaultExpiration := time.Now().Add(m.tokenCacheTTL)

	s := strings.SplitN(token, ".", 4)
	if len(s) != 3 {
		return defaultExpiration
	}

	b, err := jwt.DecodeSegment(s[1])
	if err != nil {
		return defaultExpiration
	}

	at := &jwt.StandardClaims{}
	err = json.Unmarshal(b, at)
	if err != nil || at.ExpiresAt == 0 {
		return defaultExpiration
	}

	return time.Unix(at.ExpiresAt, 0)
}

func (m oidcAuth) shouldServe(req *http.Request) bool {
	header := req.Header.Get("Authorization")

	if m.oidcIss == "" {
		return false
	}

	// todo: looks dirty, check later
	// TODO: make a PR to coreos/go-oidc for exposing userinfo endpoint on provider, see https://github.com/coreos/go-oidc/issues/248
	for _, ignoringPath := range []string{"/konnect/v1/userinfo", "/status.php"} {
		if req.URL.Path == ignoringPath {
			return false
		}
	}

	return strings.HasPrefix(header, "Bearer ")
}

func (m *oidcAuth) getProvider() OIDCProvider {
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
