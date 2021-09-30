package middleware

import (
	"net/http"
	"time"
)

// NoCache writes required cache headers to all requests.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		next.ServeHTTP(w, r)
	})
}

// Cors writes required cors headers to all requests.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept, x-requested-with")
			w.Header().Set("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")

			w.WriteHeader(http.StatusOK)
		}
	})
}

// Secure writes required access headers to all requests.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}

		next.ServeHTTP(w, r)
	})
}
