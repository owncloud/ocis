package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `yaml:"level"`
	Pretty bool   `yaml:"pretty"`
	Color  bool   `yaml:"color"`
	File   string `yaml:"file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr"`
	Token  string `yaml:"token"`
	Pprof  bool   `yaml:"pprof"`
	Zpages bool   `yaml:"zpages"`
}

// Gateway defines the available gateway configuration.
type Gateway struct {
	Port
	CommitShareToStorageGrant  bool   `yaml:"commit_share_to_storage_grant"`
	CommitShareToStorageRef    bool   `yaml:"commit_share_to_storage_ref"`
	DisableHomeCreationOnLogin bool   `yaml:"disable_home_creation_on_login"`
	ShareFolder                string `yaml:"share_folder"`
	LinkGrants                 string `yaml:"link_grants"`
	HomeMapping                string `yaml:"home_mapping"`
	EtagCacheTTL               int    `yaml:"etag_cache_ttl"`
}

// StorageRegistry defines the available storage registry configuration
type StorageRegistry struct {
	Driver string `yaml:"driver"`
	// HomeProvider is the path in the global namespace that the static storage registry uses to determine the home storage
	HomeProvider string   `yaml:"home_provider"`
	Rules        []string `yaml:"rules"`
	JSON         string   `yaml:"json"`
}

// AppRegistry defines the available app registry configuration
type AppRegistry struct {
	Driver        string `yaml:"driver"`
	MimetypesJSON string `yaml:"mime_types_json"`
}

// AppProvider defines the available app provider configuration
type AppProvider struct {
	Port
	ExternalAddr string     `yaml:"external_addr"`
	Driver       string     `yaml:"driver"`
	WopiDriver   WopiDriver `yaml:"wopi_driver"`
	AppsURL      string     `yaml:"apps_url"`
	OpenURL      string     `yaml:"open_url"`
	NewURL       string     `yaml:"new_url"`
}

type WopiDriver struct {
	AppAPIKey      string `yaml:"app_api_key"`
	AppDesktopOnly bool   `yaml:"app_desktop_only"`
	AppIconURI     string `yaml:"app_icon_uri"`
	AppInternalURL string `yaml:"app_internal_url"`
	AppName        string `yaml:"app_name"`
	AppURL         string `yaml:"app_url"`
	Insecure       bool   `yaml:"insecure"`
	IopSecret      string `yaml:"ipo_secret"`
	JWTSecret      string `yaml:"jwt_secret"`
	WopiURL        string `yaml:"wopi_url"`
}

// Sharing defines the available sharing configuration.
type Sharing struct {
	Port
	UserDriver                       string `yaml:"user_driver"`
	UserJSONFile                     string `yaml:"user_json_file"`
	CS3ProviderAddr                  string `yaml:"provider_addr"`
	CS3ServiceUser                   string `yaml:"service_user_id"`
	CS3ServiceUserIdp                string `yaml:"service_user_idp"`
	UserSQLUsername                  string `yaml:"user_sql_username"`
	UserSQLPassword                  string `yaml:"user_sql_password"`
	UserSQLHost                      string `yaml:"user_sql_host"`
	UserSQLPort                      int    `yaml:"user_sql_port"`
	UserSQLName                      string `yaml:"user_sql_name"`
	PublicDriver                     string `yaml:"public_driver"`
	PublicJSONFile                   string `yaml:"public_json_file"`
	PublicPasswordHashCost           int    `yaml:"public_password_hash_cost"`
	PublicEnableExpiredSharesCleanup bool   `yaml:"public_enable_expired_shares_cleanup"`
	PublicJanitorRunInterval         int    `yaml:"public_janitor_run_interval"`
	UserStorageMountID               string `yaml:"user_storage_mount_id"`
	Events                           Events `yaml:"events"`
}

type Events struct {
	Address   string `yaml:"address"`
	ClusterID string `yaml:"cluster_id"`
}

