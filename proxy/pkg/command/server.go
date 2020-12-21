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

	"github.com/owncloud/ocis/proxy/pkg/user/backend"

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
	"github.com/owncloud/ocis/ocis-pkg/conversions"
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

			if err := loadUserAgent(ctx, cfg); err != nil {
				return err
			}

			if err := ParseConfig(ctx, cfg); err != nil {
				return err
			}
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
	rolesClient := settings.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	revaClient, err := cs3.GetGatewayServiceClient(cfg.Reva.Address)
	var userProvider backend.UserBackend
	switch cfg.AccountBackend {
	case "accounts":
		userProvider = backend.NewAccountsServiceUserBackend(
			acc.NewAccountsService("com.owncloud.api.accounts", grpc.DefaultClient),
			rolesClient,
			cfg.OIDC.Issuer,
			l,
		)
	case "cs3":
		userProvider = backend.NewCS3UserBackend(revaClient, rolesClient, revaClient, l)
	default:
		l.Fatal().Msgf("Invalid accounts backend type '%s'", cfg.AccountBackend)
	}

	storeClient := storepb.NewStoreService("com.owncloud.api.store", grpc.DefaultClient)
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
			middleware.UserProvider(userProvider),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.CredentialsByUserAgent(cfg.Reva.Middleware.Auth.CredentialsByUserAgent),
		),
		middleware.SignedURLAuth(
			middleware.Logger(l),
			middleware.PreSignedURLConfig(cfg.PreSignedURL),
			middleware.UserProvider(userProvider),
			middleware.Store(storeClient),
		),
		middleware.AccountResolver(
			middleware.Logger(l),
			middleware.UserProvider(userProvider),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.AutoprovisionAccounts(cfg.AutoprovisionAccounts),
		),
		middleware.CreateHome(
			middleware.Logger(l),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.RevaGatewayClient(revaClient),
		),
	)
}

// loadUserAgent reads the proxy-user-agent-lock-in, since it is a string flag, and attempts to construct a map of
// "user-agent":"challenge" locks in for Reva.
// Modifies cfg. Spaces don't need to be trimmed as urfavecli takes care of it. User agents with spaces are valid. i.e:
// Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0
// This function works by relying in our format of specifying [user-agent:challenge] and the fact that the user agent
// might contain ":" (colon), so the original string is reversed, split in two parts, by the time it is split we
// have the indexes reversed and the tuple is in the format of [challenge:user-agent], then the same process is applied
// in reverse for each individual part
func loadUserAgent(c *cli.Context, cfg *config.Config) error {
	cfg.Reva.Middleware.Auth.CredentialsByUserAgent = make(map[string]string)
	locks := c.StringSlice("proxy-user-agent-lock-in")

	for _, v := range locks {
		vv := conversions.Reverse(v)
		parts := strings.SplitN(vv, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unexpected config value for user-agent lock-in: %v, expected format is user-agent:challenge", v)
		}

		cfg.Reva.Middleware.Auth.CredentialsByUserAgent[conversions.Reverse(parts[1])] = conversions.Reverse(parts[0])
	}

	return nil
}
