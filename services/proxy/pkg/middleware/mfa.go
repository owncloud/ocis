package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

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
		}
	}
}

// MultiFactorAuthentication is a authenticator that checks for mfa on specific paths
type MultiFactorAuthentication struct {
	next           http.Handler
	logger         log.Logger
	enabled        bool
	authLevelNames []string
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

	// acr is a standard OIDC claim.
	value, err := oidc.ReadStringClaim("acr", claims)
	if err != nil {
		m.logger.Error().Str("path", req.URL.Path).Interface("required", m.authLevelNames).Err(err).Interface("claims", claims).Msg("no acr claim found in access token")
		return
	}

	if !m.containsMFA(value) {
		m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("accessing path without mfa")
		return
	}

	mfa.SetHeader(req, true)
	m.logger.Debug().Str("acr", value).Str("url", req.URL.Path).Msg("mfa authenticated")
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
