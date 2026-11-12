package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/config"
)

// StoragePublicLinkConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StoragePublicLinkConfigFromStruct(cfg *config.Config) map[string]any {
	rcfg := map[string]any{
		"shared": map[string]any{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
			"grpc_client_options":       cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]any{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]any{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			"interceptors": map[string]any{
				"log": map[string]any{},
				"prometheus": map[string]any{
					"namespace": "ocis",
					"subsystem": "storage_publiclink",
				},
			},
			"services": map[string]any{
				"publicstorageprovider": map[string]any{
					"mount_id":     cfg.StorageProvider.MountID,
					"gateway_addr": cfg.Reva.Address,
				},
				"authprovider": map[string]any{
					"auth_manager": "publicshares",
					"auth_managers": map[string]any{
						"publicshares": map[string]any{
							"gateway_addr":                  cfg.Reva.Address,
							"brute_force_time_gap":          cfg.BruteForce.TimeGap,
							"brute_force_max_attempts":      cfg.BruteForce.MaxAttempts,
							"store_type":                    cfg.Store.Store,
							"store_nodes":                   cfg.Store.Nodes,
							"store_database":                cfg.Store.Database,
							"store_table":                   cfg.Store.Table,
							"store_username":                cfg.Store.AuthUsername,
							"store_password":                cfg.Store.AuthPassword,
							"store_creds_file":              cfg.Store.CredsFile,
							"store_creds":                   cfg.Store.Creds,
							"store_enable_tls":              cfg.Store.EnableTLS,
							"store_tls_insecure":            cfg.Store.TLSInsecure,
							"store_tls_root_ca_certificate": cfg.Store.TLSRootCACertificate,
						},
					},
				},
			},
		},
	}
	return rcfg
}
