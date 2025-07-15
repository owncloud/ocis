package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config is the configuration for the storage-users service
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken   bool `yaml:"skip_user_groups_in_token" env:"STORAGE_USERS_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token." introductionVersion:"pre5.0"`
	GracefulShutdownTimeout int  `yaml:"graceful_shutdown_timeout" env:"STORAGE_USERS_GRACEFUL_SHUTDOWN_TIMEOUT" desc:"The number of seconds to wait for the 'storage-users' service to shutdown cleanly before exiting with an error that gets logged. Note: This setting is only applicable when running the 'storage-users' service as a standalone service. See the text description for more details." introductionVersion:"pre5.0"`

	Driver         string  `yaml:"driver" env:"STORAGE_USERS_DRIVER" desc:"The storage driver which should be used by the service. Defaults to 'ocis', Supported values are: 'ocis', 's3ng' and 'owncloudsql'. The 'ocis' driver stores all data (blob and meta data) in an POSIX compliant volume. The 's3ng' driver stores metadata in a POSIX compliant volume and uploads blobs to the s3 bucket." introductionVersion:"pre5.0"`
	Drivers        Drivers `yaml:"drivers"`
	DataServerURL  string  `yaml:"data_server_url" env:"STORAGE_USERS_DATA_SERVER_URL" desc:"URL of the data server, needs to be reachable by the data gateway provided by the frontend service or the user if directly exposed." introductionVersion:"pre5.0"`
	DataGatewayURL string  `yaml:"data_gateway_url" env:"STORAGE_USERS_DATA_GATEWAY_URL" desc:"URL of the data gateway server" introductionVersion:"pre5.0"`

	TransferExpires   int64             `yaml:"transfer_expires" env:"STORAGE_USERS_TRANSFER_EXPIRES" desc:"The time after which the token for upload postprocessing expires" introductionVersion:"pre5.0"`
	Events            Events            `yaml:"events"`
	FilemetadataCache FilemetadataCache `yaml:"filemetadata_cache"`
	IDCache           IDCache           `yaml:"id_cache"`
	MountID           string            `yaml:"mount_id" env:"STORAGE_USERS_MOUNT_ID" desc:"Mount ID of this storage." introductionVersion:"pre5.0"`
	ExposeDataServer  bool              `yaml:"expose_data_server" env:"STORAGE_USERS_EXPOSE_DATA_SERVER" desc:"Exposes the data server directly to users and bypasses the data gateway. Ensure that the data server address is reachable by users." introductionVersion:"pre5.0"`
	ReadOnly          bool              `yaml:"readonly" env:"STORAGE_USERS_READ_ONLY" desc:"Set this storage to be read-only." introductionVersion:"pre5.0"`
	UploadExpiration  int64             `yaml:"upload_expiration" env:"STORAGE_USERS_UPLOAD_EXPIRATION" desc:"Duration in seconds after which uploads will expire. Note that when setting this to a low number, uploads could be cancelled before they are finished and return a 403 to the user." introductionVersion:"pre5.0"`
	Tasks             Tasks             `yaml:"tasks"`
	ServiceAccount    ServiceAccount    `yaml:"service_account"`

	// CLI
	RevaGatewayGRPCAddr      string `yaml:"gateway_addr" env:"OCIS_GATEWAY_GRPC_ADDR;STORAGE_USERS_GATEWAY_GRPC_ADDR" desc:"The bind address of the gateway GRPC address." introductionVersion:"5.0"`
	MachineAuthAPIKey        string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;STORAGE_USERS_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services." introductionVersion:"5.0"`
	CliMaxAttemptsRenameFile int    `yaml:"max_attempts_rename_file" env:"STORAGE_USERS_CLI_MAX_ATTEMPTS_RENAME_FILE" desc:"The maximum number of attempts to rename a file when a user restores a file to an existing destination with the same name. The minimum value is 100." introductionVersion:"5.0"`

	Context context.Context `yaml:"-"`
}

// Log configures the logging
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_USERS_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'." introductionVersion:"pre5.0"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_USERS_LOG_PRETTY" desc:"Activates pretty log output." introductionVersion:"pre5.0"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_USERS_LOG_COLOR" desc:"Activates colorized log output." introductionVersion:"pre5.0"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_USERS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set." introductionVersion:"pre5.0"`
}