// Port defines the available port configuration.
type Port struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string `yaml:"max_cpus"`
	LogLevel string `yaml:"log_level"`
	// GRPCNetwork can be tcp, udp or unix
	GRPCNetwork string `yaml:"grpc_network"`
	// GRPCAddr to listen on, hostname:port (0.0.0.0:9999 for all interfaces) or socket (/var/run/reva/sock)
	GRPCAddr string `yaml:"grpc_addr"`
	// Protocol can be grpc or http
	// HTTPNetwork can be tcp, udp or unix
	HTTPNetwork string `yaml:"http_network"`
	// HTTPAddr to listen on, hostname:port (0.0.0.0:9100 for all interfaces) or socket (/var/run/reva/sock)
	HTTPAddr string `yaml:"http_addr"`
	// Protocol can be grpc or http
	Protocol string `yaml:"protocol"`
	// Endpoint is used by the gateway and registries (eg localhost:9100 or cloud.example.com)
	Endpoint string `yaml:"endpoint"`
	// DebugAddr for the debug endpoint to bind to
	DebugAddr string `yaml:"debug_addr"`
	// Services can be used to give a list of services that should be started on this port
	Services []string `yaml:"services"`
	// Config can be used to configure the reva instance.
	// Services and Protocol will be ignored if this is used
	Config map[string]interface{} `yaml:"config"`

	// Context allows for context cancellation and propagation
	Context context.Context

	// Supervised is used when running under an oCIS runtime supervision tree
	Supervised bool // deprecated // TODO: delete me
}

// Users defines the available users configuration.
type Users struct {
	Port
	Driver                    string `yaml:"driver"`
	JSON                      string `yaml:"json"`
	UserGroupsCacheExpiration int    `yaml:"user_groups_cache_expiration"`
}

// AuthMachineConfig defines the available configuration for the machine auth driver.
type AuthMachineConfig struct {
	MachineAuthAPIKey string `yaml:"machine_auth_api_key"`
}

// Groups defines the available groups configuration.
type Groups struct {
	Port
	Driver                      string `yaml:"driver"`
	JSON                        string `yaml:"json"`
	GroupMembersCacheExpiration int    `yaml:"group_members_cache_expiration"`
}

// FrontendPort defines the available frontend configuration.
type FrontendPort struct {
	Port

	AppProviderInsecure        bool       `yaml:"app_provider_insecure"`
	AppProviderPrefix          string     `yaml:"app_provider_prefix"`
	ArchiverInsecure           bool       `yaml:"archiver_insecure"`
	ArchiverPrefix             string     `yaml:"archiver_prefix"`
	DatagatewayPrefix          string     `yaml:"data_gateway_prefix"`
	Favorites                  bool       `yaml:"favorites"`
	ProjectSpaces              bool       `yaml:"project_spaces"`
	OCDavInsecure              bool       `yaml:"ocdav_insecure"`
	OCDavPrefix                string     `yaml:"ocdav_prefix"`
	OCSPrefix                  string     `yaml:"ocs_prefix"`
	OCSSharePrefix             string     `yaml:"ocs_share_prefix"`
	OCSHomeNamespace           string     `yaml:"ocs_home_namespace"`
	PublicURL                  string     `yaml:"public_url"`
	OCSCacheWarmupDriver       string     `yaml:"ocs_cache_warmup_driver"`
	OCSAdditionalInfoAttribute string     `yaml:"ocs_additional_info_attribute"`
	OCSResourceInfoCacheTTL    int        `yaml:"ocs_resource_info_cache_ttl"`
	Middleware                 Middleware `yaml:"middleware"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agenr"`
}

// DataGatewayPort has a public url
type DataGatewayPort struct {
	Port
	PublicURL string `yaml:""`
}

type DataProvider struct {
	Insecure bool `yaml:"insecure"`
}

