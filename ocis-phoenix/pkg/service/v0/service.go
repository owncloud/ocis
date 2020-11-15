package svc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis/ocis-phoenix/pkg/assets"
	"github.com/owncloud/ocis/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

var (
	// ErrConfigInvalid is returned when the config parse is invalid.
	ErrConfigInvalid = `Invalid or missing config`
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Config(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Phoenix{
		logger: options.Logger,
		config: options.Config,
		mux:    m,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/config.json", svc.Config)
		r.Mount("/", svc.Static(options.Config.HTTP.CacheTTL))
	})

	return svc
}

// Phoenix defines implements the business logic for Service.
type Phoenix struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (p Phoenix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p Phoenix) getPayload() (payload []byte, err error) {

	if p.config.Phoenix.Path == "" {
		// render dynamically using config

		// provide default ocis-web options
		if p.config.Phoenix.Config.Options == nil {
			p.config.Phoenix.Config.Options = make(map[string]interface{})
			p.config.Phoenix.Config.Options["hideSearchBar"] = true
		}

		if p.config.Phoenix.Config.ExternalApps == nil {
			p.config.Phoenix.Config.ExternalApps = []config.ExternalApp{
				{
					ID:   "settings",
					Path: "/settings.js",
				},
				{
					ID:   "accounts",
					Path: "/accounts.js",
				},
				{
					ID:   "onlyoffice",
					Path: "/onlyoffice.js",
				},
			}
		}

		// make apps render as empty array if it is empty
		// TODO remove once https://github.com/golang/go/issues/27589 is fixed
		if len(p.config.Phoenix.Config.Apps) == 0 {
			p.config.Phoenix.Config.Apps = make([]string, 0)
		}

		return json.Marshal(p.config.Phoenix.Config)
	}

	// try loading from file
	if _, err = os.Stat(p.config.Phoenix.Path); os.IsNotExist(err) {
		p.logger.Fatal().
			Err(err).
			Str("config", p.config.Phoenix.Path).
			Msg("phoenix config doesn't exist")
	}

	payload, err = ioutil.ReadFile(p.config.Phoenix.Path)

	if err != nil {
		p.logger.Fatal().
			Err(err).
			Str("config", p.config.Phoenix.Path).
			Msg("failed to read custom config")
	}
	return
}

// Config implements the Service interface.
func (p Phoenix) Config(w http.ResponseWriter, r *http.Request) {

	payload, err := p.getPayload()
	if err != nil {
		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// Static simply serves all static files.
func (p Phoenix) Static(ttl int) http.HandlerFunc {
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

	// we don't have a last modification date of the static assets, so we use the service start date
	lastModified := time.Now().UTC().Format(http.TimeFormat)
	expires := time.Now().Add(time.Second * time.Duration(ttl)).UTC().Format(http.TimeFormat)

	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, must-revalidate", strconv.Itoa(ttl)))
		w.Header().Set("Expires", expires)
		w.Header().Set("Last-Modified", lastModified)

		static.ServeHTTP(w, r)
	}
}
