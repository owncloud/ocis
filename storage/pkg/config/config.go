package config

import (
	"context"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `ocisConfig:"level"`
	Pretty bool   `ocisConfig:"pretty"`
	Color  bool   `ocisConfig:"color"`
	File   string `ocisConfig:"file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// Gateway defines the available gateway configuration.
type Gateway struct {
	Port
	CommitShareToStorageGrant  bool   `ocisConfig:"commit_share_to_storage_grant"`
	CommitShareToStorageRef    bool   `ocisConfig:"commit_share_to_storage_ref"`
	DisableHomeCreationOnLogin bool   `ocisConfig:"disable_home_creation_on_login"`
	ShareFolder                string `ocisConfig:"share_folder"`
	LinkGrants                 string `ocisConfig:"link_grants"`
	HomeMapping                string `ocisConfig:"home_mapping"`
	EtagCacheTTL               int    `ocisConfig:"etag_cache_ttl"`
}

// StorageRegistry defines the available storage registry configuration
type StorageRegistry struct {
	Driver string `ocisConfig:"driver"`
	// HomeProvider is the path in the global namespace that the static storage registry uses to determine the home storage
	HomeProvider string   `ocisConfig:"home_provider"`
	Rules        []string `ocisConfig:"rules"`
	JSON         string   `ocisConfig:"json"`
}

// AppRegistry defines the available app registry configuration
type AppRegistry struct {
	Driver        string `ocisConfig:"driver"`
	MimetypesJSON string `ocisConfig:"mime_types_json"`
}

// AppProvider defines the available app provider configuration
type AppProvider struct {
	Port
	ExternalAddr string     `ocisConfig:"external_addr"`
	Driver       string     `ocisConfig:"driver"`
	WopiDriver   WopiDriver `ocisConfig:"wopi_driver"`
	AppsURL      string     `ocisConfig:"apps_url"`
	OpenURL      string     `ocisConfig:"open_url"`
}

type WopiDriver struct {
	AppAPIKey      string `ocisConfig:"app_api_key"`
	AppDesktopOnly bool   `ocisConfig:"app_desktop_only"`
	AppIconURI     string `ocisConfig:"app_icon_uri"`
	AppInternalURL string `ocisConfig:"app_internal_url"`
	AppName        string `ocisConfig:"app_name"`
	AppURL         string `ocisConfig:"app_url"`
	Insecure       bool   `ocisConfig:"insecure"`
	IopSecret      string `ocisConfig:"ipo_secret"`
	JWTSecret      string `ocisConfig:"jwt_secret"`
	WopiURL        string `ocisConfig:"wopi_url"`
}

// Sharing defines the available sharing configuration.
type Sharing struct {
	Port
	UserDriver                       string `ocisConfig:"user_driver"`
	UserJSONFile                     string `ocisConfig:"user_json_file"`
	UserSQLUsername                  string `ocisConfig:"user_sql_username"`
	UserSQLPassword                  string `ocisConfig:"user_sql_password"`
	UserSQLHost                      string `ocisConfig:"user_sql_host"`
	UserSQLPort                      int    `ocisConfig:"user_sql_port"`
	UserSQLName                      string `ocisConfig:"user_sql_name"`
	PublicDriver                     string `ocisConfig:"public_driver"`
	PublicJSONFile                   string `ocisConfig:"public_json_file"`
	PublicPasswordHashCost           int    `ocisConfig:"public_password_hash_cost"`
	PublicEnableExpiredSharesCleanup bool   `ocisConfig:"public_enable_expired_shares_cleanup"`
	PublicJanitorRunInterval         int    `ocisConfig:"public_janitor_run_interval"`
	UserStorageMountID               string `ocisConfig:"user_storage_mount_id"`
}

// Port defines the available port configuration.
type Port struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string `ocisConfig:"max_cpus"`
	LogLevel string `ocisConfig:"log_level"`
	// GRPCNetwork can be tcp, udp or unix
	GRPCNetwork string `ocisConfig:"grpc_network"`
	// GRPCAddr to listen on, hostname:port (0.0.0.0:9999 for all interfaces) or socket (/var/run/reva/sock)
	GRPCAddr string `ocisConfig:"grpc_addr"`
	// Protocol can be grpc or http
	// HTTPNetwork can be tcp, udp or unix
	HTTPNetwork string `ocisConfig:"http_network"`
	// HTTPAddr to listen on, hostname:port (0.0.0.0:9100 for all interfaces) or socket (/var/run/reva/sock)
	HTTPAddr string `ocisConfig:"http_addr"`
	// Protocol can be grpc or http
	Protocol string `ocisConfig:"protocol"`
	// Endpoint is used by the gateway and registries (eg localhost:9100 or cloud.example.com)
	Endpoint string `ocisConfig:"endpoint"`
	// DebugAddr for the debug endpoint to bind to
	DebugAddr string `ocisConfig:"debug_addr"`
	// Services can be used to give a list of services that should be started on this port
	Services []string `ocisConfig:"services"`
	// Config can be used to configure the reva instance.
	// Services and Protocol will be ignored if this is used
	Config map[string]interface{} `ocisConfig:"config"`

	// Context allows for context cancellation and propagation
	Context context.Context

	// Supervised is used when running under an oCIS runtime supervision tree
	Supervised bool // deprecated
}

// Users defines the available users configuration.
type Users struct {
	Port
	Driver                    string `ocisConfig:"driver"`
	JSON                      string `ocisConfig:"json"`
	UserGroupsCacheExpiration int    `ocisConfig:"user_groups_cache_expiration"`
}

// AuthMachineConfig defines the available configuration for the machine auth driver.
type AuthMachineConfig struct {
	MachineAuthAPIKey string `ocisConfig:"machine_auth_api_key"`
}

// Groups defines the available groups configuration.
type Groups struct {
	Port
	Driver                      string `ocisConfig:"driver"`
	JSON                        string `ocisConfig:"json"`
	GroupMembersCacheExpiration int    `ocisConfig:"group_members_cache_expiration"`
}

// FrontendPort defines the available frontend configuration.
type FrontendPort struct {
	Port

	AppProviderInsecure        bool       `ocisConfig:"app_provider_insecure"`
	AppProviderPrefix          string     `ocisConfig:"app_provider_prefix"`
	ArchiverInsecure           bool       `ocisConfig:"archiver_insecure"`
	ArchiverPrefix             string     `ocisConfig:"archiver_prefix"`
	DatagatewayPrefix          string     `ocisConfig:"data_gateway_prefix"`
	Favorites                  bool       `ocisConfig:"favorites"`
	OCDavInsecure              bool       `ocisConfig:"ocdav_insecure"`
	OCDavPrefix                string     `ocisConfig:"ocdav_prefix"`
	OCSPrefix                  string     `ocisConfig:"ocs_prefix"`
	OCSSharePrefix             string     `ocisConfig:"ocs_share_prefix"`
	OCSHomeNamespace           string     `ocisConfig:"ocs_home_namespace"`
	PublicURL                  string     `ocisConfig:"public_url"`
	OCSCacheWarmupDriver       string     `ocisConfig:"ocs_cache_warmup_driver"`
	OCSAdditionalInfoAttribute string     `ocisConfig:"ocs_additional_info_attribute"`
	OCSResourceInfoCacheTTL    int        `ocisConfig:"ocs_resource_info_cache_ttl"`
	Middleware                 Middleware `ocisConfig:"middleware"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `ocisConfig:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `ocisConfig:"credentials_by_user_agenr"`
}

// DataGatewayPort has a public url
type DataGatewayPort struct {
	Port
	PublicURL string `ocisConfig:""`
}

type DataProvider struct {
	Insecure bool `ocisConfig:"insecure"`
}

// StoragePort defines the available storage configuration.
type StoragePort struct {
	Port
	Driver           string `ocisConfig:"driver"`
	MountPath        string `ocisConfig:"mount_path"`
	MountID          string `ocisConfig:"mount_id"`
	AlternativeID    string `ocisConfig:"alternative_id"`
	ExposeDataServer bool   `ocisConfig:"expose_data_server"`
	// url the data gateway will use to route requests
	DataServerURL string `ocisConfig:"data_server_url"`

	// for HTTP ports with only one http service
	HTTPPrefix   string       `ocisConfig:"http_prefix"`
	TempFolder   string       `ocisConfig:"temp_folder"`
	ReadOnly     bool         `ocisConfig:"read_only"`
	DataProvider DataProvider `ocisConfig:"data_provider"`
}

// PublicStorage configures a public storage provider
type PublicStorage struct {
	StoragePort

	PublicShareProviderAddr string `ocisConfig:"public_share_provider_addr"`
	UserProviderAddr        string `ocisConfig:"user_provider_addr"`
}

// StorageConfig combines all available storage driver configuration parts.
type StorageConfig struct {
	EOS         DriverEOS         `ocisConfig:"eos"`
	Local       DriverCommon      `ocisConfig:"local"`
	OwnCloud    DriverOwnCloud    `ocisConfig:"owncloud"`
	OwnCloudSQL DriverOwnCloudSQL `ocisConfig:"owncloud_sql"`
	S3          DriverS3          `ocisConfig:"s3"`
	S3NG        DriverS3NG        `ocisConfig:"s3ng"`
	OCIS        DriverOCIS        `ocisConfig:"ocis"`
}

// DriverCommon defines common driver configuration options.
type DriverCommon struct {
	// Root is the absolute path to the location of the data
	Root string `ocisConfig:"root"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `ocisConfig:"share_folder"`
	// UserLayout contains the template used to construct
	// the internal path, eg: `{{substr 0 1 .Username}}/{{.Username}}`
	UserLayout string `ocisConfig:"user_layout"`
	// EnableHome enables the creation of home directories.
	EnableHome bool `ocisConfig:"enable_home"`
}

