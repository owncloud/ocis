package middleware

import (
	"context"
	"net/http"

	mclient "github.com/micro/go-micro/v2/client"
	acc "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
)

// AccountUUID fetches the ocis account uuid from the oidc standard claims
func AccountUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ClaimsKey)
		if claims == nil {
			next.ServeHTTP(w, r)
			return
		}

		entry, err := svcCache.Get(AccountsKey, claims.(ocisoidc.StandardClaims).Email)
		if err != nil {
			c := acc.NewSettingsService("com.owncloud.accounts", mclient.DefaultClient) // TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
			resp, err := c.Get(context.Background(), &acc.Query{
				Key: "200~a54bf154-e6a5-4e96-851b-a56c9f6c1fce",
				// Email: claims.Email // depends on https://github.com/owncloud/ocis-accounts/pull/28
			})
			if err != nil {
				// placeholder. Add more meaningful response
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = svcCache.Set(AccountsKey, claims.(ocisoidc.StandardClaims).Email, resp.Payload.Account.Uuid)
			if err != nil {
				// placeholder. Add more meaningful response
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: build JWT and set it, instead of the uuid on that header.
			w.Header().Set("x-ocis-accounts-uuid", resp.Payload.Account.Uuid)
		}

		uuid, ok := entry.V.(string)
		if !ok {
			// placeholder. Add more meaningful response
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: build JWT and set it, instead of the uuid on that header.
		w.Header().Set("x-ocis-accounts-uuid", uuid)

		next.ServeHTTP(w, r)
	})
}