// Service holds general service configuration
type Service struct {
	Name string `yaml:"-" env:"STORAGE_USERS_SERVICE_NAME" desc:"Service name to use. Change this when starting an additional storage provider with a custom configuration to prevent it from colliding with the default 'storage-users' service." introductionVersion:"7.0.0"`
}

// Debug is the configuration for the debug server
type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_USERS_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0"`
	Token  string `yaml:"token" env:"STORAGE_USERS_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_USERS_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_USERS_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0"`
}

// GRPCConfig is the configuration for the grpc server
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_USERS_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;STORAGE_USERS_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service." introductionVersion:"pre5.0"`
}

// HTTPConfig is the configuration for the http server
type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_USERS_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"pre5.0"`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_USERS_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service." introductionVersion:"pre5.0"`
	Prefix    string
	CORS      CORS `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;STORAGE_USERS_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;STORAGE_USERS_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;STORAGE_USERS_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;STORAGE_USERS_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials." introductionVersion:"pre5.0"`
	ExposedHeaders   []string `yaml:"expose_headers" env:"OCIS_CORS_EXPOSE_HEADERS;STORAGE_USERS_CORS_EXPOSE_HEADERS" desc:"A list of exposed CORS headers. See following chapter for more details: *Access-Control-Expose-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	MaxAge           uint     `yaml:"max_age" env:"OCIS_CORS_MAX_AGE;STORAGE_USERS_CORS_MAX_AGE" desc:"The max cache duration of preflight headers. See following chapter for more details: *Access-Control-Max-Age* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
}

// Drivers combine all storage driver configurations
type Drivers struct {
	OCIS        OCISDriver        `yaml:"ocis"`
	S3NG        S3NGDriver        `yaml:"s3ng"`
	OwnCloudSQL OwnCloudSQLDriver `yaml:"owncloudsql"`
	Posix       PosixDriver       `yaml:"posix"`

	S3    S3Driver    `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
	EOS   EOSDriver   `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
	Local LocalDriver `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
}

// AsyncPropagatorOptions configures the async propagator
type AsyncPropagatorOptions struct {
	PropagationDelay time.Duration `yaml:"propagation_delay" env:"STORAGE_USERS_ASYNC_PROPAGATOR_PROPAGATION_DELAY" desc:"The delay between a change made to a tree and the propagation start on treesize and treetime. Multiple propagations are computed to a single one. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
}

