package revaconfig

import (
	"path"
	"strconv"

	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config"
)

// FrontendConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func FrontendConfigFromStruct(cfg *config.Config) map[string]interface{} {
	archivers := []map[string]interface{}{
		{
			"enabled":       true,
			"version":       "2.0.0",
			"formats":       []string{"tar", "zip"},
			"archiver_url":  path.Join("/", cfg.Archiver.Prefix),
			"max_num_files": strconv.FormatInt(cfg.Archiver.MaxNumFiles, 10),
			"max_size":      strconv.FormatInt(cfg.Archiver.MaxSize, 10),
		},
	}

	appProviders := []map[string]interface{}{
		{
			"enabled":  true,
			"version":  "1.0.0",
			"apps_url": "/app/list",
			"open_url": "/app/open",
			"new_url":  "/app/new",
		},
	}

	filesCfg := map[string]interface{}{
		"private_links":     false,
		"bigfilechunking":   false,
		"blacklisted_files": []string{},
		"undelete":          true,
		"versioning":        true,
		"archivers":         archivers,
		"app_providers":     appProviders,
		"favorites":         cfg.EnableFavorites,
	}

	if cfg.DefaultUploadProtocol == "tus" {
		filesCfg["tus_support"] = map[string]interface{}{
			"version":              "1.0.0",
			"resumable":            "1.0.0",
			"extension":            "creation,creation-with-upload",
			"http_method_override": cfg.UploadHTTPMethodOverride,
			"max_chunk_size":       cfg.UploadMaxChunkSize,
		}
	}

	return map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address, // Todo or address?
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			"middlewares": map[string]interface{}{
				"cors": map[string]interface{}{
					"allow_credentials": true,
				},
				"auth": map[string]interface{}{
					"credentials_by_user_agent": cfg.Middleware.Auth.CredentialsByUserAgent,
					"credential_chain":          []string{"bearer"},
				},
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				// this reva service called "appprovider" comes from
				// `internal/http/services/appprovider` and is a translation
				// layer from the grpc app registry to http, used by eg. ownCloud Web
				// It should not be confused with `internal/grpc/services/appprovider`
				// which is currently only has only the driver for the CS3org WOPI server
				"appprovider": map[string]interface{}{
					"prefix":                 cfg.AppHandler.Prefix,
					"transfer_shared_secret": cfg.TransferSecret,
					"timeout":                86400,
					"insecure":               cfg.AppHandler.Insecure,
				},
				"archiver": map[string]interface{}{
					"prefix":        cfg.Archiver.Prefix,
					"timeout":       86400,
					"insecure":      cfg.Archiver.Insecure,
					"max_num_files": cfg.Archiver.MaxNumFiles,
					"max_size":      cfg.Archiver.MaxSize,
				},
				"datagateway": map[string]interface{}{
					"prefix":                 cfg.DataGateway.Prefix,
					"transfer_shared_secret": cfg.TransferSecret,
					"timeout":                86400,
					"insecure":               true,
				},
				"ocs": map[string]interface{}{
					"storage_registry_svc":      cfg.Reva.Address,
					"share_prefix":              cfg.OCS.SharePrefix,
					"home_namespace":            cfg.OCS.HomeNamespace,
					"resource_info_cache_ttl":   cfg.OCS.ResourceInfoCacheTTL,
					"prefix":                    cfg.OCS.Prefix,
					"additional_info_attribute": cfg.OCS.AdditionalInfoAttribute,
					"machine_auth_apikey":       cfg.MachineAuthAPIKey,
					"cache_warmup_driver":       cfg.OCS.CacheWarmupDriver,
					"cache_warmup_drivers": map[string]interface{}{
						"cbox": map[string]interface{}{
							"db_username": cfg.OCS.CacheWarmupDrivers.CBOX.DBUsername,
							"db_password": cfg.OCS.CacheWarmupDrivers.CBOX.DBPassword,
							"db_host":     cfg.OCS.CacheWarmupDrivers.CBOX.DBHost,
							"db_port":     cfg.OCS.CacheWarmupDrivers.CBOX.DBPort,
							"db_name":     cfg.OCS.CacheWarmupDrivers.CBOX.DBName,
							"namespace":   cfg.OCS.CacheWarmupDrivers.CBOX.Namespace,
							"gatewaysvc":  cfg.Reva.Address,
						},
					},
					"config": map[string]interface{}{
						"version": "1.7",
						"website": "ownCloud",
						"host":    cfg.PublicURL,
						"contact": "",
						"ssl":     "false",
					},
					"default_upload_protocol": cfg.DefaultUploadProtocol,
					"capabilities": map[string]interface{}{
						"capabilities": map[string]interface{}{
							"core": map[string]interface{}{
								"poll_interval": 60,
								"webdav_root":   "remote.php/webdav",
								"status": map[string]interface{}{
									"installed":      true,
									"maintenance":    false,
									"needsDbUpgrade": false,
									"version":        version.Legacy,
									"versionstring":  version.LegacyString,
									"edition":        "Community",
									"productname":    "Infinite Scale",
									"product":        "Infinite Scale",
									"productversion": version.GetString(),
									"hostname":       "",
								},
								"support_url_signing": true,
							},
							"checksums": map[string]interface{}{
								"supported_types":       cfg.Checksums.SupportedTypes,
								"preferred_upload_type": cfg.Checksums.PreferredUploadType,
							},
							"files": filesCfg,
							"dav": map[string]interface{}{
								"reports": []string{"search-files"},
							},
							"files_sharing": map[string]interface{}{
								"api_enabled":                       true,
								"resharing":                         cfg.EnableResharing,
								"group_sharing":                     true,
								"auto_accept_share":                 true,
								"share_with_group_members_only":     true,
								"share_with_membership_groups_only": true,
								"default_permissions":               22,
								"search_min_length":                 cfg.SearchMinLength,
								"public": map[string]interface{}{
									"alias":                      false,
									"enabled":                    true,
									"send_mail":                  true,
									"defaultPublicLinkShareName": "Public link",
									"social_share":               true,
									"upload":                     true,
									"multiple":                   true,
									"supports_upload_only":       true,
									"password": map[string]interface{}{
										"enforced": false,
										"enforced_for": map[string]interface{}{
											"read_only":   false,
											"read_write":  false,
											"upload_only": false,
										},
									},
									"expire_date": map[string]interface{}{
										"enabled": true,
									},
									"can_edit": true,
								},
								"user": map[string]interface{}{
									"send_mail":       true,
									"profile_picture": false,
									"settings": []map[string]interface{}{
										{
											"enabled": true,
											"version": "1.0.0",
										},
									},
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
							"spaces": map[string]interface{}{
								"version":    "0.0.1",
								"enabled":    cfg.EnableProjectSpaces || cfg.EnableShareJail,
								"projects":   cfg.EnableProjectSpaces,
								"share_jail": cfg.EnableShareJail,
							},
						},
						"version": map[string]interface{}{
							"product":        "Infinite Scale",
							"edition":        "Community",
							"major":          version.ParsedLegacy().Major(),
							"minor":          version.ParsedLegacy().Minor(),
							"micro":          version.ParsedLegacy().Patch(),
							"string":         version.LegacyString,
							"productversion": version.GetString(),
						},
					},
				},
			},
		},
	}
}
