package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// Version writes the current version to the headers.
func Version(name, version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				fmt.Sprintf("X-%s-VERSION", strings.ToUpper(name)),
				version,
			)

			next.ServeHTTP(w, r)
		})
	}
}
