package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Authentication is a higher level authentication middleware.
func Authentication(opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)

	oidc := OIDCAuth(
		Logger(options.Logger),
		OIDCProviderFunc(options.OIDCProviderFunc),
		HTTPClient(options.HTTPClient),
		OIDCIss(options.OIDCIss),
		TokenCacheSize(options.UserinfoCacheSize),
		TokenCacheTTL(time.Second*time.Duration(options.UserinfoCacheTTL)),
	)

	basic := BasicAuth(
		Logger(options.Logger),
		EnableBasicAuth(options.EnableBasicAuth),
		AccountsClient(options.AccountsClient),
		OIDCIss(options.OIDCIss),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// here we multiplex depending on the use agent
			userAgent := r.Header.Get("User-Agent")
			fmt.Printf("\n\nUser-Agent:\t%s\n\n", userAgent)
			switch userAgent {
			case "a":
				oidc(next).ServeHTTP(w, r)
				return
			case "b":
				basic(next).ServeHTTP(w, r)
				return
			default:
				oidc(next).ServeHTTP(w, r)
				basic(next).ServeHTTP(w, r)
				return
			}
		})
	}
}
