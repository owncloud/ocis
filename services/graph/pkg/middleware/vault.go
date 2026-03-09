package middleware

import (
	"context"
	"net/http"
)

type key int

const vaultModeKey key = iota

func SetVaultMode(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, vaultModeKey, enabled)
}

func IsVaultMode(ctx context.Context) bool {
	val, ok := ctx.Value(vaultModeKey).(bool)
	return val && ok
}

func VaultModeMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(SetVaultMode(r.Context(), true)))
		})
	}
}
