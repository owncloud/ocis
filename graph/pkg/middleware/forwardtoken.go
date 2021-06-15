package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/token"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"google.golang.org/grpc/metadata"
)

// ForwardToken provides a middleware that adds a received x-access-token to the outgoung context
func ForwardToken(opts ...account.Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := r.Header.Get(token.TokenHeader)

			ctx := metadata.AppendToOutgoingContext(r.Context(), token.TokenHeader, t)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
