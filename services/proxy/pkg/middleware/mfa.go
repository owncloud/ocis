package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc/checkers"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

var (
	// ResponseHeaderBase is the prefix for all auth related response headers.
	ResponseHeaderBase = "X-OCIS-AUTH-"

	// The list of paths that require mfa if no $search query is present
	// we use a map here for easier lookups
	_protectedPaths = map[string]struct{}{
		"/graph/v1.0/users":     struct{}{},
		"/graph/v1.0/groups":    struct{}{},
		"/graph/v1beta1/drives": struct{}{},
	}
)

// MultiFactor returns a middleware that checks requests for mfa
func MultiFactor(cfg config.MFAConfig, opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		return &MultiFactorAuthentication{
			next:          next,
			logger:        logger,
			enabled:       cfg.Enabled,
			authLevelName: cfg.AuthLevelName,
			claimsChecker: checkers.NewAcrChecker(cfg.AuthLevelName), // ?
		}
	}
}

// MultiFactorAuthentication is a authenticator that checks for mfa on specific paths
type MultiFactorAuthentication struct {
	next          http.Handler
	logger        log.Logger
	enabled       bool
	authLevelName string
	claimsChecker checkers.Checker // ?
}

// Authenticate implenents the authenticator interface and checks the access token for the correct acr claim
func (mfa MultiFactorAuthentication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !mfa.shouldCheckClaims(req) {
		mfa.next.ServeHTTP(w, req)
		return
	}

	log := mfa.logger.Error().Str("path", req.URL.Path).Str("required", mfa.authLevelName)
	claims := oidc.FromContext(req.Context())

	// either we use the claims checker here:
	if false {

		err := mfa.claimsChecker.CheckClaims(claims)
		if err == nil {
			// acr claim is correct
			mfa.next.ServeHTTP(w, req)
			return
		}

		log.Err(err).Interface("checker", mfa.claimsChecker.RequireMap()).Msg("can't access protected path without valid claims")
	}

	// or we read the acr claim directly here:
	if true {

		// acr is a standard OIDC claim.
		value, err := oidc.ReadStringClaim("acr", claims)
		if err != nil {
			log.Err(err).Interface("claims", claims).Msg("no acr claim found in access token")
			w.Header().Add(ResponseHeaderBase+"Requires-Claim", "acr")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if value == mfa.authLevelName {
			// acr claim is corrct
			mfa.next.ServeHTTP(w, req)
			return
		}

		log.Err(err).Str("acr", value).Msg("can't access protected path without valid claims")
	}

	w.Header().Add(ResponseHeaderBase+"Requires-AuthLevel", mfa.authLevelName)
	w.WriteHeader(http.StatusUnauthorized)
	return
}

// shouldCheckClaims returns true if we should check the claims for the provided request.
func (mfa MultiFactorAuthentication) shouldCheckClaims(r *http.Request) bool {
	if !mfa.enabled {
		return false
	}

	if _, protected := _protectedPaths[r.URL.Path]; !protected {
		return false
	}

	q := r.URL.Query()

	// We need to be careful here. We don't want to block access if this is a search query as this can be done without mfa.
	// But we don't want to allow bypassing mfa by just adding an empty (or ignored) $search parameter.
	// We should only check for the presence of the $search parameter if the endpoint is actually using it.
	if q.Get("$search") != "" { // if $query isn't present, it will return the empty string
		return false
	}

	return true
}
