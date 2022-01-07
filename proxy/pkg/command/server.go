package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
	"github.com/oklog/run"
	acc "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	pkgmiddleware "github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/config/parser"
	"github.com/owncloud/ocis/proxy/pkg/cs3"
	"github.com/owncloud/ocis/proxy/pkg/logging"
	"github.com/owncloud/ocis/proxy/pkg/metrics"
	"github.com/owncloud/ocis/proxy/pkg/middleware"
	"github.com/owncloud/ocis/proxy/pkg/proxy"
	"github.com/owncloud/ocis/proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis/proxy/pkg/server/http"
	"github.com/owncloud/ocis/proxy/pkg/tracing"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	storepb "github.com/owncloud/ocis/store/pkg/proto/v0"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			err := tracing.Configure(cfg)
			if err != nil {
				return err
			}

			var (
				m = metrics.New()
			)

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			defer cancel()

			m.BuildInfo.WithLabelValues(version.String).Set(1)

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
					proxyHTTP.Metrics(metrics.New()),
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
					logger.Error().Err(err).Str("server", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

func loadMiddlewares(ctx context.Context, logger log.Logger, cfg *config.Config) alice.Chain {
	rolesClient := settings.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	revaClient, err := cs3.GetGatewayServiceClient(cfg.Reva.Address)
	var userProvider backend.UserBackend
	switch cfg.AccountBackend {
	case "accounts":
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  cfg.TokenManager.JWTSecret,
			"expires": int64(24 * 60 * 60),
		})
		if err != nil {
			logger.Error().Err(err).
				Msg("Failed to create token manager")
		}
		userProvider = backend.NewAccountsServiceUserBackend(
			acc.NewAccountsService("com.owncloud.api.accounts", grpc.DefaultClient),
			rolesClient,
			cfg.OIDC.Issuer,
			tokenManager,
			logger,
		)
	case "cs3":
		userProvider = backend.NewCS3UserBackend(rolesClient, revaClient, cfg.MachineAuthAPIKey, logger)
	default:
		logger.Fatal().Msgf("Invalid accounts backend type '%s'", cfg.AccountBackend)
	}

	storeClient := storepb.NewStoreService("com.owncloud.api.store", grpc.DefaultClient)
	if err != nil {
		logger.Error().Err(err).
			Str("gateway", cfg.Reva.Address).
			Msg("Failed to create reva gateway service client")
	}

	var oidcHTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.OIDC.Insecure, //nolint:gosec
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 10,
	}

	return alice.New(
		// first make sure we log all requests and redirect to https if necessary
		pkgmiddleware.TraceContext,
		chimiddleware.RealIP,
		chimiddleware.RequestID,
		middleware.AccessLog(logger),
		middleware.HTTPSRedirect,

		// now that we established the basics, on with authentication middleware
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
			middleware.Logger(logger),
			middleware.EnableBasicAuth(cfg.EnableBasicAuth),
			middleware.UserProvider(userProvider),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.UserOIDCClaim(cfg.UserOIDCClaim),
			middleware.UserCS3Claim(cfg.UserCS3Claim),
			middleware.CredentialsByUserAgent(cfg.AuthMiddleware.CredentialsByUserAgent),
		),
		middleware.SignedURLAuth(
			middleware.Logger(logger),
			middleware.PreSignedURLConfig(cfg.PreSignedURL),
			middleware.UserProvider(userProvider),
			middleware.Store(storeClient),
		),
		middleware.AccountResolver(
			middleware.Logger(logger),
			middleware.UserProvider(userProvider),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.UserOIDCClaim(cfg.UserOIDCClaim),
			middleware.UserCS3Claim(cfg.UserCS3Claim),
			middleware.AutoprovisionAccounts(cfg.AutoprovisionAccounts),
		),

		middleware.SelectorCookie(
			middleware.Logger(logger),
			middleware.UserProvider(userProvider),
			middleware.PolicySelectorConfig(*cfg.PolicySelector),
		),

		// finally, trigger home creation when a user logs in
		middleware.CreateHome(
			middleware.Logger(logger),
			middleware.TokenManagerConfig(cfg.TokenManager),
			middleware.RevaGatewayClient(revaClient),
		),
		middleware.PublicShareAuth(
			middleware.Logger(logger),
			middleware.RevaGatewayClient(revaClient),
		),
	)
}
