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
			c := acc.NewAccountsService("com.owncloud.accounts", mclient.DefaultClient) // TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
			resp, err := c.Get(context.Background(), &acc.GetRequest{
				Email: claims.(ocisoidc.StandardClaims).Email,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = svcCache.Set(AccountsKey, claims.(ocisoidc.StandardClaims).Email, resp.Payload.Account.Uuid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: build JWT and set it, instead of the uuid on that header.
			w.Header().Set("x-ocis-accounts-uuid", resp.Payload.Account.Uuid)
		} else {
			uuid, ok := entry.V.(string)
			if !ok {
				// placeholder. Add more meaningful response
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO: build JWT and set it, instead of the uuid on that header.
			w.Header().Set("x-ocis-accounts-uuid", uuid)
		}

		next.ServeHTTP(w, r)
	})
}