// OCISDriver is the storage driver configuration when using 'ocis' storage driver
type OCISDriver struct {
	Propagator             string                 `yaml:"propagator" env:"OCIS_DECOMPOSEDFS_PROPAGATOR;STORAGE_USERS_OCIS_PROPAGATOR" desc:"The propagator used for decomposedfs. At the moment, only 'sync' is fully supported, 'async' is available as an experimental option." introductionVersion:"pre5.0"`
	AsyncPropagatorOptions AsyncPropagatorOptions `yaml:"async_propagator_options"`
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root" env:"STORAGE_USERS_OCIS_ROOT" desc:"The directory where the filesystem storage will store blobs and metadata. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users." introductionVersion:"pre5.0"`
	UserLayout          string `yaml:"user_layout" env:"STORAGE_USERS_OCIS_USER_LAYOUT" desc:"Template string for the user storage layout in the user directory." introductionVersion:"pre5.0"`
	PermissionsEndpoint string `yaml:"permissions_endpoint" env:"STORAGE_USERS_PERMISSION_ENDPOINT;STORAGE_USERS_OCIS_PERMISSIONS_ENDPOINT" desc:"Endpoint of the permissions service. The endpoints can differ for 'ocis' and 's3ng'." introductionVersion:"pre5.0"`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template" env:"STORAGE_USERS_OCIS_PERSONAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct personal space aliases." introductionVersion:"pre5.0"`
	PersonalSpacePathTemplate  string `yaml:"personalspacepath_template" env:"STORAGE_USERS_OCIS_PERSONAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the personal space roots." introductionVersion:"6.0.0"`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template" env:"STORAGE_USERS_OCIS_GENERAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct general space aliases." introductionVersion:"pre5.0"`
	GeneralSpacePathTemplate  string `yaml:"generalspacepath_template" env:"STORAGE_USERS_OCIS_GENERAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the projects space roots." introductionVersion:"6.0.0"`
	// ShareFolder defines the name of the folder jailing all shares
	ShareFolder             string `yaml:"share_folder" env:"STORAGE_USERS_OCIS_SHARE_FOLDER" desc:"Name of the folder jailing all shares." introductionVersion:"pre5.0"`
	MaxAcquireLockCycles    int    `yaml:"max_acquire_lock_cycles" env:"STORAGE_USERS_OCIS_MAX_ACQUIRE_LOCK_CYCLES" desc:"When trying to lock files, ocis will try this amount of times to acquire the lock before failing. After each try it will wait for an increasing amount of time. Values of 0 or below will be ignored and the default value will be used." introductionVersion:"pre5.0"`
	LockCycleDurationFactor int    `yaml:"lock_cycle_duration_factor" env:"STORAGE_USERS_OCIS_LOCK_CYCLE_DURATION_FACTOR" desc:"When trying to lock files, ocis will multiply the cycle with this factor and use it as a millisecond timeout. Values of 0 or below will be ignored and the default value will be used." introductionVersion:"pre5.0"`
	MaxConcurrency          int    `yaml:"max_concurrency" env:"OCIS_MAX_CONCURRENCY;STORAGE_USERS_OCIS_MAX_CONCURRENCY" desc:"Maximum number of concurrent go-routines. Higher values can potentially get work done faster but will also cause more load on the system. Values of 0 or below will be ignored and the default value will be used." introductionVersion:"pre5.0"`
	AsyncUploads            bool   `yaml:"async_uploads" env:"OCIS_ASYNC_UPLOADS" desc:"Enable asynchronous file uploads." introductionVersion:"pre5.0"`
	MaxQuota                uint64 `yaml:"max_quota" env:"OCIS_SPACES_MAX_QUOTA;STORAGE_USERS_OCIS_MAX_QUOTA" desc:"Set a global max quota for spaces in bytes. A value of 0 equals unlimited. If not using the global OCIS_SPACES_MAX_QUOTA, you must define the FRONTEND_MAX_QUOTA in the frontend service." introductionVersion:"pre5.0"`
	DisableVersioning       bool   `yaml:"disable_versioning" env:"OCIS_DISABLE_VERSIONING" desc:"Disables versioning of files. When set to true, new uploads with the same filename will overwrite existing files instead of creating a new version." introductionVersion:"7.0.0"`
}

