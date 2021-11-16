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
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// Gateway defines the available gateway configuration.
type Gateway struct {
	Port
	CommitShareToStorageGrant  bool   `mapstructure:"commit_share_to_storage_grant"`
	CommitShareToStorageRef    bool   `mapstructure:"commit_share_to_storage_ref"`
	DisableHomeCreationOnLogin bool   `mapstructure:"disable_home_creation_on_login"`
	ShareFolder                string `mapstructure:"share_folder"`
	LinkGrants                 string `mapstructure:"link_grants"`
	HomeMapping                string `mapstructure:"home_mapping"`
	EtagCacheTTL               int    `mapstructure:"etag_cache_ttl"`
}

// StorageRegistry defines the available storage registry configuration
type StorageRegistry struct {
	Driver string `mapstructure:"driver"`
	// HomeProvider is the path in the global namespace that the static storage registry uses to determine the home storage
	HomeProvider string   `mapstructure:"home_provider"`
	Rules        []string `mapstructure:"rules"`
	JSON         string   `mapstructure:"json"`
}

// AppRegistry defines the available app registry configuration
type AppRegistry struct {
	Driver        string `mapstructure:"driver"`
	MimetypesJSON string `mapstructure:"mime_types_json"`
}

// AppProvider defines the available app provider configuration
type AppProvider struct {
	Port
	ExternalAddr string     `mapstructure:"external_addr"`
	Driver       string     `mapstructure:"driver"`
	WopiDriver   WopiDriver `mapstructure:"wopi_driver"`
	AppsURL      string     `mapstructure:"apps_url"`
	OpenURL      string     `mapstructure:"open_url"`
}

type WopiDriver struct {
	AppAPIKey      string `mapstructure:"app_api_key"`
	AppDesktopOnly bool   `mapstructure:"app_desktop_only"`
	AppIconURI     string `mapstructure:"app_icon_uri"`
	AppInternalURL string `mapstructure:"app_internal_url"`
	AppName        string `mapstructure:"app_name"`
	AppURL         string `mapstructure:"app_url"`
	Insecure       bool   `mapstructure:"insecure"`
	IopSecret      string `mapstructure:"ipo_secret"`
	JWTSecret      string `mapstructure:"jwt_secret"`
	WopiURL        string `mapstructure:"wopi_url"`
}

// Sharing defines the available sharing configuration.
type Sharing struct {
	Port
	UserDriver                       string `mapstructure:"user_driver"`
	UserJSONFile                     string `mapstructure:"user_json_file"`
	UserSQLUsername                  string `mapstructure:"user_sql_username"`
	UserSQLPassword                  string `mapstructure:"user_sql_password"`
	UserSQLHost                      string `mapstructure:"user_sql_host"`
	UserSQLPort                      int    `mapstructure:"user_sql_port"`
	UserSQLName                      string `mapstructure:"user_sql_name"`
	PublicDriver                     string `mapstructure:"public_driver"`
	PublicJSONFile                   string `mapstructure:"public_json_file"`
	PublicPasswordHashCost           int    `mapstructure:"public_password_hash_cost"`
	PublicEnableExpiredSharesCleanup bool   `mapstructure:"public_enable_expired_shares_cleanup"`
	PublicJanitorRunInterval         int    `mapstructure:"public_janitor_run_interval"`
	UserStorageMountID               string `mapstructure:"user_storage_mount_id"`
}

