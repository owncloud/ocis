package command

import (
	"context"
	"fmt"
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

// Frontend is the entrypoint for the frontend command.
func Frontend(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "frontend",
		Usage: "Start reva frontend service",
		Flags: flagset.FrontendWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Frontend.Services = c.StringSlice("service")

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

				// pregenerate list of valid localhost ports for the desktop redirect_uri
				// TODO use custom scheme like "owncloud://localhost/user/callback" tracked in
				var desktopRedirectURIs [65535 - 1024]string
				for port := 0; port < len(desktopRedirectURIs); port++ {
					desktopRedirectURIs[port] = fmt.Sprintf("http://localhost:%d", (port + 1024))
				}

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.Users.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": "frontend",
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
						"gatewaysvc": cfg.Reva.Gateway.URL, // Todo or address?
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.Frontend.Network,
						"address": cfg.Reva.Frontend.Addr,
						"middlewares": map[string]interface{}{
							"cors": map[string]interface{}{
								"allow_credentials": true,
							},
						},
						// TODO build services dynamically
						"services": map[string]interface{}{
							"datagateway": map[string]interface{}{
								"transfer_shared_secret": cfg.Reva.TransferSecret,
							},
							"ocdav": map[string]interface{}{
								"prefix":           "",
								"chunk_folder":     "/var/tmp/reva/chunks",
								"files_namespace":  cfg.Reva.OCDav.DavFilesNamespace,
								"webdav_namespace": cfg.Reva.OCDav.WebdavNamespace,
							},
							"ocs": map[string]interface{}{
								"config": map[string]interface{}{
									"version": "1.8",
									"website": "reva",
									"host":    "http://" + cfg.Reva.Frontend.URL, // TODO URLs should include the protocol
									"contact": "admin@localhost",
									"ssl":     "false",
								},
								"capabilities": map[string]interface{}{
									"capabilities": map[string]interface{}{
										"core": map[string]interface{}{
											"poll_interval": 60,
											"webdav_root":   "remote.php/webdav",
											"status": map[string]interface{}{
												"installed":      true,
												"maintenance":    false,
												"needsDbUpgrade": false,
												"version":        "10.0.11.5",
												"versionstring":  "10.0.11",
												"edition":        "community",
												"productname":    "reva",
												"hostname":       "",
											},
										},
										"checksums": map[string]interface{}{
											"supported_types":       []string{"SHA256"},
											"preferred_upload_type": "SHA256",
										},
										"files": map[string]interface{}{
											"private_links":     false,
											"bigfilechunking":   false,
											"blacklisted_files": []string{},
											"undelete":          true,
											"versioning":        true,
										},
										"dav": map[string]interface{}{},
										"files_sharing": map[string]interface{}{
											"api_enabled":                       true,
											"resharing":                         true,
											"group_sharing":                     true,
											"auto_accept_share":                 true,
											"share_with_group_members_only":     true,
											"share_with_membership_groups_only": true,
											"default_permissions":               22,
											"search_min_length":                 3,
											"public": map[string]interface{}{
												"enabled":              true,
												"send_mail":            true,
												"social_share":         true,
												"upload":               true,
												"multiple":             true,
												"supports_upload_only": true,
												"password": map[string]interface{}{
													"enforced": true,
													"enforced_for": map[string]interface{}{
														"read_only":   true,
														"read_write":  true,
														"upload_only": true,
													},
												},
												"expire_date": map[string]interface{}{
													"enabled": true,
												},
											},
											"user": map[string]interface{}{
												"send_mail": true,
											},
											"user_enumeration": map[string]interface{}{
												"enabled":            true,
												"group_members_only": true,
											},
											"federation": map[string]interface{}{
												"outgoing": true,
												"incoming": true,
											},
										},
										"notifications": map[string]interface{}{
											"endpoints": []string{"list", "get", "delete"},
										},
									},
									"version": map[string]interface{}{
										"edition": "reva",
										"major":   10,
										"minor":   0,
										"micro":   11,
										"string":  "10.0.11",
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
					debug.Addr(cfg.Reva.Frontend.DebugAddr),
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
