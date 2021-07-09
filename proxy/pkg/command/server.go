package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
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
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}
			cfg.PreSignedURL.AllowedHTTPMethods = ctx.StringSlice("presignedurl-allow-method")

			if err := loadUserAgent(ctx, cfg); err != nil {
				return err
			}

			// beforeOverride contains cfg with values parsed by urfavecli,
			// this should take precedence when merging as they are more explicit.
			// beforeOverride has the highest priority, as they are inherited values.
			beforeOverride := config.Config{}
			if err := copier.Copy(&beforeOverride, cfg); err != nil {
				return err
			}
			defaultConfig := config.DefaultConfig()

			// By the time we unmarshal viper parsed values onto cfg, any value having been set by the cli framework
			// will get overridden, in order to ensure that this values are accounted for we have to perform a 3-way merge:
			// 1. merge viper onto cfg
			// 2. merge defaults onto cfg
			// 3. merge parsed flags onto cfg
			// the result of this is the same order of precedence as the cli framework claims, except a new "artificial"
			// source which accounts for structured configuration. This all goes to the moon when the extension is running
			// in supervised mode, this is because in such case we want the single config file to take precedence over
			// the global ocis.yaml config file. This is happening because in supervised mode, sending commands to a hot
			// runtime, flags forwarding is not possible, because the process is probably running in a machine elsewhere.
			// It is not impossible to do, it just needs design.
			if !cfg.Supervised {
				if err := ParseConfig(ctx, cfg); err != nil {
					return err
				}
			}

			fromProxyConfigFile := config.Config{}
			if err := ParseConfig(ctx, &fromProxyConfigFile); err != nil {
				return err
			}

			if err := mergo.Merge(cfg, defaultConfig); err != nil {
				panic(err)
			}

			// When an extension is running supervised, we have the use case where executing `ocis run extension`
			// we want to ONLY take into consideration fhe existing config file.
			if !reflect.DeepEqual(fromProxyConfigFile, config.Config{}) {
				if err := mergo.Merge(cfg, fromProxyConfigFile); err != nil {
					panic(err)
				}
				return nil
			}

			if err := mergo.Merge(cfg, fromProxyConfigFile); err != nil {
				panic(err)
			}

			// preserves the original order from inherited values. This has the drawback that also persists values
			// inherited from an ocis.yaml global config file, with the side effect of these global values overriding
			// concrete values from a closer-to-the-process proxy.yaml file.
			if err := mergo.Merge(cfg, beforeOverride, mergo.WithOverride); err != nil {
				panic(err)
			}

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
