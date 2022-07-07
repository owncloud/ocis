package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	gOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/sync"
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

			// inject claims to the request context for the account_resolver middleware.
			next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), claims)))
		})
	}
}

type oidcAuth struct {
	logger        log.Logger
	provider      OIDCProvider
	jwks          *keyfunc.JWKS
	providerFunc  func() (OIDCProvider, error)
	httpClient    *http.Client
	oidcIss       string
	tokenCache    *sync.Cache
	tokenCacheTTL time.Duration
}

func (m oidcAuth) getClaims(token string, req *http.Request) (claims map[string]interface{}, status int) {
	hit := m.tokenCache.Load(token)
	if hit == nil {
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

		if err := userInfo.Claims(&claims); err != nil {
			m.logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
			status = http.StatusInternalServerError
			return
		}

		expiration := m.extractExpiration(token)
		m.tokenCache.Store(token, claims, expiration)

		m.logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Time("expiration", expiration.UTC()).Msg("unmarshalled and cached userinfo")
		return
	}

	var ok bool
	if claims, ok = hit.V.(map[string]interface{}); !ok {
		status = http.StatusInternalServerError
		return
	}
	m.logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
	return
}

// extractExpiration tries to extract the expriration time from the access token
// It tries so by parsing (and verifying the signature) the access_token as JWT.
// If it is a valid JWT the `exp` claim will be used that the token expiration time.
// If it is not a valid JWT	we fallback to the configured cache TTL.
// This could still be enhanced by trying a to use the introspection endpoint (RFC7662),
// to validate the token. If it exists.
func (m oidcAuth) extractExpiration(token string) time.Time {
	defaultExpiration := time.Now().Add(m.tokenCacheTTL)
	jwks := m.getKeyfunc()
	if jwks == nil {
		return defaultExpiration
	}

	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, jwks.Keyfunc)
	if err != nil {
		m.logger.Info().Err(err).Msg("Error parsing access_token as JWT")
		return defaultExpiration
	}
	if claims.ExpiresAt != nil {
		m.logger.Debug().Str("exp", claims.ExpiresAt.String()).Msg("Expiration Time from access_token")
		return claims.ExpiresAt.Time
	}
	return defaultExpiration
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

type jwksJSON struct {
	JWKSURL string `json:"jwks_uri"`
}

func (m *oidcAuth) getKeyfunc() *keyfunc.JWKS {
	if m.jwks == nil {
		wellKnown := strings.TrimSuffix(m.oidcIss, "/") + "/.well-known/openid-configuration"
		resp, err := m.httpClient.Get(wellKnown)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			m.logger.Error().Err(err).Msg("unable to read discovery response body")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			m.logger.Error().Str("status", resp.Status).Str("body", string(body)).Msg("error requesting openid-configuration")
			return nil
		}

		var j jwksJSON
		err = json.
			Unmarshal(body, &j)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to decode provider discovered openid-configuration")
			return nil
		}
		m.logger.Debug().Str("jwks", j.JWKSURL).Msg("discovered jwks endpoint")
		// FIXME: make configurable
		options := keyfunc.Options{
			RefreshErrorHandler: func(err error) {
				m.logger.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
			},
			RefreshInterval:   time.Hour,
			RefreshRateLimit:  time.Minute * 5,
			RefreshTimeout:    time.Second * 10,
			RefreshUnknownKID: true,
		}
		m.jwks, err = keyfunc.Get(j.JWKSURL, options)
		if err != nil {
			m.jwks = nil
			m.logger.Error().Err(err).Msg("Failed to create JWKS from resource at the given URL.")
			return nil
		}
	}
	return m.jwks
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
