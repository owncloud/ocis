package middleware

import (
	"net"
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

			if parsedUri, err := url.Parse(uriString); err == nil {
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

// urisMatch will check if both URL match. In case of localhost, only the port
// will be ignored. Note that we need to deal with URLs like
// "[scheme:][//[userinfo@]host][/]path[?query][#fragment]"
// and also like "scheme:opaque[?query][#fragment]"
func urisMatch(u1, u2 *url.URL) bool {
	if isLocalhost(u1) && isLocalhost(u2) {
		copyu1 := *u1
		copyu2 := *u2
		// if there is a port, remove it from the URL
		if copyu1.Port() != "" {
			host1, _, _ := net.SplitHostPort(copyu1.Host)
			copyu1.Host = host1
		}
		if copyu2.Port() != "" {
			host2, _, _ := net.SplitHostPort(copyu2.Host)
			copyu2.Host = host2
		}
		return copyu1.String() == copyu2.String()
	}
	return u1.String() == u2.String()
}

func isLocalhost(u1 *url.URL) bool {
	hostname := u1.Hostname()
	return hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1"
}

// isValidRedirect check if the u1 URL is part of the configured redirect URIs.
// This will use the urisMatch for matching.
func isValidRedirect(u1 *url.URL, cfg *config.Config) bool {
	for _, client := range cfg.Clients {
		for _, redirect_uri := range client.RedirectURIs {
			if parsedRedirect, err := url.Parse(redirect_uri); err == nil {
				if urisMatch(u1, parsedRedirect) {
					return true
				}
			}
		}
	}
	return false
}
