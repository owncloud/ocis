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

type OIDCAuthenticator struct {
	Logger                  log.Logger
	HTTPClient              *http.Client
	OIDCIss                 string
	TokenCache              *osync.Cache
	TokenCacheTTL           time.Duration
	ProviderFunc            func() (OIDCProvider, error)
	AccessTokenVerifyMethod string
	JWKSOptions             config.JWKS

	providerLock *sync.Mutex
	provider     OIDCProvider

	jwksLock *sync.Mutex
	JWKS     *keyfunc.JWKS
}

func (m OIDCAuthenticator) getClaims(token string, req *http.Request) (claims map[string]interface{}, status int) {
	hit := m.TokenCache.Load(token)
	if hit == nil {
		aClaims, err := m.verifyAccessToken(token)
		if err != nil {
			m.Logger.Error().Err(err).Msg("Failed to verify access token")
			status = http.StatusUnauthorized
			return
		}

		oauth2Token := &oauth2.Token{
			AccessToken: token,
		}

		userInfo, err := m.getProvider().UserInfo(
			context.WithValue(req.Context(), oauth2.HTTPClient, m.HTTPClient),
			oauth2.StaticTokenSource(oauth2Token),
		)
		if err != nil {
			m.Logger.Error().Err(err).Msg("Failed to get userinfo")
			status = http.StatusUnauthorized
			return
		}

		if err := userInfo.Claims(&claims); err != nil {
			m.Logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
			status = http.StatusInternalServerError
			return
		}

		expiration := m.extractExpiration(aClaims)
		m.TokenCache.Store(token, claims, expiration)

		m.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Time("expiration", expiration.UTC()).Msg("unmarshalled and cached userinfo")
		return
	}

	var ok bool
	if claims, ok = hit.V.(map[string]interface{}); !ok {
		status = http.StatusInternalServerError
		return
	}
	m.Logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
	return
}

func (m OIDCAuthenticator) verifyAccessToken(token string) (jwt.RegisteredClaims, error) {
	switch m.AccessTokenVerifyMethod {
	case config.AccessTokenVerificationJWT:
		return m.verifyAccessTokenJWT(token)
	case config.AccessTokenVerificationNone:
		m.Logger.Debug().Msg("Access Token verification disabled")
		return jwt.RegisteredClaims{}, nil
	default:
		m.Logger.Error().Str("access_token_verify_method", m.AccessTokenVerifyMethod).Msg("Unknown Access Token verification setting")
		return jwt.RegisteredClaims{}, errors.New("Unknown Access Token Verification method")
	}
}

// verifyAccessTokenJWT tries to parse and verify the access token as a JWT.
func (m OIDCAuthenticator) verifyAccessTokenJWT(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	jwks := m.getKeyfunc()
	if jwks == nil {
		return claims, errors.New("Error initializing jwks keyfunc")
	}

	_, err := jwt.ParseWithClaims(token, &claims, jwks.Keyfunc)
	m.Logger.Debug().Interface("access token", &claims).Msg("parsed access token")
	if err != nil {
		m.Logger.Info().Err(err).Msg("Failed to parse/verify the access token.")
		return claims, err
	}

	if !claims.VerifyIssuer(m.OIDCIss, true) {
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
func (m OIDCAuthenticator) extractExpiration(aClaims jwt.RegisteredClaims) time.Time {
	defaultExpiration := time.Now().Add(m.TokenCacheTTL)
	if aClaims.ExpiresAt != nil {
		m.Logger.Debug().Str("exp", aClaims.ExpiresAt.String()).Msg("Expiration Time from access_token")
		return aClaims.ExpiresAt.Time
	}
	return defaultExpiration
}

func (m OIDCAuthenticator) shouldServe(req *http.Request) bool {
	header := req.Header.Get("Authorization")

	if m.OIDCIss == "" {
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

func (m *OIDCAuthenticator) getKeyfunc() *keyfunc.JWKS {
	m.jwksLock.Lock()
	defer m.jwksLock.Unlock()
	if m.JWKS == nil {
		wellKnown := strings.TrimSuffix(m.OIDCIss, "/") + "/.well-known/openid-configuration"

		resp, err := m.HTTPClient.Get(wellKnown)
		if err != nil {
			m.Logger.Error().Err(err).Msg("Failed to set request for .well-known/openid-configuration")
			return nil
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			m.Logger.Error().Err(err).Msg("unable to read discovery response body")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			m.Logger.Error().Str("status", resp.Status).Str("body", string(body)).Msg("error requesting openid-configuration")
			return nil
		}

		var j jwksJSON
		err = json.Unmarshal(body, &j)
		if err != nil {
			m.Logger.Error().Err(err).Msg("failed to decode provider openid-configuration")
			return nil
		}
		m.Logger.Debug().Str("jwks", j.JWKSURL).Msg("discovered jwks endpoint")
		options := keyfunc.Options{
			Client: m.HTTPClient,
			RefreshErrorHandler: func(err error) {
				m.Logger.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
			},
			RefreshInterval:   time.Minute * time.Duration(m.JWKSOptions.RefreshInterval),
			RefreshRateLimit:  time.Second * time.Duration(m.JWKSOptions.RefreshRateLimit),
			RefreshTimeout:    time.Second * time.Duration(m.JWKSOptions.RefreshTimeout),
			RefreshUnknownKID: m.JWKSOptions.RefreshUnknownKID,
		}
		m.JWKS, err = keyfunc.Get(j.JWKSURL, options)
		if err != nil {
			m.JWKS = nil
			m.Logger.Error().Err(err).Msg("Failed to create JWKS from resource at the given URL.")
			return nil
		}
	}
	return m.JWKS
}

func (m *OIDCAuthenticator) getProvider() OIDCProvider {
	m.providerLock.Lock()
	defer m.providerLock.Unlock()
	if m.provider == nil {
		// Lazily initialize a provider

		// provider needs to be cached as when it is created
		// it will fetch the keys from the issuer using the .well-known
		// endpoint
		provider, err := m.ProviderFunc()
		if err != nil {
			m.Logger.Error().Err(err).Msg("could not initialize oidcAuth provider")
			return nil
		}

		m.provider = provider
	}
	return m.provider
}

func (m OIDCAuthenticator) Authenticate(r *http.Request) (*http.Request, bool) {
	// there is no bearer token on the request,
	if !m.shouldServe(r) {
		// // oidc supported but token not present, add header and handover to the next middleware.
		// userAgentAuthenticateLockIn(w, r, options.CredentialsByUserAgent, "bearer")
		// next.ServeHTTP(w, r)
		return nil, false
	}

	if m.getProvider() == nil {
		// w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}
	// Force init of jwks keyfunc if needed (contacts the .well-known and jwks endpoints on first call)
	if m.AccessTokenVerifyMethod == config.AccessTokenVerificationJWT && m.getKeyfunc() == nil {
		return nil, false
	}

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	claims, status := m.getClaims(token, r)
	if status != 0 {
		// w.WriteHeader(status)
		// TODO log
		return nil, false
	}
	return r.WithContext(oidc.NewContext(r.Context(), claims)), true
}
