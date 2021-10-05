package command

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// AppProvider is the entrypoint for the app provider command.
func AppProvider(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "app-provider",
		Usage: "Start appprovider for providing apps",
		Flags: flagset.AppProviderWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.AppProvider.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
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
				debug.Addr(cfg.Reva.AppProvider.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.AppProvider.Supervised {
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
			"max_cpus":             cfg.Reva.AppProvider.MaxCPUs,
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret": cfg.Reva.JWTSecret,
			"gatewaysvc": cfg.Reva.Gateway.Endpoint,
		},
		"grpc": map[string]interface{}{
			"network": cfg.Reva.AppProvider.GRPCNetwork,
			"address": cfg.Reva.AppProvider.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"appprovider": map[string]interface{}{
					"gatewaysvc":       cfg.Reva.Gateway.Endpoint,
					"app_provider_url": cfg.Reva.AppProvider.ExternalAddr,
					"driver":           cfg.Reva.AppProvider.Driver,
					"drivers": map[string]interface{}{
						"wopi": map[string]interface{}{
							"app_api_key":          cfg.Reva.AppProvider.WopiDriver.AppAPIKey,
							"app_desktop_only":     cfg.Reva.AppProvider.WopiDriver.AppDesktopOnly,
							"app_icon_uri":         cfg.Reva.AppProvider.WopiDriver.AppIconURI,
							"app_int_url":          cfg.Reva.AppProvider.WopiDriver.AppInternalURL,
							"app_name":             cfg.Reva.AppProvider.WopiDriver.AppName,
							"app_url":              cfg.Reva.AppProvider.WopiDriver.AppURL,
							"insecure_connections": cfg.Reva.AppProvider.WopiDriver.Insecure,
							"iop_secret":           cfg.Reva.AppProvider.WopiDriver.IopSecret,
							"jwt_secret":           cfg.Reva.AppProvider.WopiDriver.JWTSecret,
							"wopi_url":             cfg.Reva.AppProvider.WopiDriver.WopiURL,
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
	if cfg.Mode == 0 {
		cfg.Storage.Reva.AppProvider.Supervised = true
	}
	return AppProviderSutureService{
		cfg: cfg.Storage,
	}
}

func (s AppProviderSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.AppProvider.Context = ctx
	f := &flag.FlagSet{}
	for k := range AppProvider(s.cfg).Flags {
		if err := AppProvider(s.cfg).Flags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if AppProvider(s.cfg).Before != nil {
		if err := AppProvider(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := AppProvider(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
