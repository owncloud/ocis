package middleware

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	ocisoidc "github.com/owncloud/ocis-pkg/oidc"
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
	var oidcMetadata *ocisoidc.ProviderMetadata

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
					InsecureSkipVerify: opt.Insecure,
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
				metadata := &ocisoidc.ProviderMetadata{}
				if err := provider.Claims(metadata); err != nil {
					opt.Logger.Error().Err(err).Msg("could not not unmarshal provider metadata")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				oidcMetadata = metadata
			}
			provider := oidcProvider

			// The claims we want to have
			var claims ocisoidc.StandardClaims

			if oidcMetadata.IntrospectionEndpoint == "" {

				opt.Logger.Debug().Msg("no introspection endpoint, trying to decode access token as jwt")
				//maybe our access token is a jwt token
				c := &oidc.Config{
					ClientID:             opt.Audience,
					SupportedSigningAlgs: opt.SigningAlgs,
				}
				if opt.SkipChecks { // not safe but only way for simplesamlphp to work with an almost compliant oidc (for now)
					c.SkipClientIDCheck = true
					c.SkipIssuerCheck = true
				}
				verifier := provider.Verifier(c)
				idToken, err := verifier.Verify(customCtx, token)
				if err != nil {
					opt.Logger.Error().Err(err).Str("token", token).Msg("could not verify jwt")
					w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, opt.Realm))
					http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
					return
				}
				if err := idToken.Claims(&claims); err != nil {
					opt.Logger.Error().Err(err).Str("token", token).Interface("id_token", idToken).Msg("failed to parse claims")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

			} else {

				// we need to lookup the id token with the access token we got
				// see oidc IDToken.Verifytoken

				data := fmt.Sprintf("token=%s&token_type_hint=access_token", token)
				req, err := http.NewRequest("POST", oidcMetadata.IntrospectionEndpoint, strings.NewReader(data))
				if err != nil {
					opt.Logger.Error().Err(err).Msg("could not create introspection request")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				// we follow https://tools.ietf.org/html/rfc7662
				req.Header.Set("Accept", "application/json")
				if opt.ClientID != "" {
					req.SetBasicAuth(opt.ClientID, opt.ClientSecret)
				}

				res, err := customHTTPClient.Do(req)
				if err != nil {
					opt.Logger.Error().Err(err).Str("token", token).Msg("could not introspect access token")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					opt.Logger.Error().Err(err).Msg("could not read introspection response body")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				opt.Logger.Debug().Str("body", string(body)).Msg("body")
				switch strings.Split(res.Header.Get("Content-Type"), ";")[0] {
				// application/jwt is in draft https://tools.ietf.org/html/draft-ietf-oauth-jwt-introspection-response-03
				case "application/jwt":
					// verify the jwt
					// TODO this is a yet untested verification of jwt encoded introspection response

					verifier := provider.Verifier(&oidc.Config{ClientID: opt.Audience})
					idToken, err := verifier.Verify(customCtx, string(body))
					if err != nil {
						opt.Logger.Error().Err(err).Str("token", string(body)).Msg("could not verify jwt")
						w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, opt.Realm))
						http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
						return
					}

					if err := idToken.Claims(&claims); err != nil {
						opt.Logger.Error().Err(err).Str("token", string(body)).Interface("id_token", idToken).Msg("failed to parse claims")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				case "application/json":
					var ir ocisoidc.IntrospectionResponse
					// parse json
					if err := json.Unmarshal(body, &ir); err != nil {
						opt.Logger.Error().Err(err).Str("token", string(body)).Msg("failed to parse introspection response")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					// verify the auth token is still active
					if !ir.Active {
						opt.Logger.Error().Interface("ir", ir).Str("body", string(body)).Msg("token no longer active")
						w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, opt.Realm))
						http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
						return
					}
					// resolve user info here? cache it?
					oauth2Token := &oauth2.Token{
						AccessToken: token,
					}
					userInfo, err := provider.UserInfo(customCtx, oauth2.StaticTokenSource(oauth2Token))
					if err != nil {
						opt.Logger.Error().Err(err).Str("token", string(body)).Msg("Failed to get userinfo")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if err := userInfo.Claims(&claims); err != nil {
						opt.Logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					claims.Iss = ir.Iss
					opt.Logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Msg("unmarshalled userinfo")

				default:
					opt.Logger.Error().Str("content-type", res.Header.Get("Content-Type")).Msg("unknown content type")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			// store claims in context
			// uses the original context, not the one with probably reduced security
			nr := r.WithContext(ocisoidc.NewContext(r.Context(), &claims))

			next.ServeHTTP(w, nr)
		})
	}
}
