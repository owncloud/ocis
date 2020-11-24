package middleware

import (
	"net/http"

	"go.uber.org/ratelimit"
)

// RateLimit creates a simple rate limiter middleware
func RateLimit(limit int) func(http.Handler) http.Handler {
	var rl ratelimit.Limiter

	if limit > 0 {
		rl = ratelimit.New(limit)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				if rl != nil {
					rl.Take()
				}

				next.ServeHTTP(res, req)
			},
		)
	}
}
