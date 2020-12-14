package middleware

import (
	"net/http"
	"path"
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

	// TODO: investigate broken caching - https://github.com/owncloud/ocis/issues/1094
	// we don't have a last modification date of the static assets, so we use the service start date
	//lastModified := time.Now().UTC().Format(http.TimeFormat)
	//expires := time.Now().Add(time.Second * time.Duration(ttl)).UTC().Format(http.TimeFormat)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, path.Join(root, "api")) {
				next.ServeHTTP(w, r)
			} else {
				// TODO: investigate broken caching - https://github.com/owncloud/ocis/issues/1094
				//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, must-revalidate", strconv.Itoa(ttl)))
				//w.Header().Set("Expires", expires)
				//w.Header().Set("Last-Modified", lastModified)
				w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
				w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
				w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				static.ServeHTTP(w, r)
			}
		})
	}
}
