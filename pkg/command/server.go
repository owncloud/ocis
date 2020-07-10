package command

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/justinas/alice"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-proxy/pkg/cs3"
	"github.com/owncloud/ocis-proxy/pkg/middleware"
	"golang.org/x/oauth2"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/micro/cli/v2"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	acc "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-proxy/pkg/config"
	"github.com/owncloud/ocis-proxy/pkg/flagset"
	"github.com/owncloud/ocis-proxy/pkg/metrics"
	"github.com/owncloud/ocis-proxy/pkg/proxy"
	"github.com/owncloud/ocis-proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis-proxy/pkg/server/http"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			l := NewLogger(cfg)
			l.Debug().Str("tracing", strconv.FormatBool(cfg.Tracing.Enabled)).Msg("init: before")
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			httpNamespace := c.String("http-namespace")

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					exporter, err := ocagent.NewExporter(
						ocagent.WithReconnectionPeriod(5*time.Second),
						ocagent.WithAddress(cfg.Tracing.Endpoint),
						ocagent.WithServiceName(cfg.Tracing.Service),
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create agent tracing")

						return err
					}

					trace.RegisterExporter(exporter)
					view.RegisterExporter(exporter)

				case "jaeger":
					exporter, err := jaeger.NewExporter(
						jaeger.Options{
							AgentEndpoint:     cfg.Tracing.Endpoint,
							CollectorEndpoint: cfg.Tracing.Collector,
							ServiceName:       cfg.Tracing.Service,
						},
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create jaeger tracing")

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create zipkin tracing")

						return err
					}

					exporter := zipkin.NewExporter(
						zipkinhttp.NewReporter(
							cfg.Tracing.Collector,
						),
						endpoint,
					)

					trace.RegisterExporter(exporter)

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				metrics     = metrics.New()
			)

			defer cancel()

			rp := proxy.NewMultiHostReverseProxy(
				proxy.Logger(logger),
				proxy.Config(cfg),
			)

			{
				server, err := proxyHTTP.Server(
					proxyHTTP.Handler(rp),
					proxyHTTP.Logger(logger),
					proxyHTTP.Namespace(httpNamespace),
					proxyHTTP.Context(ctx),
					proxyHTTP.Config(cfg),
					proxyHTTP.Metrics(metrics),
					proxyHTTP.Flags(flagset.RootWithConfig(config.New())),
					proxyHTTP.Flags(flagset.ServerWithConfig(config.New())),
					proxyHTTP.Middlewares(loadMiddlewares(ctx, logger, cfg)),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(_ error) {
					logger.Info().
						Str("server", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Error().
							Err(err).
							Str("server", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", "debug").
							Msg("Shutting down server")
					}
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

func loadMiddlewares(ctx context.Context, l log.Logger, cfg *config.Config) alice.Chain {
	if cfg.OIDC.Issuer != "" {
		l.Info().Msg("Loading OIDC-Middleware")
		l.Debug().Interface("oidc_config", cfg.OIDC).Msg("OIDC-Config")

		var oidcHTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.OIDC.Insecure,
				},
				DisableKeepAlives: true,
			},
			Timeout: time.Second * 10,
		}

		customCtx := context.WithValue(ctx, oauth2.HTTPClient, oidcHTTPClient)

		// Initialize a provider by specifying the issuer URL.
		// it will fetch the keys from the issuer using the .well-known
		// endpoint
		provider := func() (middleware.OIDCProvider, error) {
			return oidc.NewProvider(customCtx, cfg.OIDC.Issuer)
		}

		oidcMW := middleware.OpenIDConnect(
			middleware.Logger(l),
			middleware.HTTPClient(oidcHTTPClient),
			middleware.OIDCProviderFunc(provider),
		)

		// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
		// https://github.com/owncloud/ocis-proxy/issues/38
		accounts := acc.NewAccountsService("com.owncloud.api.accounts", mclient.DefaultClient)

		uuidMW := middleware.AccountUUID(
			middleware.Logger(l),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.AccountsClient(accounts),
		)

		// the connection will be established in a non blocking fashion
		sc, err := cs3.GetGatewayServiceClient(cfg.Reva.Address)
		if err != nil {
			l.Error().Err(err).
				Str("gateway", cfg.Reva.Address).
				Msg("Failed to create reva gateway service client")
		}

		chMW := middleware.CreateHome(
			middleware.Logger(l),
			middleware.RevaGatewayClient(sc),
			middleware.AccountsClient(accounts),
		)

		return alice.New(middleware.RedirectToHTTPS, oidcMW, uuidMW, chMW)
	}

	return alice.New(middleware.RedirectToHTTPS)
}
