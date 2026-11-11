package revaconfig

import (
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
)

// AuthAppConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AuthAppConfigFromStruct(cfg *config.Config) map[string]any {
	appAuthJSON := filepath.Join(defaults.BaseDataPath(), "appauth.json")

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
			"services": map[string]any{
				"authprovider": map[string]any{
					"auth_manager": "appauth",
					"auth_managers": map[string]any{
						"appauth": map[string]any{
							"gateway_addr": cfg.Reva.Address,
						},
					},
				},
				"applicationauth": map[string]any{
					"driver": "json",
					"drivers": map[string]any{
						"json": map[string]any{
							"file": appAuthJSON,
						},
					},
				},
			},
			"interceptors": map[string]any{
				"prometheus": map[string]any{
					"namespace": "ocis",
					"subsystem": "auth_app",
				},
			},
		},
	}
	return rcfg
}
