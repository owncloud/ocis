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

// StorageHome is the entrypoint for the storage-home command.
func StorageHome(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "storage-home",
		Usage: "Start reva storage-home service",
		Flags: flagset.StorageHomeWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.StorageHome.Services = c.StringSlice("service")

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
						"max_cpus": cfg.Reva.StorageHome.MaxCPUs,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.StorageHome.Network,
						"address": cfg.Reva.StorageHome.Addr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"driver": cfg.Reva.StorageHome.Driver,
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
										"enable_home":            true,
										"user_layout":            cfg.Reva.Storages.EOS.Layout,
									},
									"local": map[string]interface{}{
										"root": cfg.Reva.Storages.Local.Root,
									},
									"owncloud": map[string]interface{}{
										"datadirectory": cfg.Reva.Storages.OwnCloud.Datadirectory,
										"scan":          cfg.Reva.Storages.OwnCloud.Scan,
										"redis":         cfg.Reva.Storages.OwnCloud.Redis,
										"enable_home":   true,
										"user_layout":   cfg.Reva.Storages.OwnCloud.Layout,
									},
									"s3": map[string]interface{}{
										"region":     cfg.Reva.Storages.S3.Region,
										"access_key": cfg.Reva.Storages.S3.AccessKey,
										"secret_key": cfg.Reva.Storages.S3.SecretKey,
										"endpoint":   cfg.Reva.Storages.S3.Endpoint,
										"bucket":     cfg.Reva.Storages.S3.Bucket,
										"prefix":     cfg.Reva.Storages.S3.Prefix,
									},
								},
								"mount_path":         cfg.Reva.StorageHome.MountPath,
								"mount_id":           cfg.Reva.StorageHome.MountID,
								"expose_data_server": cfg.Reva.StorageHome.ExposeDataServer,
								// TODO use cfg.Reva.StorageHomeData.URL, ?
								"data_server_url":      cfg.Reva.StorageHome.DataServerURL,
								"enable_home_creation": cfg.Reva.StorageHome.EnableHomeCreation,
							},
						},
					},
				}

				gr.Add(func() error {
					runtime.Run(rcfg, pidFile)
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
					debug.Addr(cfg.Reva.StorageHome.DebugAddr),
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
