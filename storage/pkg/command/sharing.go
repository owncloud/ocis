package command

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
)

// Sharing is the entrypoint for the sharing command.
func Sharing(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "sharing",
		Usage: "Start sharing service",
		Flags: flagset.SharingWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Sharing.Services = c.StringSlice("service")

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
				//metrics     = metrics.New()
			)

			defer cancel()

			// precreate folders
			if cfg.Reva.Sharing.UserDriver == "json" && cfg.Reva.Sharing.UserJSONFile != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Sharing.UserJSONFile), os.FileMode(0700)); err != nil {
					return err
				}
			}
			if cfg.Reva.Sharing.PublicDriver == "json" && cfg.Reva.Sharing.PublicJSONFile != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Sharing.PublicJSONFile), os.FileMode(0700)); err != nil {
					return err
				}
			}

			{

				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.Sharing.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Sharing.GRPCNetwork,
						"address": cfg.Reva.Sharing.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"usershareprovider": map[string]interface{}{
								"driver": cfg.Reva.Sharing.UserDriver,
								"drivers": map[string]interface{}{
									"json": map[string]interface{}{
										"file": cfg.Reva.Sharing.UserJSONFile,
									},
									"sql": map[string]interface{}{
										"db_username": cfg.Reva.Sharing.UserSQLUsername,
										"db_password": cfg.Reva.Sharing.UserSQLPassword,
										"db_host":     cfg.Reva.Sharing.UserSQLHost,
										"db_port":     cfg.Reva.Sharing.UserSQLPort,
										"db_name":     cfg.Reva.Sharing.UserSQLName,
									},
								},
							},
							"publicshareprovider": map[string]interface{}{
								"driver": cfg.Reva.Sharing.PublicDriver,
								"drivers": map[string]interface{}{
									"json": map[string]interface{}{
										"file": cfg.Reva.Sharing.PublicJSONFile,
									},
								},
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
					debug.Addr(cfg.Reva.Sharing.DebugAddr),
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

// SharingSutureService allows for the storage-sharing command to be embedded and supervised by a suture supervisor tree.
type SharingSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewSharingSutureService creates a new store.SharingSutureService
func NewSharing(ctx context.Context) SharingSutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg := config.New()
	cfg.Context = sctx
	return SharingSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s SharingSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range Sharing(s.cfg).Flags {
		if err := Sharing(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if Sharing(s.cfg).Before != nil {
		if err := Sharing(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := Sharing(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s SharingSutureService) Stop() {
	s.cancel()
}
