package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	gOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	osync "github.com/owncloud/ocis/v2/ocis-pkg/sync"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"golang.org/x/oauth2"
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

// OIDCAuth provides a middleware to check access secured by a static token.
func OIDCAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	tokenCache := osync.NewCache(options.UserinfoCacheSize)

	h := oidcAuth{
		logger:                  options.Logger,
		providerFunc:            options.OIDCProviderFunc,
		httpClient:              options.HTTPClient,
		oidcIss:                 options.OIDCIss,
		tokenCache:              &tokenCache,
		tokenCacheTTL:           options.UserinfoCacheTTL,
		accessTokenVerifyMethod: options.AccessTokenVerifyMethod,
		jwksOptions:             options.JWKS,
		jwksLock:                &sync.Mutex{},
		providerLock:            &sync.Mutex{},
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

			// Force init of jwks keyfunc if needed (contacts the .well-known and jwks endpoints on first call)
			if h.accessTokenVerifyMethod == config.AccessTokenVerificationJWT && h.getKeyfunc() == nil {
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
	logger                  log.Logger
	provider                OIDCProvider
	providerLock            *sync.Mutex
	jwksOptions             config.JWKS
	jwks                    *keyfunc.JWKS
	jwksLock                *sync.Mutex
	providerFunc            func() (OIDCProvider, error)
	httpClient              *http.Client
	oidcIss                 string
	tokenCache              *osync.Cache
	tokenCacheTTL           time.Duration
	accessTokenVerifyMethod string
}

func (m oidcAuth) getClaims(token string, req *http.Request) (claims map[string]interface{}, status int) {
	hit := m.tokenCache.Load(token)
	if hit == nil {
		aClaims, err := m.verifyAccessToken(token)
		if err != nil {
			m.logger.Error().Err(err).Msg("Failed to verify access token")
			status = http.StatusUnauthorized
			return
		}

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

		expiration := m.extractExpiration(aClaims)
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

func (m oidcAuth) verifyAccessToken(token string) (jwt.RegisteredClaims, error) {
	switch m.accessTokenVerifyMethod {
	case config.AccessTokenVerificationJWT:
		return m.verifyAccessTokenJWT(token)
	case config.AccessTokenVerificationNone:
		m.logger.Debug().Msg("Access Token verification disabled")
		return jwt.RegisteredClaims{}, nil
	default:
		m.logger.Error().Str("access_token_verify_method", m.accessTokenVerifyMethod).Msg("Unknown Access Token verification setting")
		return jwt.RegisteredClaims{}, errors.New("Unknown Access Token Verification method")
	}
}

// verifyAccessTokenJWT tries to parse and verify the access token as a JWT.
func (m oidcAuth) verifyAccessTokenJWT(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	jwks := m.getKeyfunc()
	if jwks == nil {
		return claims, errors.New("Error initializing jwks keyfunc")
	}

	_, err := jwt.ParseWithClaims(token, &claims, jwks.Keyfunc)
	m.logger.Debug().Interface("access token", &claims).Msg("parsed access token")
	if err != nil {
		m.logger.Info().Err(err).Msg("Failed to parse/verify the access token.")
		return claims, err
	}

	if !claims.VerifyIssuer(m.oidcIss, true) {
		vErr := jwt.ValidationError{}
		vErr.Inner = jwt.ErrTokenInvalidIssuer
		vErr.Errors |= jwt.ValidationErrorIssuer
		return claims, vErr
	}

	return claims, nil
}

// extractExpiration tries to extract the expriration time from the access token
// If the access token does not have an exp claim it will fallback to the configured
// default expiration
func (m oidcAuth) extractExpiration(aClaims jwt.RegisteredClaims) time.Time {
	defaultExpiration := time.Now().Add(m.tokenCacheTTL)
	if aClaims.ExpiresAt != nil {
		m.logger.Debug().Str("exp", aClaims.ExpiresAt.String()).Msg("Expiration Time from access_token")
		return aClaims.ExpiresAt.Time
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
	m.jwksLock.Lock()
	defer m.jwksLock.Unlock()
	if m.jwks == nil {
		wellKnown := strings.TrimSuffix(m.oidcIss, "/") + "/.well-known/openid-configuration"

		resp, err := m.httpClient.Get(wellKnown)
		if err != nil {
			m.logger.Error().Err(err).Msg("Failed to set request for .well-known/openid-configuration")
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
		err = json.Unmarshal(body, &j)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to decode provider openid-configuration")
			return nil
		}
		m.logger.Debug().Str("jwks", j.JWKSURL).Msg("discovered jwks endpoint")
		options := keyfunc.Options{
			Client: m.httpClient,
			RefreshErrorHandler: func(err error) {
				m.logger.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
			},
			RefreshInterval:   time.Minute * time.Duration(m.jwksOptions.RefreshInterval),
			RefreshRateLimit:  time.Second * time.Duration(m.jwksOptions.RefreshRateLimit),
			RefreshTimeout:    time.Second * time.Duration(m.jwksOptions.RefreshTimeout),
			RefreshUnknownKID: m.jwksOptions.RefreshUnknownKID,
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
	m.providerLock.Lock()
	defer m.providerLock.Unlock()
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
