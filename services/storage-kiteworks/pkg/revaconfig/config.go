// Package revaconfig contains the config for the reva service
package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
)

// StorageKiteworksConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StorageKiteworksConfigFromStruct(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"graceful_shutdown_timeout": cfg.GracefulShutdownTimeout,
		},
		"shared": map[string]interface{}{
			"jwt_secret":          cfg.TokenManager.JWTSecret,
			"gatewaysvc":          cfg.Reva.Address,
			"grpc_client_options": cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]interface{}{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			"interceptors": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "storage_kiteworks",
				},
			},
			"services": map[string]interface{}{
				"storageprovider": map[string]interface{}{
					"driver": "kiteworks",
					"drivers": map[string]interface{}{
						"kiteworks": map[string]interface{}{
							"endpoint":   cfg.Driver.Endpoint,
							"insecure":   cfg.Driver.Insecure,
							"chunk_size": cfg.Driver.ChunkSize,
						},
					},
					"mount_id": cfg.MountID,
				},
			},
		},
	}
}
