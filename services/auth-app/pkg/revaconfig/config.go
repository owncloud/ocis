package revaconfig

import (
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
)

// AuthAppConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AuthAppConfigFromStruct(cfg *config.Config) map[string]interface{} {
	appAuthJSON := filepath.Join(defaults.BaseDataPath(), "appauth.json")

	rcfg := map[string]interface{}{
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
			"grpc_client_options":       cfg.Reva.GetGRPCClientConfig(),
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
				"authprovider": map[string]interface{}{
					"auth_manager": "appauth",
					"auth_managers": map[string]interface{}{
						"appauth": map[string]interface{}{
							"gateway_addr": cfg.Reva.Address,
						},
					},
				},
				"applicationauth": map[string]interface{}{
					"driver": "json",
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": appAuthJSON,
						},
					},
				},
			},
			"interceptors": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "auth_app",
				},
			},
		},
	}
	return rcfg
}
