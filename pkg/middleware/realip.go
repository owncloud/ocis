package middleware

import (
	"net/http"

	"github.com/tomasen/realip"
)

// RealIP is a middleware that sets a http.Request RemoteAddr.
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip := realip.RealIP(r); ip != "" {
			r.RemoteAddr = ip
		}

		next.ServeHTTP(w, r)
	})
}
