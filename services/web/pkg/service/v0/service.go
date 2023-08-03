package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/assetsfs"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/web"
	"github.com/owncloud/ocis/v2/services/web/pkg/assets"
	"github.com/owncloud/ocis/v2/services/web/pkg/config"
	"github.com/riandyrn/otelchi"
)

// ErrConfigInvalid is returned when the config parse is invalid.
var ErrConfigInvalid = `Invalid or missing config`

// Service defines the service handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Config(http.ResponseWriter, *http.Request)
	UploadLogo(http.ResponseWriter, *http.Request)
	ResetLogo(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	m.Use(
		otelchi.Middleware(
			"web",
			otelchi.WithChiRoutes(m),
			otelchi.WithTracerProvider(options.TraceProvider),
			otelchi.WithPropagators(tracing.GetPropagator()),
		),
	)
	svc := Web{
		logger:          options.Logger,
		config:          options.Config,
		mux:             m,
		fs:              assetsfs.New(web.Assets, options.Config.Asset.Path, options.Logger),
		gatewaySelector: options.GatewaySelector,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/config.json", svc.Config)
		r.Route("/branding/logo", func(r chi.Router) {
			r.Use(middleware.ExtractAccountUUID(
				account.Logger(options.Logger),
				account.JWTSecret(options.Config.TokenManager.JWTSecret),
			))
			r.Post("/", svc.UploadLogo)
			r.Delete("/", svc.ResetLogo)
		})
		r.Mount("/", svc.Static(options.Config.HTTP.CacheTTL))
	})

	_ = chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	return svc
}

// Web defines implements the business logic for Service.
type Web struct {
	logger          log.Logger
	config          *config.Config
	mux             *chi.Mux
	fs              *assetsfs.FileSystem
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// ServeHTTP implements the Service interface.
func (p Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p Web) getPayload() (payload []byte, err error) {
	if p.config.Web.Path == "" {
		// render dynamically using config

		// build theme url
		if themeServer, err := url.Parse(p.config.Web.ThemeServer); err == nil {
			p.config.Web.Config.Theme = themeServer.String() + p.config.Web.ThemePath
		} else {
			p.config.Web.Config.Theme = p.config.Web.ThemePath
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

	payload, err = os.ReadFile(p.config.Web.Path)

	if err != nil {
		p.logger.Fatal().
			Err(err).
			Str("config", p.config.Web.Path).
			Msg("failed to read custom config")
	}
	return
}

// Config implements the Service interface.
func (p Web) Config(w http.ResponseWriter, _ *http.Request) {
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

	static := http.StripPrefix(
		rootWithSlash,
		assets.FileServer(p.fs),
	)

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

		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s", strconv.Itoa(ttl)))
		w.Header().Set("Expires", expires)
		w.Header().Set("Last-Modified", lastModified)
		w.Header().Set("SameSite", "Strict")

		static.ServeHTTP(w, r)
	}
}