// DriverEOS defines the available EOS driver configuration.
type DriverEOS struct {
	DriverCommon

	// ShadowNamespace for storing shadow data
	ShadowNamespace string `ocisConfig:"shadow_namespace"`

	// UploadsNamespace for storing upload data
	UploadsNamespace string `ocisConfig:"uploads_namespace"`

	// Location of the eos binary.
	// Default is /usr/bin/eos.
	EosBinary string `ocisConfig:"eos_binary"`

	// Location of the xrdcopy binary.
	// Default is /usr/bin/xrdcopy.
	XrdcopyBinary string `ocisConfig:"xrd_copy_binary"`

	// URL of the Master EOS MGM.
	// Default is root://eos-example.org
	MasterURL string `ocisConfig:"master_url"`

	// URI of the EOS MGM grpc server
	// Default is empty
	GrpcURI string `ocisConfig:"grpc_uri"`

	// URL of the Slave EOS MGM.
	// Default is root://eos-example.org
	SlaveURL string `ocisConfig:"slave_url"`

	// Location on the local fs where to store reads.
	// Defaults to os.TempDir()
	CacheDirectory string `ocisConfig:"cache_directory"`

	// Enables logging of the commands executed
	// Defaults to false
	EnableLogging bool `ocisConfig:"enable_logging"`

	// ShowHiddenSysFiles shows internal EOS files like
	// .sys.v# and .sys.a# files.
	ShowHiddenSysFiles bool `ocisConfig:"shadow_hidden_files"`

	// ForceSingleUserMode will force connections to EOS to use SingleUsername
	ForceSingleUserMode bool `ocisConfig:"force_single_user_mode"`

	// UseKeyTabAuth changes will authenticate requests by using an EOS keytab.
	UseKeytab bool `ocisConfig:"user_keytab"`

	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string `ocisConfig:"sec_protocol"`

	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string `ocisConfig:"keytab"`

	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string `ocisConfig:"single_username"`

	// gateway service to use for uid lookups
	GatewaySVC string `ocisConfig:"gateway_svc"`
}

// DriverOCIS defines the available oCIS storage driver configuration.
type DriverOCIS struct {
	DriverCommon

	ServiceUserUUID string `ocisConfig:"service_user_uuid"`
}

// DriverOwnCloud defines the available ownCloud storage driver configuration.
type DriverOwnCloud struct {
	DriverCommon

	UploadInfoDir string `ocisConfig:"upload_info_dir"`
	Redis         string `ocisConfig:"redis"`
	Scan          bool   `ocisConfig:"scan"`
}

// DriverOwnCloudSQL defines the available ownCloudSQL storage driver configuration.
type DriverOwnCloudSQL struct {
	DriverCommon

	UploadInfoDir string `ocisConfig:"upload_info_dir"`
	DBUsername    string `ocisConfig:"db_username"`
	DBPassword    string `ocisConfig:"db_password"`
	DBHost        string `ocisConfig:"db_host"`
	DBPort        int    `ocisConfig:"db_port"`
	DBName        string `ocisConfig:"db_name"`
}

// DriverS3 defines the available S3 storage driver configuration.
type DriverS3 struct {
	DriverCommon

	Region    string `ocisConfig:"region"`
	AccessKey string `ocisConfig:"access_key"`
	SecretKey string `ocisConfig:"secret_key"`
	Endpoint  string `ocisConfig:"endpoint"`
	Bucket    string `ocisConfig:"bucket"`
}

