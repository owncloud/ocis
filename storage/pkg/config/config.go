package config

import "context"

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// Gateway defines the available gateway configuration.
type Gateway struct {
	Port
	CommitShareToStorageGrant  bool
	CommitShareToStorageRef    bool
	DisableHomeCreationOnLogin bool
	ShareFolder                string
	LinkGrants                 string
	HomeMapping                string
	EtagCacheTTL               int
}

// StorageRegistry defines the available storage registry configuration
type StorageRegistry struct {
	Driver string
	// HomeProvider is the path in the global namespace that the static storage registry uses to determine the home storage
	HomeProvider string
	Rules        []string
	JSON         string
}

// Sharing defines the available sharing configuration.
type Sharing struct {
	Port
	UserDriver      string
	UserJSONFile    string
	UserSQLUsername string
	UserSQLPassword string
	UserSQLHost     string
	UserSQLPort     int
	UserSQLName     string
	PublicDriver    string
	PublicJSONFile  string
}

// Port defines the available port configuration.
type Port struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string
	LogLevel string
	// GRPCNetwork can be tcp, udp or unix
	GRPCNetwork string
	// GRPCAddr to listen on, hostname:port (0.0.0.0:9999 for all interfaces) or socket (/var/run/reva/sock)
	GRPCAddr string
	// Protocol can be grpc or http
	// HTTPNetwork can be tcp, udp or unix
	HTTPNetwork string
	// HTTPAddr to listen on, hostname:port (0.0.0.0:9100 for all interfaces) or socket (/var/run/reva/sock)
	HTTPAddr string
	// Protocol can be grpc or http
	Protocol string
	// Endpoint is used by the gateway and registries (eg localhost:9100 or cloud.example.com)
	Endpoint string
	// DebugAddr for the debug endpoint to bind to
	DebugAddr string
	// Services can be used to give a list of services that should be started on this port
	Services []string
	// Config can be used to configure the reva instance.
	// Services and Protocol will be ignored if this is used
	Config map[string]interface{}

	// Context allows for context cancellation and propagation
	Context context.Context

	// Supervised is used when running under an oCIS runtime supervision tree
	Supervised bool
}

// Users defines the available users configuration.
type Users struct {
	Port
	Driver                    string
	JSON                      string
	UserGroupsCacheExpiration int
}

// Groups defines the available groups configuration.
type Groups struct {
	Port
	Driver                      string
	JSON                        string
	GroupMembersCacheExpiration int
}

// FrontendPort defines the available frontend configuration.
type FrontendPort struct {
	Port

	DatagatewayPrefix string
	OCDavPrefix       string
	OCSPrefix         string
	OCSSharePrefix    string
	OCSHomeNamespace  string
	PublicURL         string
	Middleware        Middleware
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string
}

// DataGatewayPort has a public url
type DataGatewayPort struct {
	Port
	PublicURL string
}

// StoragePort defines the available storage configuration.
type StoragePort struct {
	Port
	Driver           string
	MountPath        string
	MountID          string
	ExposeDataServer bool
	// url the data gateway will use to route requests
	DataServerURL string

	// for HTTP ports with only one http service
	HTTPPrefix string
	TempFolder string
}

// PublicStorage configures a public storage provider
type PublicStorage struct {
	StoragePort

	PublicShareProviderAddr string
	UserProviderAddr        string
}

// StorageConfig combines all available storage driver configuration parts.
type StorageConfig struct {
	Home     DriverCommon
	EOS      DriverEOS
	Local    DriverCommon
	OwnCloud DriverOwnCloud
	S3       DriverS3
	Common   DriverCommon
	// TODO checksums ... figure out what that is supposed to do
}

// DriverCommon defines common driver configuration options.
type DriverCommon struct {
	// Root is the absolute path to the location of the data
	Root string
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string
	// UserLayout contains the template used to construct
	// the internal path, eg: `{{substr 0 1 .Username}}/{{.Username}}`
	UserLayout string
	// EnableHome enables the creation of home directories.
	EnableHome bool
}

