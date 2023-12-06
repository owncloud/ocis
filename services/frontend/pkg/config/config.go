package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	HTTP HTTPConfig `yaml:"http"`

	// JWTSecret used to verify reva access token

	TransferSecret string `yaml:"transfer_secret" env:"OCIS_TRANSFER_SECRET" desc:"Transfer secret for signing file up- and download requests."`

	TokenManager      *TokenManager `yaml:"token_manager"`
	Reva              *shared.Reva  `yaml:"reva"`
	MachineAuthAPIKey string        `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;FRONTEND_MACHINE_AUTH_API_KEY" desc:"The machine auth API key used to validate internal requests necessary to access resources from other services."`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"FRONTEND_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	EnableFavorites                bool   `yaml:"enable_favorites" env:"FRONTEND_ENABLE_FAVORITES" desc:"Enables the support for favorites in the clients."`
	EnableProjectSpaces            bool   `yaml:"enable_project_spaces" env:"FRONTEND_ENABLE_PROJECT_SPACES" desc:"Changing this value is NOT supported. Indicates to clients that project spaces are supposed to be made available."`
	EnableShareJail                bool   `yaml:"enable_share_jail" env:"FRONTEND_ENABLE_SHARE_JAIL" desc:"Changing this value is NOT supported. Indicates to clients that the share jail is supposed to be used."`
	MaxQuota                       uint64 `yaml:"max_quota" env:"OCIS_SPACES_MAX_QUOTA;FRONTEND_MAX_QUOTA" desc:"Set the global max quota value in bytes. A value of 0 equals unlimited. The value is provided via capabilities."`
	UploadMaxChunkSize             int    `yaml:"upload_max_chunk_size" env:"FRONTEND_UPLOAD_MAX_CHUNK_SIZE" desc:"Sets the max chunk sizes in bytes for uploads via the clients."`
	UploadHTTPMethodOverride       string `yaml:"upload_http_method_override" env:"FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE" desc:"Advise TUS to replace PATCH requests by POST requests."`
	DefaultUploadProtocol          string `yaml:"default_upload_protocol" env:"FRONTEND_DEFAULT_UPLOAD_PROTOCOL" desc:"The default upload protocol to use in clients. Currently only 'tus' is avaliable. See the developer API documentation for more details about TUS."`
	EnableResharing                bool   `yaml:"enable_resharing" env:"OCIS_ENABLE_RESHARING;FRONTEND_ENABLE_RESHARING" desc:"Changing this value is NOT supported. Enables the support for resharing in the clients."`
	EnableFederatedSharingIncoming bool   `yaml:"enable_federated_sharing_incoming" env:"FRONTEND_ENABLE_FEDERATED_SHARING_INCOMING" desc:"Changing this value is NOT supported. Enables support for incoming federated sharing for clients. The backend behaviour is not changed."`
	EnableFederatedSharingOutgoing bool   `yaml:"enable_federated_sharing_outgoing" env:"FRONTEND_ENABLE_FEDERATED_SHARING_OUTGOING" desc:"Changing this value is NOT supported. Enables support for outgoing federated sharing for clients. The backend behaviour is not changed."`
	SearchMinLength                int    `yaml:"search_min_length" env:"FRONTEND_SEARCH_MIN_LENGTH" desc:"Minimum number of characters to enter before a client should start a search for Share receivers. This setting can be used to customize the user experience if e.g too many results are displayed."`
	Edition                        string `yaml:"edition" env:"OCIS_EDITION;FRONTEND_EDITION"`
	DisableSSE                     bool   `yaml:"disable_sse" env:"OCIS_DISABLE_SSE;FRONTEND_DISABLE_SSE" desc:"When set to true, clients are informed that the Server-Sent Events endpoint is not accessible."`
	DefaultLinkPermissions         int    `yaml:"default_link_permissions" env:"FRONTEND_DEFAULT_LINK_PERMISSIONS" desc:"Defines the default permissions a link is being created with. Possible values are 0 (= internal link, for instance members only) and 1 (= public link with viewer permissions). Defaults to 1."`

	PublicURL string `yaml:"public_url" env:"OCIS_URL;FRONTEND_PUBLIC_URL" desc:"The public facing URL of the oCIS frontend."`

	AppHandler             AppHandler  `yaml:"app_handler"`
	Archiver               Archiver    `yaml:"archiver"`
	DataGateway            DataGateway `yaml:"data_gateway"`
	OCS                    OCS         `yaml:"ocs"`
	Checksums              Checksums   `yaml:"checksums"`
	ReadOnlyUserAttributes []string    `yaml:"read_only_user_attributes" env:"FRONTEND_READONLY_USER_ATTRIBUTES" desc:"A list of user attributes to indicate as read-only. Supported values: 'user.onPremisesSamAccountName' (username), 'user.displayName', 'user.mail', 'user.passwordProfile' (password), 'user.appRoleAssignments' (role), 'user.memberOf' (groups), 'user.accountEnabled' (login allowed), 'drive.quota' (quota). See the Environment Variable Types description for more details."`
	LDAPServerWriteEnabled bool        `yaml:"ldap_server_write_enabled" env:"OCIS_LDAP_SERVER_WRITE_ENABLED;FRONTEND_LDAP_SERVER_WRITE_ENABLED" desc:"Allow creating, modifying and deleting LDAP users via the GRAPH API. This can only be set to 'true' when keeping default settings for the LDAP user and group attribute types (the 'OCIS_LDAP_USER_SCHEMA_* and 'OCIS_LDAP_GROUP_SCHEMA_* variables)."`
	FullTextSearch         bool        `yaml:"full_text_search" env:"FRONTEND_FULL_TEXT_SEARCH_ENABLED" descr:"Set to true to signal the web client that full-text search is enabled."`

	Middleware Middleware `yaml:"middleware"`

	Events           Events                `yaml:"events"`
	GRPCClientTLS    *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	AutoAcceptShares bool                  `yaml:"auto_accept_shares" env:"FRONTEND_AUTO_ACCEPT_SHARES" desc:"Defines if shares should be auto accepted by default. Users can change this setting individually in their profile."`
	ServiceAccount   ServiceAccount        `yaml:"service_account"`

	PasswordPolicy PasswordPolicy `yaml:"password_policy"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;FRONTEND_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;FRONTEND_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;FRONTEND_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;FRONTEND_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"FRONTEND_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"FRONTEND_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"FRONTEND_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"FRONTEND_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"FRONTEND_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"FRONTEND_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
	Prefix    string `yaml:"prefix" env:"FRONTEND_HTTP_PREFIX" desc:"The Path prefix where the frontend can be accessed (defaults to /)."`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;FRONTEND_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details."`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;FRONTEND_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details."`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;FRONTEND_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details."`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;FRONTEND_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials."`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent"`
}

