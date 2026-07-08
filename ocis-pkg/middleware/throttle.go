package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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

// LimiterPerEndpoint rate limits the number of requests per endpoint.
func LimiterPerEndpoint(requestLimit int, windowLength time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitBy(requestLimit, windowLength, httprate.KeyByEndpoint)
}
