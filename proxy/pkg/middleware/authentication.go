package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var SupportedAuthStrategies []string

// ProxyWwwAuthenticate is a list of endpoints that do not rely on reva underlying authentication, such as ocs.
// services that fallback to reva authentication are declared in the "frontend" command on OCIS.
// TODO this should be a regexp, or it can be confused with routes that contain "/ocs" somewhere along the URI
var ProxyWwwAuthenticate = []string{"ocs"}

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
		CredentialsByUserAgent(options.CredentialsByUserAgent),
	)

	basic := BasicAuth(
		Logger(options.Logger),
		EnableBasicAuth(options.EnableBasicAuth),
		AccountsClient(options.AccountsClient),
		OIDCIss(options.OIDCIss),
		CredentialsByUserAgent(options.CredentialsByUserAgent),
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

func removeSuperfluousAuthenticate(w http.ResponseWriter) {
	w.Header().Del("Www-Authenticate")
}

// userAgentAuthenticateLockIn sets Www-Authenticate according to configured user agents. This is useful for the case of
// legacy clients that do not support protocols like OIDC or OAuth and want to lock a given user agent to a challenge,
// such as basic. More info in https://github.com/cs3org/reva/pull/1350
func userAgentAuthenticateLockIn(w http.ResponseWriter, req *http.Request, creds map[string]string, fallback string) {
	for i := 0; i < len(ProxyWwwAuthenticate); i++ {
		if strings.Contains(req.RequestURI, fmt.Sprintf("/%v/", ProxyWwwAuthenticate[i])) {
			for k, v := range creds {
				if strings.Contains(k, req.UserAgent()) {
					removeSuperfluousAuthenticate(w)
					w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(v), req.Host))
					return
				}
			}
			w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(fallback), req.Host))
		}
	}
}