// S3NGDriver is the storage driver configuration when using 's3ng' storage driver
type S3NGDriver struct {
	Propagator             string                 `yaml:"propagator" env:"OCIS_DECOMPOSEDFS_PROPAGATOR;STORAGE_USERS_S3NG_PROPAGATOR" desc:"The propagator used for decomposedfs. At the moment, only 'sync' is fully supported, 'async' is available as an experimental option." introductionVersion:"pre5.0"`
	AsyncPropagatorOptions AsyncPropagatorOptions `yaml:"async_propagator_options"`
	// Root is the absolute path to the location of the data
	Root                  string `yaml:"root" env:"STORAGE_USERS_S3NG_ROOT" desc:"The directory where the filesystem storage will store metadata for blobs. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users." introductionVersion:"pre5.0"`
	UserLayout            string `yaml:"user_layout" env:"STORAGE_USERS_S3NG_USER_LAYOUT" desc:"Template string for the user storage layout in the user directory." introductionVersion:"pre5.0"`
	PermissionsEndpoint   string `yaml:"permissions_endpoint" env:"STORAGE_USERS_PERMISSION_ENDPOINT;STORAGE_USERS_S3NG_PERMISSIONS_ENDPOINT" desc:"Endpoint of the permissions service. The endpoints can differ for 'ocis' and 's3ng'." introductionVersion:"pre5.0"`
	Region                string `yaml:"region" env:"STORAGE_USERS_S3NG_REGION" desc:"Region of the S3 bucket." introductionVersion:"pre5.0"`
	AccessKey             string `yaml:"access_key" env:"STORAGE_USERS_S3NG_ACCESS_KEY" desc:"Access key for the S3 bucket." introductionVersion:"pre5.0"`
	SecretKey             string `yaml:"secret_key" env:"STORAGE_USERS_S3NG_SECRET_KEY" desc:"Secret key for the S3 bucket." introductionVersion:"pre5.0"`
	Endpoint              string `yaml:"endpoint" env:"STORAGE_USERS_S3NG_ENDPOINT" desc:"Endpoint for the S3 bucket." introductionVersion:"pre5.0"`
	Bucket                string `yaml:"bucket" env:"STORAGE_USERS_S3NG_BUCKET" desc:"Name of the S3 bucket." introductionVersion:"pre5.0"`
	DisableContentSha256  bool   `yaml:"put_object_disable_content_sha254" env:"STORAGE_USERS_S3NG_PUT_OBJECT_DISABLE_CONTENT_SHA256" desc:"Disable sending content sha256 when copying objects to S3." introductionVersion:"5.0"`
	DisableMultipart      bool   `yaml:"put_object_disable_multipart" env:"STORAGE_USERS_S3NG_PUT_OBJECT_DISABLE_MULTIPART" desc:"Disable multipart uploads when copying objects to S3" introductionVersion:"5.0"`
	SendContentMd5        bool   `yaml:"put_object_send_content_md5" env:"STORAGE_USERS_S3NG_PUT_OBJECT_SEND_CONTENT_MD5" desc:"Send a Content-MD5 header when copying objects to S3." introductionVersion:"5.0"`
	ConcurrentStreamParts bool   `yaml:"put_object_concurrent_stream_parts" env:"STORAGE_USERS_S3NG_PUT_OBJECT_CONCURRENT_STREAM_PARTS" desc:"Always precreate parts when copying objects to S3." introductionVersion:"5.0"`
	NumThreads            uint   `yaml:"put_object_num_threads" env:"STORAGE_USERS_S3NG_PUT_OBJECT_NUM_THREADS" desc:"Number of concurrent uploads to use when copying objects to S3." introductionVersion:"5.0"`
	PartSize              uint64 `yaml:"put_object_part_size" env:"STORAGE_USERS_S3NG_PUT_OBJECT_PART_SIZE" desc:"Part size for concurrent uploads to S3. If no value or 0 is set, the library's default value of 16MB is used. The value range is min 5MB and max 5GB." introductionVersion:"5.0"`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template" env:"STORAGE_USERS_S3NG_PERSONAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct personal space aliases." introductionVersion:"pre5.0"`
	PersonalSpacePathTemplate  string `yaml:"personalspacepath_template" env:"STORAGE_USERS_S3NG_PERSONAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the personal space roots." introductionVersion:"6.0.0"`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template" env:"STORAGE_USERS_S3NG_GENERAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct general space aliases." introductionVersion:"pre5.0"`
	GeneralSpacePathTemplate  string `yaml:"generalspacepath_template" env:"STORAGE_USERS_S3NG_GENERAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the projects space roots." introductionVersion:"6.0.0"`
	// ShareFolder defines the name of the folder jailing all shares
	ShareFolder             string `yaml:"share_folder" env:"STORAGE_USERS_S3NG_SHARE_FOLDER" desc:"Name of the folder jailing all shares." introductionVersion:"pre5.0"`
	MaxAcquireLockCycles    int    `yaml:"max_acquire_lock_cycles" env:"STORAGE_USERS_S3NG_MAX_ACQUIRE_LOCK_CYCLES" desc:"When trying to lock files, ocis will try this amount of times to acquire the lock before failing. After each try it will wait for an increasing amount of time. Values of 0 or below will be ignored and the default value of 20 will be used." introductionVersion:"pre5.0"`
	LockCycleDurationFactor int    `yaml:"lock_cycle_duration_factor" env:"STORAGE_USERS_S3NG_LOCK_CYCLE_DURATION_FACTOR" desc:"When trying to lock files, ocis will multiply the cycle with this factor and use it as a millisecond timeout. Values of 0 or below will be ignored and the default value of 30 will be used." introductionVersion:"pre5.0"`
	MaxConcurrency          int    `yaml:"max_concurrency" env:"OCIS_MAX_CONCURRENCY;STORAGE_USERS_S3NG_MAX_CONCURRENCY" desc:"Maximum number of concurrent go-routines. Higher values can potentially get work done faster but will also cause more load on the system. Values of 0 or below will be ignored and the default value of 100 will be used." introductionVersion:"pre5.0"`
	DisableVersioning       bool   `yaml:"disable_versioning" env:"OCIS_DISABLE_VERSIONING" desc:"Disables versioning of files. When set to true, new uploads with the same filename will overwrite existing files instead of creating a new version." introductionVersion:"7.0.0"`
}

