package middleware

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

// Static is a middleware that serves static assets.
func Static(root string, fs http.FileSystem, ttl int) func(http.Handler) http.Handler {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	static := http.StripPrefix(
		root,
		http.FileServer(
			fs,
		),
	)

	// we don't have a last modification date of the static assets, so we use the service start date
	lastModified := time.Now().UTC().Format(http.TimeFormat)
	expires := time.Now().Add(time.Second * time.Duration(ttl)).UTC().Format(http.TimeFormat)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, path.Join(root, "api")) {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", strconv.Itoa(ttl)))
				w.Header().Set("Expires", expires)
				w.Header().Set("Last-Modified", lastModified)
				static.ServeHTTP(w, r)
			}
		})
	}
}
