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

// AuthBasic is the entrypoint for the auth-basic command.
func AuthBasic(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-basic",
		Usage: "Start authprovider for basic auth",
		Flags: flagset.AuthBasicWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.AuthBasic.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Storage only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring storage to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Storage only supports the jaeger tracing backend")

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
			if cfg.Reva.AuthProvider.Driver == "json" && cfg.Reva.AuthProvider.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.AuthProvider.JSON), os.ModeExclusive); err != nil {
					return err
				}
			}

			{

				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.AuthBasic.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.AuthBasic.GRPCNetwork,
						"address": cfg.Reva.AuthBasic.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"authprovider": map[string]interface{}{
								"auth_manager": cfg.Reva.AuthProvider.Driver,
								"auth_managers": map[string]interface{}{
									"json": map[string]interface{}{
										"users": cfg.Reva.AuthProvider.JSON,
									},
									"ldap": map[string]interface{}{
										"hostname":      cfg.Reva.LDAP.Hostname,
										"port":          cfg.Reva.LDAP.Port,
										"base_dn":       cfg.Reva.LDAP.BaseDN,
										"loginfilter":   cfg.Reva.LDAP.LoginFilter,
										"bind_username": cfg.Reva.LDAP.BindDN,
										"bind_password": cfg.Reva.LDAP.BindPassword,
										"idp":           cfg.Reva.LDAP.IDP,
										"gatewaysvc":    cfg.Reva.Gateway.Endpoint,
										"schema": map[string]interface{}{
											"dn":          "dn",
											"uid":         cfg.Reva.LDAP.UserSchema.UID,
											"mail":        cfg.Reva.LDAP.UserSchema.Mail,
											"displayName": cfg.Reva.LDAP.UserSchema.DisplayName,
											"cn":          cfg.Reva.LDAP.UserSchema.CN,
										},
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
					debug.Addr(cfg.Reva.AuthBasic.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", "debug").
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
							Str("server", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", "debug").
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

// AuthBasicSutureService allows for the storage-authbasic command to be embedded and supervised by a suture supervisor tree.
type AuthBasicSutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewAuthBasicSutureService creates a new store.AuthBasicSutureService
func NewAuthBasic(ctx context.Context, cfg *config.Config) AuthBasicSutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Context = sctx
	return AuthBasicSutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s AuthBasicSutureService) Serve() {
	f := &flag.FlagSet{}
	for k := range AuthBasic(s.cfg).Flags {
		if err := AuthBasic(s.cfg).Flags[k].Apply(f); err != nil {
			return
		}
	}
	ctx := cli.NewContext(nil, f, nil)
	if AuthBasic(s.cfg).Before != nil {
		if err := AuthBasic(s.cfg).Before(ctx); err != nil {
			return
		}
	}
	if err := AuthBasic(s.cfg).Action(ctx); err != nil {
		return
	}
}

func (s AuthBasicSutureService) Stop() {
	s.cancel()
}
