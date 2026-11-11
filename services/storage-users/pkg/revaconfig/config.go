// Package revaconfig contains the config for the reva service
package revaconfig

import (
	"strconv"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
)

// StorageUsersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StorageUsersConfigFromStruct(cfg *config.Config) map[string]any {
	// Credentials must not be combined with a wildcard origin.
	corsAllowCredentials := cfg.HTTP.CORS.AllowCredentials
	if corsAllowCredentials && cors.AllowsAnyOrigin(cfg.HTTP.CORS.AllowedOrigins) {
		corsAllowCredentials = false
	}

	rcfg := map[string]any{
		"core": map[string]any{
			"graceful_shutdown_timeout": cfg.GracefulShutdownTimeout,
		},
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
			// TODO build services dynamically
			"services": map[string]any{
				"storageprovider": map[string]any{
					"driver":             cfg.Driver,
					"drivers":            StorageProviderDrivers(cfg),
					"mount_id":           cfg.MountID,
					"expose_data_server": cfg.ExposeDataServer,
					"data_server_url":    cfg.DataServerURL,
					"upload_expiration":  cfg.UploadExpiration,
					"events": map[string]any{
						"nats_address":     cfg.Events.Endpoint,
						"nats_clusterid":   cfg.Events.Cluster,
						"tls_insecure":     cfg.Events.TLSInsecure,
						"tls_root_ca_cert": cfg.Events.TLSRootCACertificate,
						"nats_enable_tls":  cfg.Events.EnableTLS,
						"nats_username":    cfg.Events.AuthUsername,
						"nats_password":    cfg.Events.AuthPassword,
						"nats_creds_file":  cfg.Events.CredsFile,
						"nats_creds":       cfg.Events.Creds,
					},
				},
			},
			"interceptors": map[string]any{
				"eventsmiddleware": map[string]any{
					"group":            "sharing",
					"type":             "nats",
					"address":          cfg.Events.Endpoint,
					"clusterID":        cfg.Events.Cluster,
					"tls-insecure":     cfg.Events.TLSInsecure,
					"tls-root-ca-cert": cfg.Events.TLSRootCACertificate,
					"enable-tls":       cfg.Events.EnableTLS,
					"name":             generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus),
					"username":         cfg.Events.AuthUsername,
					"password":         cfg.Events.AuthPassword,
					"creds-file":       cfg.Events.CredsFile,
					"creds":            cfg.Events.Creds,
				},
				"prometheus": map[string]any{
					"namespace": "ocis",
					"subsystem": "storage_users",
				},
			},
		},
		"http": map[string]any{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			"middlewares": map[string]any{
				"requestid": map[string]any{},
			},
			// TODO build services dynamically
			"services": map[string]any{
				"dataprovider": map[string]any{
					"prefix":                 cfg.HTTP.Prefix,
					"driver":                 cfg.Driver,
					"drivers":                DataProviderDrivers(cfg),
					"nats_address":           cfg.Events.Endpoint,
					"nats_clusterID":         cfg.Events.Cluster,
					"mount_id":               cfg.MountID,
					"nats_tls_insecure":      cfg.Events.TLSInsecure,
					"nats_root_ca_cert_path": cfg.Events.TLSRootCACertificate,
					"nats_enable_tls":        cfg.Events.EnableTLS,
					"nats_username":          cfg.Events.AuthUsername,
					"nats_password":          cfg.Events.AuthPassword,
					"nats_creds_file":        cfg.Events.CredsFile,
					"nats_creds":             cfg.Events.Creds,
					"data_txs": map[string]any{
						"simple": map[string]any{
							"cache_store":    "noop",
							"cache_database": "system",
							"cache_table":    "stat",
						},
						"spaces": map[string]any{
							"cache_store":    "noop",
							"cache_database": "system",
							"cache_table":    "stat",
						},
						"tus": map[string]any{
							"cache_store":    "noop",
							"cache_database": "system",
							"cache_table":    "stat",
							"cors_enabled":   true,
							// allow_origin is configured as a regex in tusd, so we concatenate the configured values into a regex
							"cors_allow_origin":      "(" + strings.ReplaceAll(strings.Join(cfg.HTTP.CORS.AllowedOrigins, "|"), "*", ".*") + ")",
							"cors_allow_credentials": corsAllowCredentials,
							"cors_allow_methods":     strings.Join(cfg.HTTP.CORS.AllowedMethods, ","),
							"cors_allow_headers":     strings.Join(cfg.HTTP.CORS.AllowedHeaders, ","),
							"cors_max_age":           strconv.FormatUint(uint64(cfg.HTTP.CORS.MaxAge), 10),
							"cors_expose_headers":    strings.Join(cfg.HTTP.CORS.ExposedHeaders, ","),
						},
					},
				},
			},
		},
	}
	gcfg := rcfg["grpc"].(map[string]any)
	if cfg.ReadOnly {
		// Replace all interceptors with readonly when the storage is read-only.
		// eventsmiddleware and prometheus are intentionally dropped in this mode.
		gcfg["interceptors"] = map[string]any{
			"readonly": map[string]any{},
		}
	}
	if cfg.EnableVaultMode {
		// Set mfa_enabled inside the auth interceptor config so that all gRPC
		// calls to this vault storage-users instance require MFA authentication.
		interceptors := gcfg["interceptors"].(map[string]any)
		if authCfg, ok := interceptors["auth"].(map[string]any); ok {
			authCfg["mfa_enabled"] = true
		} else {
			interceptors["auth"] = map[string]any{
				"mfa_enabled": true,
			}
		}
	}
	return rcfg
}
