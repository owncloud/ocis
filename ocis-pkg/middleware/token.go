package middleware

import (
	"errors"
	"net/http"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")
)

// Token provides a middleware to check access secured by a static token.
func Token(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")

			if header == "" {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			if header != "Bearer "+token {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
