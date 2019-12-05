package middleware

import (
	"net/http"
	"strings"
)

// Static is a middleware that serves static assets.
func Static(fs http.FileSystem) func(http.Handler) http.Handler {
	static := http.StripPrefix(
		"/",
		http.FileServer(
			fs,
		),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				if strings.HasSuffix(r.URL.Path, "/") {
					http.NotFound(w, r)
				} else {
					static.ServeHTTP(w, r)
				}
			}
		})
	}
}