// DriverEOS defines the available EOS driver configuration.
type DriverEOS struct {
	DriverCommon

	// ShadowNamespace for storing shadow data
	ShadowNamespace string

	// UploadsNamespace for storing upload data
	UploadsNamespace string

	// Location of the eos binary.
	// Default is /usr/bin/eos.
	EosBinary string

	// Location of the xrdcopy binary.
	// Default is /usr/bin/xrdcopy.
	XrdcopyBinary string

	// URL of the Master EOS MGM.
	// Default is root://eos-example.org
	MasterURL string

	// URI of the EOS MGM grpc server
	// Default is empty
	GrpcURI string

	// URL of the Slave EOS MGM.
	// Default is root://eos-example.org
	SlaveURL string

	// Location on the local fs where to store reads.
	// Defaults to os.TempDir()
	CacheDirectory string

	// Enables logging of the commands executed
	// Defaults to false
	EnableLogging bool

	// ShowHiddenSysFiles shows internal EOS files like
	// .sys.v# and .sys.a# files.
	ShowHiddenSysFiles bool

	// ForceSingleUserMode will force connections to EOS to use SingleUsername
	ForceSingleUserMode bool

	// UseKeyTabAuth changes will authenticate requests by using an EOS keytab.
	UseKeytab bool

	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string

	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string

	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string

	// gateway service to use for uid lookups
	GatewaySVC string
}

// DriverOwnCloud defines the available ownCloud storage driver configuration.
type DriverOwnCloud struct {
	DriverCommon

	UploadInfoDir string
	Redis         string
	Scan          bool
}

// DriverS3 defines the available S3 storage driver configuration.
type DriverS3 struct {
	DriverCommon

	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
}

// OIDC defines the available OpenID Connect configuration.
type OIDC struct {
	Issuer   string
	Insecure bool
	IDClaim  string
	UIDClaim string
	GIDClaim string
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname             string
	Port                 int
	BaseDN               string
	LoginFilter          string
	UserFilter           string
	UserAttributeFilter  string
	UserFindFilter       string
	UserGroupFilter      string
	GroupFilter          string
	GroupAttributeFilter string
	GroupFindFilter      string
	GroupMemberFilter    string
	BindDN               string
	BindPassword         string
	IDP                  string
	UserSchema           LDAPUserSchema
	GroupSchema          LDAPGroupSchema
}

// UserGroupRest defines the REST driver specification for user and group resolution.
type UserGroupRest struct {
	ClientID          string
	ClientSecret      string
	RedisAddress      string
	RedisUsername     string
	RedisPassword     string
	IDProvider        string
	APIBaseURL        string
	OIDCTokenEndpoint string
	TargetAPI         string
}

// LDAPUserSchema defines the available ldap user schema configuration.
type LDAPUserSchema struct {
	UID         string
	Mail        string
	DisplayName string
	CN          string
	UIDNumber   string
	GIDNumber   string
}

// LDAPGroupSchema defines the available ldap group schema configuration.
type LDAPGroupSchema struct {
	GID         string
	Mail        string
	DisplayName string
	CN          string
	GIDNumber   string
}

// OCDav defines the available ocdav configuration.
type OCDav struct {
	WebdavNamespace   string
	DavFilesNamespace string
}

// Reva defines the available reva configuration.
type Reva struct {
	// JWTSecret used to sign jwt tokens between services
	JWTSecret       string
	TransferSecret  string
	TransferExpires int
	OIDC            OIDC
	LDAP            LDAP
	UserGroupRest   UserGroupRest
	OCDav           OCDav
	Storages        StorageConfig
	// Ports are used to configure which services to start on which port
	Frontend          FrontendPort
	DataGateway       DataGatewayPort
	Gateway           Gateway
	StorageRegistry   StorageRegistry
	Users             Users
	Groups            Groups
	AuthProvider      Users
	AuthBasic         Port
	AuthBearer        Port
	Sharing           Sharing
	StorageHome       StoragePort
	StorageUsers      StoragePort
	StoragePublicLink PublicStorage
	StorageMetadata   StoragePort
	// Configs can be used to configure the reva instance.
	// Services and Ports will be ignored if this is used
	Configs map[string]interface{}
	// chunking and resumable upload config (TUS)
	UploadMaxChunkSize       int
	UploadHTTPMethodOverride string
	// checksumming capabilities
	ChecksumSupportedTypes      []string
	ChecksumPreferredUploadType string
	DefaultUploadProtocol       string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	Reva    Reva
	Tracing Tracing
	Asset   Asset
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
