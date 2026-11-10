package revaconfig

import (
	"math"
	"net/url"
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config"
)

// OCMConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func OCMConfigFromStruct(cfg *config.Config, logger log.Logger) map[string]any {

	// Construct the ocm provider domain from the oCIS URL
	providerDomain := ""
	u, err := url.Parse(cfg.Commons.OcisURL)
	switch {
	case err != nil:
		logger.Error().Err(err).Msg("could not parse oCIS URL")
	case u.Host == "":
		logger.Error().Msg("oCIS URL has no host")
	default:
		providerDomain = u.Host
	}

	corsAllowCredentials := cfg.HTTP.CORS.AllowCredentials
	if corsAllowCredentials && cors.AllowsAnyOrigin(cfg.HTTP.CORS.AllowedOrigins) {
		logger.Warn().
			Strs("allowed_origins", cfg.HTTP.CORS.AllowedOrigins).
			Msg("cors: refusing to allow credentials together with a wildcard origin, disabling allow_credentials")
		corsAllowCredentials = false
	}

	return map[string]any{
		"shared": map[string]any{
			"jwt_secret":          cfg.TokenManager.JWTSecret,
			"gatewaysvc":          cfg.Reva.Address, // Todo or address?
			"grpc_client_options": cfg.Reva.GetGRPCClientConfig(),
		},
		"http": map[string]any{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			"middlewares": map[string]any{
				"cors": map[string]any{
					"allowed_origins":   cfg.HTTP.CORS.AllowedOrigins,
					"allowed_methods":   cfg.HTTP.CORS.AllowedMethods,
					"allowed_headers":   cfg.HTTP.CORS.AllowedHeaders,
					"allow_credentials": corsAllowCredentials,
					// currently unused
					//"options_passthrough": ,
					//"debug": ,
					//"max_age": ,
					//"priority": ,
					//"exposed_headers": ,
				},
				"auth": map[string]any{
					"credentials_by_user_agent": cfg.Middleware.Auth.CredentialsByUserAgent,
				},
				"prometheus": map[string]any{
					"namespace": "ocis",
					"subsystem": "ocm",
				},
				"requestid": map[string]any{},
			},
			// TODO build services dynamically
			"services": map[string]any{
				"wellknown": map[string]any{
					"prefix": ".well-known",
					"ocmprovider": map[string]any{
						"ocm_prefix":    cfg.OCMD.Prefix,
						"endpoint":      cfg.Commons.OcisURL,
						"provider":      "oCIS",
						"webdav_root":   "/dav/ocm",
						"webapp_root":   cfg.ScienceMesh.Prefix,
						"enable_webapp": false,
						"enable_datatx": false,
					},
				},
				"sciencemesh": map[string]any{
					"prefix":             cfg.ScienceMesh.Prefix,
					"smtp_credentials":   map[string]string{},
					"gatewaysvc":         cfg.Reva.Address,
					"mesh_directory_url": cfg.ScienceMesh.MeshDirectoryURL,
					"provider_domain":    providerDomain,
					"events": map[string]any{
						"natsaddress":          cfg.Events.Endpoint,
						"natsclusterid":        cfg.Events.Cluster,
						"tlsinsecure":          cfg.Events.TLSInsecure,
						"tlsrootcacertificate": cfg.Events.TLSRootCACertificate,
						"authusername":         cfg.Events.AuthUsername,
						"authpassword":         cfg.Events.AuthPassword,
						"credsfile":            cfg.Events.CredsFile,
						"creds":                cfg.Events.Creds,
					},
				},
				"ocmd": map[string]any{
					"prefix":                        cfg.OCMD.Prefix,
					"gatewaysvc":                    cfg.Reva.Address,
					"expose_recipient_display_name": cfg.OCMD.ExposeRecipientDisplayName,
					"service_account_id":            cfg.ServiceAccount.ID,
					"service_account_secret":        cfg.ServiceAccount.Secret,
				},
				"dataprovider": map[string]any{
					"prefix":           "data",
					"driver":           "ocmreceived",
					"upload_directory": filepath.Join(cfg.OCMStorageProvider.StorageRoot, "uploads"),
					"drivers": map[string]any{
						"ocmreceived": map[string]any{
							"insecure":               cfg.OCMStorageProvider.Insecure,
							"storage_root":           cfg.OCMStorageProvider.StorageRoot,
							"service_account_id":     cfg.ServiceAccount.ID,
							"service_account_secret": cfg.ServiceAccount.Secret,
						},
					},
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
						},
					},
				},
			},
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
				"ocminvitemanager": map[string]any{
					"driver": cfg.OCMInviteManager.Driver,
					"drivers": map[string]any{
						"json": map[string]any{
							"file": cfg.OCMInviteManager.Drivers.JSON.File,
						},
					},
					"provider_domain":  providerDomain,
					"token_expiration": cfg.OCMInviteManager.TokenExpiration.String(),
					"ocm_timeout":      int(math.Round(cfg.OCMInviteManager.Timeout.Seconds())),
					"ocm_insecure":     cfg.OCMInviteManager.Insecure,
				},
				"ocmproviderauthorizer": map[string]any{
					"driver": cfg.OCMProviderAuthorizerDriver,
					"drivers": map[string]any{
						"json": map[string]any{
							"providers": cfg.OCMProviderAuthorizerDrivers.JSON.Providers,
						},
					},
				},
				"ocmshareprovider": map[string]any{
					"driver": cfg.OCMShareProvider.Driver,
					"drivers": map[string]any{
						"json": map[string]any{
							"file": cfg.OCMShareProvider.Drivers.JSON.File,
						},
					},
					"gatewaysvc":          cfg.Reva.Address,
					"provider_domain":     providerDomain,
					"webdav_endpoint":     cfg.Commons.OcisURL,
					"webapp_template":     cfg.OCMShareProvider.WebappTemplate,
					"client_insecure":     cfg.OCMShareProvider.Insecure,
					"enable_user_sharing": cfg.OCMShareProvider.EnableUserSharing,
				},
				"ocmcore": map[string]any{
					"driver": cfg.OCMCore.Driver,
					"drivers": map[string]any{
						"json": map[string]any{
							"file": cfg.OCMCore.Drivers.JSON.File,
						},
					},
					"events": map[string]any{
						"natsaddress":          cfg.Events.Endpoint,
						"natsclusterid":        cfg.Events.Cluster,
						"tlsinsecure":          cfg.Events.TLSInsecure,
						"tlsrootcacertificate": cfg.Events.TLSRootCACertificate,
						"authusername":         cfg.Events.AuthUsername,
						"authpassword":         cfg.Events.AuthPassword,
						"credsfile":            cfg.Events.CredsFile,
						"creds":                cfg.Events.Creds,
					},
				},
				"storageprovider": map[string]any{
					"driver":           "ocmreceived",
					"upload_directory": filepath.Join(cfg.OCMStorageProvider.StorageRoot, "uploads"),
					"drivers": map[string]any{
						"ocmreceived": map[string]any{
							"insecure":     cfg.OCMStorageProvider.Insecure,
							"storage_root": cfg.OCMStorageProvider.StorageRoot,
						},
					},
					"data_server_url": cfg.OCMStorageProvider.DataServerURL,
				},
				"authprovider": map[string]any{
					"auth_manager": "ocmshares",
					"auth_managers": map[string]any{
						"ocmshares": map[string]any{
							"gatewaysvc": cfg.Reva.Address,
						},
					},
				},
			},
		},
	}
}