// OwnCloudSQLDriver is the storage driver configuration when using 'owncloudsql' storage driver
type OwnCloudSQLDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root" env:"STORAGE_USERS_OWNCLOUDSQL_DATADIR" desc:"The directory where the filesystem storage will store SQL migration data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/owncloud." introductionVersion:"pre5.0"`
	// ShareFolder defines the name of the folder jailing all shares
	ShareFolder           string `yaml:"share_folder" env:"STORAGE_USERS_OWNCLOUDSQL_SHARE_FOLDER" desc:"Name of the folder jailing all shares." introductionVersion:"pre5.0"`
	UserLayout            string `yaml:"user_layout" env:"STORAGE_USERS_OWNCLOUDSQL_LAYOUT" desc:"Path layout to use to navigate into a users folder in an owncloud data directory" introductionVersion:"pre5.0"`
	UploadInfoDir         string `yaml:"upload_info_dir" env:"STORAGE_USERS_OWNCLOUDSQL_UPLOADINFO_DIR" desc:"The directory where the filesystem will store uploads temporarily. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/uploadinfo." introductionVersion:"pre5.0"`
	DBUsername            string `yaml:"db_username" env:"STORAGE_USERS_OWNCLOUDSQL_DB_USERNAME" desc:"Username for the database." introductionVersion:"pre5.0"`
	DBPassword            string `yaml:"db_password" env:"STORAGE_USERS_OWNCLOUDSQL_DB_PASSWORD" desc:"Password for the database." introductionVersion:"pre5.0"`
	DBHost                string `yaml:"db_host" env:"STORAGE_USERS_OWNCLOUDSQL_DB_HOST" desc:"Hostname or IP of the database server." introductionVersion:"pre5.0"`
	DBPort                int    `yaml:"db_port" env:"STORAGE_USERS_OWNCLOUDSQL_DB_PORT" desc:"Port that the database server is listening on." introductionVersion:"pre5.0"`
	DBName                string `yaml:"db_name" env:"STORAGE_USERS_OWNCLOUDSQL_DB_NAME" desc:"Name of the database to be used." introductionVersion:"pre5.0"`
	UsersProviderEndpoint string `yaml:"users_provider_endpoint" env:"STORAGE_USERS_OWNCLOUDSQL_USERS_PROVIDER_ENDPOINT" desc:"Endpoint of the users provider." introductionVersion:"pre5.0"`
}

