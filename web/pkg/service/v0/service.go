package svc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/web/pkg/assets"
	"github.com/owncloud/ocis/web/pkg/config"
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

	svc := Web{
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

// Web defines implements the business logic for Service.
type Web struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (p Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p Web) getPayload() (payload []byte, err error) {

	if p.config.Web.Path == "" {
		// render dynamically using config

		// provide default ocis-web options
		if p.config.Web.Config.Options == nil {
			p.config.Web.Config.Options = make(map[string]interface{})
			p.config.Web.Config.Options["hideSearchBar"] = true
		}

		// build theme url
		if themeServer, err := url.Parse(p.config.Web.ThemeServer); err == nil {
			p.config.Web.Config.Theme = themeServer.String() + p.config.Web.ThemePath
		} else {
			p.config.Web.Config.Theme = p.config.Web.ThemePath
		}

		if p.config.Web.Config.ExternalApps == nil {
			p.config.Web.Config.ExternalApps = []config.ExternalApp{
				{
					ID:   "settings",
					Path: "/settings.js",
				},
				{
					ID:   "accounts",
					Path: "/accounts.js",
				},
			}
		}

		// make apps render as empty array if it is empty
		// TODO remove once https://github.com/golang/go/issues/27589 is fixed
		if len(p.config.Web.Config.Apps) == 0 {
			p.config.Web.Config.Apps = make([]string, 0)
		}

		return json.Marshal(p.config.Web.Config)
	}

	// try loading from file
	if _, err = os.Stat(p.config.Web.Path); os.IsNotExist(err) {
		p.logger.Fatal().
			Err(err).
			Str("config", p.config.Web.Path).
			Msg("web config doesn't exist")
	}

	payload, err = ioutil.ReadFile(p.config.Web.Path)

	if err != nil {
		p.logger.Fatal().
			Err(err).
			Str("config", p.config.Web.Path).
			Msg("failed to read custom config")
	}
	return
}

// Config implements the Service interface.
func (p Web) Config(w http.ResponseWriter, r *http.Request) {

	payload, err := p.getPayload()
	if err != nil {
		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		p.logger.Error().Err(err).Msg("could not write config response")
	}
}

// Static simply serves all static files.
func (p Web) Static(ttl int) http.HandlerFunc {
	rootWithSlash := p.config.HTTP.Root

	if !strings.HasSuffix(rootWithSlash, "/") {
		rootWithSlash = rootWithSlash + "/"
	}
	assets := assets.New(
		assets.Logger(p.logger),
		assets.Config(p.config),
	)

	notFoundFunc := func(w http.ResponseWriter, r *http.Request) {
		// TODO: replace the redirect with a not found page containing a link to the Web UI
		http.Redirect(w, r, rootWithSlash, http.StatusTemporaryRedirect)
	}

	static := http.StripPrefix(
		rootWithSlash,
		interceptNotFound(
			http.FileServer(assets),
			notFoundFunc,
		),
	)

	// TODO: investigate broken caching - https://github.com/owncloud/ocis/issues/1094
	// we don't have a last modification date of the static assets, so we use the service start date
	//lastModified := time.Now().UTC().Format(http.TimeFormat)
	//expires := time.Now().Add(time.Second * time.Duration(ttl)).UTC().Format(http.TimeFormat)

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
			notFoundFunc(w, r)
			return
		}

		// TODO: investigate broken caching - https://github.com/owncloud/ocis/issues/1094
		//w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, must-revalidate", strconv.Itoa(ttl)))
		//w.Header().Set("Expires", expires)
		//w.Header().Set("Last-Modified", lastModified)
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("SameSite", "Strict")

		static.ServeHTTP(w, r)
	}
}

func interceptNotFound(h http.Handler, notFoundFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notFoundInterceptor := &NotFoundInterceptor{ResponseWriter: w}
		h.ServeHTTP(notFoundInterceptor, r)
		if notFoundInterceptor.status == http.StatusNotFound {
			notFoundFunc(w, r)
		}
	}
}

type NotFoundInterceptor struct {
	http.ResponseWriter
	status int
}

func (w *NotFoundInterceptor) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundInterceptor) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil
}
