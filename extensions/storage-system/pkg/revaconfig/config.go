package revaconfig

import (
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/extensions/storage-system/pkg/config"
)

// StorageMetadataFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StorageMetadataFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					// registries are located on the gateway
					"authregistrysvc":    cfg.GRPC.Addr,
					"storageregistrysvc": cfg.GRPC.Addr,
					// user metadata is located on the users services
					"userprovidersvc":  cfg.GRPC.Addr,
					"groupprovidersvc": cfg.GRPC.Addr,
					"permissionssvc":   cfg.GRPC.Addr,
					// other
					"disable_home_creation_on_login": true, // metadata manually creates a space
					// metadata always uses the simple upload, so no transfer secret or datagateway needed
				},
				"userprovider": map[string]interface{}{
					"driver": "memory",
					"drivers": map[string]interface{}{
						"memory": map[string]interface{}{
							"users": map[string]interface{}{
								"serviceuser": map[string]interface{}{
									"id": map[string]interface{}{
										"opaqueId": cfg.SystemUserID,
										"idp":      "internal",
										"type":     userpb.UserType_USER_TYPE_PRIMARY,
									},
									"username":     "serviceuser",
									"display_name": "System User",
								},
							},
						},
					},
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"machine": cfg.GRPC.Addr,
							},
						},
					},
				},
				"authprovider": map[string]interface{}{
					"auth_manager": "machine",
					"auth_managers": map[string]interface{}{
						"machine": map[string]interface{}{
							"api_key":      cfg.SystemUserAPIKey,
							"gateway_addr": cfg.GRPC.Addr,
						},
					},
				},
				"permissions": map[string]interface{}{
					"driver": "demo",
					"drivers": map[string]interface{}{
						"demo": map[string]interface{}{},
					},
				},
				"storageregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"/": map[string]interface{}{
									"address": cfg.GRPC.Addr,
								},
							},
						},
					},
				},
				"storageprovider": map[string]interface{}{
					"driver":          cfg.Driver,
					"drivers":         metadataDrivers(cfg),
					"data_server_url": cfg.DataServerURL,
					"tmp_folder":      cfg.TempFolder,
				},
			},
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			// no datagateway needed as the metadata clients directly talk to the dataprovider with the simple protocol
			"services": map[string]interface{}{
				"dataprovider": map[string]interface{}{
					"prefix":      "data",
					"driver":      cfg.Driver,
					"drivers":     metadataDrivers(cfg),
					"timeout":     86400,
					"insecure":    cfg.DataProviderInsecure,
					"disable_tus": true,
				},
			},
		},
	}
	return rcfg
}

func metadataDrivers(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"ocis": map[string]interface{}{
			"root":                cfg.Drivers.OCIS.Root,
			"user_layout":         "{{.Id.OpaqueId}}",
			"treetime_accounting": false,
			"treesize_accounting": false,
			"permissionssvc":      cfg.GRPC.Addr,
		},
	}
}
