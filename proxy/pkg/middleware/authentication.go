package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var SupportedAuthStrategies []string

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Authentication is a higher order authentication middleware.
func Authentication(opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)
	if options.OIDCIss != "" {
		SupportedAuthStrategies = append(SupportedAuthStrategies, "bearer")
	}

	if options.EnableBasicAuth {
		SupportedAuthStrategies = append(SupportedAuthStrategies, "basic")
	}

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
			if options.OIDCIss != "" && options.EnableBasicAuth {
				oidc(basic(next)).ServeHTTP(w, r)
			}

			if options.OIDCIss != "" && !options.EnableBasicAuth {
				oidc(next).ServeHTTP(w, r)
			}

			if options.OIDCIss == "" && options.EnableBasicAuth {
				basic(next).ServeHTTP(w, r)
			}

		})
	}
}

func writeSupportedAuthenticateHeader(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(SupportedAuthStrategies); i++ {
		w.Header().Add("WWW-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(SupportedAuthStrategies[i]), r.Host))
	}
}
