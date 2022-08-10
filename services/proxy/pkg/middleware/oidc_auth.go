package middleware

import (
	"context"
	"encoding/json"
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
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	_headerAuthorization = "Authorization"
	_bearerPrefix        = "Bearer "
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

func NewOIDCAuthenticator(logger log.Logger, tokenCacheTTL int, oidcHTTPClient *http.Client, oidcIss string, providerFunc func() (OIDCProvider, error),
	jwksOptions config.JWKS, accessTokenVerifyMethod string) OIDCAuthenticator {
	tokenCache := osync.NewCache(tokenCacheTTL)
	return OIDCAuthenticator{
		Logger:                  logger,
		tokenCache:              &tokenCache,
		TokenCacheTTL:           time.Duration(tokenCacheTTL),
		HTTPClient:              oidcHTTPClient,
		OIDCIss:                 oidcIss,
		ProviderFunc:            providerFunc,
		JWKSOptions:             jwksOptions,
		AccessTokenVerifyMethod: accessTokenVerifyMethod,
		providerLock:            &sync.Mutex{},
		jwksLock:                &sync.Mutex{},
	}
}

type OIDCAuthenticator struct {
	Logger                  log.Logger
	HTTPClient              *http.Client
	OIDCIss                 string
	tokenCache              *osync.Cache
	TokenCacheTTL           time.Duration
	ProviderFunc            func() (OIDCProvider, error)
	AccessTokenVerifyMethod string
	JWKSOptions             config.JWKS

	providerLock *sync.Mutex
	provider     OIDCProvider

	jwksLock *sync.Mutex
	JWKS     *keyfunc.JWKS
}

func (m OIDCAuthenticator) getClaims(token string, req *http.Request) (map[string]interface{}, error) {
	var claims map[string]interface{}
	hit := m.tokenCache.Load(token)
	if hit == nil {
		aClaims, err := m.verifyAccessToken(token)
		if err != nil {
			return nil, errors.Wrap(err, "failed to verify access token")
		}

		oauth2Token := &oauth2.Token{
			AccessToken: token,
		}

		userInfo, err := m.getProvider().UserInfo(
			context.WithValue(req.Context(), oauth2.HTTPClient, m.HTTPClient),
			oauth2.StaticTokenSource(oauth2Token),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get userinfo")
		}
		if err := userInfo.Claims(&claims); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal userinfo claims")
		}

		expiration := m.extractExpiration(aClaims)
		m.tokenCache.Store(token, claims, expiration)

		m.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Time("expiration", expiration.UTC()).Msg("unmarshalled and cached userinfo")
		return claims, nil
	}

	var ok bool
	if claims, ok = hit.V.(map[string]interface{}); !ok {
		return nil, errors.New("failed to cast claims from the cache")
	}
	m.Logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
	return claims, nil
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
	if m.OIDCIss == "" {
		return false
	}

	header := req.Header.Get(_headerAuthorization)
	return strings.HasPrefix(header, _bearerPrefix)
}

type jwksJSON struct {
	JWKSURL string `json:"jwks_uri"`
}

func (m OIDCAuthenticator) getKeyfunc() *keyfunc.JWKS {
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

func (m OIDCAuthenticator) getProvider() OIDCProvider {
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
		return nil, false
	}

	if m.getProvider() == nil {
		return nil, false
	}

	// Force init of jwks keyfunc if needed (contacts the .well-known and jwks endpoints on first call)
	if m.AccessTokenVerifyMethod == config.AccessTokenVerificationJWT && m.getKeyfunc() == nil {
		return nil, false
	}
	token := strings.TrimPrefix(r.Header.Get(_headerAuthorization), _bearerPrefix)

	claims, err := m.getClaims(token, r)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "oidc").
			Str("path", r.URL.Path).
			Msg("failed to authenticate the request")
		return nil, false
	}
	m.Logger.Debug().
		Str("authenticator", "oidc").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	return r.WithContext(oidc.NewContext(r.Context(), claims)), true
}
