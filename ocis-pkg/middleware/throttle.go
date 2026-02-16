package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Throttle limits the number of concurrent requests.
func Throttle(limit int) func(http.Handler) http.Handler {
	if limit > 0 {
		return middleware.Throttle(limit)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