// StoragePort defines the available storage configuration.
type StoragePort struct {
	Port
	Driver           string `yaml:"driver"`
	MountID          string `yaml:"mount_id"`
	AlternativeID    string `yaml:"alternative_id"`
	ExposeDataServer bool   `yaml:"expose_data_server"`
	// url the data gateway will use to route requests
	DataServerURL string `yaml:"data_server_url"`

	// for HTTP ports with only one http service
	HTTPPrefix      string       `yaml:"http_prefix"`
	TempFolder      string       `yaml:"temp_folder"`
	ReadOnly        bool         `yaml:"read_only"`
	DataProvider    DataProvider `yaml:"data_provider"`
	GatewayEndpoint string       `yaml:"gateway_endpoint"`
}

// PublicStorage configures a public storage provider
type PublicStorage struct {
	StoragePort

	PublicShareProviderAddr string `yaml:"public_share_provider_addr"`
	UserProviderAddr        string `yaml:"user_provider_addr"`
}

// StorageConfig combines all available storage driver configuration parts.
type StorageConfig struct {
	EOS         DriverEOS         `yaml:"eos"`
	Local       DriverCommon      `yaml:"local"`
	OwnCloudSQL DriverOwnCloudSQL `yaml:"owncloud_sql"`
	S3          DriverS3          `yaml:"s3"`
	S3NG        DriverS3NG        `yaml:"s3ng"`
	OCIS        DriverOCIS        `yaml:"ocis"`
}

