package middleware

import (
	"net/http"
)

// SilentRefresh allows the oidc client lib to silently refresh the token in an iframe
func SilentRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		next.ServeHTTP(w, r)
	})
}
