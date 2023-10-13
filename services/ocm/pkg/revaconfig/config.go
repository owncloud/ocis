package revaconfig

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config"
)

// OCMConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func OCMConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]interface{} {
	return map[string]interface{}{
		"shared": map[string]interface{}{
			"gatewaysvc":          cfg.Reva.Address, // Todo or address?
			"grpc_client_options": cfg.Reva.GetGRPCClientConfig(),
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			"middlewares": map[string]interface{}{
				"cors": map[string]interface{}{
					"allowed_origins":   cfg.HTTP.CORS.AllowedOrigins,
					"allowed_methods":   cfg.HTTP.CORS.AllowedMethods,
					"allowed_headers":   cfg.HTTP.CORS.AllowedHeaders,
					"allow_credentials": cfg.HTTP.CORS.AllowCredentials,
					// currently unused
					//"options_passthrough": ,
					//"debug": ,
					//"max_age": ,
					//"priority": ,
					//"exposed_headers": ,
				},
				"auth": map[string]interface{}{
					"credentials_by_user_agent": cfg.Middleware.Auth.CredentialsByUserAgent,
				},
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "ocm",
				},
				"requestid": map[string]interface{}{},
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				"sciencemesh": map[string]interface{}{
					"prefix":             cfg.ScienceMesh.Prefix,
					"smtp_credentials":   map[string]string{},
					"gatewaysvc":         cfg.Reva.Address,
					"mesh_directory_url": cfg.Commons.OcisURL,
					"provider_domain":    cfg.Commons.OcisURL,
				},
				"ocmd": map[string]interface{}{
					"prefix":                        cfg.OCMD.Prefix,
					"gatewaysvc":                    cfg.Reva.Address,
					"expose_recipient_display_name": cfg.OCMD.ExposeRecipientDisplayName,
				},
			},
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
				"ocminvitemanager": map[string]interface{}{
					"driver": cfg.OCMInviteManager.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": cfg.OCMInviteManager.Drivers.JSON.File,
						},
					},
					"provider_domain": cfg.Commons.OcisURL,
					"ocm_insecure":    cfg.OCMInviteManager.Insecure,
				},
				"ocmproviderauthorizer": map[string]interface{}{
					"driver": cfg.OCMProviderAuthorizerDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"providers":               cfg.OCMProviderAuthorizerDrivers.JSON.Providers,
							"verify_request_hostname": cfg.OCMProviderAuthorizerDrivers.JSON.VerifyRequestHostname,
						},
					},
				},
				"ocmshareprovider": map[string]interface{}{
					"driver": cfg.OCMShareProvider.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": cfg.OCMShareProvider.Drivers.JSON.File,
						},
					},
					"gatewaysvc":      cfg.Reva.Address,
					"provider_domain": cfg.Commons.OcisURL,
					"webdav_endpoint": cfg.Commons.OcisURL,
					"client_insecure": cfg.OCMShareProvider.Insecure,
				},
				"ocmcore": map[string]interface{}{
					"driver": cfg.OCMCore.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file": cfg.OCMCore.Drivers.JSON.File,
						},
					},
				},
				"authprovider": map[string]interface{}{
					"auth_manager": "ocmshares",
					"auth_managers": map[string]interface{}{
						"ocmshares": map[string]interface{}{
							"gatewaysvc": cfg.Reva.Address,
						},
					},
				},
			},
		},
	}
}
