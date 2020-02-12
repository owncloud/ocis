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
					"grpc": map[string]interface{}{
						"network": cfg.Reva.StorageHome.Network,
						"address": cfg.Reva.StorageHome.Addr,
						// TODO extract interceptor config, which is the same for all grpc services
						"interceptors": map[string]interface{}{
							"auth": map[string]interface{}{
								"token_manager": "jwt",
								"token_managers": map[string]interface{}{
									"jwt": map[string]interface{}{
										"secret": cfg.Reva.JWTSecret,
									},
								},
							},
						},
						// TODO build services dynamically
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"driver": cfg.Reva.StorageHome.Driver,
								"drivers": map[string]interface{}{
									"eos": map[string]interface{}{
										"namespace":              cfg.Reva.Storages.EOS.Namespace,
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
									},
									"local": map[string]interface{}{
										"root": cfg.Reva.Storages.Local.Root,
									},
									"owncloud": map[string]interface{}{
										"datadirectory": cfg.Reva.Storages.OwnCloud.Datadirectory,
										"scan":          cfg.Reva.Storages.OwnCloud.Scan,
										"autocreate":    cfg.Reva.Storages.OwnCloud.Autocreate,
										"redis":         cfg.Reva.Storages.OwnCloud.Redis,
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
								"path_wrapper": cfg.Reva.StorageHome.PathWrapper,
								"path_wrappers": map[string]interface{}{
									"context": map[string]interface{}{
										"prefix": cfg.Reva.StorageHome.PathWrapperContext.Prefix,
									},
								},
								"mount_path":         cfg.Reva.StorageHome.MountPath,
								"mount_id":           cfg.Reva.StorageHome.MountID,
								"expose_data_server": cfg.Reva.StorageHome.ExposeDataServer,
								// TODO use cfg.Reva.StorageHomeData.URL, ?
								"data_server_url": cfg.Reva.StorageHome.DataServerURL,
								"available_checksums": map[string]interface{}{
									"md5":   100,
									"unset": 1000,
								},
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
