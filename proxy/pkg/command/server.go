package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/justinas/alice"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	acc "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/conversions"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/cs3"
	"github.com/owncloud/ocis/proxy/pkg/flagset"
	"github.com/owncloud/ocis/proxy/pkg/metrics"
	"github.com/owncloud/ocis/proxy/pkg/middleware"
	"github.com/owncloud/ocis/proxy/pkg/proxy"
	"github.com/owncloud/ocis/proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis/proxy/pkg/server/http"
	"github.com/owncloud/ocis/proxy/pkg/tracing"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	storepb "github.com/owncloud/ocis/store/pkg/proto/v0"
	"golang.org/x/oauth2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: append(flagset.ServerWithConfig(cfg), flagset.RootWithConfig(cfg)...),
		Before: func(ctx *cli.Context) error {
			logger := NewLogger(cfg)
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}
			cfg.PreSignedURL.AllowedHTTPMethods = ctx.StringSlice("presignedurl-allow-method")

			if err := loadUserAgent(ctx, cfg); err != nil {
				return err
			}

			if !cfg.Supervised {
				return ParseConfig(ctx, cfg)
			}
			logger.Debug().Str("service", "ocs").Msg("ignoring config file parsing when running supervised")
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if err := tracing.Configure(cfg, logger); err != nil {
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

			m.BuildInfo.WithLabelValues(cfg.Service.Version).Set(1)

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

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
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
				InsecureSkipVerify: cfg.OIDC.Insecure, //nolint:gosec
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 10,
	}

	return alice.New(
		// first make sure we log all requests and redirect to https if necessary
		chimiddleware.RealIP,
		chimiddleware.RequestID,
		middleware.AccessLog(l),
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
			middleware.UserOIDCClaim(cfg.UserOIDCClaim),
			middleware.UserCS3Claim(cfg.UserCS3Claim),
			middleware.AutoprovisionAccounts(cfg.AutoprovisionAccounts),
		),

		middleware.SelectorCookie(
			middleware.Logger(l),
			middleware.UserProvider(userProvider),
			middleware.PolicySelectorConfig(*cfg.PolicySelector),
		),

		// finally, trigger home creation when a user logs in
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