type AppHandler struct {
	Prefix   string `yaml:"-"`
	Insecure bool   `yaml:"insecure" env:"OCIS_INSECURE;FRONTEND_APP_HANDLER_INSECURE" desc:"Allow insecure connections to the frontend."`
}

type Archiver struct {
	MaxNumFiles int64  `yaml:"max_num_files" env:"FRONTEND_ARCHIVER_MAX_NUM_FILES" desc:"Max number of files that can be packed into an archive."`
	MaxSize     int64  `yaml:"max_size" env:"FRONTEND_ARCHIVER_MAX_SIZE" desc:"Max size in bytes of the zip archive the archiver can create."`
	Prefix      string `yaml:"-"`
	Insecure    bool   `yaml:"insecure" env:"OCIS_INSECURE;FRONTEND_ARCHIVER_INSECURE" desc:"Allow insecure connections to the archiver."`
}

type DataGateway struct {
	Prefix string `yaml:"prefix" env:"FRONTEND_DATA_GATEWAY_PREFIX" desc:"Path prefix for the data gateway."`
}

type OCS struct {
	Prefix                               string             `yaml:"prefix" env:"FRONTEND_OCS_PREFIX" desc:"URL path prefix for the OCS service. Note that the string must not start with '/'."`
	SharePrefix                          string             `yaml:"share_prefix" env:"FRONTEND_OCS_SHARE_PREFIX" desc:"Path prefix for shares as part of an ocis resource. Note that the path must start with '/'."`
	HomeNamespace                        string             `yaml:"home_namespace" env:"FRONTEND_OCS_PERSONAL_NAMESPACE" desc:"Homespace namespace identifier."`
	AdditionalInfoAttribute              string             `yaml:"additional_info_attribute" env:"FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE" desc:"Additional information attribute for the user like {{.Mail}}."`
	StatCacheType                        string             `yaml:"stat_cache_type" env:"OCIS_CACHE_STORE;FRONTEND_OCS_STAT_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	StatCacheNodes                       []string           `yaml:"stat_cache_nodes" env:"OCIS_CACHE_STORE_NODES;FRONTEND_OCS_STAT_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	StatCacheDatabase                    string             `yaml:"stat_cache_database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use."`
	StatCacheTable                       string             `yaml:"stat_cache_table" env:"FRONTEND_OCS_STAT_CACHE_TABLE" desc:"The database table the store should use."`
	StatCacheTTL                         time.Duration      `yaml:"stat_cache_ttl" env:"OCIS_CACHE_TTL;FRONTEND_OCS_STAT_CACHE_TTL" desc:"Default time to live for user info in the cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details."`
	StatCacheSize                        int                `yaml:"stat_cache_size" env:"OCIS_CACHE_SIZE;FRONTEND_OCS_STAT_CACHE_SIZE" desc:"Max number of entries to hold in the cache."`
	CacheWarmupDriver                    string             `yaml:"cache_warmup_driver,omitempty"`  // not supported by the oCIS product, therefore not part of docs
	CacheWarmupDrivers                   CacheWarmupDrivers `yaml:"cache_warmup_drivers,omitempty"` // not supported by the oCIS product, therefore not part of docs
	EnableDenials                        bool               `yaml:"enable_denials" env:"FRONTEND_OCS_ENABLE_DENIALS" desc:"EXPERIMENTAL: enable the feature to deny access on folders."`
	ListOCMShares                        bool               `yaml:"list_ocm_shares" env:"FRONTEND_OCS_LIST_OCM_SHARES" desc:"Include OCM shares when listing shares. See the OCM service documentation for more details."`
	PublicShareMustHavePassword          bool               `yaml:"public_sharing_share_must_have_password" env:"OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD;FRONTEND_OCS_PUBLIC_SHARE_MUST_HAVE_PASSWORD" desc:"Set this to true if you want to enforce passwords on all public shares."`
	WriteablePublicShareMustHavePassword bool               `yaml:"public_sharing_writeableshare_must_have_password" env:"OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD;FRONTEND_OCS_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" desc:"Set this to true if you want to enforce passwords on Uploader, Editor or Contributor shares."`
}

