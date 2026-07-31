package middleware

import (
	"context"
	"net/http"
)

type key int

const vaultModeKey key = iota

// SetVaultMode sets the vault mode in the context.
func SetVaultMode(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, vaultModeKey, enabled)
}

// IsVaultMode checks if the vault mode is enabled in the context.
func IsVaultMode(ctx context.Context) bool {
	val, ok := ctx.Value(vaultModeKey).(bool)
	return val && ok
}

// VaultModeMiddleware is a middleware that sets the vault mode in the context.
func VaultModeMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(SetVaultMode(r.Context(), true)))
		})
	}
}
