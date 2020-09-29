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
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
	"github.com/owncloud/ocis/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/ocis-reva/pkg/server/debug"
)

// StorageMetadata the entrypoint for the reva-storage-metadata command.
//
// It provides a ocis-specific storage store metadata (shares,account,settings...)
func StorageMetadata(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "reva-storage-metadata",
		Usage:    "Start reva storage-metadata service",
		Flags:    flagset.StorageMetadata(cfg),
		Category: "Extensions",
		Before: func(c *cli.Context) error {
			storageRoot := c.String("storage-root")

			cfg.Reva.Storages.OwnCloud.Root = storageRoot
			cfg.Reva.Storages.EOS.Root = storageRoot
			cfg.Reva.Storages.Local.Root = storageRoot
			cfg.Reva.Storages.S3.Root = storageRoot
			cfg.Reva.Storages.Home.Root = storageRoot

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

				// Disable home because the metadata is stored independently
				// of the user. This also means that a valid-token without any user-id
				// is allowed to write to the metadata-storage.
				cfg.Reva.Storages.Common.EnableHome = false
				cfg.Reva.Storages.EOS.EnableHome = false
				cfg.Reva.Storages.Local.EnableHome = false
				cfg.Reva.Storages.OwnCloud.EnableHome = false
				cfg.Reva.Storages.S3.EnableHome = false

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             "100",
						"tracing_enabled":      false,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": "storage-metadata",
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.StorageMetadata.Network,
						"address": cfg.Reva.StorageMetadata.Addr,
						"interceptors": map[string]interface{}{
							"log": map[string]interface{}{},
						},
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"mount_path":      "/meta",
								"data_server_url": cfg.Reva.StorageMetadataData.URL,
								"driver":          cfg.Reva.StorageMetadata.Driver,
								"drivers":         drivers(cfg),
							},
						},
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.StorageMetadataData.Network,
						"address": cfg.Reva.StorageMetadataData.Addr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"dataprovider": map[string]interface{}{
								"prefix":      "data",
								"driver":      cfg.Reva.StorageMetadataData.Driver,
								"drivers":     drivers(cfg),
								"timeout":     86400,
								"insecure":    true,
								"disable_tus": true,
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
					debug.Addr(cfg.Reva.StorageMetadata.DebugAddr),
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
