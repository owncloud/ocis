package middleware

import (
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

var (
	wellKnownPath = "/.well-known/openid-configuration"
)

// OIDCWellKnownRewrite is a middleware that rewrites the /.well-known/openid-configuration endpoint for external IDPs.
func OIDCWellKnownRewrite(logger log.Logger, oidcISS string, rewrite bool, oidcClient *http.Client) func(http.Handler) http.Handler {

	oidcURL, _ := url.Parse(oidcISS)
	oidcURL.Path = path.Join(oidcURL.Path, wellKnownPath)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rewrite && path.Clean(r.URL.Path) == wellKnownPath {

				wellKnownRes, err := oidcClient.Get(oidcURL.String())
				if err != nil {
					logger.Error().
						Err(err).
						Str("middleware", "oidc wellknown rewrite").
						Str("url", oidcURL.String()).
						Msg("get information from url failed")
					w.WriteHeader(http.StatusInternalServerError)
				}
				defer wellKnownRes.Body.Close()

				copyHeader(w.Header(), wellKnownRes.Header)
				w.WriteHeader(wellKnownRes.StatusCode)
				io.Copy(w, wellKnownRes.Body)

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
