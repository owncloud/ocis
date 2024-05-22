package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func retryAfterFn(ctxDone bool) time.Duration {
	if ctxDone {
		return time.Minute
	}
	return time.Minute * 5
}

// Throttle limits the number of concurrent requests.
func Throttle(limit int) func(http.Handler) http.Handler {
	if limit > 0 {
		opts := middleware.ThrottleOpts{
			RetryAfterFn:   retryAfterFn,
			Limit:          limit,
			BacklogLimit:   0,
			BacklogTimeout: time.Minute,
		}
		return middleware.ThrottleWithOpts(opts)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
