package middleware

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	mclient "github.com/micro/go-micro/v2/client"
	acc "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
	"github.com/owncloud/ocis-proxy/pkg/cache"
	"golang.org/x/oauth2"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")

	// svcCache caches requests for given services to prevent round trips to the service
	svcCache = cache.NewCache()

	accountSvc = "com.owncloud.accounts"

	// UUIDKey works as a context key
	UUIDKey interface{} = "uuid"

	// ClaimsKey works as a context key for user claims
	ClaimsKey interface{} = "claims"
)

// newOIDCOptions initializes the available default options.
func newOIDCOptions(opts ...ocisoidc.Option) ocisoidc.Options {
	opt := ocisoidc.Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// OpenIDConnect provides a middleware to check access secured by a static token.
func OpenIDConnect(opts ...ocisoidc.Option) func(next http.Handler) http.Handler {
	opt := newOIDCOptions(opts...)

	// set defaults
	if opt.Realm == "" {
		opt.Realm = opt.Endpoint
	}
	if len(opt.SigningAlgs) < 1 {
		opt.SigningAlgs = []string{"RS256", "PS256"}
	}

	var oidcProvider *oidc.Provider

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			path := r.URL.Path

			// Ignore request to "/konnect/v1/userinfo" as this will cause endless loop when getting userinfo
			// needs a better idea on how to not hardcode this
			if header == "" || !strings.HasPrefix(header, "Bearer ") || path == "/konnect/v1/userinfo" {
				next.ServeHTTP(w, r)
				return
			}

			token := header[7:]
			customHTTPClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: opt.Insecure,
					},
				},
				Timeout: time.Second * 10,
			}

			customCtx := context.WithValue(r.Context(), oauth2.HTTPClient, customHTTPClient)

			// use cached provider
			if oidcProvider == nil {
				// Initialize a provider by specifying the issuer URL.
				// provider needs to be cached as when it is created
				// it will fetch the keys from the issuer using the .well-known
				// endpoint
				provider, err := oidc.NewProvider(customCtx, opt.Endpoint)
				if err != nil {
					opt.Logger.Error().Err(err).Msg("could not initialize oidc provider")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				oidcProvider = provider
			}

			// The claims we want to have
			var claims ocisoidc.StandardClaims

			// TODO cache userinfo for access token if we can determine the expiry (which works in case it is a jwt based access token)
			oauth2Token := &oauth2.Token{
				AccessToken: token,
			}
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

			// inject claims to the request context for the account_uuid middleware.
			ctxWithClaims := context.WithValue(r.Context(), ClaimsKey, claims)
			r = r.WithContext(ctxWithClaims)

			// add UUID to the request context for the handler to deal with
			_, err = uuidFromClaims(claims)
			if err != nil {
				opt.Logger.Error().Err(err).Interface("account uuid", userInfo).Msg("failed to unmarshal userinfo claims")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

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

// from the user claims we need to get the uuid from the accounts service
func uuidFromClaims(claims ocisoidc.StandardClaims) (string, error) {
	entry, err := svcCache.Get(AccountsKey, claims.Email)
	if err != nil {
		c := acc.NewAccountsService("com.owncloud.accounts", mclient.DefaultClient) // TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
		resp, err := c.Get(context.Background(), &acc.GetRequest{
			Uuid: "200~a54bf154-e6a5-4e96-851b-a56c9f6c1fce",
			// Email: claims.Email // depends on https://github.com/owncloud/ocis-accounts/pull/28
		})
		if err != nil {
			return "", err
		}

		// TODO add logging info. Round trip has been made to the accounts service.
		err = svcCache.Set(AccountsKey, claims.Email, resp.Payload.Account.Uuid)
		if err != nil {
			return "", err
		}

		return resp.Key, nil
	}

	uuid, ok := entry.V.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type on cache entry value. Expected string type")
	}

	// TODO add logging info. Read from cache.
	return uuid, nil
}
