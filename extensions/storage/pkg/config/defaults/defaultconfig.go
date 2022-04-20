package defaults

import (
	"os"
	"path"

	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

const (
	defaultPublicURL                  = "https://localhost:9200"
	defaultShareFolder                = "/Shares"
	defaultStorageNamespace           = "/users/{{.Id.OpaqueId}}"
	defaultGatewayAddr                = "127.0.0.1:9142"
	defaultUserLayout                 = "{{.Id.OpaqueId}}"
	defaultPersonalSpaceAliasTemplate = "{{.SpaceType}}/{{.User.Username | lower}}"
	defaultGeneralSpaceAliasTemplate  = "{{.SpaceType}}/{{.SpaceName | replace \" \" \"-\" | lower}}"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)
	Sanitize(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		// log is inherited
		Debug: config.Debug{
			Addr: "127.0.0.1:9109",
		},
		Reva: config.Reva{
			JWTSecret:             "Pive-Fumkiu4",
			SkipUserGroupsInToken: false,
			TransferSecret:        "replace-me-with-a-transfer-secret",
			TransferExpires:       24 * 60 * 60,
			OIDC: config.OIDC{
				Issuer:   defaultPublicURL,
				Insecure: false,
				IDClaim:  "preferred_username",
			},
			LDAP: config.LDAP{
				URI:              "ldaps://localhost:9126",
				CACert:           path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
				Insecure:         false,
				UserBaseDN:       "dc=ocis,dc=test",
				GroupBaseDN:      "dc=ocis,dc=test",
				UserScope:        "sub",
				GroupScope:       "sub",
				LoginAttributes:  []string{"cn", "mail"},
				UserFilter:       "",
				GroupFilter:      "",
				UserObjectClass:  "posixAccount",
				GroupObjectClass: "posixGroup",
				BindDN:           "cn=reva,ou=sysusers,dc=ocis,dc=test",
				BindPassword:     "reva",
				IDP:              defaultPublicURL,
				UserSchema: config.LDAPUserSchema{
					ID:          "ownclouduuid",
					Mail:        "mail",
					DisplayName: "displayname",
					Username:    "cn",
					UIDNumber:   "uidnumber",
					GIDNumber:   "gidnumber",
				},
				GroupSchema: config.LDAPGroupSchema{
					ID:          "cn",
					Mail:        "mail",
					DisplayName: "cn",
					Groupname:   "cn",
					Member:      "cn",
					GIDNumber:   "gidnumber",
				},
			},
			UserGroupRest: config.UserGroupRest{
				RedisAddress: "localhost:6379",
			},
			UserOwnCloudSQL: config.UserOwnCloudSQL{
				DBUsername:         "owncloud",
				DBPassword:         "secret",
				DBHost:             "mysql",
				DBPort:             3306,
				DBName:             "owncloud",
				Idp:                defaultPublicURL,
				Nobody:             90,
				JoinUsername:       false,
				JoinOwnCloudUUID:   false,
				EnableMedialSearch: false,
			},
			Archiver: config.Archiver{
				MaxNumFiles: 10000,
				MaxSize:     1073741824,
				ArchiverURL: "/archiver",
			},
			UserStorage: config.StorageConfig{
				EOS: config.DriverEOS{
					DriverCommon: config.DriverCommon{
						Root:        "/eos/dockertest/reva",
						ShareFolder: defaultShareFolder,
						UserLayout:  "{{substr 0 1 .Username}}/{{.Username}}",
					},
					ShadowNamespace:  "", // Defaults to path.Join(c.Namespace, ".shadow")
					UploadsNamespace: "", // Defaults to path.Join(c.Namespace, ".uploads")
					EosBinary:        "/usr/bin/eos",
					XrdcopyBinary:    "/usr/bin/xrdcopy",
					MasterURL:        "root://eos-mgm1.eoscluster.cern.ch:1094",
					SlaveURL:         "root://eos-mgm1.eoscluster.cern.ch:1094",
					CacheDirectory:   os.TempDir(),
					GatewaySVC:       defaultGatewayAddr,
				},
				Local: config.DriverCommon{
					Root:        path.Join(defaults.BaseDataPath(), "storage", "local", "users"),
					ShareFolder: defaultShareFolder,
					UserLayout:  "{{.Username}}",
					EnableHome:  false,
				},
				OwnCloudSQL: config.DriverOwnCloudSQL{
					DriverCommon: config.DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "owncloud"),
						ShareFolder: defaultShareFolder,
						UserLayout:  "{{.Username}}",
						EnableHome:  false,
					},
					UploadInfoDir: path.Join(defaults.BaseDataPath(), "storage", "uploadinfo"),
					DBUsername:    "owncloud",
					DBPassword:    "owncloud",
					DBHost:        "",
					DBPort:        3306,
					DBName:        "owncloud",
				},
				S3: config.DriverS3{
					DriverCommon: config.DriverCommon{},
					Region:       "default",
					AccessKey:    "",
					SecretKey:    "",
					Endpoint:     "",
					Bucket:       "",
				},
				S3NG: config.DriverS3NG{
					DriverCommon: config.DriverCommon{
						Root:                       path.Join(defaults.BaseDataPath(), "storage", "users"),
						ShareFolder:                defaultShareFolder,
						UserLayout:                 defaultUserLayout,
						PersonalSpaceAliasTemplate: defaultPersonalSpaceAliasTemplate,
						GeneralSpaceAliasTemplate:  defaultGeneralSpaceAliasTemplate,
						EnableHome:                 false,
					},
					Region:    "default",
					AccessKey: "",
					SecretKey: "",
					Endpoint:  "",
					Bucket:    "",
				},
				OCIS: config.DriverOCIS{
					DriverCommon: config.DriverCommon{
						Root:                       path.Join(defaults.BaseDataPath(), "storage", "users"),
						ShareFolder:                defaultShareFolder,
						UserLayout:                 defaultUserLayout,
						PersonalSpaceAliasTemplate: defaultPersonalSpaceAliasTemplate,
						GeneralSpaceAliasTemplate:  defaultGeneralSpaceAliasTemplate,
					},
				},
			},
			MetadataStorage: config.StorageConfig{
				EOS: config.DriverEOS{
					DriverCommon: config.DriverCommon{
						Root:        "/eos/dockertest/reva",
						ShareFolder: defaultShareFolder,
						UserLayout:  "{{substr 0 1 .Username}}/{{.Username}}",
						EnableHome:  false,
					},
					ShadowNamespace:     "",
					UploadsNamespace:    "",
					EosBinary:           "/usr/bin/eos",
					XrdcopyBinary:       "/usr/bin/xrdcopy",
					MasterURL:           "root://eos-mgm1.eoscluster.cern.ch:1094",
					GrpcURI:             "",
					SlaveURL:            "root://eos-mgm1.eoscluster.cern.ch:1094",
					CacheDirectory:      os.TempDir(),
					EnableLogging:       false,
					ShowHiddenSysFiles:  false,
					ForceSingleUserMode: false,
					UseKeytab:           false,
					SecProtocol:         "",
					Keytab:              "",
					SingleUsername:      "",
					GatewaySVC:          defaultGatewayAddr,
				},
				Local: config.DriverCommon{
					Root: path.Join(defaults.BaseDataPath(), "storage", "local", "metadata"),
				},
				OwnCloudSQL: config.DriverOwnCloudSQL{},
				S3: config.DriverS3{
					DriverCommon: config.DriverCommon{},
					Region:       "default",
				},
				S3NG: config.DriverS3NG{
					DriverCommon: config.DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  defaultUserLayout,
						EnableHome:  false,
					},
					Region:    "default",
					AccessKey: "",
					SecretKey: "",
					Endpoint:  "",
					Bucket:    "",
				},
				OCIS: config.DriverOCIS{
					DriverCommon: config.DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  defaultUserLayout,
						EnableHome:  false,
					},
				},
			},
			Frontend: config.FrontendPort{
				Port: config.Port{
					MaxCPUs:     "",
					LogLevel:    "",
					GRPCNetwork: "",
					GRPCAddr:    "",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9140",
					Protocol:    "",
					Endpoint:    "",
					DebugAddr:   "127.0.0.1:9141",
					Services:    []string{"datagateway", "ocs", "appprovider"},
					Config:      nil,
					Context:     nil,
					Supervised:  false,
				},
				AppProviderInsecure:        false,
				AppProviderPrefix:          "",
				ArchiverInsecure:           false,
				ArchiverPrefix:             "archiver",
				DatagatewayPrefix:          "data",
				Favorites:                  false,
				ProjectSpaces:              true,
				OCSPrefix:                  "ocs",
				OCSSharePrefix:             defaultShareFolder,
				OCSHomeNamespace:           defaultStorageNamespace,
				PublicURL:                  defaultPublicURL,
				OCSCacheWarmupDriver:       "",
				OCSAdditionalInfoAttribute: "{{.Mail}}",
				OCSResourceInfoCacheTTL:    0,
				Middleware:                 config.Middleware{},
			},
			DataGateway: config.DataGatewayPort{
				Port:      config.Port{},
				PublicURL: "",
			},
			Gateway: config.Gateway{
				Port: config.Port{
					Endpoint:    defaultGatewayAddr,
					DebugAddr:   "127.0.0.1:9143",
					GRPCNetwork: "tcp",
					GRPCAddr:    defaultGatewayAddr,
				},
				CommitShareToStorageGrant:  true,
				CommitShareToStorageRef:    true,
				DisableHomeCreationOnLogin: true,
				ShareFolder:                "Shares",
				LinkGrants:                 "",
				HomeMapping:                "",
				EtagCacheTTL:               0,
			},
			StorageRegistry: config.StorageRegistry{
				Driver:       "static",
				HomeProvider: "/home", // unused for spaces, static currently not supported
				JSON:         "",
			},
			AppRegistry: config.AppRegistry{
				Driver:        "static",
				MimetypesJSON: "",
			},
			Users: config.Users{
				Port: config.Port{
					Endpoint:    "localhost:9144",
					DebugAddr:   "127.0.0.1:9145",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9144",
					Services:    []string{"userprovider"},
				},
				Driver:                    "ldap",
				UserGroupsCacheExpiration: 5,
			},
			Groups: config.Groups{
				Port: config.Port{
					Endpoint:    "localhost:9160",
					DebugAddr:   "127.0.0.1:9161",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9160",
					Services:    []string{"groupprovider"},
				},
				Driver:                      "ldap",
				GroupMembersCacheExpiration: 5,
			},
			AuthProvider: config.Users{
				Port:                      config.Port{},
				Driver:                    "ldap",
				UserGroupsCacheExpiration: 0,
			},
			AuthBasic: config.Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9146",
				DebugAddr:   "127.0.0.1:9147",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9146",
			},
			AuthBearer: config.Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9148",
				DebugAddr:   "127.0.0.1:9149",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9148",
			},
			AuthMachine: config.Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9166",
				DebugAddr:   "127.0.0.1:9167",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9166",
			},
			AuthMachineConfig: config.AuthMachineConfig{
				MachineAuthAPIKey: "change-me-please",
			},
			Sharing: config.Sharing{
				Port: config.Port{
					Endpoint:    "localhost:9150",
					DebugAddr:   "127.0.0.1:9151",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9150",
					Services:    []string{"usershareprovider", "publicshareprovider"},
				},
				CS3ProviderAddr:                  "127.0.0.1:9215",
				CS3ServiceUser:                   "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
				CS3ServiceUserIdp:                "https://localhost:9200",
				UserDriver:                       "json",
				UserJSONFile:                     path.Join(defaults.BaseDataPath(), "storage", "shares.json"),
				UserSQLUsername:                  "",
				UserSQLPassword:                  "",
				UserSQLHost:                      "",
				UserSQLPort:                      1433,
				UserSQLName:                      "",
				PublicDriver:                     "json",
				PublicJSONFile:                   path.Join(defaults.BaseDataPath(), "storage", "publicshares.json"),
				PublicPasswordHashCost:           11,
				PublicEnableExpiredSharesCleanup: true,
				PublicJanitorRunInterval:         60,
				UserStorageMountID:               "",
				Events: config.Events{
					Address:   "127.0.0.1:9233",
					ClusterID: "ocis-cluster",
				},
			},
			StorageShares: config.StoragePort{
				Port: config.Port{
					Endpoint:    "localhost:9154",
					DebugAddr:   "127.0.0.1:9156",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9154",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9155",
				},
				ReadOnly:      false,
				AlternativeID: "1284d238-aa92-42ce-bdc4-0b0000009154",
				MountID:       "1284d238-aa92-42ce-bdc4-0b0000009157",
			},
			StorageUsers: config.StoragePort{
				Port: config.Port{
					Endpoint:    "localhost:9157",
					DebugAddr:   "127.0.0.1:9159",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9157",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9158",
				},
				MountID:       "1284d238-aa92-42ce-bdc4-0b0000009157",
				Driver:        "ocis",
				DataServerURL: "http://localhost:9158/data",
				HTTPPrefix:    "data",
				TempFolder:    path.Join(defaults.BaseDataPath(), "tmp", "users"),
			},
			StoragePublicLink: config.PublicStorage{
				StoragePort: config.StoragePort{
					Port: config.Port{
						Endpoint:    "localhost:9178",
						DebugAddr:   "127.0.0.1:9179",
						GRPCNetwork: "tcp",
						GRPCAddr:    "127.0.0.1:9178",
					},
					MountID: "7993447f-687f-490d-875c-ac95e89a62a4",
				},
				PublicShareProviderAddr: "",
				UserProviderAddr:        "",
			},
			StorageMetadata: config.StoragePort{
				Port: config.Port{
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9215",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9216",
					DebugAddr:   "127.0.0.1:9217",
				},
				MountID:          "0dba9855-3ab1-432f-ace7-e01224fe2c65",
				Driver:           "ocis",
				ExposeDataServer: false,
				DataServerURL:    "http://localhost:9216/data",
				TempFolder:       path.Join(defaults.BaseDataPath(), "tmp", "metadata"),
				DataProvider:     config.DataProvider{},
			},
			AppProvider: config.AppProvider{
				Port: config.Port{
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9164",
					DebugAddr:   "127.0.0.1:9165",
					Endpoint:    "localhost:9164",
					Services:    []string{"appprovider"},
				},
				ExternalAddr: "127.0.0.1:9164",
				WopiDriver:   config.WopiDriver{},
				AppsURL:      "/app/list",
				OpenURL:      "/app/open",
				NewURL:       "/app/new",
			},
			Permissions: config.Port{
				Endpoint: "localhost:9191",
			},
			Configs:                     nil,
			UploadMaxChunkSize:          1e+8,
			UploadHTTPMethodOverride:    "",
			ChecksumSupportedTypes:      []string{"sha1", "md5", "adler32"},
			ChecksumPreferredUploadType: "",
			DefaultUploadProtocol:       "tus",
		},
		// TODO move ocdav config to a separate service
		OCDav: config.OCDav{
			Addr:            "127.0.0.1:0", // :0 to pick any local free port
			DebugAddr:       "127.0.0.1:9163",
			WebdavNamespace: defaultStorageNamespace,
			FilesNamespace:  defaultStorageNamespace,
			SharesNamespace: defaultShareFolder,
			PublicURL:       defaultPublicURL,
			Prefix:          "",
			GatewaySVC:      defaultGatewayAddr,
			Insecure:        false, // true?
			Timeout:         84300,
			JWTSecret:       "Pive-Fumkiu4",
		},
		Tracing: config.Tracing{
			Service: "storage",
			Type:    "jaeger",
		},
		Asset: config.Asset{},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// TODO: IMPLEMENT ME!
}

func Sanitize(cfg *config.Config) {
	// TODO: IMPLEMENT ME!
}
