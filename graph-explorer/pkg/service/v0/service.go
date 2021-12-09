package svc

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/graph-explorer/pkg/assets"
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	ConfigJs(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := GraphExplorer{
		logger: options.Logger,
		config: options.Config,
		mux:    m,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/config.js", svc.ConfigJs)
		r.Mount("/", svc.Static())
	})

	return svc
}

// GraphExplorer defines implements the business logic for Service.
type GraphExplorer struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (p GraphExplorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

// ConfigJs implements the Service interface.
func (p GraphExplorer) ConfigJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, fmt.Sprintf("window.ClientId = \"%v\";", p.config.GraphExplorer.ClientID)); err != nil {
		p.logger.Error().Err(err).Msg("Could not write to response writer")
	}
	if _, err := io.WriteString(w, fmt.Sprintf("window.Iss = \"%v\";", p.config.GraphExplorer.Issuer)); err != nil {
		p.logger.Error().Err(err).Msg("Could not write to response writer")
	}
	if _, err := io.WriteString(w, fmt.Sprintf("window.GraphUrl = \"%v\";", p.config.GraphExplorer.GraphURLBase+p.config.GraphExplorer.GraphURLPath)); err != nil {
		p.logger.Error().Err(err).Msg("Could not write to response writer")
	}
}

// Static simply serves all static files.
func (p GraphExplorer) Static() http.HandlerFunc {
	rootWithSlash := p.config.HTTP.Root

	if !strings.HasSuffix(rootWithSlash, "/") {
		rootWithSlash = rootWithSlash + "/"
	}

	static := http.StripPrefix(
		rootWithSlash,
		http.FileServer(
			assets.New(
				assets.Logger(p.logger),
				assets.Config(p.config),
			),
		),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rootWithSlash != "/" && r.URL.Path == p.config.HTTP.Root {
			http.Redirect(
				w,
				r,
				rootWithSlash,
				http.StatusMovedPermanently,
			)

			return
		}

		if r.URL.Path != rootWithSlash && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(
				w,
				r,
			)

			return
		}

		static.ServeHTTP(w, r)
	})
}