// Port defines the available port configuration.
type Port struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string `mapstructure:"max_cpus"`
	LogLevel string `mapstructure:"log_level"`
	// GRPCNetwork can be tcp, udp or unix
	GRPCNetwork string `mapstructure:"grpc_network"`
	// GRPCAddr to listen on, hostname:port (0.0.0.0:9999 for all interfaces) or socket (/var/run/reva/sock)
	GRPCAddr string `mapstructure:"grpc_addr"`
	// Protocol can be grpc or http
	// HTTPNetwork can be tcp, udp or unix
	HTTPNetwork string `mapstructure:"http_network"`
	// HTTPAddr to listen on, hostname:port (0.0.0.0:9100 for all interfaces) or socket (/var/run/reva/sock)
	HTTPAddr string `mapstructure:"http_addr"`
	// Protocol can be grpc or http
	Protocol string `mapstructure:"protocol"`
	// Endpoint is used by the gateway and registries (eg localhost:9100 or cloud.example.com)
	Endpoint string `mapstructure:"endpoint"`
	// DebugAddr for the debug endpoint to bind to
	DebugAddr string `mapstructure:"debug_addr"`
	// Services can be used to give a list of services that should be started on this port
	Services []string `mapstructure:"services"`
	// Config can be used to configure the reva instance.
	// Services and Protocol will be ignored if this is used
	Config map[string]interface{} `mapstructure:"config"`

	// Context allows for context cancellation and propagation
	Context context.Context

	// Supervised is used when running under an oCIS runtime supervision tree
	Supervised bool // deprecated
}

// Users defines the available users configuration.
type Users struct {
	Port
	Driver                    string `mapstucture:"driver"`
	JSON                      string `mapstucture:"json"`
	UserGroupsCacheExpiration int    `mapstucture:"user_groups_cache_expiration"`
}

// AuthMachineConfig defines the available configuration for the machine auth driver.
type AuthMachineConfig struct {
	MachineAuthAPIKey string `mapstucture:"machine_auth_api_key"`
}

// Groups defines the available groups configuration.
type Groups struct {
	Port
	Driver                      string `mapstucture:"driver"`
	JSON                        string `mapstucture:"json"`
	GroupMembersCacheExpiration int    `mapstucture:"group_members_cache_expiration"`
}

// FrontendPort defines the available frontend configuration.
type FrontendPort struct {
	Port

	AppProviderInsecure        bool       `mapstucture:"app_provider_insecure"`
	AppProviderPrefix          string     `mapstucture:"app_provider_prefix"`
	ArchiverInsecure           bool       `mapstucture:"archiver_insecure"`
	ArchiverPrefix             string     `mapstucture:"archiver_prefix"`
	DatagatewayPrefix          string     `mapstucture:"data_gateway_prefix"`
	Favorites                  bool       `mapstucture:"favorites"`
	OCDavInsecure              bool       `mapstucture:"ocdav_insecure"`
	OCDavPrefix                string     `mapstucture:"ocdav_prefix"`
	OCSPrefix                  string     `mapstucture:"ocs_prefix"`
	OCSSharePrefix             string     `mapstucture:"ocs_share_prefix"`
	OCSHomeNamespace           string     `mapstucture:"ocs_home_namespace"`
	PublicURL                  string     `mapstucture:"public_url"`
	OCSCacheWarmupDriver       string     `mapstucture:"ocs_cache_warmup_driver"`
	OCSAdditionalInfoAttribute string     `mapstucture:"ocs_additional_info_attribute"`
	OCSResourceInfoCacheTTL    int        `mapstucture:"ocs_resource_info_cache_ttl"`
	Middleware                 Middleware `mapstucture:"middleware"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `mapstructure:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `mapstructure:"credentials_by_user_agenr"`
}

// DataGatewayPort has a public url
type DataGatewayPort struct {
	Port
	PublicURL string `mapstructure:""`
}

type DataProvider struct {
	Insecure bool `mapstructure:"insecure"`
}

// StoragePort defines the available storage configuration.
type StoragePort struct {
	Port
	Driver           string `mapstructure:"driver"`
	MountPath        string `mapstructure:"mount_path"`
	MountID          string `mapstructure:"mount_id"`
	ExposeDataServer bool   `mapstructure:"expose_data_server"`
	// url the data gateway will use to route requests
	DataServerURL string `mapstructure:"data_server_url"`

	// for HTTP ports with only one http service
	HTTPPrefix   string       `mapstructure:"http_prefix"`
	TempFolder   string       `mapstructure:"temp_folder"`
	ReadOnly     bool         `mapstructure:"read_only"`
	DataProvider DataProvider `mapstructure:"data_provider"`
}

