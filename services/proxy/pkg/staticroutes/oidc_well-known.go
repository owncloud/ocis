package staticroutes

import (
	"io"
	"net/http"
)

var (
	wellKnownPath = "/.well-known/openid-configuration"
)

// OIDCWellKnownRewrite is a handler that rewrites the /.well-known/openid-configuration endpoint for external IDPs.
func (s *StaticRouteHandler) oIDCWellKnownRewrite(w http.ResponseWriter, r *http.Request) {
	wellKnownRes, err := s.OidcHttpClient.Get(s.oidcURL.String())
	if err != nil {
		s.Logger.Error().
			Err(err).
			Str("handler", "oidc wellknown rewrite").
			Str("url", s.oidcURL.String()).
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

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
