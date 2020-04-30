package middleware

import (
	"net/http"
)

// Logger undocummented
func Logger() M {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do some logging logic here
			next.ServeHTTP(w, r)
		})
	}
}
