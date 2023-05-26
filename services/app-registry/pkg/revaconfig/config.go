package revaconfig

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"

	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/ocis/v2/services/app-registry/pkg/config"
)

// AppRegistryConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AppRegistryConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_exporter":     cfg.Tracing.Type,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
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
			"services": map[string]interface{}{
				"appregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"mime_types": mimetypes(cfg, logger),
						},
					},
				},
			},
			"interceptors": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "app_registry",
				},
			},
		},
	}
	return rcfg
}

func mimetypes(cfg *config.Config, logger log.Logger) []map[string]interface{} {
	var m []map[string]interface{}
	if err := mapstructure.Decode(cfg.AppRegistry.MimeTypeConfig, &m); err != nil {
		logger.Error().Err(err).Msg("Failed to decode appregistry mimetypes to mapstructure")
		return nil
	}
	return m
}
