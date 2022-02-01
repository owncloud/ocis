package config

import (
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

const (
	defaultPublicURL        = "https://localhost:9200"
	defaultShareFolder      = "/Shares"
	defaultStorageNamespace = "/users/{{.Id.OpaqueId}}"
	defaultGatewayAddr      = "127.0.0.1:9142"
	defaultUserLayout       = "{{.Id.OpaqueId}}"
	defaultServiceUserUUID  = "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad"
)

func DefaultConfig() *Config {
	return &Config{
		// log is inherited
		Debug: Debug{
			Addr: "127.0.0.1:9109",
		},
		Reva: Reva{
			JWTSecret:             "Pive-Fumkiu4",
			SkipUserGroupsInToken: false,
			TransferSecret:        "replace-me-with-a-transfer-secret",
			TransferExpires:       24 * 60 * 60,
			OIDC: OIDC{
				Issuer:   defaultPublicURL,
				Insecure: false,
				IDClaim:  "preferred_username",
			},
			LDAP: LDAP{
				Hostname:             "localhost",
				Port:                 9126,
				CACert:               path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
				Insecure:             false,
				BaseDN:               "dc=ocis,dc=test",
				LoginFilter:          "(&(objectclass=posixAccount)(|(cn={{login}})(mail={{login}})))",
				UserFilter:           "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
				UserAttributeFilter:  "(&(objectclass=posixAccount)({{attr}}={{value}}))",
				UserFindFilter:       "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
				UserGroupFilter:      "(&(objectclass=posixGroup)(cn={{query}}*))",
				GroupFilter:          "(&(objectclass=posixGroup)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
				GroupAttributeFilter: "(&(objectclass=posixGroup)({{attr}}={{value}}))",
				GroupFindFilter:      "(&(objectclass=posixGroup)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
				GroupMemberFilter:    "(&(objectclass=posixAccount)(ownclouduuid={{.OpaqueId}}*))",
				BindDN:               "cn=reva,ou=sysusers,dc=ocis,dc=test",
				BindPassword:         "reva",
				IDP:                  defaultPublicURL,
				UserSchema: LDAPUserSchema{
					UID:         "ownclouduuid",
					Mail:        "mail",
					DisplayName: "displayname",
					CN:          "cn",
					UIDNumber:   "uidnumber",
					GIDNumber:   "gidnumber",
				},
				GroupSchema: LDAPGroupSchema{
					GID:         "cn",
					Mail:        "mail",
					DisplayName: "cn",
					CN:          "cn",
					GIDNumber:   "gidnumber",
				},
			},
			UserGroupRest: UserGroupRest{
				RedisAddress: "localhost:6379",
			},
			UserOwnCloudSQL: UserOwnCloudSQL{
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
			OCDav: OCDav{
				WebdavNamespace:   defaultStorageNamespace,
				DavFilesNamespace: defaultStorageNamespace,
			},
			Archiver: Archiver{
				MaxNumFiles: 10000,
				MaxSize:     1073741824,
				ArchiverURL: "/archiver",
			},
			UserStorage: StorageConfig{
				EOS: DriverEOS{
					DriverCommon: DriverCommon{
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
				Local: DriverCommon{
					Root:        path.Join(defaults.BaseDataPath(), "storage", "local", "users"),
					ShareFolder: defaultShareFolder,
					UserLayout:  "{{.Username}}",
					EnableHome:  false,
				},
				OwnCloudSQL: DriverOwnCloudSQL{
					DriverCommon: DriverCommon{
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
				S3: DriverS3{
					DriverCommon: DriverCommon{},
					Region:       "default",
					AccessKey:    "",
					SecretKey:    "",
					Endpoint:     "",
					Bucket:       "",
				},
				S3NG: DriverS3NG{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "users"),
						ShareFolder: defaultShareFolder,
						UserLayout:  defaultUserLayout,
						EnableHome:  false,
					},
					ServiceUserUUID: defaultServiceUserUUID,
					Region:          "default",
					AccessKey:       "",
					SecretKey:       "",
					Endpoint:        "",
					Bucket:          "",
				},
				OCIS: DriverOCIS{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "users"),
						ShareFolder: defaultShareFolder,
						UserLayout:  defaultUserLayout,
					},
					ServiceUserUUID: defaultServiceUserUUID,
				},
			},
			MetadataStorage: StorageConfig{
				EOS: DriverEOS{
					DriverCommon: DriverCommon{
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
				Local: DriverCommon{
					Root: path.Join(defaults.BaseDataPath(), "storage", "local", "metadata"),
				},
				OwnCloudSQL: DriverOwnCloudSQL{},
				S3: DriverS3{
					DriverCommon: DriverCommon{},
					Region:       "default",
				},
				S3NG: DriverS3NG{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  defaultUserLayout,
						EnableHome:  false,
					},
					ServiceUserUUID: defaultServiceUserUUID,
					Region:          "default",
					AccessKey:       "",
					SecretKey:       "",
					Endpoint:        "",
					Bucket:          "",
				},
				OCIS: DriverOCIS{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  defaultUserLayout,
						EnableHome:  false,
					},
					ServiceUserUUID: defaultServiceUserUUID,
				},
			},
			Frontend: FrontendPort{
				Port: Port{
					MaxCPUs:     "",
					LogLevel:    "",
					GRPCNetwork: "",
					GRPCAddr:    "",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9140",
					Protocol:    "",
					Endpoint:    "",
					DebugAddr:   "127.0.0.1:9141",
					Services:    []string{"datagateway", "ocdav", "ocs", "appprovider"},
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
				OCDavInsecure:              false, // true?
				OCDavPrefix:                "",
				OCSPrefix:                  "ocs",
				OCSSharePrefix:             defaultShareFolder,
				OCSHomeNamespace:           defaultStorageNamespace,
				PublicURL:                  defaultPublicURL,
				OCSCacheWarmupDriver:       "",
				OCSAdditionalInfoAttribute: "{{.Mail}}",
				OCSResourceInfoCacheTTL:    0,
				Middleware:                 Middleware{},
			},
			DataGateway: DataGatewayPort{
				Port:      Port{},
				PublicURL: "",
			},
			Gateway: Gateway{
				Port: Port{
					Endpoint:    defaultGatewayAddr,
					DebugAddr:   "127.0.0.1:9143",
					GRPCNetwork: "tcp",
					GRPCAddr:    defaultGatewayAddr,
				},
				CommitShareToStorageGrant:  true,
				CommitShareToStorageRef:    true,
				DisableHomeCreationOnLogin: false,
				ShareFolder:                "Shares",
				LinkGrants:                 "",
				HomeMapping:                "",
				EtagCacheTTL:               0,
			},
			StorageRegistry: StorageRegistry{
				Driver:       "spaces",
				HomeProvider: "/home", // unused for spaces, static currently not supported
				JSON:         "",
			},
			AppRegistry: AppRegistry{
				Driver:        "static",
				MimetypesJSON: "",
			},
			Users: Users{
				Port: Port{
					Endpoint:    "localhost:9144",
					DebugAddr:   "127.0.0.1:9145",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9144",
					Services:    []string{"userprovider"},
				},
				Driver:                    "ldap",
				UserGroupsCacheExpiration: 5,
			},
			Groups: Groups{
				Port: Port{
					Endpoint:    "localhost:9160",
					DebugAddr:   "127.0.0.1:9161",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9160",
					Services:    []string{"groupprovider"},
				},
				Driver:                      "ldap",
				GroupMembersCacheExpiration: 5,
			},
			AuthProvider: Users{
				Port:                      Port{},
				Driver:                    "ldap",
				UserGroupsCacheExpiration: 0,
			},
			AuthBasic: Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9146",
				DebugAddr:   "127.0.0.1:9147",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9146",
			},
			AuthBearer: Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9148",
				DebugAddr:   "127.0.0.1:9149",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9148",
			},
			AuthMachine: Port{
				GRPCNetwork: "tcp",
				GRPCAddr:    "127.0.0.1:9166",
				DebugAddr:   "127.0.0.1:9167",
				Services:    []string{"authprovider"},
				Endpoint:    "localhost:9166",
			},
			AuthMachineConfig: AuthMachineConfig{
				MachineAuthAPIKey: "change-me-please",
			},
			Sharing: Sharing{
				Port: Port{
					Endpoint:    "localhost:9150",
					DebugAddr:   "127.0.0.1:9151",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9150",
					Services:    []string{"usershareprovider", "publicshareprovider"},
				},
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
			},
			StorageShares: StoragePort{
				Port: Port{
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
			StorageUsers: StoragePort{
				Port: Port{
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
			StoragePublicLink: PublicStorage{
				StoragePort: StoragePort{
					Port: Port{
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
			StorageMetadata: StoragePort{
				Port: Port{
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9215",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9216",
					DebugAddr:   "127.0.0.1:9217",
				},
				Driver:           "ocis",
				ExposeDataServer: false,
				DataServerURL:    "http://localhost:9216/data",
				TempFolder:       path.Join(defaults.BaseDataPath(), "tmp", "metadata"),
				DataProvider:     DataProvider{},
			},
			AppProvider: AppProvider{
				Port: Port{
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9164",
					DebugAddr:   "127.0.0.1:9165",
					Endpoint:    "localhost:9164",
					Services:    []string{"appprovider"},
				},
				ExternalAddr: "127.0.0.1:9164",
				WopiDriver:   WopiDriver{},
				AppsURL:      "/app/list",
				OpenURL:      "/app/open",
				NewURL:       "/app/new",
			},
			Permissions: Port{
				Endpoint: "localhost:9191",
			},
			Configs:                     nil,
			UploadMaxChunkSize:          1e+8,
			UploadHTTPMethodOverride:    "",
			ChecksumSupportedTypes:      []string{"sha1", "md5", "adler32"},
			ChecksumPreferredUploadType: "",
			DefaultUploadProtocol:       "tus",
		},
		Tracing: Tracing{
			Service: "storage",
			Type:    "jaeger",
		},
		Asset: Asset{},
	}
}
