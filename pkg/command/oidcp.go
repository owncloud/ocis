package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/command"
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis-reva/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// OIDCPCommand is the entrypoint for the revaoidcp command.
func OIDCPCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "oidcp",
		Usage:    "Start openid connect provider",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Reva),
		Action: func(c *cli.Context) error {
			scfg := configureOIDCProvider(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureOIDCProvider(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Reva.Configs = map[string]interface{}{
		"oidcprovider": map[string]interface{}{
			"core": map[string]interface{}{
				"max_cpus":             cfg.Reva.Reva.MaxCPUs,
				"tracing_enabled":      cfg.Reva.Tracing.Enabled,
				"tracing_endpoint":     cfg.Reva.Tracing.Endpoint,
				"tracing_collector":    cfg.Reva.Tracing.Collector,
				"tracing_service_name": cfg.Reva.Tracing.Service,
			},
			"log": map[string]interface{}{
				"level": cfg.Reva.Reva.LogLevel,
				//TODO mode = "console" # "console" or "json"
				//TODO output = "./standalone.log"
			},
			"http": map[string]interface{}{
				"network": cfg.Reva.Reva.HTTP.Network,
				"address": cfg.Reva.Reva.HTTP.Addr,
				"enabled_services": []string{
					"oidcprovider",
					"wellknown",
					"prometheus",
					"ocs", // TODO remove when phoenix no longer tries to fetch the capabilities
				},
				"enabled_middlewares": []string{
					"cors",
					"auth",
				},
				"middlewares": map[string]interface{}{
					"auth": map[string]interface{}{
						"gateway":          cfg.Reva.Reva.GRPC.Addr,
						"credential_chain": []string{"basic", "bearer"},
						"token_strategy":   "header",
						"token_writer":     "header",
						"token_manager":    "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
						"skip_methods": []string{
							"/favicon.ico",
							"/oauth2",
							"/.well-known",
							"/metrics", // for prometheus metrics
						},
					},
					"cors": map[string]interface{}{
						"allowed_origins": []string{"*"},
						"allowed_methods": []string{
							"OPTIONS",
							"GET",
							"PUT",
							"POST",
							"DELETE",
							"MKCOL",
							"PROPFIND",
							"PROPPATCH",
							"MOVE",
							"COPY",
							"REPORT",
							"SEARCH",
						},
						"allowed_headers": []string{
							"Origin",
							"Accept",
							"Depth",
							"Content-Type",
							"X-Requested-With",
							"Authorization",
							"Ocs-Apirequest",
							"If-Match",
							"If-None-Match",
							"Destination",
							"Overwrite",
						},
						"allow_credentials":   true,
						"options_passthrough": false,
					},
				},
				"services": map[string]interface{}{
					// TODO investigate: service must be here as well, otherwise eg wellknown won't get started
					// TODO we hardcoded the url because the config option is used to tell the server which address to bind to,
					//      which is 0.0.0.0:9135 by default, but the iss needs to use a hostname
					"wellknown": map[string]interface{}{
						"prefix":                 ".well-known",
						"issuer":                 "http://localhost:9135",
						"authorization_endpoint": "http://localhost:9135/oauth2/auth",
						"token_endpoint":         "http://localhost:9135/oauth2/token",
						"revocation_endpoint":    "http://localhost:9135/oauth2/auth",
						"introspection_endpoint": "http://localhost:9135/oauth2/introspect",
						"userinfo_endpoint":      "http://localhost:9135/oauth2/userinfo",
					},
					"oidcprovider": map[string]interface{}{
						"prefix":    "oauth2",
						"gateway":   cfg.Reva.Reva.GRPC.Addr,
						"auth_type": "basic",
						// TODO we hardcoded the url because the config option is used to tell the server which address to bind to,
						//      which is 0.0.0.0:9135 by default, but the iss needs to use a hostname
						"issuer": "http://localhost:9135",
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
					"ocs": map[string]interface{}{
						"prefix":  "ocs",
						"gateway": cfg.Reva.Reva.GRPC.Addr,
					},
				},
			},
			"grpc": map[string]interface{}{
				"network": cfg.Reva.Reva.GRPC.Network,
				"address": cfg.Reva.Reva.GRPC.Addr,
				"enabled_services": []string{
					"authprovider", // provides basic auth
					"userprovider", // provides user matadata (used to look up email, displayname etc after a login)
					"gateway",      // to lookup services and authenticate requests
					"authregistry", // used by the gateway to look up auth providers
				},
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
						"skip_methods": []string{
							// we need to allow calls that happen during authentication
							"/cs3.auth.registry.v1beta1.RegistryAPI/GetAuthProvider",
							"/cs3.auth.provider.v1beta1.ProviderAPI/Authenticate",
							"/cs3.gateway.v1beta1.GatewayAPI/Authenticate",
							"/cs3.identity.user.v1beta1.UserAPI/GetUser",
							"/cs3.gateway.v1beta1.GatewayAPI/GetUser",
						},
					},
				},
				"services": map[string]interface{}{
					"gateway": map[string]interface{}{
						"authregistrysvc":               cfg.Reva.Reva.GRPC.Addr,
						"storageregistrysvc":            cfg.Reva.Reva.GRPC.Addr,
						"appregistrysvc":                cfg.Reva.Reva.GRPC.Addr,
						"preferencessvc":                cfg.Reva.Reva.GRPC.Addr,
						"usershareprovidersvc":          cfg.Reva.Reva.GRPC.Addr,
						"publicshareprovidersvc":        cfg.Reva.Reva.GRPC.Addr,
						"ocmshareprovidersvc":           cfg.Reva.Reva.GRPC.Addr,
						"userprovidersvc":               cfg.Reva.Reva.GRPC.Addr,
						"commit_share_to_storage_grant": true,
						"datagateway":                   "http://" + cfg.Reva.Reva.HTTP.Addr + "/data",
						"transfer_shared_secret":        "replace-me-with-a-transfer-secret",
						"transfer_expires":              6, // give it a moment
						"token_manager":                 "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
					"authregistry": map[string]interface{}{
						"driver": "static",
						"drivers": map[string]interface{}{
							"static": map[string]interface{}{
								"rules": map[string]interface{}{
									"basic":  cfg.Reva.Reva.GRPC.Addr,
									"bearer": "localhost:9138",
								},
							},
						},
					},
					"authprovider": map[string]interface{}{
						"auth_manager":    "demo",
						"userprovidersvc": cfg.Reva.Reva.GRPC.Addr,
					},
					"userprovider": map[string]interface{}{
						"driver": "demo",
					},
				},
			},
		},
		"oidcauthprovider": map[string]interface{}{
			"core": map[string]interface{}{
				"max_cpus":             cfg.Reva.Reva.MaxCPUs,
				"tracing_enabled":      cfg.Reva.Tracing.Enabled,
				"tracing_endpoint":     cfg.Reva.Tracing.Endpoint,
				"tracing_collector":    cfg.Reva.Tracing.Collector,
				"tracing_service_name": cfg.Reva.Tracing.Service,
			},
			"log": map[string]interface{}{
				"level": cfg.Reva.Reva.LogLevel,
				//TODO mode = "console" # "console" or "json"
				//TODO output = "./standalone.log"
			},
			"grpc": map[string]interface{}{
				"network": cfg.Reva.Reva.GRPC.Network,
				"address": "localhost:9138", // use another port
				"enabled_services": []string{
					"authprovider", // provides oidc auth
				},
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
						"skip_methods": []string{
							"/cs3.auth.registry.v1beta1.RegistryAPI/GetAuthProvider",
							"/cs3.auth.provider.v1beta1.ProviderAPI/Authenticate",
						},
					},
				},
				"services": map[string]interface{}{
					"authprovider": map[string]interface{}{
						"auth_manager":    "oidc",
						"userprovidersvc": cfg.Reva.Reva.GRPC.Addr,
						"auth_managers": map[string]interface{}{
							"oidc": map[string]interface{}{
								// TODO we hardcoded the url because the config option is used to tell the server which address to bind to,
								//      which is 0.0.0.0:9135 by default, but the iss needs to use a hostname
								"issuer": "http://localhost:9135",
							},
						},
					},
				},
			},
		},
	}

	return cfg.Reva
}

func init() {
	register.AddCommand(OIDCPCommand)
}
