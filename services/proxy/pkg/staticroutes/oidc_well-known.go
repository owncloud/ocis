package staticroutes

import (
	"io"
	"net/http"
)

var (
	wellKnownPath = "/.well-known/openid-configuration"
)

//oidcURL, _ := url.Parse(oidcISS)
//oidcURL.Path = path.Join(oidcURL.Path, wellKnownPath)

// OIDCWellKnownRewrite is a middleware that rewrites the /.well-known/openid-configuration endpoint for external IDPs.
func (s *StaticRouteHandler) OIDCWellKnownRewrite(w http.ResponseWriter, r *http.Request) {
	wellKnownRes, err := s.OidcHttpClient.Get(s.oidcURL.String())
	if err != nil {
		s.Logger.Error().
			Err(err).
			Str("middleware", "oidc wellknown rewrite").
			Str("url", s.oidcURL.String()).
			Msg("get information from url failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer wellKnownRes.Body.Close()

	copyHeader(w.Header(), wellKnownRes.Header)
	w.WriteHeader(wellKnownRes.StatusCode)
	io.Copy(w, wellKnownRes.Body)

	return
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