// PosixDriver is the storage driver configuration when using 'posix' storage driver
type PosixDriver struct {
	// Root is the absolute path to the location of the data
	Root                      string        `yaml:"root" env:"STORAGE_USERS_POSIX_ROOT" desc:"The directory where the filesystem storage will store its data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users." introductionVersion:"6.0.0"`
	PersonalSpacePathTemplate string        `yaml:"personalspacepath_template" env:"STORAGE_USERS_POSIX_PERSONAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the personal space roots." introductionVersion:"6.0.0"`
	GeneralSpacePathTemplate  string        `yaml:"generalspacepath_template" env:"STORAGE_USERS_POSIX_GENERAL_SPACE_PATH_TEMPLATE" desc:"Template string to construct the paths of the projects space roots." introductionVersion:"6.0.0"`
	PermissionsEndpoint       string        `yaml:"permissions_endpoint" env:"STORAGE_USERS_PERMISSION_ENDPOINT;STORAGE_USERS_POSIX_PERMISSIONS_ENDPOINT" desc:"Endpoint of the permissions service. The endpoints can differ for 'ocis', 'posix' and 's3ng'." introductionVersion:"6.0.0"`
	AsyncUploads              bool          `yaml:"async_uploads" env:"OCIS_ASYNC_UPLOADS" desc:"Enable asynchronous file uploads." introductionVersion:"pre5.0"`
	ScanDebounceDelay         time.Duration `yaml:"scan_debounce_delay" env:"STORAGE_USERS_POSIX_SCAN_DEBOUNCE_DELAY" desc:"The time in milliseconds to wait before scanning the filesystem for changes after a change has been detected." introductionVersion:"6.0.0"`

	UseSpaceGroups bool `yaml:"use_space_groups" env:"STORAGE_USERS_POSIX_USE_SPACE_GROUPS" desc:"Use space groups to manage permissions on spaces." introductionVersion:"6.0.0"`

	WatchType               string `yaml:"watch_type" env:"STORAGE_USERS_POSIX_WATCH_TYPE" desc:"Type of the watcher to use for getting notified about changes to the filesystem. Currently available options are 'inotifywait' (default), 'gpfswatchfolder' and 'gpfsfileauditlogging'." introductionVersion:"6.0.0"`
	WatchPath               string `yaml:"watch_path" env:"STORAGE_USERS_POSIX_WATCH_PATH" desc:"Path to the watch directory/file. Only applies to the 'gpfsfileauditlogging' and 'inotifywait' watcher, in which case it is the path of the file audit log file/base directory to watch." introductionVersion:"6.0.0"`
	WatchFolderKafkaBrokers string `yaml:"watch_folder_kafka_hosts" env:"STORAGE_USERS_POSIX_WATCH_FOLDER_KAFKA_BROKERS" desc:"Comma-separated list of kafka brokers to read the watchfolder events from." introductionVersion:"6.0.0"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Addr              string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;STORAGE_USERS_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture." introductionVersion:"pre5.0"`
	ClusterID         string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;STORAGE_USERS_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system." introductionVersion:"pre5.0"`
	TLSInsecure       bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;STORAGE_USERS_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates." introductionVersion:"pre5.0"`
	TLSRootCaCertPath string `yaml:"tls_root_ca_cert_path" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;STORAGE_USERS_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided STORAGE_USERS_EVENTS_TLS_INSECURE will be seen as false." introductionVersion:"pre5.0"`
	EnableTLS         bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;STORAGE_USERS_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"pre5.0"`
	NumConsumers      int    `yaml:"num_consumers" env:"STORAGE_USERS_EVENTS_NUM_CONSUMERS" desc:"The amount of concurrent event consumers to start. Event consumers are used for post-processing files. Multiple consumers increase parallelisation, but will also increase CPU and memory demands. The setting has no effect when the OCIS_ASYNC_UPLOADS is set to false. The default and minimum value is 1." introductionVersion:"pre5.0"`
	AuthUsername      string `yaml:"username" env:"OCIS_EVENTS_AUTH_USERNAME;STORAGE_USERS_EVENTS_AUTH_USERNAME" desc:"The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
	AuthPassword      string `yaml:"password" env:"OCIS_EVENTS_AUTH_PASSWORD;STORAGE_USERS_EVENTS_AUTH_PASSWORD" desc:"The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
}

// FilemetadataCache holds cache config
type FilemetadataCache struct {
	Store              string        `yaml:"store" env:"OCIS_CACHE_STORE;STORAGE_USERS_FILEMETADATA_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	Nodes              []string      `yaml:"nodes" env:"OCIS_CACHE_STORE_NODES;STORAGE_USERS_FILEMETADATA_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Database           string        `yaml:"database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	TTL                time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;STORAGE_USERS_FILEMETADATA_CACHE_TTL" desc:"Default time to live for user info in the user info cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	DisablePersistence bool          `yaml:"disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;STORAGE_USERS_FILEMETADATA_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	AuthUsername       string        `yaml:"username" env:"OCIS_CACHE_AUTH_USERNAME;STORAGE_USERS_FILEMETADATA_CACHE_AUTH_USERNAME" desc:"The username to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	AuthPassword       string        `yaml:"password" env:"OCIS_CACHE_AUTH_PASSWORD;STORAGE_USERS_FILEMETADATA_CACHE_AUTH_PASSWORD" desc:"The password to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
}

// IDCache holds cache config
type IDCache struct {
	Store              string        `yaml:"store" env:"OCIS_CACHE_STORE;STORAGE_USERS_ID_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	Nodes              []string      `yaml:"nodes" env:"OCIS_CACHE_STORE_NODES;STORAGE_USERS_ID_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Database           string        `yaml:"database" env:"OCIS_CACHE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	TTL                time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;STORAGE_USERS_ID_CACHE_TTL" desc:"Default time to live for user info in the user info cache. Only applied when access tokens have no expiration. Defaults to 300s which is derived from the underlaying package though not explicitly set as default. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	DisablePersistence bool          `yaml:"disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;STORAGE_USERS_ID_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	AuthUsername       string        `yaml:"username" env:"OCIS_CACHE_AUTH_USERNAME;STORAGE_USERS_ID_CACHE_AUTH_USERNAME" desc:"The username to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	AuthPassword       string        `yaml:"password" env:"OCIS_CACHE_AUTH_PASSWORD;STORAGE_USERS_ID_CACHE_AUTH_PASSWORD" desc:"The password to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
}

// S3Driver is the storage driver configuration when using 's3' storage driver
type S3Driver struct {
	// Root is the absolute path to the location of the data
	Root      string `yaml:"root"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
}

