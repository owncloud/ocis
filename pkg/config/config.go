package config

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
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
	ShareFolder                string
	DisableHomeCreationOnLogin bool
}

// Port defines the available port configuration.
type Port struct {
	// MaxCPUs can be a number or a percentage
	MaxCPUs  string
	LogLevel string
	// Network can be tcp, udp or unix
	Network string
	// Addr to listen on, hostname:port (0.0.0.0:9999 for all interfaces) or socket (/var/run/reva/sock)
	Addr string
	// Protocol can be grpc or http
	Protocol string
	// URL is used by the gateway and registries (eg http://localhost:9100 or https://cloud.example.com)
	URL string
	// DebugAddr for the debug endpoint to bind to
	DebugAddr string
	// Services can be used to give a list of services that should be started on this port
	Services []string
	// Config can be used to configure the reva instance.
	// Services and Protocol will be ignored if this is used
	Config map[string]interface{}
}

// Users defines the available users configuration.
type Users struct {
	Port
	Driver string
	JSON   string
}

// PathWrapperContext defines the available PathWrapperContext configuration.
type PathWrapperContext struct {
	Prefix string
}

// StoragePort defines the available storage configuration.
type StoragePort struct {
	Port
	Driver             string
	PathWrapper        string
	PathWrapperContext PathWrapperContext
	MountPath          string
	MountID            string
	ExposeDataServer   bool
	DataServerURL      string
	EnableHomeCreation bool

	// for HTTP ports with only one http service
	Prefix     string
	TempFolder string
}

// StorageConfig combines all available storage driver configuration parts.
type StorageConfig struct {
	EOS      DriverEOS
	Local    DriverLocal
	OwnCloud DriverOwnCloud
	S3       DriverS3
	// TODO checksums ... figure out what that is supposed to do
}

// DriverEOS defines the available EOS driver configuration.
type DriverEOS struct {
	// Namespace for metadata operations
	Namespace string

	// Location of the eos binary.
	// Default is /usr/bin/eos.
	EosBinary string

	// Location of the xrdcopy binary.
	// Default is /usr/bin/xrdcopy.
	XrdcopyBinary string

	// URL of the Master EOS MGM.
	// Default is root://eos-example.org
	MasterURL string

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

	// EnableHome enables the creation of home directories.
	EnableHome bool

	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string

	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string

	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string

	// Layout of the users home dir path
	Layout string
}

// DriverLocal defines the available local storage driver configuration.
type DriverLocal struct {
	Root string
}

// DriverOwnCloud defines the available ownCloud storage driver configuration.
type DriverOwnCloud struct {
	Datadirectory string
	Scan          bool
	Redis         string
	Layout        string
}

// DriverS3 defines the available S3 storage driver configuration.
type DriverS3 struct {
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
	Prefix    string
}

// OIDC defines the available OpenID Connect configuration.
type OIDC struct {
	Issuer   string
	Insecure bool
	IDClaim  string
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname     string
	Port         int
	BaseDN       string
	UserFilter   string
	GroupFilter  string
	BindDN       string
	BindPassword string
	IDP          string
	Schema       LDAPSchema
}

// LDAPSchema defines the available ldap schema configuration.
type LDAPSchema struct {
	UID         string
	Mail        string
	DisplayName string
	CN          string
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
	OCDav           OCDav
	Storages        StorageConfig
	// Ports are used configure which services to start on which port
	Frontend          Port
	Gateway           Gateway
	Users             Users
	AuthBasic         Port
	AuthBearer        Port
	Sharing           Port
	StorageRoot       StoragePort
	StorageHome       StoragePort
	StorageHomeData   StoragePort
	StorageEOS        StoragePort
	StorageEOSData    StoragePort
	StorageOC         StoragePort
	StorageOCData     StoragePort
	StorageS3         StoragePort
	StorageS3Data     StoragePort
	StorageWND        StoragePort
	StorageWNDData    StoragePort
	StorageCustom     StoragePort
	StorageCustomData StoragePort
	// Configs can be used to configure the reva instance.
	// Services and Ports will be ignored if this is used
	Configs map[string]interface{}
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
