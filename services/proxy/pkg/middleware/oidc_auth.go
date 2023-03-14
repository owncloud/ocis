package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	ocstore "github.com/owncloud/ocis/v2/ocis-pkg/store"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"

	"github.com/MicahParks/keyfunc"
	gOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/shamaton/msgpack/v2"
	store "go-micro.dev/v4/store"
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

// NewOIDCAuthenticator returns a ready to use authenticator which can handle OIDC authentication.
func NewOIDCAuthenticator(opts ...Option) *OIDCAuthenticator {
	options := newOptions(opts...)

	// if no cache is configured use noop cache
	if options.Cache == nil {
		options.Cache = ocstore.Create(ocstore.Type("noop"))
	}
	return &OIDCAuthenticator{
		Logger:                  options.Logger,
		tokenCache:              options.Cache,
		DefaultTokenCacheTTL:    options.DefaultAccessTokenTTL,
		HTTPClient:              options.HTTPClient,
		OIDCIss:                 options.OIDCIss,
		ProviderFunc:            options.OIDCProviderFunc,
		JWKSOptions:             options.JWKS,
		AccessTokenVerifyMethod: options.AccessTokenVerifyMethod,
		providerLock:            &sync.Mutex{},
		jwksLock:                &sync.Mutex{},
	}
}

// OIDCAuthenticator is an authenticator responsible for OIDC authentication.
type OIDCAuthenticator struct {
	Logger                  log.Logger
	HTTPClient              *http.Client
	OIDCIss                 string
	tokenCache              store.Store
	DefaultTokenCacheTTL    time.Duration
	ProviderFunc            func() (OIDCProvider, error)
	AccessTokenVerifyMethod string
	JWKSOptions             config.JWKS

	providerLock *sync.Mutex
	provider     OIDCProvider

	jwksLock *sync.Mutex
	JWKS     *keyfunc.JWKS
}

func (m *OIDCAuthenticator) getClaims(token string, req *http.Request) (map[string]interface{}, error) {
	var claims map[string]interface{}
	record, _ := m.tokenCache.Read(token) // TODO log error in debug?
	if len(record) < 1 {
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
		d, err := msgpack.Marshal(claims)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal claims for cache")
		}

		err = m.tokenCache.Write(&store.Record{
			Key:    token,
			Value:  d,
			Expiry: time.Until(expiration),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to write to cache") // TODO log if cache does not work, but continue
		}

		m.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Time("expiration", expiration.UTC()).Msg("unmarshalled and cached userinfo")
		return claims, nil
	}

	if err := msgpack.Unmarshal(record[0].Value, &claims); err != nil {
		return nil, errors.New("failed to unmarshal claims from the cache") // TODO log if cache does not work, but continue
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
	defaultExpiration := time.Now().Add(m.DefaultTokenCacheTTL)
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

func (m *OIDCAuthenticator) getKeyfunc() *keyfunc.JWKS {
	m.jwksLock.Lock()
	defer m.jwksLock.Unlock()
	if m.JWKS == nil {
		oidcMetadata, err := oidc.GetIDPMetadata(m.Logger, m.HTTPClient, m.OIDCIss)
		if err != nil {
			m.Logger.Error().Err(err).Msg("failed to decode provider openid-configuration")
			return nil
		}
		m.Logger.Debug().Str("jwks", oidcMetadata.JwksURI).Msg("discovered jwks endpoint")
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
		m.JWKS, err = keyfunc.Get(oidcMetadata.JwksURI, options)
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

// Authenticate implements the authenticator interface to authenticate requests via oidc auth.
func (m *OIDCAuthenticator) Authenticate(r *http.Request) (*http.Request, bool) {
	// there is no bearer token on the request,
	if !m.shouldServe(r) || isPublicPath(r.URL.Path) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
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
