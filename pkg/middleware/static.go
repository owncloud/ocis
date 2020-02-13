package middleware

import (
	"net/http"
	"strings"
)

// Static is a middleware that serves static assets.
func Static(root string, fs http.FileSystem) func(http.Handler) http.Handler {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	static := http.StripPrefix(
		root,
		http.FileServer(
			fs,
		),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// serve the static assets for the identifier web app
			if strings.HasPrefix(r.URL.Path, "/signin/v1/static/") {
				if strings.HasSuffix(r.URL.Path, "/") {
					// but no listing of folders
					http.NotFound(w, r)
				} else {
					r.URL.Path = strings.Replace(r.URL.Path, "/signin/v1/static/", "/signin/v1/identifier/static/", 1)
					static.ServeHTTP(w, r)
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
