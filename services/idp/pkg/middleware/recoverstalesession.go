package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

const (
	// logonCookieName is the name lico uses for the encrypted logon session
	// cookie (see the vendored identifier bootstrap). It is cleared here to
	// recover from a stale session pointing at a deleted user.
	logonCookieName = "__Secure-KKT"
	// logonCookiePath mirrors the path lico sets the logon cookie with
	// (pathPrefix + "/identifier/_/"); it must match for the browser to drop it.
	logonCookiePath = "/signin/v1/identifier/_/"
	// authorizeSuffix is the path suffix of the OIDC authorization endpoint as
	// mounted by lico.
	authorizeSuffix = "/identifier/_/authorize"
)

// RecoverStaleSession breaks the infinite authorization redirect loop that
// happens when the logon session cookie points at a user that no longer exists
// (OCISDEV-128). When the current user has been deleted, lico still accepts its
// encrypted logon cookie but cannot resolve the user, and the sign-in SPA ends
// up re-issuing the authorization request with only an OAuth2 `error` parameter.
// lico answers that with a relative, self-referential redirect back to the same
// endpoint, so the browser follows `?error=...` forever and the user is stuck
// until the browser data is cleared manually.
//
// An `error` parameter is an OAuth2 response field; a legitimate inbound
// authorization request never carries one. So its presence on the authorize
// endpoint is an unambiguous signal that the flow is stuck. In that case we
// expire the stale logon cookie (the server-side equivalent of clearing the
// browser data) and send the browser back to the issuer to start a clean login,
// where the now-missing session makes the identity provider show the sign-in
// form again and a different user can log in.
func RecoverStaleSession(cfg *config.Config, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasSuffix(r.URL.Path, authorizeSuffix) || r.URL.Query().Get("error") == "" {
				next.ServeHTTP(w, r)
				return
			}

			logger.Warn().
				Str("error", r.URL.Query().Get("error")).
				Msg("recovering from stale logon session: clearing logon cookie and restarting login")

			// Expire the logon cookie using the same name/path/attributes lico
			// set it with, so the browser drops the stale (deleted-user) session.
			http.SetCookie(w, &http.Cookie{
				Name:     logonCookieName,
				Value:    "",
				Path:     logonCookiePath,
				Secure:   true,
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
			})

			http.Redirect(w, r, cfg.IDP.Iss, http.StatusFound)
		})
	}
}
