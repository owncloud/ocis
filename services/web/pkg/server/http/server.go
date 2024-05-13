package http

import (
	"fmt"
	"path"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go-micro.dev/v4"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
	"github.com/owncloud/ocis/v2/services/web"
	"github.com/owncloud/ocis/v2/services/web/pkg/apps"
	webmid "github.com/owncloud/ocis/v2/services/web/pkg/middleware"
	svc "github.com/owncloud/ocis/v2/services/web/pkg/service/v0"
)

var (
	// _customAppsEndpoint path is used to make app artifacts available by the web service.
	_customAppsEndpoint = "/assets/apps"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service, err := http.NewService(
		http.TLSConfig(options.Config.HTTP.TLS),
		http.Logger(options.Logger),
		http.Namespace(options.Namespace),
		http.Name("web"),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
		http.TraceProvider(options.TraceProvider),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	gatewaySelector, err := pool.GatewaySelector(
		options.Config.GatewayAddress,
		pool.WithRegistry(registry.GetRegistry()),
		pool.WithTracerProvider(options.TraceProvider),
	)
	if err != nil {
		return http.Service{}, err
	}

	coreFS := fsx.NewFallbackFS(
		fsx.NewBasePathFs(fsx.NewOsFs(), options.Config.Asset.CorePath),
		fsx.NewBasePathFs(fsx.FromIOFS(web.Assets), "assets/core"),
	)
	appsFS := fsx.NewFallbackFS(
		fsx.NewReadOnlyFs(fsx.NewBasePathFs(fsx.NewOsFs(), options.Config.Asset.AppsPath)),
		fsx.NewBasePathFs(fsx.FromIOFS(web.Assets), "assets/apps"),
	)

	// build and inject the list of applications into the config
	for _, application := range apps.List(options.Logger, options.Config.Apps, appsFS.Secondary().IOFS(), appsFS.Primary().IOFS()) {
		options.Config.Web.Config.ExternalApps = append(
			options.Config.Web.Config.ExternalApps,
			application.ToExternal(path.Join(options.Config.HTTP.Root, _customAppsEndpoint)),
		)
	}

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.CoreFS(coreFS),
		svc.AppFS(appsFS.IOFS()),
		svc.AppsHTTPEndpoint(_customAppsEndpoint),
		svc.Config(options.Config),
		svc.GatewaySelector(gatewaySelector),
		svc.Middleware(
			chimiddleware.RealIP,
			chimiddleware.RequestID,
			middleware.NoCache,
			webmid.SilentRefresh,
			middleware.Version(
				options.Config.Service.Name,
				version.GetString(),
			),
			middleware.Logger(
				options.Logger,
			),
			middleware.Cors(
				cors.Logger(options.Logger),
				cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
				cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
				cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
				cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
			),
		),
		svc.TraceProvider(options.TraceProvider),
	)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}
