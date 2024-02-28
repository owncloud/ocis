package staticroutes

import (
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	microstore "go-micro.dev/v4/store"
	"net/http"
)

// StaticRouteHandler defines a Route Handler for static routes
type StaticRouteHandler struct {
	Prefix        string
	Proxy         http.Handler
	UserInfoCache microstore.Store
	Logger        log.Logger
	Config        config.Config
	OidcClient    oidc.OIDCClient
}

type jse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *StaticRouteHandler) Handler() http.Handler {
	m := chi.NewMux()
	m.Route(s.Prefix, func(r chi.Router) {
		// Wrapper for backchannel logout
		r.Post("/backchannel_logout", s.backchannelLogout)

		// TODO: migrate oidc well knowns here in a second wrapper

		// Send all requests to the proxy handler
		r.HandleFunc("/*", s.Proxy.ServeHTTP)
	})

	// Also send requests for methods unknown to chi to the proxy handler as well
	m.MethodNotAllowed(s.Proxy.ServeHTTP)

	return m
}