// DriverS3NG defines the available s3ng storage driver configuration.
type DriverS3NG struct {
	DriverCommon

	Region    string `ocisConfig:"region"`
	AccessKey string `ocisConfig:"access_key"`
	SecretKey string `ocisConfig:"secret_key"`
	Endpoint  string `ocisConfig:"endpoint"`
	Bucket    string `ocisConfig:"bucket"`
}

// OIDC defines the available OpenID Connect configuration.
type OIDC struct {
	Issuer   string `ocisConfig:"issuer"`
	Insecure bool   `ocisConfig:"insecure"`
	IDClaim  string `ocisConfig:"id_claim"`
	UIDClaim string `ocisConfig:"uid_claim"`
	GIDClaim string `ocisConfig:"gid_claim"`
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname             string          `ocisConfig:"hostname"`
	Port                 int             `ocisConfig:"port"`
	CACert               string          `ocisConfig:"ca_cert"`
	Insecure             bool            `ocisConfig:"insecure"`
	BaseDN               string          `ocisConfig:"base_dn"`
	LoginFilter          string          `ocisConfig:"login_filter"`
	UserFilter           string          `ocisConfig:"user_filter"`
	UserAttributeFilter  string          `ocisConfig:"user_attribute_filter"`
	UserFindFilter       string          `ocisConfig:"user_find_filter"`
	UserGroupFilter      string          `ocisConfig:"user_group_filter"`
	GroupFilter          string          `ocisConfig:"group_filter"`
	GroupAttributeFilter string          `ocisConfig:"group_attribute_filter"`
	GroupFindFilter      string          `ocisConfig:"group_finder_filter"`
	GroupMemberFilter    string          `ocisConfig:"group_member_filter"`
	BindDN               string          `ocisConfig:"bind_dn"`
	BindPassword         string          `ocisConfig:"bind_password"`
	IDP                  string          `ocisConfig:"idp"`
	UserSchema           LDAPUserSchema  `ocisConfig:"user_schema"`
	GroupSchema          LDAPGroupSchema `ocisConfig:"group_schema"`
}

// UserGroupRest defines the REST driver specification for user and group resolution.
type UserGroupRest struct {
	ClientID          string `ocisConfig:"client_id"`
	ClientSecret      string `ocisConfig:"client_secret"`
	RedisAddress      string `ocisConfig:"redis_address"`
	RedisUsername     string `ocisConfig:"redis_username"`
	RedisPassword     string `ocisConfig:"redis_password"`
	IDProvider        string `ocisConfig:"idp_provider"`
	APIBaseURL        string `ocisConfig:"api_base_url"`
	OIDCTokenEndpoint string `ocisConfig:"oidc_token_endpoint"`
	TargetAPI         string `ocisConfig:"target_api"`
}

// UserOwnCloudSQL defines the available ownCloudSQL user provider configuration.
type UserOwnCloudSQL struct {
	DBUsername         string `ocisConfig:"db_username"`
	DBPassword         string `ocisConfig:"db_password"`
	DBHost             string `ocisConfig:"db_host"`
	DBPort             int    `ocisConfig:"db_port"`
	DBName             string `ocisConfig:"db_name"`
	Idp                string `ocisConfig:"idp"`
	Nobody             int64  `ocisConfig:"nobody"`
	JoinUsername       bool   `ocisConfig:"join_username"`
	JoinOwnCloudUUID   bool   `ocisConfig:"join_owncloud_uuid"`
	EnableMedialSearch bool   `ocisConfig:"enable_medial_search"`
}

// LDAPUserSchema defines the available ldap user schema configuration.
type LDAPUserSchema struct {
	UID         string `ocisConfig:"uid"`
	Mail        string `ocisConfig:"mail"`
	DisplayName string `ocisConfig:"display_name"`
	CN          string `ocisConfig:"cn"`
	UIDNumber   string `ocisConfig:"uid_number"`
	GIDNumber   string `ocisConfig:"gid_number"`
}

// LDAPGroupSchema defines the available ldap group schema configuration.
type LDAPGroupSchema struct {
	GID         string `ocisConfig:"gid"`
	Mail        string `ocisConfig:"mail"`
	DisplayName string `ocisConfig:"display_name"`
	CN          string `ocisConfig:"cn"`
	GIDNumber   string `ocisConfig:"gid_number"`
}

// OCDav defines the available ocdav configuration.
type OCDav struct {
	WebdavNamespace   string `ocisConfig:"webdav_namespace"`
	DavFilesNamespace string `ocisConfig:"dav_files_namespace"`
}

// Archiver defines the available archiver configuration.
type Archiver struct {
	MaxNumFiles int64  `ocisConfig:"max_num_files"`
	MaxSize     int64  `ocisConfig:"max_size"`
	ArchiverURL string `ocisConfig:"archiver_url"`
}

