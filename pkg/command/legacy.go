package command

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis-reva/pkg/server/debug"
)

// Legacy is the entrypoint for the legacy command.
func Legacy(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "legacy",
		Usage: "Start legacy server mimicking oc10",
		Flags: flagset.ServerWithConfig(cfg),
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
			)

			defer cancel()

			// Flags have to be injected all the way down to the go-micro service
			{

				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": cfg.Tracing.Service,
					},
					"log": map[string]interface{}{
						"level": cfg.Reva.LogLevel,
						//TODO mode": ""console" # "console" or "json"
						//TODO output": ""./standalone.log"
					},
					"http": map[string]interface{}{
						"network": cfg.Reva.HTTP.Network,
						"address": cfg.Reva.HTTP.Addr,
						"enabled_services": []string{
							"dataprovider",
							"ocdav",
							"prometheus",
							"wellknown",
							"oidcprovider",
							"ocs",
						},
						"enabled_middlewares": []string{
							//"cors",
							"auth",
						},
						"middlewares": map[string]interface{}{
							"auth": map[string]interface{}{
								"gateway":             cfg.Reva.GRPC.Addr,
								"auth_type":           "oidc",
								"credential_strategy": "oidc",
								"token_strategy":      "header",
								"token_writer":        "header",
								"token_manager":       "jwt",
								"token_managers": map[string]interface{}{
									"jwt": map[string]interface{}{
										"secret": cfg.Reva.JWTSecret,
									},
								},
								"skip_methods": []string{
									"/favicon.ico",
									"/status.php",
									"/oauth2",
									"/oauth2/auth",
									"/oauth2/token",
									// TODO protect the introspection endpoint from external requests.
									// should only be reachable by internal services, which is why the
									// oidc-provider.toml has clientid and secret that are used for a basic auth
									//"/oauth2/introspect", // no longer used, oidc auth checks access token using the userinfo endpoint
									"/oauth2/userinfo",
									"/oauth2/sessions",
									"/.well-known/openid-configuration",
									"/metrics", // for prometheus metrics
								},
							},
						},
						"services": map[string]interface{}{
							"oidcprovider": map[string]interface{}{
								"prefix":    "oauth2",
								"gateway":   cfg.Reva.GRPC.Addr,
								"auth_type": "basic",
								"issuer":    cfg.Reva.HTTP.Addr,
								"clients": map[string]interface{}{
									"phoenix": map[string]interface{}{
										"id": "phoenix",
										// use ocis port range for phoenix
										// TODO should use the micro / ocis http gateway, but then it would no longer be able to run standalone
										// IMO the ports should be fetched from the ocis registry anyway
										"redirect_uris":  []string{"http://localhost:9100/oidc-callback.html", "http://localhost:9100/"},
										"grant_types":    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
										"response_types": []string{"code"}, // use authorization code flow, see https://developer.okta.com/blog/2019/05/01/is-the-oauth-implicit-flow-dead for details
										"scopes":         []string{"openid", "profile", "email", "offline"},
										"public":         true, // force PKCS for public clients
									},
									"cli": map[string]interface{}{
										"id":            "cli",
										"client_secret": "$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO", // = "foobar"
										// use hardcoded port credentials for cli
										"redirect_uris":  []string{"http://localhost:18080/callback"},
										"grant_types":    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
										"response_types": []string{"code"}, // use authorization code flow, see https://developer.okta.com/blog/2019/05/01/is-the-oauth-implicit-flow-dead for details
										"scopes":         []string{"openid", "profile", "email", "offline"},
									},
								},
							},
							"dataprovider": map[string]interface{}{
								"driver":     "owncloud",
								"prefix":     "data",
								"tmp_folder": "/var/tmp/",
								"drivers": map[string]interface{}{
									"owncloud": map[string]interface{}{
										"datadirectory": "/var/tmp/reva/data",
									},
								},
							},
							"ocdav": map[string]interface{}{
								"prefix":       "",
								"chunk_folder": "/var/tmp/revad/chunks",
								"gateway":      cfg.Reva.GRPC.Addr,
							},
							"ocs": map[string]interface{}{
								"gateway": cfg.Reva.GRPC.Addr,
								"config": map[string]interface{}{
									"version": "1.8",
									"website": "ocis",
									"host":    cfg.Reva.HTTP.Addr, // TODO should be read from registry
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
										"dav": map[string]interface{}{
											"chunking": "1.0",
										},
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
										"edition": "ocis",
										"major":   11,
										"minor":   0,
										"micro":   0,
										"string":  "11.0.0",
									},
								},
							},
						},
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.GRPC.Network,
						"address": cfg.Reva.GRPC.Addr,
						"enabled_services": []string{
							"authprovider",      // provides basic auth
							"storageprovider",   // handles storage metadata
							"usershareprovider", // provides user shares
							"userprovider",      // provides user matadata (used to look up email, displayname etc after a login)
							"preferences",       // provides user preferences
							"gateway",           // to lookup services and authenticate requests
							"authregistry",      // used by the gateway to look up auth providers
							"storageregistry",   // used by the gateway to look up storage providers
						},
						"interceptors": map[string]interface{}{
							"auth": map[string]interface{}{
								"token_manager": "jwt",
								"token_managers": map[string]interface{}{
									"jwt": map[string]interface{}{
										"secret": cfg.Reva.JWTSecret,
									},
								},
								"skip_methods": []string{
									// we need to allow calls that happen during authentication
									"/cs3.gatewayv0alpha.GatewayService/Authenticate",
									"/cs3.gatewayv0alpha.GatewayService/WhoAmI",
									"/cs3.gatewayv0alpha.GatewayService/GetUser",
									"/cs3.gatewayv0alpha.GatewayService/ListAuthProviders",
									"/cs3.authregistryv0alpha.AuthRegistryService/ListAuthProviders",
									"/cs3.authregistryv0alpha.AuthRegistryService/GetAuthProvider",
									"/cs3.authproviderv0alpha.AuthProviderService/Authenticate",
									"/cs3.userproviderv0alpha.UserProviderService/GetUser",
								},
							},
						},
						"services": map[string]interface{}{
							"gateway": map[string]interface{}{
								"authregistrysvc":               cfg.Reva.GRPC.Addr,
								"storageregistrysvc":            cfg.Reva.GRPC.Addr,
								"appregistrysvc":                cfg.Reva.GRPC.Addr,
								"preferencessvc":                cfg.Reva.GRPC.Addr,
								"usershareprovidersvc":          cfg.Reva.GRPC.Addr,
								"publicshareprovidersvc":        cfg.Reva.GRPC.Addr,
								"ocmshareprovidersvc":           cfg.Reva.GRPC.Addr,
								"userprovidersvc":               cfg.Reva.GRPC.Addr,
								"commit_share_to_storage_grant": true,
								"datagateway":                   "http://" + cfg.Reva.HTTP.Addr + "/data",
								"transfer_shared_secret":        "replace-me-with-a-transfer-secret",
								"transfer_expires":              6, // give it a moment
								"token_manager":                 "jwt",
								"token_managers": map[string]interface{}{
									"jwt": map[string]interface{}{
										"secret": cfg.Reva.JWTSecret,
									},
								},
							},
							"authregistry": map[string]interface{}{
								"driver": "static",
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"rules": map[string]interface{}{
											//"basic": "localhost:9999",
											"oidc": cfg.Reva.GRPC.Addr,
										},
									},
								},
							},
							"storageregistry": map[string]interface{}{
								"driver": "static",
								"drivers": map[string]interface{}{
									"static": map[string]interface{}{
										"rules": map[string]interface{}{
											"/":                                    cfg.Reva.GRPC.Addr,
											"123e4567-e89b-12d3-a456-426655440000": cfg.Reva.GRPC.Addr,
										},
									},
								},
							},
							"authprovider": map[string]interface{}{
								"auth_manager": "oidc",
								"auth_managers": map[string]interface{}{
									"oidc": map[string]interface{}{
										"provider": cfg.AuthProvider.Provider,
										"insecure": cfg.AuthProvider.Insecure,
									},
								},
								"userprovidersvc": cfg.Reva.GRPC.Addr,
							},
							"userprovider": map[string]interface{}{
								"driver": "demo", // TODO use graph api
								/*
									"drivers": map[string]interface{}{
										"graph": map[string]interface{}{
											"provider": cfg.AuthProvider.Provider,
											"insecure": cfg.AuthProvider.Insecure,
										},
									},
								*/
							},
							"usershareprovider": map[string]interface{}{
								"driver": "memory",
							},
							"storageprovider": map[string]interface{}{
								"mount_path":         "/",
								"mount_id":           "123e4567-e89b-12d3-a456-426655440000",
								"data_server_url":    "http://" + cfg.Reva.HTTP.Addr + "/data",
								"expose_data_server": true,
								"available_checksums": map[string]interface{}{
									"md5":   100,
									"unset": 1000,
								},
								"driver": "owncloud",
								"drivers": map[string]interface{}{
									"owncloud": map[string]interface{}{
										"datadirectory": "/var/tmp/reva/data",
									},
								},
							},
						},
					},
				}
				gr.Add(func() error {
					// TODO micro knows nothing about reva
					runtime.Run(rcfg, pidFile)
					return nil
				}, func(_ error) {
					logger.Info().
						Str("server", "reva").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
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
