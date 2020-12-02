package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/coreos/go-oidc"
	"github.com/justinas/alice"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	acc "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/cs3"
	"github.com/owncloud/ocis/proxy/pkg/flagset"
	"github.com/owncloud/ocis/proxy/pkg/metrics"
	"github.com/owncloud/ocis/proxy/pkg/middleware"
	"github.com/owncloud/ocis/proxy/pkg/proxy"
	"github.com/owncloud/ocis/proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis/proxy/pkg/server/http"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	storepb "github.com/owncloud/ocis/store/pkg/proto/v0"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"golang.org/x/oauth2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}
			cfg.PreSignedURL.AllowedHTTPMethods = ctx.StringSlice("presignedurl-allow-method")

			cfg.Reva.Middleware.Auth.CredentialsByUserAgent = make(map[string]string, 0)
			uaw := ctx.StringSlice("proxy-user-agent-whitelist")
			for _, v := range uaw {
				parts := strings.Split(v, ":")
				if len(parts) != 2 {
					return fmt.Errorf("unexpected config value for user-agent whitelist: %v, expected format is userAgent:challenge", v)
				}

				cfg.Reva.Middleware.Auth.CredentialsByUserAgent[parts[0]] = parts[1]
			}

			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

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
							Process: jaeger.Process{
								ServiceName: cfg.Tracing.Service,
							},
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

			metrics.BuildInfo.WithLabelValues(cfg.Service.Version).Set(1)

			rp := proxy.NewMultiHostReverseProxy(
				proxy.Logger(logger),
				proxy.Config(cfg),
			)

			{
				server, err := proxyHTTP.Server(
					proxyHTTP.Handler(rp),
					proxyHTTP.Logger(logger),
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
	accountsClient := acc.NewAccountsService("com.owncloud.api.accounts", grpc.DefaultClient)
	rolesClient := settings.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	storeClient := storepb.NewStoreService("com.owncloud.api.store", grpc.DefaultClient)
	revaClient, err := cs3.GetGatewayServiceClient(cfg.Reva.Address)
	if err != nil {
		l.Error().Err(err).
			Str("gateway", cfg.Reva.Address).
			Msg("Failed to create reva gateway service client")
	}

	var oidcHTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.OIDC.Insecure,
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 10,
	}

	return alice.New(
		middleware.HTTPSRedirect,
		middleware.Authentication(
			// OIDC Options
			middleware.OIDCProviderFunc(func() (middleware.OIDCProvider, error) {
				// Initialize a provider by specifying the issuer URL.
				// it will fetch the keys from the issuer using the .well-known
				// endpoint
				return oidc.NewProvider(
					context.WithValue(ctx, oauth2.HTTPClient, oidcHTTPClient),
					cfg.OIDC.Issuer,
				)
			}),
			middleware.HTTPClient(oidcHTTPClient),
			middleware.TokenCacheSize(cfg.OIDC.UserinfoCache.Size),
			middleware.TokenCacheTTL(time.Second*time.Duration(cfg.OIDC.UserinfoCache.TTL)),

			// basic Options
			middleware.Logger(l),
			middleware.EnableBasicAuth(cfg.EnableBasicAuth),
			middleware.AccountsClient(accountsClient),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.CredentialsByUserAgent(cfg.Reva.Middleware.Auth.CredentialsByUserAgent),
		),
		middleware.SignedURLAuth(
			middleware.Logger(l),
			middleware.PreSignedURLConfig(cfg.PreSignedURL),
			middleware.AccountsClient(accountsClient),
			middleware.Store(storeClient),
		),
		middleware.AccountResolver(
			middleware.Logger(l),
			middleware.AccountsClient(accountsClient),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.AutoprovisionAccounts(cfg.AutoprovisionAccounts),
			middleware.SettingsRoleService(rolesClient),
		),
		middleware.CreateHome(
			middleware.Logger(l),
			middleware.AccountsClient(accountsClient),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.RevaGatewayClient(revaClient),
		),
	)
}
