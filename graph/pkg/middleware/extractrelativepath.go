package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
)

var relativePathRegex = regexp.MustCompile(":/[^:]+:?")

// ExtractRelativePath provides a middleware that adds a relativepath chi parameter to the context
// for graph drives urls
func ExtractRelativePath() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// path.Clean / path.Join etc will strip any . which is used to distinguish the reference types in the gateway
			// so we need to start it with a .
			relPath := "."

			r.URL.Path = relativePathRegex.ReplaceAllStringFunc(r.URL.Path, func(s string) string {
				// get rid of all leading and ending :
				// this will leave a string looking like an absolute path, prefix it with a . for proper routing as a relative reference
				relPath = "." + strings.Trim(s, ":")
				return ""
			})

			rctx := chi.RouteContext(r.Context())
			rctx.URLParams.Add("relative-path", relPath) // prefix path with a .

			next.ServeHTTP(w, r)
		})
	}
}
