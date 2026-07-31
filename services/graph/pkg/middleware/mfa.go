package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
)

// RequireMFA middleware is used to require the user in context to have MFA satisfied
func RequireMFA(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !revactx.HasMFA(r.Context()) {
				l := logger.SubloggerWithRequestID(r.Context())
				l.Error().Str("path", r.URL.Path).Msg("MFA required but not satisfied")
				mfa.SetRequiredStatus(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
