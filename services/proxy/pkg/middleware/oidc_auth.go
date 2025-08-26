package middleware

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc/checkers"
	"github.com/pkg/errors"
	"github.com/shamaton/msgpack/v2"
	store "go-micro.dev/v4/store"
	"golang.org/x/crypto/sha3"
	"golang.org/x/oauth2"
)

const (
	_headerAuthorization = "Authorization"
	_bearerPrefix        = "Bearer "
)

// NewOIDCAuthenticator returns a ready to use authenticator which can handle OIDC authentication.
func NewOIDCAuthenticator(opts ...Option) *OIDCAuthenticator {
	options := newOptions(opts...)

	return &OIDCAuthenticator{
		Logger:                  options.Logger,
		userInfoCache:           options.UserInfoCache,
		DefaultTokenCacheTTL:    options.DefaultAccessTokenTTL,
		HTTPClient:              options.HTTPClient,
		OIDCIss:                 options.OIDCIss,
		oidcClient:              options.OIDCClient,
		claimsChecker:           options.ClaimsChecker,
		AccessTokenVerifyMethod: options.AccessTokenVerifyMethod,
		skipUserInfo:            options.SkipUserInfo,
		TimeFunc:                time.Now,
	}
}

// OIDCAuthenticator is an authenticator responsible for OIDC authentication.
type OIDCAuthenticator struct {
	Logger                  log.Logger
	HTTPClient              *http.Client
	OIDCIss                 string
	userInfoCache           store.Store
	DefaultTokenCacheTTL    time.Duration
	oidcClient              oidc.OIDCClient
	claimsChecker           checkers.Checker
	AccessTokenVerifyMethod string
	skipUserInfo            bool
	TimeFunc                func() time.Time
}

func (m *OIDCAuthenticator) getClaims(token string, req *http.Request) (map[string]interface{}, bool, error) {
	var claims map[string]interface{}

	// use a 64 bytes long hash to have 256-bit collision resistance.
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, []byte(token))
	encodedHash := base64.URLEncoding.EncodeToString(hash)

	record, err := m.userInfoCache.Read(encodedHash)
	if err != nil && err != store.ErrNotFound {
		m.Logger.Error().Err(err).Msg("could not read from userinfo cache")
	}
	if len(record) > 0 {
		if err = msgpack.UnmarshalAsMap(record[0].Value, &claims); err == nil {
			m.Logger.Debug().Interface("claims", claims).Msg("cache hit for userinfo")
			if ok := verifyExpiresAt(claims, m.TimeFunc()); !ok {
				return nil, false, jwt.ErrTokenExpired
			}
			return claims, false, nil
		}
		m.Logger.Error().Err(err).Msg("could not unmarshal userinfo")
	}

	aClaims, claims, err := m.oidcClient.VerifyAccessToken(req.Context(), token)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to verify access token")
	}

	if !m.skipUserInfo {
		oauth2Token := &oauth2.Token{
			AccessToken: token,
		}

		userInfo, err := m.oidcClient.UserInfo(
			context.WithValue(req.Context(), oauth2.HTTPClient, m.HTTPClient),
			oauth2.StaticTokenSource(oauth2Token),
		)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to get userinfo")
		}
		if err := userInfo.Claims(&claims); err != nil {
			return nil, false, errors.Wrap(err, "failed to unmarshal userinfo claims")
		}
	}

	expiration := m.extractExpiration(aClaims)
	// always set an exp claim
	claims["exp"] = expiration.Unix()
	go func() {
		if d, err := msgpack.MarshalAsMap(claims); err != nil {
			m.Logger.Error().Err(err).Msg("failed to marshal claims for userinfo cache")
		} else {
			err = m.userInfoCache.Write(&store.Record{
				Key:    encodedHash,
				Value:  d,
				Expiry: time.Until(expiration),
			})
			if err != nil {
				m.Logger.Error().Err(err).Msg("failed to write to userinfo cache")
			}

			if sid := aClaims.SessionID; sid != "" {
				// reuse user cache for session id lookup
				err = m.userInfoCache.Write(&store.Record{
					Key:    sid,
					Value:  []byte(encodedHash),
					Expiry: time.Until(expiration),
				})
				if err != nil {
					m.Logger.Error().Err(err).Msg("failed to write session lookup cache")
				}
			}
		}
	}()

	// If we get here this was a new login (or a renewal of the token)
	// add a flag about that to the claims, to be able to distinguish
	// it in the accountresolver middleware

	m.Logger.Debug().Interface("claims", claims).Msg("extracted claims")
	return claims, true, nil
}

