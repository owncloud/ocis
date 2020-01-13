// +build !simple

package command

import (
	svcconfig "github.com/owncloud/ocis-reva/pkg/config"
	"github.com/owncloud/ocis/pkg/config"
)

func configureReva(cfg *config.Config) *svcconfig.Config {
	cfg.Reva.Log.Level = cfg.Log.Level
	cfg.Reva.Log.Pretty = cfg.Log.Pretty
	cfg.Reva.Log.Color = cfg.Log.Color

	// reva will use a global namespace, for that to work we need to start three storage providers
	// - / for the root listing
	// - /home which always points to the users home dir
	// - /oc for a global view
	// in addition to that we need a frontend with
	// - an openid connect provider
	// - the webdav api
	// - the ocs api
	// - the .well-know endpoint for oidc configuration
	cfg.Reva.Reva.Configs = map[string]interface{}{
		"frontend": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:20099",
				"services": map[string]interface{}{
					"authprovider": map[string]interface{}{
						"auth_manager": "oidc",
						"auth_managers": map[string]interface{}{
							"oidc": map[string]interface{}{
								"issuer": "http://localhost:20080",
							},
						},
					},
				},
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
			},
			"http": map[string]interface{}{
				"address": "0.0.0.0:20080",
				"middlewares": map[string]interface{}{
					"auth": map[string]interface{}{
						"gateway":          "localhost:19000",
						"credential_chain": []string{"basic", "bearer"},
						"token_strategy":   "header",
						"token_writer":     "header",
						"token_manager":    "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
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
					"wellknown": map[string]interface{}{
						"issuer":                 "http://localhost:20080",
						"authorization_endpoint": "http://localhost:20080/oauth2/auth",
						"token_endpoint":         "http://localhost:20080/oauth2/token",
						"revocation_endpoint":    "http://localhost:20080/oauth2/auth",
						"introspection_endpoint": "http://localhost:20080/oauth2/introspect",
						"userinfo_endpoint":      "http://localhost:20080/oauth2/userinfo",
					},
					"oidcprovider": map[string]interface{}{
						"prefix":  "oauth2",
						"gateway": "localhost:19000",
						"issuer":  "http://localhost:20080",
						"clients": map[string]interface{}{
							"phoenix": map[string]interface{}{
								"id":             "phoenix",
								"redirect_uris":  []string{"http://localhost:9100/oidc-callback.html", "http://localhost:9100/"},
								"grant_types":    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
								"response_types": []string{"code"}, // use authorization code flow, see https://developer.okta.com/blog/2019/05/01/is-the-oauth-implicit-flow-dead for details
								"scopes":         []string{"openid", "profile", "email", "offline"},
								"public":         true, // force PKCS for public clients
							},
							// TODO add cli command for token fetching
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
					"ocdav": map[string]interface{}{
						"prefix":           "",
						"chunk_folder":     "/var/tmp/revad/chunks",
						"gateway":          "localhost:19000",
						"files_namespace":  "/",
						"webdav_namespace": "/",
					},
					"ocs": map[string]interface{}{
						"gateway": "localhost:19000",
						"config": map[string]interface{}{
							"version": "1.8",
							"website": "reva",
							"host":    "http://localhost:20080",
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
		},
		"gateway": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:19000",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"gateway": map[string]interface{}{
						// registries
						"authregistrysvc":    "localhost:19000",
						"storageregistrysvc": "localhost:19000",
						"appregistrysvc":     "localhost:19000",
						// user metadata
						"preferencessvc":  "localhost:18000",
						"userprovidersvc": "localhost:18000",
						// sharing
						"usershareprovidersvc":          "localhost:17000",
						"publicshareprovidersvc":        "localhost:17000",
						"ocmshareprovidersvc":           "localhost:17000",
						"commit_share_to_storage_grant": true,
						// other
						"datagateway":            "http://localhost:19001/data",
						"transfer_shared_secret": "replace-me-with-a-transfer-secret",
						"transfer_expires":       6, // give it a moment
						"token_manager":          "jwt",
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
									"basic":  "localhost:18000",
									"bearer": "localhost:20099",
								},
							},
						},
					},
					"storageregistry": map[string]interface{}{
						"driver": "static",
						"drivers": map[string]interface{}{
							"static": map[string]interface{}{
								"rules": map[string]interface{}{
									"/home":                                "localhost:12000",
									"/oc":                                  "localhost:11000",
									"123e4567-e89b-12d3-a456-426655440000": "localhost:11000",
									"/":                                    "localhost:11100",
									"123e4567-e89b-12d3-a456-426655440001": "localhost:11100",
								},
							},
						},
					},
				},
			},
			"http": map[string]interface{}{
				"address": "0.0.0.0:19001",
				"middlewares": map[string]interface{}{
					"auth": map[string]interface{}{
						"gateway":          "localhost:19000",
						"credential_chain": []string{"basic", "bearer"},
						"token_strategy":   "header",
						"token_writer":     "header",
						"token_manager":    "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"datagateway": map[string]interface{}{
						"prefix":                 "data",
						"gateway":                "",
						"transfer_shared_secret": "replace-me-with-a-transfer-secret",
					},
				},
			},
		},
		"shares": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:17000",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"usershareprovider": map[string]interface{}{
						"driver": "memory",
					},
					"publicshareprovider": map[string]interface{}{
						"driver": "memory",
					},
				},
			},
		},
		"storage-home": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:12000",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"storageprovider": map[string]interface{}{
						"driver": "owncloud",
						"drivers": map[string]interface{}{
							"owncloud": map[string]interface{}{
								"datadirectory": "/var/tmp/reva/data",
							},
						},
						"path_wrapper": "context",
						"path_wrappers": map[string]interface{}{
							"context": map[string]interface{}{
								"prefix": "",
							},
						},
						"mount_path":         "/home",
						"mount_id":           "123e4567-e89b-12d3-a456-426655440000",
						"expose_data_server": true,
						"data_server_url":    "http://localhost:12001/data",
						"available_checksums": map[string]interface{}{
							"md5":   100,
							"unset": 1000,
						},
					},
				},
			},
			"http": map[string]interface{}{
				"address": "0.0.0.0:12001",
				"middlewares": map[string]interface{}{
					"auth": map[string]interface{}{
						"gateway":          "localhost:19000",
						"credential_chain": []string{"basic", "bearer"},
						"token_strategy":   "header",
						"token_writer":     "header",
						"token_manager":    "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"dataprovider": map[string]interface{}{
						"prefix": "data",
						"driver": "owncloud",
						"drivers": map[string]interface{}{
							"owncloud": map[string]interface{}{
								"datadirectory": "/var/tmp/reva/data",
							},
						},
						"temp_folder": "/var/tmp/",
					},
				},
			},
		},
		"storage-oc": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:11000",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"storageprovider": map[string]interface{}{
						"driver": "owncloud",
						"drivers": map[string]interface{}{
							"owncloud": map[string]interface{}{
								"datadirectory": "/var/tmp/reva/data",
							},
						},
						"mount_path":         "/oc",
						"mount_id":           "123e4567-e89b-12d3-a456-426655440000",
						"expose_data_server": true,
						"data_server_url":    "http://localhost:11001/data",
						"available_checksums": map[string]interface{}{
							"md5":   100,
							"unset": 1000,
						},
					},
				},
			},
			"http": map[string]interface{}{
				"address": "0.0.0.0:11001",
				"middlewares": map[string]interface{}{
					"auth": map[string]interface{}{
						"gateway":          "localhost:19000",
						"credential_chain": []string{"basic", "bearer"},
						"token_strategy":   "header",
						"token_writer":     "header",
						"token_manager":    "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"dataprovider": map[string]interface{}{
						"prefix": "data",
						"driver": "owncloud",
						"drivers": map[string]interface{}{
							"owncloud": map[string]interface{}{
								"datadirectory": "/var/tmp/reva/data",
							},
						},
						"temp_folder": "/var/tmp/",
					},
				},
			},
		},
		"storage-root": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:11100",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"storageprovider": map[string]interface{}{
						"driver": "local",
						"drivers": map[string]interface{}{
							"local": map[string]interface{}{
								"root": "/var/tmp/reva/root",
							},
						},
						"mount_path": "/",
						"mount_id":   "123e4567-e89b-12d3-a456-426655440001",
						"available_checksums": map[string]interface{}{
							"md5":   100,
							"unset": 1000,
						},
					},
				},
			},
		},
		"users": map[string]interface{}{
			"grpc": map[string]interface{}{
				"address": "0.0.0.0:18000",
				"interceptors": map[string]interface{}{
					"auth": map[string]interface{}{
						"token_manager": "jwt",
						"token_managers": map[string]interface{}{
							"jwt": map[string]interface{}{
								"secret": cfg.Reva.Reva.JWTSecret,
							},
						},
					},
				},
				"services": map[string]interface{}{
					"authprovider": map[string]interface{}{
						"auth_manager": "demo",
					},
					"userprovider": map[string]interface{}{
						"driver": "demo",
					},
				},
			},
		},
	}

	return cfg.Reva
}
