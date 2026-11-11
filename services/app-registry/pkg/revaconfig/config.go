package revaconfig

import (
	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/app-registry/pkg/config"
)

// AppRegistryConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AppRegistryConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]any {
	rcfg := map[string]any{
		"core": map[string]any{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_exporter":     cfg.Tracing.Type,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]any{
			"jwt_secret":          cfg.TokenManager.JWTSecret,
			"gatewaysvc":          cfg.Reva.Address,
			"grpc_client_options": cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]any{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]any{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			"services": map[string]any{
				"appregistry": map[string]any{
					"driver": "static",
					"drivers": map[string]any{
						"static": map[string]any{
							"mime_types": mimetypes(cfg, logger),
						},
					},
				},
			},
			"interceptors": map[string]any{
				"prometheus": map[string]any{
					"namespace": "ocis",
					"subsystem": "app_registry",
				},
			},
		},
	}
	return rcfg
}

func mimetypes(cfg *config.Config, logger log.Logger) []map[string]any {
	var m []map[string]any
	if err := mapstructure.Decode(cfg.AppRegistry.MimeTypeConfig, &m); err != nil {
		logger.Error().Err(err).Msg("Failed to decode appregistry mimetypes to mapstructure")
		return nil
	}
	return m
}
