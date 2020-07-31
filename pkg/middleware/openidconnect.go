package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
	"github.com/owncloud/ocis-proxy/pkg/cache"
	"golang.org/x/oauth2"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")

	// svcCache caches requests for given services to prevent round trips to the service
	svcCache = cache.NewCache(
		cache.Size(256),
	)
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*oidc.UserInfo, error)
}

// OpenIDConnect provides a middleware to check access secured by a static token.
func OpenIDConnect(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {

		var oidcProvider OIDCProvider
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			path := r.URL.Path

			// Ignore request to "/konnect/v1/userinfo" as this will cause endless loop when getting userinfo
			// needs a better idea on how to not hardcode this
			if header == "" || !strings.HasPrefix(header, "Bearer ") || path == "/konnect/v1/userinfo" {
				next.ServeHTTP(w, r)
				return
			}

			customCtx := context.WithValue(r.Context(), oauth2.HTTPClient, opt.HTTPClient)

			// check if oidc provider is initialized
			if oidcProvider == nil {
				// Lazily initialize a provider

				// provider needs to be cached as when it is created
				// it will fetch the keys from the issuer using the .well-known
				// endpoint
				var err error
				oidcProvider, err = opt.OIDCProviderFunc()
				if err != nil {
					opt.Logger.Error().Err(err).Msg("could not initialize oidc provider")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			token := strings.TrimPrefix(header, "Bearer ")

			// TODO cache userinfo for access token if we can determine the expiry (which works in case it is a jwt based access token)
			oauth2Token := &oauth2.Token{
				AccessToken: token,
			}

			// The claims we want to have
			var claims ocisoidc.StandardClaims
			userInfo, err := oidcProvider.UserInfo(customCtx, oauth2.StaticTokenSource(oauth2Token))
			if err != nil {
				opt.Logger.Error().Err(err).Str("token", token).Msg("Failed to get userinfo")
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			if err := userInfo.Claims(&claims); err != nil {
				opt.Logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			//TODO: This should be read from the token instead of config
			claims.Iss = opt.OIDCIss

			// inject claims to the request context for the account_uuid middleware.
			ctxWithClaims := ocisoidc.NewContext(r.Context(), &claims)
			r = r.WithContext(ctxWithClaims)

			opt.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Msg("unmarshalled userinfo")
			// store claims in context
			// uses the original context, not the one with probably reduced security
			nr := r.WithContext(ocisoidc.NewContext(r.Context(), &claims))

			next.ServeHTTP(w, nr)
		})
	}
}

// AccountsCacheEntry stores a request to the accounts service on the cache.
// this type declaration should be on each respective service.
type AccountsCacheEntry struct {
	Email string
	UUID  string
}

const (
	// AccountsKey declares the svcKey for the Accounts service.
	AccountsKey = "accounts"

	// NodeKey declares the key that will be used to store the node address.
	// It is shared between services.
	NodeKey = "node"
)
