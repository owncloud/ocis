package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// HTTPSRedirect creates middleware that redirects insecure requests to HTTPS using a trusted base URL.
func HTTPSRedirect(trustedBaseURL string) func(http.Handler) http.Handler {
	var trustedHost string
	if trustedBaseURL != "" {
		if u, err := url.Parse(trustedBaseURL); err == nil {
			trustedHost = u.Host
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			proto := req.Header.Get("x-forwarded-proto")
			if proto == "http" || proto == "HTTP" {
				if strings.TrimSpace(trustedHost) != "" {
					target := &url.URL{
						Scheme:   "https",
						Host:     trustedHost,
						Path:     req.URL.Path,
						RawQuery: req.URL.RawQuery,
					}
					http.Redirect(res, req, target.String(), http.StatusPermanentRedirect)
					return
				}
				// No trusted host configured; do not perform a redirect
			}

			next.ServeHTTP(res, req)
		})
	}
}
