package revaconfig

import "github.com/owncloud/ocis/v2/services/auth-bearer/pkg/config"

// AuthBearerConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AuthBearerConfigFromStruct(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_exporter":     cfg.Tracing.Type,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
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
					"auth_manager": "oidc",
					"auth_managers": map[string]interface{}{
						"oidc": map[string]interface{}{
							"issuer":    cfg.OIDC.Issuer,
							"insecure":  cfg.OIDC.Insecure,
							"id_claim":  cfg.OIDC.IDClaim,
							"uid_claim": cfg.OIDC.UIDClaim,
							"gid_claim": cfg.OIDC.GIDClaim,
						},
					},
				},
			},
			"interceptors": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "auth_bearer",
				},
			},
		},
	}
}
