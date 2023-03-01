package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	stdhttp "net/http"
	"os"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	graphMiddleware "github.com/owncloud/ocis/v2/services/graph/pkg/middleware"
	svc "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/pkg/errors"
	"go-micro.dev/v4"
	"go-micro.dev/v4/events"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service, err := http.NewService(
		http.TLSConfig(options.Config.HTTP.TLS),
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("graph"),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	var publisher events.Stream

	if options.Config.Events.Endpoint != "" {
		var err error

		var tlsConf *tls.Config
		if options.Config.Events.EnableTLS {
			var rootCAPool *x509.CertPool
			if options.Config.Events.TLSRootCACertificate != "" {
				rootCrtFile, err := os.Open(options.Config.Events.TLSRootCACertificate)
				if err != nil {
					return http.Service{}, err
				}

				rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
				if err != nil {
					return http.Service{}, err
				}
				options.Config.Events.TLSInsecure = false
			}

			tlsConf = &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: options.Config.Events.TLSInsecure, //nolint:gosec
				RootCAs:            rootCAPool,
			}
		}
		publisher, err = stream.Nats(
			natsjs.TLSConfig(tlsConf),
			natsjs.Address(options.Config.Events.Endpoint),
			natsjs.ClusterID(options.Config.Events.Cluster),
		)
		if err != nil {
			options.Logger.Error().
				Err(err).
				Msg("Error initializing events publisher")
			return http.Service{}, errors.Wrap(err, "could not initialize events publisher")
		}
	}

	middlewares := []func(stdhttp.Handler) stdhttp.Handler{
		middleware.TraceContext,
		chimiddleware.RequestID,
		middleware.Version(
			"graph",
			version.GetString(),
		),
		middleware.Logger(
			options.Logger,
		),
	}
	// how do we secure the api?
	var requireAdminMiddleware func(stdhttp.Handler) stdhttp.Handler
	var roleService svc.RoleService
	var gatewayClient svc.GatewayClient
	if options.Config.HTTP.APIToken == "" {
		middlewares = append(middlewares,
			graphMiddleware.Auth(
				account.Logger(options.Logger),
				account.JWTSecret(options.Config.TokenManager.JWTSecret),
			))
		roleService = settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient())
		gatewayClient, err = pool.GetGatewayServiceClient(options.Config.Reva.Address, options.Config.Reva.GetRevaOptions()...)
		if err != nil {
			return http.Service{}, errors.Wrap(err, "could not initialize gateway client")
		}
	} else {
		middlewares = append(middlewares, graphMiddleware.Token(options.Config.HTTP.APIToken))
		// use a dummy admin middleware for the chi router
		requireAdminMiddleware = func(next stdhttp.Handler) stdhttp.Handler {
			return next
		}
		// no gatewayclient needed
	}

	var handle svc.Service
	handle, err = svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(middlewares...),
		svc.EventsPublisher(publisher),
		svc.WithRoleService(roleService),
		svc.WithRequireAdminMiddleware(requireAdminMiddleware),
		svc.WithGatewayClient(gatewayClient),
		svc.WithSearchService(searchsvc.NewSearchProviderService("com.owncloud.api.search", grpc.DefaultClient())),
	)

	if err != nil {
		return http.Service{}, errors.New("could not initialize graph service")
	}

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}
