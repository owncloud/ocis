package revaconfig

import (
	"github.com/owncloud/ocis/v2/extensions/storage-users/pkg/config"
)

// StorageUsersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StorageUsersConfigFromStruct(cfg *config.Config) map[string]interface{} {
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
				"storageprovider": map[string]interface{}{
					"driver":             cfg.Driver,
					"drivers":            UserDrivers(cfg),
					"mount_id":           cfg.MountID,
					"expose_data_server": cfg.ExposeDataServer,
					"data_server_url":    cfg.DataServerURL,
					"tmp_folder":         cfg.TempFolder,
				},
			},
			"interceptors": map[string]interface{}{
				"eventsmiddleware": map[string]interface{}{
					"group":     "sharing",
					"type":      "nats",
					"address":   cfg.Events.Addr,
					"clusterID": cfg.Events.ClusterID,
				},
			},
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"dataprovider": map[string]interface{}{
					"prefix":      cfg.HTTP.Prefix,
					"driver":      cfg.Driver,
					"drivers":     UserDrivers(cfg),
					"timeout":     86400,
					"insecure":    cfg.DataProviderInsecure,
					"disable_tus": false,
				},
			},
		},
	}
	if cfg.ReadOnly {
		gcfg := rcfg["grpc"].(map[string]interface{})
		gcfg["interceptors"] = map[string]interface{}{
			"readonly": map[string]interface{}{},
		}
	}
	return rcfg
}
