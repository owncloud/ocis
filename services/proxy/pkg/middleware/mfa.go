package middleware

import (
	"net/http"
	"time"

	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// mfaStoreTTL is how long a verified MFA status is remembered for non-OIDC
// requests (e.g. signed-URL archiver downloads). It should be at least as
// long as the signed-URL expiry (OC-Expires). Default: 1 hour.
const mfaStoreTTL = time.Hour

// MultiFactor returns a middleware that checks requests for mfa
func MultiFactor(cfg config.MFAConfig, opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		return &MultiFactorAuthentication{
			next:           next,
			logger:         logger,
			enabled:        cfg.Enabled,
			authLevelNames: cfg.AuthLevelNames,
			store:          options.MFAStore,
		}
	}
}

// MultiFactorAuthentication is a authenticator that checks for mfa on specific paths
type MultiFactorAuthentication struct {
	next           http.Handler
	logger         log.Logger
	enabled        bool
	authLevelNames []string
	// store persists verified MFA status so that non-OIDC requests (e.g.
	// signed-URL archiver downloads) can inherit it from the user's most
	// recent OIDC session. Nil when no store is configured.
	store microstore.Store
}

// ServeHTTP adds the mfa header if the request contains a valid mfa token
func (m MultiFactorAuthentication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer m.next.ServeHTTP(w, req)

	if !m.enabled {
		// if mfa is disabled we always set the header to true.
		// this allows all other services to assume mfa is always active.
		// this should reduce code and configuration complexity in other services.
		mfa.SetHeader(req, true)
		return
	}

	// overwrite the mfa header to avoid passing on wrong information
	mfa.SetHeader(req, false)

	claims := oidc.FromContext(req.Context())
	// TODO:
	// Should we set the "acr" to the claims "X-Access-Token" for the "aud":["reva"]?
	// Should we get the claims from the "X-Access-Token"?

	if claims == nil {
		// No OIDC claims — request was authenticated via a non-OIDC method
		// (e.g. signed URL, basic auth, app token). MFA cannot be determined
		// from claims directly.
		//
		// Fall back to the persisted MFA status from the user's most recent
		// OIDC-authenticated session. This allows, for example, a signed-URL
		// archiver download to succeed when the user has recently proven MFA
		// in their browser session.
		if m.store != nil {
			if u, ok := revactx.ContextGetUser(req.Context()); ok && u.GetId().GetOpaqueId() != "" {
				if m.readMFAFromStore(u.GetId().GetOpaqueId()) {
					mfa.SetHeader(req, true)
					m.logger.Debug().Str("path", req.URL.Path).Msg("MFA status restored from store for non-OIDC request")
					return
				}
			}
		}
		m.logger.Debug().Str("path", req.URL.Path).Msg("no OIDC claims in context, skipping MFA check")
		return
	}

	// acr is a standard OIDC claim.
	value, err := oidc.ReadStringClaim("acr", claims)
	if err != nil {
		m.logger.Debug().Str("path", req.URL.Path).Interface("required", m.authLevelNames).Err(err).Msg("acr claim not set in access token")
		return
	}

	if !m.containsMFA(value) {
		m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("accessing path without mfa")
		return
	}

	mfa.SetHeader(req, true)
	m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("mfa authenticated")

	// Persist the verified MFA status so that subsequent non-OIDC requests
	// (e.g. signed-URL archiver downloads) can inherit it. The entry is
	// refreshed on every successful OIDC MFA verification and expires after
	// mfaStoreTTL if no further OIDC requests are made.
	if m.store != nil {
		if u, ok := revactx.ContextGetUser(req.Context()); ok && u.GetId().GetOpaqueId() != "" {
			m.writeMFAToStore(u.GetId().GetOpaqueId())
		}
	}
}

func (m MultiFactorAuthentication) readMFAFromStore(userID string) bool {
	records, err := m.store.Read("mfa:" + userID)
	if err != nil || len(records) == 0 {
		return false
	}
	return string(records[0].Value) == "true"
}

func (m MultiFactorAuthentication) writeMFAToStore(userID string) {
	if err := m.store.Write(&microstore.Record{
		Key:    "mfa:" + userID,
		Value:  []byte("true"),
		Expiry: mfaStoreTTL,
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
