package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/owncloud/reva/v2/pkg/autoprop"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

const defaultMFASessionDuration = 3600

// mfaSessionDuration returns the configured MFA session duration or the default.
func mfaSessionDuration(cfg config.MFAConfig) time.Duration {
	if cfg.SessionDuration <= 0 {
		return defaultMFASessionDuration * time.Second
	}
	return time.Duration(cfg.SessionDuration) * time.Second
}

// MultiFactor returns a middleware that checks requests for mfa.
//
// It must run after the AccountResolver, which adds the user to the context.
// The user is needed to store and to read the MFA status.
func MultiFactor(cfg config.MFAConfig, opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		return &MultiFactorAuthentication{
			next:            next,
			logger:          logger,
			enabled:         cfg.Enabled,
			authLevelNames:  cfg.AuthLevelNames,
			store:           options.MFAStore,
			sessionDuration: mfaSessionDuration(cfg),
		}
	}
}

// MultiFactorAuthentication is a authenticator that checks for mfa on specific paths
type MultiFactorAuthentication struct {
	next            http.Handler
	logger          log.Logger
	enabled         bool
	authLevelNames  []string
	sessionDuration time.Duration
	// store holds the time of the last request with MFA of a user, so that requests
	// without claims (e.g. signed-URL downloads) can inherit the MFA status. Nil when
	// no store is configured.
	store microstore.Store
}

// ServeHTTP adds the mfa header if the request contains a valid mfa token
func (m MultiFactorAuthentication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Incoming requests must not have autopropagation headers
	for key, _ := range req.Header {
		if strings.HasPrefix(key, autoprop.HTTPAutoPropPrefix) || strings.HasPrefix(key, autoprop.MicroAutoPropPrefix) {
			req.Header.Del(key)
		}
	}

	if !m.enabled {
		// if mfa is disabled we always set the header to true.
		// this allows all other services to assume mfa is always active.
		// this should reduce code and configuration complexity in other services.
		req = req.WithContext(revactx.SetMFA(req.Context()))
		m.next.ServeHTTP(w, req)
		return
	}

	// overwrite the mfa header to avoid passing on wrong information
	ctx := revactx.RemoveMFA(req.Context())

	claims := oidc.FromContext(ctx)

	if claims == nil {
		// No OIDC claims, e.g. the request has a signed URL. MFA cannot be determined
		// from claims directly.
		//
		// Fall back to the stored MFA status from the most recent request with MFA
		// of the user. This allows, for example, a signed-URL archiver download to
		// succeed when the user has recently proven MFA in their browser session.
		if m.readMFAFromStore(ctx) {
			ctx = revactx.SetMFA(ctx)
		}

		m.logger.Debug().Str("path", req.URL.Path).Bool("mfaStatus", revactx.HasMFA(ctx)).Msg("no OIDC claims in context")
		m.next.ServeHTTP(w, req.WithContext(ctx)) // ensure the request has the right context
		return
	}

	// acr is a standard OIDC claim.
	value, err := oidc.ReadStringClaim("acr", claims)
	if err != nil {
		m.logger.Debug().Str("path", req.URL.Path).Interface("required", m.authLevelNames).Err(err).Msg("acr claim not set in access token")
	} else if !m.containsMFA(value) {
		m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("accessing path without mfa")
	} else {
		m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("mfa authenticated")
		ctx = revactx.SetMFA(ctx)
		m.writeMFAToStore(ctx)
	}

	// MFA status will only be true if the acr claim contains the proper value,
	// otherwise it wll be false (removed from the context early)
	m.next.ServeHTTP(w, req.WithContext(ctx))
}

// readMFAFromStore checks if the user had MFA in an OIDC request within the session duration.
// The age is checked here because stores like nats-js-kv ignore the expiry of a record.
func (m MultiFactorAuthentication) readMFAFromStore(ctx context.Context) bool {
	userID := mfaUserID(ctx)
	if m.store == nil || userID == "" {
		return false
	}
	records, err := m.store.Read(key(userID))
	if err != nil || len(records) == 0 {
		return false
	}
	verifiedAt, err := strconv.ParseInt(string(records[0].Value), 10, 64)
	if err != nil {
		return false
	}
	return time.Since(time.Unix(verifiedAt, 0)) < m.sessionDuration
}

// writeMFAToStore stores the time of the last request with MFA of the user, so that requests
// without claims (e.g. signed-URL downloads) can inherit the MFA status, see readMFAFromStore.
func (m MultiFactorAuthentication) writeMFAToStore(ctx context.Context) {
	userID := mfaUserID(ctx)
	if m.store == nil || userID == "" {
		return
	}
	if err := m.store.Write(&microstore.Record{
		Key:   key(userID),
		Value: []byte(strconv.FormatInt(time.Now().Unix(), 10)),
		// redis honors the expiry, nats-js-kv uses the TTL of the bucket instead
		Expiry: m.sessionDuration,
	}); err != nil {
		m.logger.Error().Err(err).Str("userID", userID).Msg("failed to write MFA status to store")
	}
}

// containsMFA checks if the given value is in the list of authentication level names
func (m MultiFactorAuthentication) containsMFA(value string) bool {
	for _, v := range m.authLevelNames {
		if v == value {
			return true
		}
	}
	return false
}

// mfaUserID returns the id of the user in the context, or "" if there is no user.
func mfaUserID(ctx context.Context) string {
	u, _ := revactx.ContextGetUser(ctx)
	return u.GetId().GetOpaqueId()
}

func key(userID string) string {
	return "mfa:" + userID
}
