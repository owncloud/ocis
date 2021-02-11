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
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
)

// StorageMetadata the entrypoint for the storage-storage-metadata command.
//
// It provides a ocis-specific storage store metadata (shares,account,settings...)
func StorageMetadata(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-metadata",
		Usage:    "Start storage-metadata service",
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
						Msg("configuring storage to use the jaeger tracing backend")

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
						"max_cpus":             cfg.Reva.StorageMetadata.MaxCPUs,
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
						"network": cfg.Reva.StorageMetadata.GRPCNetwork,
						"address": cfg.Reva.StorageMetadata.GRPCAddr,
						"interceptors": map[string]interface{}{
							"log": map[string]interface{}{},
						},
						"services": map[string]interface{}{
							"storageprovider": map[string]interface{}{
								"mount_path":      "/meta",
								"driver":          cfg.Reva.StorageMetadata.Driver,
								"drivers":         drivers(cfg),
								"data_server_url": cfg.Reva.StorageMetadata.DataServerURL,
								"tmp_folder":      cfg.Reva.StorageMetadata.TempFolder,
							},
						},
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.StorageMetadata.HTTPNetwork,
						"address": cfg.Reva.StorageMetadata.HTTPAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"dataprovider": map[string]interface{}{
								"prefix":      "data",
								"driver":      cfg.Reva.StorageMetadata.Driver,
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

			// the defensive code is needed because sending to a nil channel blocks forever
			if cfg.C != nil {
				*cfg.C <- struct{}{}
			}
			return gr.Run()
		},
	}
}
