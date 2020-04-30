package middleware

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	acc "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
	"golang.org/x/oauth2"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")

	accountSvc = "com.owncloud.accounts"
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
func OpenIDConnect(opts ...ocisoidc.Option) M {
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

			opt.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Msg("unmarshalled userinfo")
			// store claims in context
			// uses the original context, not the one with probably reduced security
			nr := r.WithContext(ocisoidc.NewContext(r.Context(), &claims))

			next.ServeHTTP(w, nr)
		})
	}
}

// from the user claims we need to get the uuid from the accounts service
func uuidFromClaims(claims ocisoidc.StandardClaims) (string, error) {
	var node string
	// get accounts node from micro registry
	// TODO this assumes we use mdns as registry. This should be configurable for any ocis extension.
	svc, err := registry.GetService(accountSvc)
	if err != nil {
		return "", err
	}

	if len(svc) > 0 {
		node = svc[0].Nodes[0].Address
	}

	c := acc.NewSettingsService("accounts", mclient.DefaultClient)
	_, err = c.Get(context.Background(), &acc.Query{
		// TODO accounts query message needs to be updated to query for multiple fields
		// queries by key makes little sense as it is unknown.
		Key: "73912d13-32f7-4fb6-aeb2-ea2088a3a264",
	})
	if err != nil {
		return "", err
	}

	// by this point, rec.Payload contains the Account info. To include UUID, see:
	// https://github.com/owncloud/ocis-accounts/pull/22/files#diff-b425175389864c4f9218ecd9cae80223R23

	// return rec.GetPayload().Account.UUID, nil // depends on the aforementioned PR
	return node, nil
}