// Reva defines the available reva configuration.
type Reva struct {
	// JWTSecret used to sign jwt tokens between services
	JWTSecret             string          `ocisConfig:"jwt_secret"`
	SkipUserGroupsInToken bool            `ocisConfig:"skip_user_grooups_in_token"`
	TransferSecret        string          `ocisConfig:"transfer_secret"`
	TransferExpires       int             `ocisConfig:"transfer_expires"`
	OIDC                  OIDC            `ocisConfig:"oidc"`
	LDAP                  LDAP            `ocisConfig:"ldap"`
	UserGroupRest         UserGroupRest   `ocisConfig:"user_group_rest"`
	UserOwnCloudSQL       UserOwnCloudSQL `ocisConfig:"user_owncloud_sql"`
	OCDav                 OCDav           `ocisConfig:"ocdav"`
	Archiver              Archiver        `ocisConfig:"archiver"`
	UserStorage           StorageConfig   `ocisConfig:"user_storage"`
	MetadataStorage       StorageConfig   `ocisConfig:"metadata_storage"`
	// Ports are used to configure which services to start on which port
	Frontend          FrontendPort      `ocisConfig:"frontend"`
	DataGateway       DataGatewayPort   `ocisConfig:"data_gateway"`
	Gateway           Gateway           `ocisConfig:"gateway"`
	StorageRegistry   StorageRegistry   `ocisConfig:"storage_registry"`
	AppRegistry       AppRegistry       `ocisConfig:"app_registry"`
	Users             Users             `ocisConfig:"users"`
	Groups            Groups            `ocisConfig:"groups"`
	AuthProvider      Users             `ocisConfig:"auth_provider"`
	AuthBasic         Port              `ocisConfig:"auth_basic"`
	AuthBearer        Port              `ocisConfig:"auth_bearer"`
	AuthMachine       Port              `ocisConfig:"auth_machine"`
	AuthMachineConfig AuthMachineConfig `ocisConfig:"auth_machine_config"`
	Sharing           Sharing           `ocisConfig:"sharing"`
	StorageHome       StoragePort       `ocisConfig:"storage_home"`
	StorageUsers      StoragePort       `ocisConfig:"storage_users"`
	StoragePublicLink PublicStorage     `ocisConfig:"storage_public_link"`
	StorageMetadata   StoragePort       `ocisConfig:"storage_metadata"`
	AppProvider       AppProvider       `ocisConfig:"app_provider"`
	// Configs can be used to configure the reva instance.
	// Services and Ports will be ignored if this is used
	Configs map[string]interface{} `ocisConfig:"configs"`
	// chunking and resumable upload config (TUS)
	UploadMaxChunkSize       int    `ocisConfig:"uppload_max_chunk_size"`
	UploadHTTPMethodOverride string `ocisConfig:"upload_http_method_override"`
	// checksumming capabilities
	ChecksumSupportedTypes      []string `ocisConfig:"checksum_supported_types"`
	ChecksumPreferredUploadType string   `ocisConfig:"checksum_preferred_upload_type"`
	DefaultUploadProtocol       string   `ocisConfig:"default_upload_protocol"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File    string      `ocisConfig:"file"`
	Log     *shared.Log `ocisConfig:"log"`
	Debug   Debug       `ocisConfig:"debug"`
	Reva    Reva        `ocisConfig:"reva"`
	Tracing Tracing     `ocisConfig:"tracing"`
	Asset   Asset       `ocisConfig:"asset"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

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
				Issuer:   "https://localhost:9200",
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
				UserGroupFilter:      "(&(objectclass=posixGroup)(ownclouduuid={{.OpaqueId}}*))",
				GroupFilter:          "(&(objectclass=posixGroup)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
				GroupAttributeFilter: "(&(objectclass=posixGroup)({{attr}}={{value}}))",
				GroupFindFilter:      "(&(objectclass=posixGroup)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
				GroupMemberFilter:    "(&(objectclass=posixAccount)(ownclouduuid={{.OpaqueId}}*))",
				BindDN:               "cn=reva,ou=sysusers,dc=ocis,dc=test",
				BindPassword:         "reva",
				IDP:                  "https://localhost:9200",
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
				Idp:                "https://localhost:9200",
				Nobody:             90,
				JoinUsername:       false,
				JoinOwnCloudUUID:   false,
				EnableMedialSearch: false,
			},
			OCDav: OCDav{
				WebdavNamespace:   "/home/",
				DavFilesNamespace: "/users/",
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
						ShareFolder: "/Shares",
						UserLayout:  "{{substr 0 1 .Username}}/{{.Username}}",
					},
					ShadowNamespace:  "", // Defaults to path.Join(c.Namespace, ".shadow")
					UploadsNamespace: "", // Defaults to path.Join(c.Namespace, ".uploads")
					EosBinary:        "/usr/bin/eos",
					XrdcopyBinary:    "/usr/bin/xrdcopy",
					MasterURL:        "root://eos-mgm1.eoscluster.cern.ch:1094",
					SlaveURL:         "root://eos-mgm1.eoscluster.cern.ch:1094",
					CacheDirectory:   os.TempDir(),
					GatewaySVC:       "127.0.0.1:9142",
				},
				Local: DriverCommon{
					Root:        path.Join(defaults.BaseDataPath(), "storage", "local", "users"),
					ShareFolder: "/Shares",
					UserLayout:  "{{.Username}}",
					EnableHome:  false,
				},
				OwnCloud: DriverOwnCloud{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "owncloud"),
						ShareFolder: "/Shares",
						UserLayout:  "{{.Id.OpaqueId}}",
						EnableHome:  false,
					},
					UploadInfoDir: path.Join(defaults.BaseDataPath(), "storage", "uploadinfo"),
					Redis:         ":6379",
					Scan:          true,
				},
				OwnCloudSQL: DriverOwnCloudSQL{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "owncloud"),
						ShareFolder: "/Shares",
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
						ShareFolder: "/Shares",
						UserLayout:  "{{.Id.OpaqueId}}",
						EnableHome:  false,
					},
					Region:    "default",
					AccessKey: "",
					SecretKey: "",
					Endpoint:  "",
					Bucket:    "",
				},
				OCIS: DriverOCIS{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "users"),
						ShareFolder: "/Shares",
						UserLayout:  "{{.Id.OpaqueId}}",
					},
					ServiceUserUUID: "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
				},
			},
			MetadataStorage: StorageConfig{
				EOS: DriverEOS{
					DriverCommon: DriverCommon{
						Root:        "/eos/dockertest/reva",
						ShareFolder: "/Shares",
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
					GatewaySVC:          "127.0.0.1:9142",
				},
				Local: DriverCommon{
					Root: path.Join(defaults.BaseDataPath(), "storage", "local", "metadata"),
				},
				OwnCloud:    DriverOwnCloud{},
				OwnCloudSQL: DriverOwnCloudSQL{},
				S3: DriverS3{
					DriverCommon: DriverCommon{},
					Region:       "default",
				},
				S3NG: DriverS3NG{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  "{{.Id.OpaqueId}}",
						EnableHome:  false,
					},
					Region: "default",
				},
				OCIS: DriverOCIS{
					DriverCommon: DriverCommon{
						Root:        path.Join(defaults.BaseDataPath(), "storage", "metadata"),
						ShareFolder: "",
						UserLayout:  "{{.Id.OpaqueId}}",
						EnableHome:  false,
					},
					ServiceUserUUID: "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
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
				OCDavInsecure:              false,
				OCDavPrefix:                "",
				OCSPrefix:                  "ocs",
				OCSSharePrefix:             "/Shares",
				OCSHomeNamespace:           "/home",
				PublicURL:                  "https://localhost:9200",
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
					Endpoint:    "127.0.0.1:9142",
					DebugAddr:   "127.0.0.1:9143",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9142",
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
				Driver:       "static",
				HomeProvider: "/home",
				JSON:         "",
			},
			AppRegistry: AppRegistry{
				Driver:        "static",
				MimetypesJSON: "",
			},
			Users: Users{
				Port: Port{
					Endpoint:    "localhost:9144",
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
			StorageHome: StoragePort{
				Port: Port{
					Endpoint:    "localhost:9154",
					DebugAddr:   "127.0.0.1:9156",
					GRPCNetwork: "tcp",
					GRPCAddr:    "127.0.0.1:9154",
					HTTPNetwork: "tcp",
					HTTPAddr:    "127.0.0.1:9155",
				},
				Driver:        "ocis",
				ReadOnly:      false,
				MountPath:     "/home",
				AlternativeID: "1284d238-aa92-42ce-bdc4-0b0000009154",
				MountID:       "1284d238-aa92-42ce-bdc4-0b0000009157",
				DataServerURL: "http://localhost:9155/data",
				HTTPPrefix:    "data",
				TempFolder:    path.Join(defaults.BaseDataPath(), "tmp", "home"),
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
				MountPath:     "/users",
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
					MountPath: "/public",
					MountID:   "e1a73ede-549b-4226-abdf-40e69ca8230d",
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
				DataServerURL:    "http://localhost:9216",
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

// StructMappings binds a set of environment variables to a destination on cfg. Iterating over this set and editing the
// Destination value of a binding will alter the original value, as it is a pointer to its memory address. This lets
// us propagate changes easier.
func StructMappings(cfg *Config) []shared.EnvBinding {
	return structMappings(cfg)
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv(cfg *Config) []string {
	var r = make([]string, len(structMappings(cfg)))
	for i := range structMappings(cfg) {
		r = append(r, structMappings(cfg)[i].EnvVars...)
	}

	return r
}

func structMappings(cfg *Config) []shared.EnvBinding {
	return []shared.EnvBinding{
		// Shared
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL", "STORAGE_FRONTEND_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY", "STORAGE_FRONTEND_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR", "STORAGE_FRONTEND_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_METADATA_DATAPROVIDER_INSECURE"},
			Destination: &cfg.Reva.StorageMetadata.DataProvider.Insecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_HOME_DATAPROVIDER_INSECURE"},
			Destination: &cfg.Reva.StorageHome.DataProvider.Insecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_FRONTEND_APPPROVIDER_INSECURE"},
			Destination: &cfg.Reva.Frontend.AppProviderInsecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_FRONTEND_ARCHIVER_INSECURE"},
			Destination: &cfg.Reva.Frontend.ArchiverInsecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_FRONTEND_OCDAV_INSECURE"},
			Destination: &cfg.Reva.Frontend.OCDavInsecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_OIDC_INSECURE"},
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		{
			EnvVars:     []string{"OCIS_INSECURE", "STORAGE_USERS_DATAPROVIDER_INSECURE"},
			Destination: &cfg.Reva.StorageUsers.DataProvider.Insecure,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.UserStorage.Local.Root,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER"},
			Destination: &cfg.Reva.StorageUsers.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_DRIVER"},
			Destination: &cfg.Reva.StorageHome.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_DATADIR"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.Root,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_ROOT"},
			Destination: &cfg.Reva.UserStorage.OCIS.Root,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.Root,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_JSON_FILE"},
			Destination: &cfg.Reva.Sharing.UserJSONFile,
		},
		{
			EnvVars:     []string{"OCIS_URL", "STORAGE_FRONTEND_PUBLIC_URL"},
			Destination: &cfg.Reva.Frontend.PublicURL,
		},
		{
			EnvVars:     []string{"OCIS_URL", "STORAGE_OIDC_ISSUER"},
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		{
			EnvVars:     []string{"OCIS_URL", "STORAGE_LDAP_IDP"},
			Destination: &cfg.Reva.LDAP.IDP,
		},
		{
			EnvVars:     []string{"OCIS_URL", "STORAGE_USERPROVIDER_OWNCLOUDSQL_IDP"},
			Destination: &cfg.Reva.UserOwnCloudSQL.Idp,
		},
		{
			EnvVars:     []string{"STORAGE_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},

		// debug

		{
			EnvVars:     []string{"STORAGE_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"STORAGE_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"STORAGE_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},

		// app provider

		{
			EnvVars:     []string{"APP_PROVIDER_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AppProvider.DebugAddr,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_BASIC_GRPC_NETWORK"},
			Destination: &cfg.Reva.AppProvider.GRPCNetwork,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_BASIC_GRPC_ADDR"},
			Destination: &cfg.Reva.AppProvider.GRPCAddr,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_BASIC_EXTERNAL_ADDR"},
			Destination: &cfg.Reva.AppProvider.ExternalAddr,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_DRIVER"},
			Destination: &cfg.Reva.AppProvider.Driver,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_API_KEY"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppAPIKey,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_DESKTOP_ONLY"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppDesktopOnly,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_ICON_URI"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppIconURI,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_INTERNAL_URL"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppInternalURL,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_NAME"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppName,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_APP_URL"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.AppURL,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_INSECURE"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.Insecure,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_IOP_SECRET"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.IopSecret,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_WOPI_DRIVER_WOPI_URL"},
			Destination: &cfg.Reva.AppProvider.WopiDriver.WopiURL,
		},

		// authbasic
		{
			EnvVars:     []string{"STORAGE_AUTH_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBasic.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_DRIVER"},
			Destination: &cfg.Reva.AuthProvider.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_JSON"},
			Destination: &cfg.Reva.AuthProvider.JSON,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthBasic.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthBasic.GRPCAddr,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		// authbearer
		{
			EnvVars:     []string{"STORAGE_AUTH_BEARER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBearer.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_OIDC_ID_CLAIM"},
			Destination: &cfg.Reva.OIDC.IDClaim,
		},
		{
			EnvVars:     []string{"STORAGE_OIDC_UID_CLAIM"},
			Destination: &cfg.Reva.OIDC.UIDClaim,
		},
		{
			EnvVars:     []string{"STORAGE_OIDC_GID_CLAIM"},
			Destination: &cfg.Reva.OIDC.GIDClaim,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BEARER_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthBearer.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BEARER_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthBearer.GRPCAddr,
		},

		// auth-machine
		{
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthMachine.DebugAddr,
		},
		{
			EnvVars:     []string{"OCIS_MACHINE_AUTH_API_KEY", "STORAGE_AUTH_MACHINE_AUTH_API_KEY"},
			Destination: &cfg.Reva.AuthMachineConfig.MachineAuthAPIKey,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthMachine.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthMachine.GRPCAddr,
		},

		// frontend
		{
			EnvVars:     []string{"STORAGE_FRONTEND_DEBUG_ADDR"},
			Destination: &cfg.Reva.Frontend.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_TRANSFER_SECRET"},
			Destination: &cfg.Reva.TransferSecret,
		},
		{
			EnvVars:     []string{"STORAGE_CHUNK_FOLDER"},
			Destination: &cfg.Reva.OCDav.WebdavNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_WEBDAV_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.WebdavNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_DAV_FILES_NAMESPACE"},
			Destination: &cfg.Reva.OCDav.DavFilesNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_ARCHIVER_MAX_NUM_FILES"},
			Destination: &cfg.Reva.Archiver.MaxNumFiles,
		},
		{
			EnvVars:     []string{"STORAGE_ARCHIVER_MAX_SIZE"},
			Destination: &cfg.Reva.Archiver.MaxSize,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_NETWORK"},
			Destination: &cfg.Reva.Frontend.HTTPNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_HTTP_ADDR"},
			Destination: &cfg.Reva.Frontend.HTTPAddr,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_APPPROVIDER_PREFIX"},
			Destination: &cfg.Reva.Frontend.AppProviderPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_ARCHIVER_PREFIX"},
			Destination: &cfg.Reva.Frontend.ArchiverPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_DATAGATEWAY_PREFIX"},
			Destination: &cfg.Reva.Frontend.DatagatewayPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_FAVORITES"},
			Destination: &cfg.Reva.Frontend.Favorites,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCDAV_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCDavPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_SHARE_PREFIX"},
			Destination: &cfg.Reva.Frontend.OCSSharePrefix,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_HOME_NAMESPACE"},
			Destination: &cfg.Reva.Frontend.OCSHomeNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL"},
			Destination: &cfg.Reva.Frontend.OCSResourceInfoCacheTTL,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_CACHE_WARMUP_DRIVER"},
			Destination: &cfg.Reva.Frontend.OCSCacheWarmupDriver,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE"},
			Destination: &cfg.Reva.Frontend.OCSAdditionalInfoAttribute,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_DEFAULT_UPLOAD_PROTOCOL"},
			Destination: &cfg.Reva.DefaultUploadProtocol,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_MAX_CHUNK_SIZE"},
			Destination: &cfg.Reva.UploadMaxChunkSize,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE"},
			Destination: &cfg.Reva.UploadHTTPMethodOverride,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_CHECKSUM_PREFERRED_UPLOAD_TYPE"},
			Destination: &cfg.Reva.ChecksumPreferredUploadType,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_ARCHIVER_URL"},
			Destination: &cfg.Reva.Archiver.ArchiverURL,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_APP_PROVIDER_APPS_URL"},
			Destination: &cfg.Reva.AppProvider.AppsURL,
		},
		{
			EnvVars:     []string{"STORAGE_FRONTEND_APP_PROVIDER_OPEN_URL"},
			Destination: &cfg.Reva.AppProvider.OpenURL,
		},

		// gateway
		{
			EnvVars:     []string{"STORAGE_GATEWAY_DEBUG_ADDR"},
			Destination: &cfg.Reva.Gateway.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_TRANSFER_EXPIRES"},
			Destination: &cfg.Reva.TransferExpires,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_GRPC_NETWORK"},
			Destination: &cfg.Reva.Gateway.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_GRPC_ADDR"},
			Destination: &cfg.Reva.Gateway.GRPCAddr,
		},

		{
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageGrant,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF"},
			Destination: &cfg.Reva.Gateway.CommitShareToStorageRef,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_SHARE_FOLDER"},
			Destination: &cfg.Reva.Gateway.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN"},
			Destination: &cfg.Reva.Gateway.DisableHomeCreationOnLogin,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_HOME_MAPPING"},
			Destination: &cfg.Reva.Gateway.HomeMapping,
		},
		{
			EnvVars:     []string{"STORAGE_GATEWAY_ETAG_CACHE_TTL"},
			Destination: &cfg.Reva.Gateway.EtagCacheTTL,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BASIC_ENDPOINT"},
			Destination: &cfg.Reva.AuthBasic.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_BEARER_ENDPOINT"},
			Destination: &cfg.Reva.AuthBearer.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_AUTH_MACHINE_ENDPOINT"},
			Destination: &cfg.Reva.AuthMachine.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_DRIVER"},
			Destination: &cfg.Reva.StorageRegistry.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_HOME_PROVIDER"},
			Destination: &cfg.Reva.StorageRegistry.HomeProvider,
		},
		{
			EnvVars:     []string{"STORAGE_STORAGE_REGISTRY_JSON"},
			Destination: &cfg.Reva.StorageRegistry.JSON,
		},
		{
			EnvVars:     []string{"STORAGE_APP_REGISTRY_DRIVER"},
			Destination: &cfg.Reva.AppRegistry.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_APP_REGISTRY_MIMETYPES_JSON"},
			Destination: &cfg.Reva.AppRegistry.MimetypesJSON,
		},
		{
			EnvVars:     []string{"STORAGE_DATAGATEWAY_PUBLIC_URL"},
			Destination: &cfg.Reva.DataGateway.PublicURL,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Groups.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_ENDPOINT"},
			Destination: &cfg.Reva.Sharing.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_APPPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.AppProvider.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_ENDPOINT"},
			Destination: &cfg.Reva.StorageHome.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_ENDPOINT"},
			Destination: &cfg.Reva.StorageUsers.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageUsers.MountPath,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageUsers.MountID,
		},
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_ENDPOINT"},
			Destination: &cfg.Reva.StoragePublicLink.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
		},

		// groups
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.Groups.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_NETWORK"},
			Destination: &cfg.Reva.Groups.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_ADDR"},
			Destination: &cfg.Reva.Groups.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_DRIVER"},
			Destination: &cfg.Reva.Groups.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_GROUPPROVIDER_JSON"},
			Destination: &cfg.Reva.Groups.JSON,
		},
		{
			EnvVars:     []string{"STORAGE_GROUP_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.Groups.GroupMembersCacheExpiration,
		},

		// ldap
		{
			EnvVars:     []string{"STORAGE_LDAP_HOSTNAME"},
			Destination: &cfg.Reva.LDAP.Hostname,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_PORT"},
			Destination: &cfg.Reva.LDAP.Port,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_CACERT"},
			Destination: &cfg.Reva.LDAP.CACert,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_INSECURE"},
			Destination: &cfg.Reva.LDAP.Insecure,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_BASE_DN"},
			Destination: &cfg.Reva.LDAP.BaseDN,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_LOGINFILTER"},
			Destination: &cfg.Reva.LDAP.LoginFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USERATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.UserAttributeFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USERFINDFILTER"},
			Destination: &cfg.Reva.LDAP.UserFindFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USERGROUPFILTER"},
			Destination: &cfg.Reva.LDAP.UserGroupFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUPATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.GroupAttributeFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUPFINDFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFindFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUPMEMBERFILTER"},
			Destination: &cfg.Reva.LDAP.GroupMemberFilter,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_BIND_DN"},
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID"},
			Destination: &cfg.Reva.LDAP.UserSchema.UID,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.UserSchema.Mail,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.UserSchema.DisplayName,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.UserSchema.CN,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.UIDNumber,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.GIDNumber,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GID,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Mail,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.GroupSchema.DisplayName,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.GroupSchema.CN,
		},
		{
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GIDNumber,
		},

		// rest
		{
			EnvVars:     []string{"STORAGE_REST_CLIENT_ID"},
			Destination: &cfg.Reva.UserGroupRest.ClientID,
		},
		{
			EnvVars:     []string{"STORAGE_REST_CLIENT_SECRET"},
			Destination: &cfg.Reva.UserGroupRest.ClientSecret,
		},
		{
			EnvVars:     []string{"STORAGE_REST_REDIS_ADDRESS"},
			Destination: &cfg.Reva.UserGroupRest.RedisAddress,
		},
		{
			EnvVars:     []string{"STORAGE_REST_REDIS_USERNAME"},
			Destination: &cfg.Reva.UserGroupRest.RedisUsername,
		},
		{
			EnvVars:     []string{"STORAGE_REST_REDIS_PASSWORD"},
			Destination: &cfg.Reva.UserGroupRest.RedisPassword,
		},
		{
			EnvVars:     []string{"STORAGE_REST_ID_PROVIDER"},
			Destination: &cfg.Reva.UserGroupRest.IDProvider,
		},
		{
			EnvVars:     []string{"STORAGE_REST_API_BASE_URL"},
			Destination: &cfg.Reva.UserGroupRest.APIBaseURL,
		},
		{
			EnvVars:     []string{"STORAGE_REST_OIDC_TOKEN_ENDPOINT"},
			Destination: &cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
		},
		{
			EnvVars:     []string{"STORAGE_REST_TARGET_API"},
			Destination: &cfg.Reva.UserGroupRest.TargetAPI,
		},

		// secret
		{
			EnvVars:     []string{"OCIS_JWT_SECRET", "STORAGE_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},
		{
			EnvVars:     []string{"STORAGE_SKIP_USER_GROUPS_IN_TOKEN"},
			Destination: &cfg.Reva.SkipUserGroupsInToken,
		},

		// sharing
		{
			EnvVars:     []string{"STORAGE_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Sharing.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_GRPC_NETWORK"},
			Destination: &cfg.Reva.Sharing.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_GRPC_ADDR"},
			Destination: &cfg.Reva.Sharing.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_DRIVER"},
			Destination: &cfg.Reva.Sharing.UserDriver,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_DRIVER"},
			Destination: &cfg.Reva.Sharing.PublicDriver,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_JSON_FILE"},
			Destination: &cfg.Reva.Sharing.PublicJSONFile,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_PASSWORD_HASH_COST"},
			Destination: &cfg.Reva.Sharing.PublicPasswordHashCost,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_ENABLE_EXPIRED_SHARES_CLEANUP"},
			Destination: &cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_JANITOR_RUN_INTERVAL"},
			Destination: &cfg.Reva.Sharing.PublicJanitorRunInterval,
		},

		// shqringsql
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_USERNAME"},
			Destination: &cfg.Reva.Sharing.UserSQLUsername,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PASSWORD"},
			Destination: &cfg.Reva.Sharing.UserSQLPassword,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_HOST"},
			Destination: &cfg.Reva.Sharing.UserSQLHost,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PORT"},
			Destination: &cfg.Reva.Sharing.UserSQLPort,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_NAME"},
			Destination: &cfg.Reva.Sharing.UserSQLName,
		},

		// storage home
		{
			EnvVars:     []string{"STORAGE_HOME_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageHome.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageHome.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageHome.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageHome.HTTPNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageHome.HTTPAddr,
		},
		{
			EnvVars:     []string{"OCIS_STORAGE_READ_ONLY", "STORAGE_HOME_READ_ONLY"},
			Destination: &cfg.Reva.StorageHome.ReadOnly,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageHome.ExposeDataServer,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageHome.DataServerURL,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_HTTP_PREFIX"},
			Destination: &cfg.Reva.StorageHome.HTTPPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_HOME_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageHome.TempFolder,
		},

		// storage metadata
		{
			EnvVars:     []string{"STORAGE_METADATA_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_GRPC_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageMetadata.DataServerURL,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageMetadata.HTTPNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageMetadata.HTTPAddr,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageMetadata.TempFolder,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER"},
			Destination: &cfg.Reva.StorageMetadata.Driver,
		},

		// storage public link
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_DEBUG_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_GRPC_NETWORK"},
			Destination: &cfg.Reva.StoragePublicLink.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_GRPC_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.GRPCAddr,
		},

		// storage users
		{
			EnvVars:     []string{"STORAGE_USERS_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageUsers.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageUsers.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageUsers.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageUsers.HTTPNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageUsers.HTTPAddr,
		},
		{
			EnvVars:     []string{"OCIS_STORAGE_READ_ONLY", "STORAGE_USERS_READ_ONLY"},
			Destination: &cfg.Reva.StorageUsers.ReadOnly,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageUsers.ExposeDataServer,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageUsers.DataServerURL,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_HTTP_PREFIX"},
			Destination: &cfg.Reva.StorageUsers.HTTPPrefix,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageUsers.TempFolder,
		},

		// tracing
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED", "STORAGE_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE", "STORAGE_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT", "STORAGE_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR", "STORAGE_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"STORAGE_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},

		// users
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.Users.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_NETWORK"},
			Destination: &cfg.Reva.Users.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_ADDR"},
			Destination: &cfg.Reva.Users.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_DRIVER"},
			Destination: &cfg.Reva.Users.Driver,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_JSON"},
			Destination: &cfg.Reva.Users.JSON,
		},
		{
			EnvVars:     []string{"STORAGE_USER_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.Users.UserGroupsCacheExpiration,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBHOST"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBHost,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPORT"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBPort,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBNAME"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBName,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBUSER"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBUsername,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPASS"},
			Destination: &cfg.Reva.UserOwnCloudSQL.DBPassword,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_NOBODY"},
			Destination: &cfg.Reva.UserOwnCloudSQL.Nobody,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_USERNAME"},
			Destination: &cfg.Reva.UserOwnCloudSQL.JoinUsername,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_OWNCLOUDUUID"},
			Destination: &cfg.Reva.UserOwnCloudSQL.JoinOwnCloudUUID,
		},
		{
			EnvVars:     []string{"STORAGE_USERPROVIDER_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH"},
			Destination: &cfg.Reva.UserOwnCloudSQL.EnableMedialSearch,
		},

		// driver eos
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.Root,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHADOW_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.ShadowNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_UPLOADS_NAMESPACE"},
			Destination: &cfg.Reva.UserStorage.EOS.UploadsNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.EOS.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_BINARY"},
			Destination: &cfg.Reva.UserStorage.EOS.EosBinary,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_XRDCOPY_BINARY"},
			Destination: &cfg.Reva.UserStorage.EOS.XrdcopyBinary,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_MASTER_URL"},
			Destination: &cfg.Reva.UserStorage.EOS.MasterURL,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SLAVE_URL"},
			Destination: &cfg.Reva.UserStorage.EOS.SlaveURL,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_CACHE_DIRECTORY"},
			Destination: &cfg.Reva.UserStorage.EOS.CacheDirectory,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_ENABLE_LOGGING"},
			Destination: &cfg.Reva.UserStorage.EOS.EnableLogging,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SHOW_HIDDEN_SYSFILES"},
			Destination: &cfg.Reva.UserStorage.EOS.ShowHiddenSysFiles,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_FORCE_SINGLEUSER_MODE"},
			Destination: &cfg.Reva.UserStorage.EOS.ForceSingleUserMode,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_USE_KEYTAB"},
			Destination: &cfg.Reva.UserStorage.EOS.UseKeytab,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SEC_PROTOCOL"},
			Destination: &cfg.Reva.UserStorage.EOS.SecProtocol,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_KEYTAB"},
			Destination: &cfg.Reva.UserStorage.EOS.Keytab,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_SINGLE_USERNAME"},
			Destination: &cfg.Reva.UserStorage.EOS.SingleUsername,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_EOS_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.EOS.UserLayout,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.UserStorage.EOS.GatewaySVC,
		},

		// driver local
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.Local.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_USER_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.Local.UserLayout,
		},

		// driver ocis
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.OCIS.UserLayout,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.OCIS.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_SERVICE_USER_UUID"},
			Destination: &cfg.Reva.UserStorage.OCIS.ServiceUserUUID,
		},

		// driver owncloud
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_UPLOADINFO_DIR"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.UploadInfoDir,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_SCAN"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.Scan,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_REDIS_ADDR"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.Redis,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUD_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.OwnCloud.UserLayout,
		},

		// driver owncloud sql
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DATADIR"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.Root,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.UploadInfoDir,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.UserLayout,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBUSERNAME"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.DBUsername,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPASSWORD"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.DBPassword,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBHOST"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.DBHost,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPORT"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.DBPort,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBNAME"},
			Destination: &cfg.Reva.UserStorage.OwnCloudSQL.DBName,
		},

		// driver s3
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_REGION"},
			Destination: &cfg.Reva.UserStorage.S3.Region,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_ACCESS_KEY"},
			Destination: &cfg.Reva.UserStorage.S3.AccessKey,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_SECRET_KEY"},
			Destination: &cfg.Reva.UserStorage.S3.SecretKey,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_ENDPOINT"},
			Destination: &cfg.Reva.UserStorage.S3.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_BUCKET"},
			Destination: &cfg.Reva.UserStorage.S3.Bucket,
		},

		// driver s3ng
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ROOT"},
			Destination: &cfg.Reva.UserStorage.S3NG.Root,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.S3NG.UserLayout,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.S3NG.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_REGION"},
			Destination: &cfg.Reva.UserStorage.S3NG.Region,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ACCESS_KEY"},
			Destination: &cfg.Reva.UserStorage.S3NG.AccessKey,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_SECRET_KEY"},
			Destination: &cfg.Reva.UserStorage.S3NG.SecretKey,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ENDPOINT"},
			Destination: &cfg.Reva.UserStorage.S3NG.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_BUCKET"},
			Destination: &cfg.Reva.UserStorage.S3NG.Bucket,
		},

		// metadata driver eos
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_NAMESPACE"},
			Destination: &cfg.Reva.MetadataStorage.EOS.Root,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SHADOW_NAMESPACE"},
			Destination: &cfg.Reva.MetadataStorage.EOS.ShadowNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_UPLOADS_NAMESPACE"},
			Destination: &cfg.Reva.MetadataStorage.EOS.UploadsNamespace,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SHARE_FOLDER"},
			Destination: &cfg.Reva.MetadataStorage.EOS.ShareFolder,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_BINARY"},
			Destination: &cfg.Reva.MetadataStorage.EOS.EosBinary,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_XRDCOPY_BINARY"},
			Destination: &cfg.Reva.MetadataStorage.EOS.XrdcopyBinary,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_MASTER_URL"},
			Destination: &cfg.Reva.MetadataStorage.EOS.MasterURL,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SLAVE_URL"},
			Destination: &cfg.Reva.MetadataStorage.EOS.SlaveURL,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_CACHE_DIRECTORY"},
			Destination: &cfg.Reva.MetadataStorage.EOS.CacheDirectory,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_ENABLE_LOGGING"},
			Destination: &cfg.Reva.MetadataStorage.EOS.EnableLogging,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SHOW_HIDDEN_SYSFILES"},
			Destination: &cfg.Reva.MetadataStorage.EOS.ShowHiddenSysFiles,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_FORCE_SINGLEUSER_MODE"},
			Destination: &cfg.Reva.MetadataStorage.EOS.ForceSingleUserMode,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_USE_KEYTAB"},
			Destination: &cfg.Reva.MetadataStorage.EOS.UseKeytab,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SEC_PROTOCOL"},
			Destination: &cfg.Reva.MetadataStorage.EOS.SecProtocol,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_KEYTAB"},
			Destination: &cfg.Reva.MetadataStorage.EOS.Keytab,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_SINGLE_USERNAME"},
			Destination: &cfg.Reva.MetadataStorage.EOS.SingleUsername,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_EOS_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.EOS.UserLayout,
		},
		{
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.MetadataStorage.EOS.GatewaySVC,
		},

		// metadata local driver
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.Local.Root,
		},

		// metadata ocis driver
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.UserLayout,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_SERVICE_USER_UUID"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.ServiceUserUUID,
		},

		// metadata driver s3
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3_REGION"},
			Destination: &cfg.Reva.MetadataStorage.S3.Region,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3_ACCESS_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3.AccessKey,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3_SECRET_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3.SecretKey,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3_ENDPOINT"},
			Destination: &cfg.Reva.MetadataStorage.S3.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3_BUCKET"},
			Destination: &cfg.Reva.MetadataStorage.S3.Bucket,
		},

		// driver s3ng
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Root,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.UserLayout,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_REGION"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Region,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ACCESS_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.AccessKey,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_SECRET_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.SecretKey,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ENDPOINT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_BUCKET"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Bucket,
		},
	}
}
