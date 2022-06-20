package revaconfig

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"

	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/config"
)

// GatewayConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func GatewayConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]interface{} {
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
			// TODO build services dynamically
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					// registries is located on the gateway
					"authregistrysvc":    cfg.Reva.Address,
					"storageregistrysvc": cfg.Reva.Address,
					"appregistrysvc":     cfg.AppRegistryEndpoint,
					// user metadata is located on the users services
					"preferencessvc":   cfg.UsersEndpoint,
					"userprovidersvc":  cfg.UsersEndpoint,
					"groupprovidersvc": cfg.GroupsEndpoint,
					"permissionssvc":   cfg.PermissionsEndpoint,
					// sharing is located on the sharing service
					"usershareprovidersvc":          cfg.SharingEndpoint,
					"publicshareprovidersvc":        cfg.SharingEndpoint,
					"ocmshareprovidersvc":           cfg.SharingEndpoint,
					"commit_share_to_storage_grant": cfg.CommitShareToStorageGrant,
					"commit_share_to_storage_ref":   cfg.CommitShareToStorageRef,
					"share_folder":                  cfg.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.
					// other
					"disable_home_creation_on_login": cfg.DisableHomeCreationOnLogin,
					"datagateway":                    strings.TrimRight(cfg.FrontendPublicURL, "/") + "/data",
					"transfer_shared_secret":         cfg.TransferSecret,
					"transfer_expires":               cfg.TransferExpires,
					"etag_cache_ttl":                 cfg.EtagCacheTTL,
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"basic":        cfg.AuthBasicEndpoint,
								"bearer":       cfg.AuthBearerEndpoint,
								"machine":      cfg.AuthMachineEndpoint,
								"publicshares": cfg.StoragePublicLinkEndpoint,
							},
						},
					},
				},
				"storageregistry": map[string]interface{}{
					"driver": cfg.StorageRegistry.Driver,
					"drivers": map[string]interface{}{
						"spaces": map[string]interface{}{
							"providers": spacesProviders(cfg, logger),
						},
					},
				},
			},
		},
	}
	return rcfg
}

func spacesProviders(cfg *config.Config, logger log.Logger) map[string]map[string]interface{} {

	// if a list of rules is given it overrides the generated rules from below
	if len(cfg.StorageRegistry.Rules) > 0 {
		rules := map[string]map[string]interface{}{}
		for i := range cfg.StorageRegistry.Rules {
			parts := strings.SplitN(cfg.StorageRegistry.Rules[i], "=", 2)
			rules[parts[0]] = map[string]interface{}{"address": parts[1]}
		}
		return rules
	}

	// check if the rules have to be read from a json file
	if cfg.StorageRegistry.JSON != "" {
		data, err := ioutil.ReadFile(cfg.StorageRegistry.JSON)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read storage registry rules from JSON file: " + cfg.StorageRegistry.JSON)
			return nil
		}
		var rules map[string]map[string]interface{}
		if err = json.Unmarshal(data, &rules); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal storage registry rules")
			return nil
		}
		return rules
	}

	// generate rules based on default config
	return map[string]map[string]interface{}{
		cfg.StorageUsersEndpoint: {
			"providerid": "1284d238-aa92-42ce-bdc4-0b0000009157",
			"spaces": map[string]interface{}{
				"personal": map[string]interface{}{
					"mount_point":   "/users",
					"path_template": "/users/{{.Space.Owner.Id.OpaqueId}}",
				},
				"project": map[string]interface{}{
					"mount_point":   "/projects",
					"path_template": "/projects/{{.Space.Name}}",
				},
			},
		},
		cfg.StorageSharesEndpoint: {
			"providerid": utils.ShareStorageProviderID,
			"spaces": map[string]interface{}{
				"virtual": map[string]interface{}{
					// The root of the share jail is mounted here
					"mount_point": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares",
				},
				"grant": map[string]interface{}{
					// Grants are relative to a space root that the gateway will determine with a stat
					"mount_point": ".",
				},
				"mountpoint": map[string]interface{}{
					// The jail needs to be filled with mount points
					// .Space.Name is a path relative to the mount point
					"mount_point":   "/users/{{.CurrentUser.Id.OpaqueId}}/Shares",
					"path_template": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares/{{.Space.Name}}",
				},
			},
		},
		// public link storage returns the mount id of the actual storage
		cfg.StoragePublicLinkEndpoint: {
			"providerid": utils.PublicStorageProviderID,
			"spaces": map[string]interface{}{
				"grant": map[string]interface{}{
					"mount_point": ".",
				},
				"mountpoint": map[string]interface{}{
					"mount_point":   "/public",
					"path_template": "/public/{{.Space.Root.OpaqueId}}",
				},
			},
		},
		// medatada storage not part of the global namespace
	}
}
