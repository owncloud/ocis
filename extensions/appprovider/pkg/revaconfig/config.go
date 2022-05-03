package revaconfig

import (
	"github.com/owncloud/ocis/extensions/appprovider/pkg/config"
)

// AppProviderConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AppProviderConfigFromStruct(cfg *config.Config) map[string]interface{} {

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
				"appprovider": map[string]interface{}{
					"app_provider_url": cfg.ExternalAddr,
					"driver":           cfg.Driver,
					"drivers": map[string]interface{}{
						"wopi": map[string]interface{}{
							"app_api_key":          cfg.Drivers.WOPI.AppAPIKey,
							"app_desktop_only":     cfg.Drivers.WOPI.AppDesktopOnly,
							"app_icon_uri":         cfg.Drivers.WOPI.AppIconURI,
							"app_int_url":          cfg.Drivers.WOPI.AppInternalURL,
							"app_name":             cfg.Drivers.WOPI.AppName,
							"app_url":              cfg.Drivers.WOPI.AppURL,
							"insecure_connections": cfg.Drivers.WOPI.Insecure,
							"iop_secret":           cfg.Drivers.WOPI.IopSecret,
							"jwt_secret":           cfg.TokenManager.JWTSecret,
							"wopi_url":             cfg.Drivers.WOPI.WopiURL,
						},
					},
				},
			},
		},
	}
	return rcfg
}
