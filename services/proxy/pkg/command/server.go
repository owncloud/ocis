package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/cs3org/reva/v2/pkg/token/manager/jwt"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	pkgmiddleware "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/sync"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/cs3"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/logging"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis/v2/services/proxy/pkg/server/http"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/tracing"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
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
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
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

			m.BuildInfo.WithLabelValues(version.GetString()).Set(1)

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
				}, func(err error) {
					logger.Error().
						Err(err).
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
	rolesClient := settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	revaClient, err := cs3.GetGatewayServiceClient(cfg.Reva.Address)
	var userProvider backend.UserBackend
	switch cfg.AccountBackend {
	case "cs3":
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret": cfg.TokenManager.JWTSecret,
		})
		if err != nil {
			logger.Error().Err(err).
				Msg("Failed to create token manager")
		}

		userProvider = backend.NewCS3UserBackend(rolesClient, revaClient, cfg.MachineAuthAPIKey, cfg.OIDC.Issuer, tokenManager, logger)
	default:
		logger.Fatal().Msgf("Invalid accounts backend type '%s'", cfg.AccountBackend)
	}

	storeClient := storesvc.NewStoreService("com.owncloud.api.store", grpc.DefaultClient)
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

	var authenticators []middleware.Authenticator
	if cfg.EnableBasicAuth {
		logger.Warn().Msg("basic auth enabled, use only for testing or development")
		authenticators = append(authenticators, middleware.BasicAuthenticator{
			Logger:       logger,
			UserProvider: userProvider,
		})
	}
	tokenCache := sync.NewCache(cfg.OIDC.UserinfoCache.Size)
	authenticators = append(authenticators, middleware.OIDCAuthenticator{
		Logger:        logger,
		TokenCache:    &tokenCache,
		TokenCacheTTL: time.Duration(cfg.OIDC.UserinfoCache.TTL),
		HTTPClient:    oidcHTTPClient,
		OIDCIss:       cfg.OIDC.Issuer,
		ProviderFunc: func() (middleware.OIDCProvider, error) {
			// Initialize a provider by specifying the issuer URL.
			// it will fetch the keys from the issuer using the .well-known
			// endpoint
			return oidc.NewProvider(
				context.WithValue(ctx, oauth2.HTTPClient, oidcHTTPClient),
				cfg.OIDC.Issuer,
			)
		},
		JWKSOptions:             cfg.OIDC.JWKS,
		AccessTokenVerifyMethod: cfg.OIDC.AccessTokenVerifyMethod,
	})
	authenticators = append(authenticators, middleware.PublicShareAuthenticator{
		Logger:            logger,
		RevaGatewayClient: revaClient,
	})

	authenticators = append(authenticators, middleware.SignedURLAuthenticator{
		Logger:             logger,
		PreSignedURLConfig: cfg.PreSignedURL,
		UserProvider:       userProvider,
		Store:              storeClient,
	})

	return alice.New(
		// first make sure we log all requests and redirect to https if necessary
		pkgmiddleware.TraceContext,
		chimiddleware.RealIP,
		chimiddleware.RequestID,
		middleware.AccessLog(logger),
		middleware.HTTPSRedirect,
		middleware.OIDCWellKnownRewrite(
			logger, cfg.OIDC.Issuer,
			cfg.OIDC.RewriteWellKnown,
			oidcHTTPClient,
		),

		middleware.Authentication(
			authenticators,
			middleware.CredentialsByUserAgent(cfg.AuthMiddleware.CredentialsByUserAgent),
			middleware.Logger(logger),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.EnableBasicAuth(cfg.EnableBasicAuth),
		),
		middleware.AccountResolver(
			middleware.Logger(logger),
			middleware.UserProvider(userProvider),
			middleware.TokenManagerConfig(*cfg.TokenManager),
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
			middleware.TokenManagerConfig(*cfg.TokenManager),
			middleware.RevaGatewayClient(revaClient),
		),
	)
}
