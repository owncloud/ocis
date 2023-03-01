package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = "invalid or missing token"
)

// Token provides a middleware to check access secured by a static token.
func Token(token string) func(http.Handler) http.Handler {
	h := sha256.New()
	requiredTokenHash := h.Sum(([]byte("Bearer " + token)))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")

			if header == "" {
				errorcode.InvalidAuthenticationToken.Render(w, r, http.StatusUnauthorized, ErrInvalidToken)
				return
			}

			h = sha256.New()
			providedTokenHash := h.Sum([]byte(header))

			if subtle.ConstantTimeCompare(requiredTokenHash, providedTokenHash) == 0 {
				errorcode.InvalidAuthenticationToken.Render(w, r, http.StatusUnauthorized, ErrInvalidToken)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
