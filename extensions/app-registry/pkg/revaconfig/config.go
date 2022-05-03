package revaconfig

import (
	"github.com/owncloud/ocis/ocis-pkg/log"

	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/ocis/extensions/app-registry/pkg/config"
)

// AppRegistryConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AppRegistryConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret": cfg.TokenManager.JWTSecret,
			"gatewaysvc": cfg.Reva.Address,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
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
