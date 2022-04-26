package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/appprovider/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// AppProvider is the entrypoint for the app provider command.
func AppProvider(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "app-provider",
		Usage: "start appprovider for providing apps",
		Action: func(c *cli.Context) error {
			logCfg := cfg.Logging
			logger := log.NewLogger(
				log.Level(logCfg.Level),
				log.File(logCfg.File),
				log.Pretty(logCfg.Pretty),
				log.Color(logCfg.Color),
			)
			tracing.Configure(cfg.Tracing.Enabled, cfg.Tracing.Type, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			rcfg := appProviderConfigFromStruct(c, cfg)

			gr.Add(func() error {
				runtime.RunWithOptions(rcfg, pidFile, runtime.WithLogger(&logger.Logger))
				return nil
			}, func(_ error) {
				logger.Info().
					Str("server", c.Command.Name).
					Msg("Shutting down server")

				cancel()
			})

			debugServer, err := debug.Server(
				debug.Name(c.Command.Name+"-debug"),
				debug.Addr(cfg.Debug.Addr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Pprof(cfg.Debug.Pprof),
				debug.Zpages(cfg.Debug.Zpages),
				debug.Token(cfg.Debug.Token),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// appProviderConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func appProviderConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {

	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.JWTSecret,
			"gatewaysvc":                cfg.GatewayEndpoint,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"appprovider": map[string]interface{}{
					"app_provider_url": cfg.ExternalAddr,
					"driver":           cfg.Driver,
					"drivers": map[string]interface{}{
						"wopi": map[string]interface{}{
							"app_api_key":          cfg.Drivers.WOPI.AppAPIKey,
							"app_desktop_only":     cfg.Drivers.WOPI.AppDesktopOnly,
							"app_icon_uri":         cfg.Drivers.WOPI.AppIconURI,
							"app_int_url":          cfg.Drivers.WOPI.AppInternalURL,
							"app_name":             cfg.Drivers.WOPI.AppName,
							"app_url":              cfg.Drivers.WOPI.AppURL,
							"insecure_connections": cfg.Drivers.WOPI.Insecure,
							"iop_secret":           cfg.Drivers.WOPI.IopSecret,
							"jwt_secret":           cfg.JWTSecret,
							"wopi_url":             cfg.Drivers.WOPI.WopiURL,
						},
					},
				},
			},
		},
	}
	return rcfg
}

// AppProviderSutureService allows for the app-provider command to be embedded and supervised by a suture supervisor tree.
type AppProviderSutureService struct {
	cfg *config.Config
}

// NewAppProvider creates a new store.AppProviderSutureService
func NewAppProvider(cfg *ociscfg.Config) suture.Service {
	cfg.AppProvider.Commons = cfg.Commons
	return AppProviderSutureService{
		cfg: cfg.AppProvider,
	}
}

func (s AppProviderSutureService) Serve(ctx context.Context) error {
	cmd := AppProvider(s.cfg)
	f := &flag.FlagSet{}
	cmdFlags := cmd.Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if cmd.Before != nil {
		if err := cmd.Before(cliCtx); err != nil {
			return err
		}
	}
	if err := cmd.Action(cliCtx); err != nil {
		return err
	}

	return nil
}