type CacheWarmupDrivers struct {
	CBOX CBOXDriver `yaml:"cbox,omitempty"`
}

type CBOXDriver struct {
	DBUsername string `yaml:"db_username,omitempty"`
	DBPassword string `yaml:"db_password,omitempty"`
	DBHost     string `yaml:"db_host,omitempty"`
	DBPort     int    `yaml:"db_port,omitempty"`
	DBName     string `yaml:"db_name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}

type Checksums struct {
	SupportedTypes      []string `yaml:"supported_types" env:"FRONTEND_CHECKSUMS_SUPPORTED_TYPES" desc:"A list of checksum types that indicate to clients which hashes the server can use to verify upload integrity. Supported types are 'sha1', 'md5' and 'adler32'. See the Environment Variable Types description for more details."`
	PreferredUploadType string   `yaml:"preferred_upload_type" env:"FRONTEND_CHECKSUMS_PREFERRED_UPLOAD_TYPE" desc:"The supported checksum type for uploads that indicates to clients supporting multiple hash algorithms which one is preferred by the server. Must be one out of the defined list of SUPPORTED_TYPES."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;FRONTEND_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;FRONTEND_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;FRONTEND_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"FRONTEND_EVENTS_TLS_ROOT_CA_CERTIFICATE;OCS_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;FRONTEND_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ServiceAccountID     string `yaml:"service_account_id" env:"OCIS_SERVICE_ACCOUNT_ID;FRONTEND_SERVICE_ACCOUNT_ID" desc:"The ID of the service account the service should use. See the 'auth-service' service description for more details."`
	ServiceAccountSecret string `yaml:"service_account_secret" env:"OCIS_SERVICE_ACCOUNT_SECRET;FRONTEND_SERVICE_ACCOUNT_SECRET" desc:"The service account secret."`
}

// PasswordPolicy configures reva password policy
type PasswordPolicy struct {
	MinCharacters          int    `yaml:"min_characters,omitempty" env:"OCIS_PASSWORD_POLICY_MIN_CHARACTERS;FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS" desc:"Define the minimum password length. Defaults to 0 if not set."`
	MinLowerCaseCharacters int    `yaml:"min_lowercase_characters" env:"OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS;FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS" desc:"Define the minimum number of uppercase letters. Defaults to 0 if not set."`
	MinUpperCaseCharacters int    `yaml:"min_uppercase_characters" env:"OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS;FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS" desc:"Define the minimum number of lowercase letters. Defaults to 0 if not set."`
	MinDigits              int    `yaml:"min_digits" env:"OCIS_PASSWORD_POLICY_MIN_DIGITS;FRONTEND_PASSWORD_POLICY_MIN_DIGITS" desc:"Define the minimum number of digits. Defaults to 0 if not set."`
	MinSpecialCharacters   int    `yaml:"min_special_characters" env:"OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS;FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS" desc:"Define the minimum number of characters from the special characters list to be present. Defaults to 0 if not set."`
	BannedPasswordsList    string `yaml:"banned_passwords_list" env:"OCIS_PASSWORD_POLICY_BANNED_PASSWORDS_LIST;FRONTEND_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" desc:"Path to the 'banned passwords list' file. See the documentation for more details."`
}
