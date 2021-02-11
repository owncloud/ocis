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

// Users is the entrypoint for the sharing command.
func Users(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Start users service",
		Flags: flagset.UsersWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Users.Services = c.StringSlice("service")

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

			{
				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.Users.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Users.GRPCNetwork,
						"address": cfg.Reva.Users.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"userprovider": map[string]interface{}{
								"driver": cfg.Reva.Users.Driver,
								"drivers": map[string]interface{}{
									"json": map[string]interface{}{
										"users": cfg.Reva.Users.JSON,
									},
									"ldap": map[string]interface{}{
										"hostname":        cfg.Reva.LDAP.Hostname,
										"port":            cfg.Reva.LDAP.Port,
										"base_dn":         cfg.Reva.LDAP.BaseDN,
										"userfilter":      cfg.Reva.LDAP.UserFilter,
										"attributefilter": cfg.Reva.LDAP.AttributeFilter,
										"findfilter":      cfg.Reva.LDAP.FindFilter,
										"groupfilter":     cfg.Reva.LDAP.GroupFilter,
										"bind_username":   cfg.Reva.LDAP.BindDN,
										"bind_password":   cfg.Reva.LDAP.BindPassword,
										"idp":             cfg.Reva.LDAP.IDP,
										"schema": map[string]interface{}{
											"dn":          "dn",
											"uid":         cfg.Reva.LDAP.Schema.UID,
											"mail":        cfg.Reva.LDAP.Schema.Mail,
											"displayName": cfg.Reva.LDAP.Schema.DisplayName,
											"cn":          cfg.Reva.LDAP.Schema.CN,
											"uidNumber":   cfg.Reva.LDAP.Schema.UIDNumber,
											"gidNumber":   cfg.Reva.LDAP.Schema.GIDNumber,
										},
									},
									"rest": map[string]interface{}{
										"client_id":                    cfg.Reva.UserRest.ClientID,
										"client_secret":                cfg.Reva.UserRest.ClientSecret,
										"redis_address":                cfg.Reva.UserRest.RedisAddress,
										"redis_username":               cfg.Reva.UserRest.RedisUsername,
										"redis_password":               cfg.Reva.UserRest.RedisPassword,
										"user_groups_cache_expiration": cfg.Reva.UserRest.UserGroupsCacheExpiration,
										"id_provider":                  cfg.Reva.UserRest.IDProvider,
										"api_base_url":                 cfg.Reva.UserRest.APIBaseURL,
										"oidc_token_endpoint":          cfg.Reva.UserRest.OIDCTokenEndpoint,
										"target_api":                   cfg.Reva.UserRest.TargetAPI,
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
					debug.Addr(cfg.Reva.Users.DebugAddr),
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
