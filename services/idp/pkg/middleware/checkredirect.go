package middleware

import (
	"net/http"
	"net/url"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

func CheckRedirect(cfg *config.Config, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uriString := r.URL.Query().Get("redirect_uri")
			if uriString == "" {
				next.ServeHTTP(w, r)
				return
			}

			if parsedUri, ok := url.Parse(uriString); ok == nil {
				if isValidRedirect(parsedUri, cfg) {
					next.ServeHTTP(w, r)
					return
				}
			}

			logger.Warn().Str("redirect_uri", uriString).Msg("Unknown redirect uri")
			http.Error(w, "Unknown URI", http.StatusInternalServerError)
		})
	}
}

// urisMatch will check if both URL match. They match if they have the same
// scheme, hostname and opaque; other elements (in particular the port)
// are ignored. Note that we need to deal with URLs like
// "[scheme:][//[userinfo@]host][/]path[?query][#fragment]"
// and also like "scheme:opaque[?query][#fragment]"
func urisMatch(u1, u2 *url.URL) bool {
	return u1.Scheme == u2.Scheme && u1.Hostname() == u2.Hostname() && u1.Opaque == u2.Opaque
}

// isValidRedirect check if the u1 URL is part of the configured redirect URIs.
// This will use the urisMatch for matching.
func isValidRedirect(u1 *url.URL, cfg *config.Config) bool {
	for _, client := range cfg.Clients {
		for _, redirect_uri := range client.RedirectURIs {
			if parsedRedirect, ok2 := url.Parse(redirect_uri); ok2 == nil {
				if urisMatch(u1, parsedRedirect) {
					return true
				}
			}
		}
	}
	return false
}