// EOSDriver is the storage driver configuration when using 'eos' storage driver
type EOSDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root"`
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
	// URL of the Slave EOS MGM.
	// Default is root://eos-example.org
	SlaveURL string `yaml:"slave_url"`
	// Location on the local fs where to store reads.
	// Defaults to os.TempDir()
	CacheDirectory string `yaml:"cache_directory"`
	// SecProtocol specifies the xrootd security protocol to use between the server and EOS.
	SecProtocol string `yaml:"sec_protocol"`
	// Keytab specifies the location of the keytab to use to authenticate to EOS.
	Keytab string `yaml:"keytab"`
	// SingleUsername is the username to use when SingleUserMode is enabled
	SingleUsername string `yaml:"single_username"`
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
	// gateway service to use for uid lookups
	GatewaySVC string `yaml:"gateway_svc"`
	// ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	GRPCURI     string
	UserLayout  string
}

// LocalDriver is the storage driver configuration when using 'local' storage driver
type LocalDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root"`
	// ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	UserLayout  string `yaml:"user_layout"`
}

// Tasks wraps task configurations
type Tasks struct {
	PurgeTrashBin PurgeTrashBin `yaml:"purge_trash_bin"`
}

// PurgeTrashBin contains all necessary configurations to clean up the respective trash cans
type PurgeTrashBin struct {
	UserID               string        `yaml:"user_id" env:"OCIS_ADMIN_USER_ID;STORAGE_USERS_PURGE_TRASH_BIN_USER_ID" desc:"ID of the user who collects all necessary information for deletion. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand." introductionVersion:"pre5.0"`
	PersonalDeleteBefore time.Duration `yaml:"personal_delete_before" env:"STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE" desc:"Specifies the period of time in which items that have been in the personal trash-bin for longer than this value should be deleted. A value of 0 means no automatic deletion. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	ProjectDeleteBefore  time.Duration `yaml:"project_delete_before" env:"STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE" desc:"Specifies the period of time in which items that have been in the project trash-bin for longer than this value should be deleted. A value of 0 means no automatic deletion. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ServiceAccountID     string `yaml:"service_account_id" env:"OCIS_SERVICE_ACCOUNT_ID;STORAGE_USERS_SERVICE_ACCOUNT_ID" desc:"The ID of the service account the service should use. See the 'auth-service' service description for more details." introductionVersion:"5.0"`
	ServiceAccountSecret string `yaml:"service_account_secret" env:"OCIS_SERVICE_ACCOUNT_SECRET;STORAGE_USERS_SERVICE_ACCOUNT_SECRET" desc:"The service account secret." introductionVersion:"5.0"`
}