// extractExpiration tries to extract the expriration time from the access token
// If the access token does not have an exp claim it will fallback to the configured
// default expiration
func (m OIDCAuthenticator) extractExpiration(aClaims oidc.RegClaimsWithSID) time.Time {
	defaultExpiration := time.Now().Add(m.DefaultTokenCacheTTL)
	if aClaims.ExpiresAt != nil {
		m.Logger.Debug().Str("exp", aClaims.ExpiresAt.String()).Msg("Expiration Time from access_token")
		return aClaims.ExpiresAt.Time
	}
	return defaultExpiration
}

func verifyExpiresAt(claims map[string]interface{}, cmp time.Time) bool {
	var expiry time.Time
	switch v := claims["exp"].(type) {
	case nil:
		return false
	case int64:
		expiry = time.Unix(v, 0)
	case uint32:
		expiry = time.Unix(int64(v), 0)
	}
	return cmp.Before(expiry)
}

func (m OIDCAuthenticator) shouldServe(req *http.Request) bool {
	if m.OIDCIss == "" {
		return false
	}

	header := req.Header.Get(_headerAuthorization)
	return strings.HasPrefix(header, _bearerPrefix)
}

// shouldCheckClaims returns true if we should check the claims for the
// provided request.
func (m *OIDCAuthenticator) shouldCheckClaims(r *http.Request) bool {
	// the list is currently hardcoded
	protectedPaths := []string{
		"/graph/v1.0/users",
		"/graph/v1.0/groups",
		"/graph/v1beta1/drives",
	}

	for _, path := range protectedPaths {
		if r.URL.Path == path {
			q := r.URL.Query()
			// we need to check claims if the $search query is NOT present (or empty)
			if q.Get("$search") == "" { // if $query isn't present, it will return the empty string
				return true
			}
		}
	}
	return false
}

// Authenticate implements the authenticator interface to authenticate requests via oidc auth.
func (m *OIDCAuthenticator) Authenticate(r *http.Request) (*http.Request, map[string]string, bool) {
	// there is no bearer token on the request,
	if !m.shouldServe(r) {
		// The authentication of public path requests is handled by another authenticator.
		// Since we can't guarantee the order of execution of the authenticators, we better
		// implement an early return here for paths we can't authenticate in this authenticator.
		return nil, nil, false
	}
	token := strings.TrimPrefix(r.Header.Get(_headerAuthorization), _bearerPrefix)
	if token == "" {
		return nil, nil, false
	}

	claims, newSession, err := m.getClaims(token, r)
	if m.shouldCheckClaims(r) {
		if err := m.claimsChecker.CheckClaims(claims); err != nil {
			m.Logger.Error().
				Err(err).
				Str("path", r.URL.Path).
				Msg("can't access protected path without valid claims")
			return nil, m.claimsChecker.RequireMap(), false
		}
	}

	if err != nil {
		host, port, _ := net.SplitHostPort(r.RemoteAddr)
		m.Logger.Error().
			Err(err).
			Str("authenticator", "oidc").
			Str("path", r.URL.Path).
			Str("user_agent", r.UserAgent()).
			Str("client.address", r.Header.Get("X-Forwarded-For")).
			Str("network.peer.address", host).
			Str("network.peer.port", port).
			Msg("failed to authenticate the request")
		return nil, nil, false
	}
	m.Logger.Debug().
		Str("authenticator", "oidc").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")

	ctx := r.Context()
	if newSession {
		ctx = oidc.NewContextSessionFlag(ctx, true)
	}

	return r.WithContext(oidc.NewContext(ctx, claims)), nil, true
}
