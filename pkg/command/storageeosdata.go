package command

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis-reva/pkg/server/debug"
)

// StorageEOSData is the entrypoint for the storage-oc-data command.
func StorageEOSData(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-eos-data",
		Usage: "Start reva storage-eos-data service",
		Flags: flagset.StorageEOSDataWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.StorageEOSData.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring reva to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//metrics     = metrics.New()
			)

			defer cancel()

			{
				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.StorageEOSData.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": "storage-eos-data",
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
						"gatewaysvc": cfg.Reva.Gateway.URL, // Todo or address?
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.StorageEOSData.Network,
						"address": cfg.Reva.StorageEOSData.Addr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"dataprovider": map[string]interface{}{
								"prefix": cfg.Reva.StorageEOSData.Prefix,
								"driver": cfg.Reva.StorageEOSData.Driver,
								"drivers": map[string]interface{}{
									"eos": map[string]interface{}{
										"namespace":              cfg.Reva.Storages.EOS.Namespace,
										"shadow_namespace":       cfg.Reva.Storages.EOS.ShadowNamespace,
										"share_folder":           cfg.Reva.Storages.EOS.ShareFolder,
										"eos_binary":             cfg.Reva.Storages.EOS.EosBinary,
										"xrdcopy_binary":         cfg.Reva.Storages.EOS.XrdcopyBinary,
										"master_url":             cfg.Reva.Storages.EOS.MasterURL,
										"slave_url":              cfg.Reva.Storages.EOS.SlaveURL,
										"cache_directory":        cfg.Reva.Storages.EOS.CacheDirectory,
										"enable_logging":         cfg.Reva.Storages.EOS.EnableLogging,
										"show_hidden_sys_files":  cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
										"force_single_user_mode": cfg.Reva.Storages.EOS.ForceSingleUserMode,
										"use_keytab":             cfg.Reva.Storages.EOS.UseKeytab,
										"sec_protocol":           cfg.Reva.Storages.EOS.SecProtocol,
										"keytab":                 cfg.Reva.Storages.EOS.Keytab,
										"single_username":        cfg.Reva.Storages.EOS.SingleUsername,
										"user_layout":            cfg.Reva.Storages.EOS.Layout,
									},
									"eoshome": map[string]interface{}{
										"namespace":              cfg.Reva.Storages.EOS.Namespace,
										"shadow_namespace":       cfg.Reva.Storages.EOS.ShadowNamespace,
										"share_folder":           cfg.Reva.Storages.EOS.ShareFolder,
										"eos_binary":             cfg.Reva.Storages.EOS.EosBinary,
										"xrdcopy_binary":         cfg.Reva.Storages.EOS.XrdcopyBinary,
										"master_url":             cfg.Reva.Storages.EOS.MasterURL,
										"slave_url":              cfg.Reva.Storages.EOS.SlaveURL,
										"cache_directory":        cfg.Reva.Storages.EOS.CacheDirectory,
										"enable_logging":         cfg.Reva.Storages.EOS.EnableLogging,
										"show_hidden_sys_files":  cfg.Reva.Storages.EOS.ShowHiddenSysFiles,
										"force_single_user_mode": cfg.Reva.Storages.EOS.ForceSingleUserMode,
										"use_keytab":             cfg.Reva.Storages.EOS.UseKeytab,
										"sec_protocol":           cfg.Reva.Storages.EOS.SecProtocol,
										"keytab":                 cfg.Reva.Storages.EOS.Keytab,
										"single_username":        cfg.Reva.Storages.EOS.SingleUsername,
										"user_layout":            cfg.Reva.Storages.EOS.Layout,
									},
								},
								"temp_folder": cfg.Reva.StorageEOSData.TempFolder,
							},
						},
					},
				}

				gr.Add(func() error {
					runtime.RunWithOptions(
						rcfg,
						pidFile,
						runtime.WithLogger(&logger.Logger),
					)
					return nil
				}, func(_ error) {
					logger.Info().
						Str("server", c.Command.Name).
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.Reva.StorageEOSData.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", c.Command.Name+"-debug").
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
						logger.Info().
							Err(err).
							Str("server", c.Command.Name+"-debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", c.Command.Name+"-debug").
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