// PublicStorage configures a public storage provider
type PublicStorage struct {
	StoragePort

	PublicShareProviderAddr string `mapstructure:"public_share_provider_addr"`
	UserProviderAddr        string `mapstructure:"user_provider_addr"`
}

// StorageConfig combines all available storage driver configuration parts.
type StorageConfig struct {
	EOS         DriverEOS         `mapstructure:"eos"`
	Local       DriverCommon      `mapstructure:"local"`
	OwnCloud    DriverOwnCloud    `mapstructure:"owncloud"`
	OwnCloudSQL DriverOwnCloudSQL `mapstructure:"owncloud_sql"`
	S3          DriverS3          `mapstructure:"s3"`
	S3NG        DriverS3NG        `mapstructure:"s3ng"`
	OCIS        DriverOCIS        `mapstructure:"ocis"`
}

// DriverCommon defines common driver configuration options.
type DriverCommon struct {
	// Root is the absolute path to the location of the data
	Root string `mapstructure:"root"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `mapstructure:"share_folder"`
	// UserLayout contains the template used to construct
	// the internal path, eg: `{{substr 0 1 .Username}}/{{.Username}}`
	UserLayout string `mapstructure:"user_layout"`
	// EnableHome enables the creation of home directories.
	EnableHome bool `mapstructure:"enable_home"`
}

// DriverEOS defines the available EOS driver configuration.
type DriverEOS struct {
	DriverCommon

	// ShadowNamespace for storing shadow data
	ShadowNamespace string `mapstructure:"shadow_namespace"`

	// UploadsNamespace for storing upload data
	UploadsNamespace string `mapstructure:"uploads_namespace"`

	// Location of the eos binary.
	// Default is /usr/bin/eos.
	EosBinary string `mapstructure:"eos_binary"`

	// Location of the xrdcopy binary.
	// Default is /usr/bin/xrdcopy.
	XrdcopyBinary string `mapstructure:"xrd_copy_binary"`

	// URL of the Master EOS MGM.
	// Default is root://eos-example.org
	MasterURL string `mapstructure:"master_url"`

	// URI of the EOS MGM grpc server
	// Default is empty
	GrpcURI string `mapstructure:"grpc_uri"`

	// URL of the Slave EOS MGM.
	// Default is root://eos-example.org
	SlaveURL string `mapstructure:"slave_url"`

	// Location on the local fs where to store reads.
	// Defaults to os.TempDir()
	CacheDirectory string `mapstructure:"cache_directory"`

	// Enables logging of the commands executed
	// Defaults to false
	EnableLogging bool `mapstructure:"enable_logging"`

	// ShowHiddenSysFiles shows internal EOS files like
	// .sys.v# and .sys.a# files.
	ShowHiddenSysFiles bool `mapstructure:"shadow_hidden_files"`

	// ForceSingleUserMode will force connections to EOS to use SingleUsername
	ForceSingleUserMode bool `mapstructure:"force_single_user_mode"`

	// UseKeyTabAuth changes will authenticate requests by using an EOS keytab.
	UseKeytab bool `mapstructure:"user_keytab"`

	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string `mapstructure:"sec_protocol"`

	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string `mapstructure:"keytab"`

	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string `mapstructure:"single_username"`

	// gateway service to use for uid lookups
	GatewaySVC string `mapstructure:"gateway_svc"`
}

// DriverOCIS defines the available oCIS storage driver configuration.
type DriverOCIS struct {
	DriverCommon

	ServiceUserUUID string `mapstructure:"service_user_uuid"`
}

// DriverOwnCloud defines the available ownCloud storage driver configuration.
type DriverOwnCloud struct {
	DriverCommon

	UploadInfoDir string `mapstructure:"upload_info_dir"`
	Redis         string `mapstructure:"redis"`
	Scan          bool   `mapstructure:"scan"`
}

// DriverOwnCloudSQL defines the available ownCloudSQL storage driver configuration.
type DriverOwnCloudSQL struct {
	DriverCommon

	UploadInfoDir string `mapstructure:"upload_info_dir"`
	DBUsername    string `mapstructure:"db_username"`
	DBPassword    string `mapstructure:"db_password"`
	DBHost        string `mapstructure:"db_host"`
	DBPort        int    `mapstructure:"db_port"`
	DBName        string `mapstructure:"db_name"`
}

// DriverS3 defines the available S3 storage driver configuration.
type DriverS3 struct {
	DriverCommon

	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
}

// DriverS3NG defines the available s3ng storage driver configuration.
type DriverS3NG struct {
	DriverCommon

	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
}

// OIDC defines the available OpenID Connect configuration.
type OIDC struct {
	Issuer   string `mapstructure:"issuer"`
	Insecure bool   `mapstructure:"insecure"`
	IDClaim  string `mapstructure:"id_claim"`
	UIDClaim string `mapstructure:"uid_claim"`
	GIDClaim string `mapstructure:"gid_claim"`
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname             string          `mapstructure:"hostname"`
	Port                 int             `mapstructure:"port"`
	CACert               string          `mapstructure:"ca_cert"`
	Insecure             bool            `mapstructure:"insecure"`
	BaseDN               string          `mapstructure:"base_dn"`
	LoginFilter          string          `mapstructure:"login_filter"`
	UserFilter           string          `mapstructure:"user_filter"`
	UserAttributeFilter  string          `mapstructure:"user_attribute_filter"`
	UserFindFilter       string          `mapstructure:"user_find_filter"`
	UserGroupFilter      string          `mapstructure:"user_group_filter"`
	GroupFilter          string          `mapstructure:"group_filter"`
	GroupAttributeFilter string          `mapstructure:"group_attribute_filter"`
	GroupFindFilter      string          `mapstructure:"group_finder_filter"`
	GroupMemberFilter    string          `mapstructure:"group_member_filter"`
	BindDN               string          `mapstructure:"bind_dn"`
	BindPassword         string          `mapstructure:"bind_password"`
	IDP                  string          `mapstructure:"idp"`
	UserSchema           LDAPUserSchema  `mapstructure:"user_schema"`
	GroupSchema          LDAPGroupSchema `mapstructure:"group_schema"`
}

// UserGroupRest defines the REST driver specification for user and group resolution.
type UserGroupRest struct {
	ClientID          string `mapstructure:"client_id"`
	ClientSecret      string `mapstructure:"client_secret"`
	RedisAddress      string `mapstructure:"redis_address"`
	RedisUsername     string `mapstructure:"redis_username"`
	RedisPassword     string `mapstructure:"redis_password"`
	IDProvider        string `mapstructure:"idp_provider"`
	APIBaseURL        string `mapstructure:"api_base_url"`
	OIDCTokenEndpoint string `mapstructure:"oidc_token_endpoint"`
	TargetAPI         string `mapstructure:"target_api"`
}

// UserOwnCloudSQL defines the available ownCloudSQL user provider configuration.
type UserOwnCloudSQL struct {
	DBUsername         string `mapstructure:"db_username"`
	DBPassword         string `mapstructure:"db_password"`
	DBHost             string `mapstructure:"db_host"`
	DBPort             int    `mapstructure:"db_port"`
	DBName             string `mapstructure:"db_name"`
	Idp                string `mapstructure:"idp"`
	Nobody             int64  `mapstructure:"nobody"`
	JoinUsername       bool   `mapstructure:"join_username"`
	JoinOwnCloudUUID   bool   `mapstructure:"join_owncloud_uuid"`
	EnableMedialSearch bool   `mapstructure:"enable_medial_search"`
}

// LDAPUserSchema defines the available ldap user schema configuration.
type LDAPUserSchema struct {
	UID         string `mapstructure:"uid"`
	Mail        string `mapstructure:"mail"`
	DisplayName string `mapstructure:"display_name"`
	CN          string `mapstructure:"cn"`
	UIDNumber   string `mapstructure:"uid_number"`
	GIDNumber   string `mapstructure:"gid_number"`
}

// LDAPGroupSchema defines the available ldap group schema configuration.
type LDAPGroupSchema struct {
	GID         string `mapstructure:"gid"`
	Mail        string `mapstructure:"mail"`
	DisplayName string `mapstructure:"display_name"`
	CN          string `mapstructure:"cn"`
	GIDNumber   string `mapstructure:"gid_number"`
}

// OCDav defines the available ocdav configuration.
type OCDav struct {
	WebdavNamespace   string `mapstructure:"webdav_namespace"`
	DavFilesNamespace string `mapstructure:"dav_files_namespace"`
}

// Archiver defines the available archiver configuration.
type Archiver struct {
	MaxNumFiles int64  `mapstructure:"max_num_files"`
	MaxSize     int64  `mapstructure:"max_size"`
	ArchiverURL string `mapstructure:"archiver_url"`
}

// Reva defines the available reva configuration.
type Reva struct {
	// JWTSecret used to sign jwt tokens between services
	JWTSecret             string          `mapstructure:"jwt_secret"`
	SkipUserGroupsInToken bool            `mapstructure:"skip_user_grooups_in_token"`
	TransferSecret        string          `mapstructure:"transfer_secret"`
	TransferExpires       int             `mapstructure:"transfer_expires"`
	OIDC                  OIDC            `mapstructure:"oidc"`
	LDAP                  LDAP            `mapstructure:"ldap"`
	UserGroupRest         UserGroupRest   `mapstructure:"user_group_rest"`
	UserOwnCloudSQL       UserOwnCloudSQL `mapstructure:"user_owncloud_sql"`
	OCDav                 OCDav           `mapstructure:"ocdav"`
	Archiver              Archiver        `mapstructure:"archiver"`
	UserStorage           StorageConfig   `mapstructure:"user_storage"`
	MetadataStorage       StorageConfig   `mapstructure:"metadata_storage"`
	// Ports are used to configure which services to start on which port
	Frontend          FrontendPort      `mapstructure:"frontend"`
	DataGateway       DataGatewayPort   `mapstructure:"data_gateway"`
	Gateway           Gateway           `mapstructure:"gateway"`
	StorageRegistry   StorageRegistry   `mapstructure:"storage_registry"`
	AppRegistry       AppRegistry       `mapstructure:"app_registry"`
	Users             Users             `mapstructure:"users"`
	Groups            Groups            `mapstructure:"groups"`
	AuthProvider      Users             `mapstructure:"auth_provider"`
	AuthBasic         Port              `mapstructure:"auth_basic"`
	AuthBearer        Port              `mapstructure:"auth_bearer"`
	AuthMachine       Port              `mapstructure:"auth_machine"`
	AuthMachineConfig AuthMachineConfig `mapstructure:"auth_machine_config"`
	Sharing           Sharing           `mapstructure:"sharing"`
	StorageHome       StoragePort       `mapstructure:"storage_home"`
	StorageUsers      StoragePort       `mapstructure:"storage_users"`
	StoragePublicLink PublicStorage     `mapstructure:"storage_public_link"`
	StorageMetadata   StoragePort       `mapstructure:"storage_metadata"`
	AppProvider       AppProvider       `mapstructure:"app_provider"`
	// Configs can be used to configure the reva instance.
	// Services and Ports will be ignored if this is used
	Configs map[string]interface{} `mapstructure:"configs"`
	// chunking and resumable upload config (TUS)
	UploadMaxChunkSize       int    `mapstructure:"uppload_max_chunk_size"`
	UploadHTTPMethodOverride string `mapstructure:"upload_http_method_override"`
	// checksumming capabilities
	ChecksumSupportedTypes      []string `mapstructure:"checksum_supported_types"`
	ChecksumPreferredUploadType string   `mapstructure:"checksum_preferred_upload_type"`
	DefaultUploadProtocol       string   `mapstructure:"default_upload_protocol"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `mapstructure:"path"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File    string      `mapstructure:"file"`
	Log     *shared.Log `mapstructure:"log"`
	Debug   Debug       `mapstructure:"debug"`
	Reva    Reva        `mapstructure:"reva"`
	Tracing Tracing     `mapstructure:"tracing"`
	Asset   Asset       `mapstructure:"asset"`
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
	}
}
