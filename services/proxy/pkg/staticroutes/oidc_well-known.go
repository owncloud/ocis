package staticroutes

import (
	"io"
	"net/http"
	"net/url"
	"path"
)

var (
	wellKnownPath = "/.well-known/openid-configuration"
)

// OIDCWellKnownRewrite is a handler that rewrites the /.well-known/openid-configuration endpoint for external IDPs.
func (s *StaticRouteHandler) oIDCWellKnownRewrite(issuer string) http.HandlerFunc {
	oidcURL, _ := url.Parse(issuer)
	oidcURL.Path = path.Join(oidcURL.Path, wellKnownPath)
	return func(w http.ResponseWriter, r *http.Request) {
		wellKnownRes, err := s.OidcHttpClient.Get(oidcURL.String())
		if err != nil {
			s.Logger.Error().
				Err(err).
				Str("handler", "oidc wellknown rewrite").
				Str("url", oidcURL.String()).
				Msg("get information from url failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer wellKnownRes.Body.Close()

		copyHeader(w.Header(), wellKnownRes.Header)
		w.WriteHeader(wellKnownRes.StatusCode)
		_, err = io.Copy(w, wellKnownRes.Body)
		if err != nil {
			s.Logger.Error().
				Err(err).
				Str("handler", "oidc wellknown rewrite").
				Msg("copying response body failed")

		}
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