// DriverCommon defines common driver configuration options.
type DriverCommon struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	// UserLayout contains the template used to construct
	// the internal path, eg: `{{substr 0 1 .Username}}/{{.Username}}`
	UserLayout string `yaml:"user_layout"`
	// EnableHome enables the creation of home directories.
	EnableHome bool `yaml:"enable_home"`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template"`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template"`
}

// DriverEOS defines the available EOS driver configuration.
type DriverEOS struct {
	DriverCommon

	// ShadowNamespace for storing shadow data
	ShadowNamespace string `yaml:"shadow_namespace"`

	// UploadsNamespace for storing upload data
	UploadsNamespace string `yaml:"uploads_namespace"`

	// Location of the eos binary.
	// Default is /usr/bin/eos.
	EosBinary string `yaml:"eos_binary"`

	// Location of the xrdcopy binary.
	// Default is /usr/bin/xrdcopy.
	XrdcopyBinary string `yaml:"xrd_copy_binary"`

	// URL of the Master EOS MGM.
	// Default is root://eos-example.org
	MasterURL string `yaml:"master_url"`

	// URI of the EOS MGM grpc server
	// Default is empty
	GrpcURI string `yaml:"grpc_uri"`

	// URL of the Slave EOS MGM.
	// Default is root://eos-example.org
	SlaveURL string `yaml:"slave_url"`

	// Location on the local fs where to store reads.
	// Defaults to os.TempDir()
	CacheDirectory string `yaml:"cache_directory"`

	// Enables logging of the commands executed
	// Defaults to false
	EnableLogging bool `yaml:"enable_logging"`

	// ShowHiddenSysFiles shows internal EOS files like
	// .sys.v# and .sys.a# files.
	ShowHiddenSysFiles bool `yaml:"shadow_hidden_files"`

	// ForceSingleUserMode will force connections to EOS to use SingleUsername
	ForceSingleUserMode bool `yaml:"force_single_user_mode"`

	// UseKeyTabAuth changes will authenticate requests by using an EOS keytab.
	UseKeytab bool `yaml:"user_keytab"`

	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string `yaml:"sec_protocol"`

	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string `yaml:"keytab"`

	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string `yaml:"single_username"`

	// gateway service to use for uid lookups
	GatewaySVC string `yaml:"gateway_svc"`
}

// DriverOCIS defines the available oCIS storage driver configuration.
type DriverOCIS struct {
	DriverCommon
}

// DriverOwnCloudSQL defines the available ownCloudSQL storage driver configuration.
type DriverOwnCloudSQL struct {
	DriverCommon

	UploadInfoDir string `yaml:"upload_info_dir"`
	DBUsername    string `yaml:"db_username"`
	DBPassword    string `yaml:"db_password"`
	DBHost        string `yaml:"db_host"`
	DBPort        int    `yaml:"db_port"`
	DBName        string `yaml:"db_name"`
}

// DriverS3 defines the available S3 storage driver configuration.
type DriverS3 struct {
	DriverCommon

	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
}

// DriverS3NG defines the available s3ng storage driver configuration.
type DriverS3NG struct {
	DriverCommon

	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
}

// OIDC defines the available OpenID Connect configuration.
type OIDC struct {
	Issuer   string `yaml:"issuer"`
	Insecure bool   `yaml:"insecure"`
	IDClaim  string `yaml:"id_claim"`
	UIDClaim string `yaml:"uid_claim"`
	GIDClaim string `yaml:"gid_claim"`
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	URI              string          `yaml:"uri"`
	CACert           string          `yaml:"ca_cert"`
	Insecure         bool            `yaml:"insecure"`
	UserBaseDN       string          `yaml:"user_base_dn"`
	GroupBaseDN      string          `yaml:"group_base_dn"`
	UserScope        string          `yaml:"user_scope"`
	GroupScope       string          `yaml:"group_scope"`
	UserObjectClass  string          `yaml:"user_objectclass"`
	GroupObjectClass string          `yaml:"group_objectclass"`
	UserFilter       string          `yaml:"user_filter"`
	GroupFilter      string          `yaml:"group_filter"`
	LoginAttributes  []string        `yaml:"login_attributes"`
	BindDN           string          `yaml:"bind_dn"`
	BindPassword     string          `yaml:"bind_password"`
	IDP              string          `yaml:"idp"`
	UserSchema       LDAPUserSchema  `yaml:"user_schema"`
	GroupSchema      LDAPGroupSchema `yaml:"group_schema"`
}

// UserGroupRest defines the REST driver specification for user and group resolution.
type UserGroupRest struct {
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	RedisAddress      string `yaml:"redis_address"`
	RedisUsername     string `yaml:"redis_username"`
	RedisPassword     string `yaml:"redis_password"`
	IDProvider        string `yaml:"idp_provider"`
	APIBaseURL        string `yaml:"api_base_url"`
	OIDCTokenEndpoint string `yaml:"oidc_token_endpoint"`
	TargetAPI         string `yaml:"target_api"`
}

// UserOwnCloudSQL defines the available ownCloudSQL user provider configuration.
type UserOwnCloudSQL struct {
	DBUsername         string `yaml:"db_username"`
	DBPassword         string `yaml:"db_password"`
	DBHost             string `yaml:"db_host"`
	DBPort             int    `yaml:"db_port"`
	DBName             string `yaml:"db_name"`
	Idp                string `yaml:"idp"`
	Nobody             int64  `yaml:"nobody"`
	JoinUsername       bool   `yaml:"join_username"`
	JoinOwnCloudUUID   bool   `yaml:"join_owncloud_uuid"`
	EnableMedialSearch bool   `yaml:"enable_medial_search"`
}

// LDAPUserSchema defines the available ldap user schema configuration.
type LDAPUserSchema struct {
	ID              string `yaml:"id"`
	IDIsOctetString bool   `yaml:"id_is_octet_string"`
	Mail            string `yaml:"mail"`
	DisplayName     string `yaml:"display_name"`
	Username        string `yaml:"user_name"`
	UIDNumber       string `yaml:"uid_number"`
	GIDNumber       string `yaml:"gid_number"`
}

// LDAPGroupSchema defines the available ldap group schema configuration.
type LDAPGroupSchema struct {
	ID              string `yaml:"id"`
	IDIsOctetString bool   `yaml:"id_is_octet_string"`
	Mail            string `yaml:"mail"`
	DisplayName     string `yaml:"display_name"`
	Groupname       string `yaml:"group_name"`
	Member          string `yaml:"member"`
	GIDNumber       string `yaml:"gid_number"`
}

// OCDav defines the available ocdav configuration.
type OCDav struct {
	WebdavNamespace   string `yaml:"webdav_namespace"`
	DavFilesNamespace string `yaml:"dav_files_namespace"`
}

// Archiver defines the available archiver configuration.
type Archiver struct {
	MaxNumFiles int64  `yaml:"max_num_files"`
	MaxSize     int64  `yaml:"max_size"`
	ArchiverURL string `yaml:"archiver_url"`
}

// Reva defines the available reva configuration.
type Reva struct {
	// JWTSecret used to sign jwt tokens between services
	JWTSecret             string          `yaml:"jwt_secret"`
	SkipUserGroupsInToken bool            `yaml:"skip_user_grooups_in_token"`
	TransferSecret        string          `yaml:"transfer_secret"`
	TransferExpires       int             `yaml:"transfer_expires"`
	OIDC                  OIDC            `yaml:"oidc"`
	LDAP                  LDAP            `yaml:"ldap"`
	UserGroupRest         UserGroupRest   `yaml:"user_group_rest"`
	UserOwnCloudSQL       UserOwnCloudSQL `yaml:"user_owncloud_sql"`
	OCDav                 OCDav           `yaml:"ocdav"`
	Archiver              Archiver        `yaml:"archiver"`
	UserStorage           StorageConfig   `yaml:"user_storage"`
	MetadataStorage       StorageConfig   `yaml:"metadata_storage"`
	// Ports are used to configure which services to start on which port
	Frontend          FrontendPort      `yaml:"frontend"`
	DataGateway       DataGatewayPort   `yaml:"data_gateway"`
	Gateway           Gateway           `yaml:"gateway"`
	StorageRegistry   StorageRegistry   `yaml:"storage_registry"`
	AppRegistry       AppRegistry       `yaml:"app_registry"`
	Users             Users             `yaml:"users"`
	Groups            Groups            `yaml:"groups"`
	AuthProvider      Users             `yaml:"auth_provider"`
	AuthBasic         Port              `yaml:"auth_basic"`
	AuthBearer        Port              `yaml:"auth_bearer"`
	AuthMachine       Port              `yaml:"auth_machine"`
	AuthMachineConfig AuthMachineConfig `yaml:"auth_machine_config"`
	Sharing           Sharing           `yaml:"sharing"`
	StorageShares     StoragePort       `yaml:"storage_shares"`
	StorageUsers      StoragePort       `yaml:"storage_users"`
	StoragePublicLink PublicStorage     `yaml:"storage_public_link"`
	StorageMetadata   StoragePort       `yaml:"storage_metadata"`
	AppProvider       AppProvider       `yaml:"app_provider"`
	Permissions       Port              `yaml:"permissions"`
	// Configs can be used to configure the reva instance.
	// Services and Ports will be ignored if this is used
	Configs map[string]interface{} `yaml:"configs"`
	// chunking and resumable upload config (TUS)
	UploadMaxChunkSize       int    `yaml:"uppload_max_chunk_size"`
	UploadHTTPMethodOverride string `yaml:"upload_http_method_override"`
	// checksumming capabilities
	ChecksumSupportedTypes      []string `yaml:"checksum_supported_types"`
	ChecksumPreferredUploadType string   `yaml:"checksum_preferred_upload_type"`
	DefaultUploadProtocol       string   `yaml:"default_upload_protocol"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled"`
	Type      string `yaml:"type"`
	Endpoint  string `yaml:"endpoint"`
	Collector string `yaml:"collector"`
	Service   string `yaml:"service"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File    string      `yaml:"file"`
	Log     *shared.Log `yaml:"log"`
	Debug   Debug       `yaml:"debug"`
	Reva    Reva        `yaml:"reva"`
	Tracing Tracing     `yaml:"tracing"`
	Asset   Asset       `yaml:"asset"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
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
			EnvVars:     []string{"APP_PROVIDER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AppProvider.DebugAddr,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_GRPC_NETWORK"},
			Destination: &cfg.Reva.AppProvider.GRPCNetwork,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_GRPC_ADDR"},
			Destination: &cfg.Reva.AppProvider.GRPCAddr,
		},
		{
			EnvVars:     []string{"APP_PROVIDER_EXTERNAL_ADDR"},
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
			EnvVars:     []string{"STORAGE_FRONTEND_PROJECT_SPACES"},
			Destination: &cfg.Reva.Frontend.ProjectSpaces,
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
		{
			EnvVars:     []string{"STORAGE_FRONTEND_APP_PROVIDER_NEW_URL"},
			Destination: &cfg.Reva.AppProvider.NewURL,
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
			EnvVars:     []string{"STORAGE_USERS_ENDPOINT"},
			Destination: &cfg.Reva.StorageUsers.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageUsers.MountID,
		},
		{
			EnvVars:     []string{"STORAGE_SHARES_ENDPOINT"},
			Destination: &cfg.Reva.StorageShares.Endpoint,
		},
		{
			EnvVars:     []string{"STORAGE_PUBLIC_LINK_ENDPOINT"},
			Destination: &cfg.Reva.StoragePublicLink.Endpoint,
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
			EnvVars:     []string{"LDAP_URI", "STORAGE_LDAP_URI"},
			Destination: &cfg.Reva.LDAP.URI,
		},
		{
			EnvVars:     []string{"LDAP_CACERT", "STORAGE_LDAP_CACERT"},
			Destination: &cfg.Reva.LDAP.CACert,
		},
		{
			EnvVars:     []string{"LDAP_INSECURE", "STORAGE_LDAP_INSECURE"},
			Destination: &cfg.Reva.LDAP.Insecure,
		},
		{
			EnvVars:     []string{"LDAP_USER_BASE_DN", "STORAGE_LDAP_USER_BASE_DN"},
			Destination: &cfg.Reva.LDAP.UserBaseDN,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_BASE_DN", "STORAGE_LDAP_GROUP_BASE_DN"},
			Destination: &cfg.Reva.LDAP.GroupBaseDN,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCOPE", "STORAGE_LDAP_USER_SCOPE"},
			Destination: &cfg.Reva.LDAP.UserScope,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCOPE", "STORAGE_LDAP_GROUP_SCOPE"},
			Destination: &cfg.Reva.LDAP.GroupScope,
		},
		{
			EnvVars:     []string{"LDAP_USER_OBJECTCLASS", "STORAGE_LDAP_USER_OBJECTCLASS"},
			Destination: &cfg.Reva.LDAP.UserObjectClass,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_OBJECTCLASS", "STORAGE_LDAP_GROUP_OBJECTCLASS"},
			Destination: &cfg.Reva.LDAP.GroupObjectClass,
		},
		{
			EnvVars:     []string{"LDAP_LOGIN_ATTRIBUTES", "STORAGE_LDAP_LOGIN_ATTRIBUTES"},
			Destination: &cfg.Reva.LDAP.LoginAttributes,
		},
		{
			EnvVars:     []string{"LDAP_USERFILTER", "STORAGE_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		{
			EnvVars:     []string{"LDAP_GROUPFILTER", "STORAGE_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		{
			EnvVars:     []string{"LDAP_BIND_DN", "STORAGE_LDAP_BIND_DN"},
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		{
			EnvVars:     []string{"LDAP_BIND_PASSWORD", "STORAGE_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_ID", "STORAGE_LDAP_USER_SCHEMA_ID"},
			Destination: &cfg.Reva.LDAP.UserSchema.ID,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING", "STORAGE_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"},
			Destination: &cfg.Reva.LDAP.UserSchema.IDIsOctetString,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_MAIL", "STORAGE_LDAP_USER_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.UserSchema.Mail,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_DISPLAYNAME", "STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.UserSchema.DisplayName,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_USERNAME", "STORAGE_LDAP_USER_SCHEMA_USERNAME"},
			Destination: &cfg.Reva.LDAP.UserSchema.Username,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_UID_NUMBER", "STORAGE_LDAP_USER_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.UIDNumber,
		},
		{
			EnvVars:     []string{"LDAP_USER_SCHEMA_GID_NUMBER", "STORAGE_LDAP_USER_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.GIDNumber,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_ID", "STORAGE_LDAP_GROUP_SCHEMA_ID"},
			Destination: &cfg.Reva.LDAP.GroupSchema.ID,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING", "STORAGE_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"},
			Destination: &cfg.Reva.LDAP.GroupSchema.IDIsOctetString,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_MAIL", "STORAGE_LDAP_GROUP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Mail,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_DISPLAYNAME", "STORAGE_LDAP_GROUP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.GroupSchema.DisplayName,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_GROUPNAME", "STORAGE_LDAP_GROUP_SCHEMA_GROUPNAME"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Groupname,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_MEMBER", "STORAGE_LDAP_GROUP_SCHEMA_MEMBER"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Member,
		},
		{
			EnvVars:     []string{"LDAP_GROUP_SCHEMA_GID_NUMBER", "STORAGE_LDAP_GROUP_SCHEMA_GID_NUMBER"},
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

		// sharing cs3

		{
			EnvVars:     []string{"STORAGE_SHARING_CS3_PROVIDER_ADDR"},
			Destination: &cfg.Reva.Sharing.CS3ProviderAddr,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_CS3_SERVICE_USER"},
			Destination: &cfg.Reva.Sharing.CS3ServiceUser,
		},
		{
			EnvVars:     []string{"OCIS_URL", "STORAGE_SHARING_CS3_SERVICE_USER_IDP"},
			Destination: &cfg.Reva.Sharing.CS3ServiceUserIdp,
		},

		// sharingsql
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
		{
			EnvVars:     []string{"STORAGE_SHARING_EVENTS_ADDRESS"},
			Destination: &cfg.Reva.Sharing.Events.Address,
		},
		{
			EnvVars:     []string{"STORAGE_SHARING_EVENTS_CLUSTER_ID"},
			Destination: &cfg.Reva.Sharing.Events.ClusterID,
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
			EnvVars:     []string{"STORAGE_METADATA_GRPC_ADDR"},
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

		// storage shares
		{
			EnvVars:     []string{"STORAGE_SHARES_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageShares.DebugAddr,
		},
		{
			EnvVars:     []string{"STORAGE_SHARES_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageShares.GRPCNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_SHARES_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageShares.GRPCAddr,
		},
		{
			EnvVars:     []string{"STORAGE_SHARES_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageShares.HTTPNetwork,
		},
		{
			EnvVars:     []string{"STORAGE_SHARES_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageShares.HTTPAddr,
		},
		{
			EnvVars:     []string{"OCIS_STORAGE_READ_ONLY", "STORAGE_SHARES_READ_ONLY"},
			Destination: &cfg.Reva.StorageShares.ReadOnly,
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
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_PERSONAL_SPACE_ALIAS_TEMPLATE"},
			Destination: &cfg.Reva.UserStorage.OCIS.PersonalSpaceAliasTemplate,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_GENERAL_SPACE_ALIAS_TEMPLATE"},
			Destination: &cfg.Reva.UserStorage.OCIS.GeneralSpaceAliasTemplate,
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
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_PERSONAL_SPACE_ALIAS_TEMPLATE"},
			Destination: &cfg.Reva.UserStorage.S3NG.PersonalSpaceAliasTemplate,
		},
		{
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_GENERAL_SPACE_ALIAS_TEMPLATE"},
			Destination: &cfg.Reva.UserStorage.S3NG.GeneralSpaceAliasTemplate,
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

		// permissions
		{
			EnvVars:     []string{"STORAGE_PERMISSIONS_ENDPOINT"},
			Destination: &cfg.Reva.Permissions.Endpoint,
		},
	}
}
