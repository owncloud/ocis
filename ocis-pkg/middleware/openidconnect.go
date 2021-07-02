package middleware

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	ocisoidc "github.com/owncloud/ocis/ocis-pkg/oidc"
	"golang.org/x/oauth2"
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
func OpenIDConnect(opts ...ocisoidc.Option) func(http.Handler) http.Handler {
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

			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, opt.Realm))
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			token := header[7:]

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: opt.Insecure, //nolint:gosec
				},
			}
			customHTTPClient := &http.Client{
				Transport: tr,
				Timeout:   time.Second * 10,
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
				opt.Logger.Error().Err(err).Msg("Failed to get userinfo")
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// parse claims
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
