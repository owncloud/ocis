package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	// SupportedAuthStrategies stores configured challenges.
	SupportedAuthStrategies []string

	// ProxyWwwAuthenticate is a list of endpoints that do not rely on reva underlying authentication, such as ocs.
	// services that fallback to reva authentication are declared in the "frontend" command on oCIS. It is a list of strings
	// to be regexp compiled.
	ProxyWwwAuthenticate = []string{"/ocs/v[12].php/cloud/"}

	// WWWAuthenticate captures the Www-Authenticate header string.
	WWWAuthenticate = "Www-Authenticate"
)

// userAgentLocker aids in dependency injection for helper methods. The set of fields is arbitrary and the only relation
// they share is to fulfill their duty and lock a User-Agent to its correct challenge if configured.
type userAgentLocker struct {
	w        http.ResponseWriter
	r        *http.Request
	locks    map[string]string // locks represents a reva user-agent:challenge mapping.
	fallback string
}

// Authentication is a higher order authentication middleware.
func Authentication(opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)

	configureSupportedChallenges(options)
	oidc := newOIDCAuth(options)
	basic := newBasicAuth(options)

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

// configureSupportedChallenges adds known authentication challenges to the current session.
func configureSupportedChallenges(options Options) {
	if options.OIDCIss != "" {
		SupportedAuthStrategies = append(SupportedAuthStrategies, "bearer")
	}

	if options.EnableBasicAuth {
		SupportedAuthStrategies = append(SupportedAuthStrategies, "basic")
	}
}

func writeSupportedAuthenticateHeader(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(SupportedAuthStrategies); i++ {
		w.Header().Add(WWWAuthenticate, fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(SupportedAuthStrategies[i]), r.Host))
	}
}

func removeSuperfluousAuthenticate(w http.ResponseWriter) {
	w.Header().Del(WWWAuthenticate)
}

// userAgentAuthenticateLockIn sets Www-Authenticate according to configured user agents. This is useful for the case of
// legacy clients that do not support protocols like OIDC or OAuth and want to lock a given user agent to a challenge
// such as basic. For more context check https://github.com/cs3org/reva/pull/1350
func userAgentAuthenticateLockIn(w http.ResponseWriter, r *http.Request, locks map[string]string, fallback string) {
	u := userAgentLocker{
		w:        w,
		r:        r,
		locks:    locks,
		fallback: fallback,
	}

	for i := 0; i < len(ProxyWwwAuthenticate); i++ {
		evalRequestURI(&u, i)
	}
}

func evalRequestURI(l *userAgentLocker, i int) {
	r := regexp.MustCompile(ProxyWwwAuthenticate[i])
	if r.Match([]byte(l.r.RequestURI)) {
		for k, v := range l.locks {
			if strings.Contains(k, l.r.UserAgent()) {
				removeSuperfluousAuthenticate(l.w)
				l.w.Header().Add(WWWAuthenticate, fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(v), l.r.Host))
				return
			}
		}
		l.w.Header().Add(WWWAuthenticate, fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(l.fallback), l.r.Host))
	}
}

// newOIDCAuth returns a configured oidc middleware
func newOIDCAuth(options Options) func(http.Handler) http.Handler {
	return OIDCAuth(
		Logger(options.Logger),
		OIDCProviderFunc(options.OIDCProviderFunc),
		HTTPClient(options.HTTPClient),
		OIDCIss(options.OIDCIss),
		TokenCacheSize(options.UserinfoCacheSize),
		TokenCacheTTL(options.UserinfoCacheTTL),
		CredentialsByUserAgent(options.CredentialsByUserAgent),
	)
}

// newBasicAuth returns a configured basic middleware
func newBasicAuth(options Options) func(http.Handler) http.Handler {
	return BasicAuth(
		UserProvider(options.UserProvider),
		Logger(options.Logger),
		EnableBasicAuth(options.EnableBasicAuth),
		AccountsClient(options.AccountsClient),
		OIDCIss(options.OIDCIss),
		CredentialsByUserAgent(options.CredentialsByUserAgent),
	)
}
