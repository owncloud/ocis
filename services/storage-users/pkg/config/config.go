package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"STORAGE_USERS_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	Driver           string  `yaml:"driver" env:"STORAGE_USERS_DRIVER" desc:"The storage driver which should be used by the service"`
	Drivers          Drivers `yaml:"drivers"`
	DataServerURL    string  `yaml:"data_server_url" env:"STORAGE_USERS_DATA_SERVER_URL" desc:"URL of the data server, needs to be reachable by the data gateway provided by the frontend service or the user if directly exposed."`
	Events           Events  `yaml:"events"`
	MountID          string  `yaml:"mount_id" env:"STORAGE_USERS_MOUNT_ID" desc:"Mount ID of this storage."`
	ExposeDataServer bool    `yaml:"expose_data_server" env:"STORAGE_USERS_EXPOSE_DATA_SERVER" desc:"Exposes the data server directly to users and bypasses the data gateway. Ensure that the data server address is reachable by users."`
	ReadOnly         bool    `yaml:"readonly" env:"STORAGE_USERS_READ_ONLY" desc:"Set this storage to be read-only."`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_USERS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_USERS_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_USERS_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_USERS_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_USERS_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_USERS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_USERS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_USERS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_USERS_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"STORAGE_USERS_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_USERS_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"STORAGE_USERS_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_USERS_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_USERS_GRPC_PROTOCOL" desc:"The transport protocol of the GPRC service."`
}

type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"STORAGE_USERS_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"STORAGE_USERS_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service."`
	Prefix    string
}

type Drivers struct {
	OCIS        OCISDriver        `yaml:"ocis"`
	S3NG        S3NGDriver        `yaml:"s3ng"`
	OwnCloudSQL OwnCloudSQLDriver `yaml:"owncloudsql"`

	S3    S3Driver    `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
	EOS   EOSDriver   `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
	Local LocalDriver `yaml:",omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type OCISDriver struct {
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root" env:"STORAGE_USERS_OCIS_ROOT" desc:"Path for the persistence directory."`
	UserLayout          string `yaml:"user_layout" env:"STORAGE_USERS_OCIS_USER_LAYOUT" desc:"Template string for the user storage layout in the persistence directory."`
	PermissionsEndpoint string `yaml:"permissions_endpoint" env:"STORAGE_USERS_PERMISSION_ENDPOINT,STORAGE_USERS_OCIS_PERMISSIONS_ENDPOINT" desc:"Endpoint of the permissions service."`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template" env:"STORAGE_USERS_OCIS_PERSONAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct personal space aliases."`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template" env:"STORAGE_USERS_OCIS_GENERAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct general space aliases."`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder" env:"STORAGE_USERS_OCIS_SHARE_FOLDER" desc:"Name of the folder jailing all shares."`
}

type S3NGDriver struct {
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root" env:"STORAGE_USERS_S3NG_ROOT" desc:"Path for the persistence directory."`
	UserLayout          string `yaml:"user_layout" env:"STORAGE_USERS_S3NG_USER_LAYOUT" desc:"Template string for the user storage layout in the persistence directory."`
	PermissionsEndpoint string `yaml:"permissions_endpoint" env:"STORAGE_USERS_PERMISSION_ENDPOINT;STORAGE_USERS_S3NG_PERMISSIONS_ENDPOINT" desc:"Endpoint of the permissions service."`
	Region              string `yaml:"region" env:"STORAGE_USERS_S3NG_REGION" desc:"Region of the S3 bucket."`
	AccessKey           string `yaml:"access_key" env:"STORAGE_USERS_S3NG_ACCESS_KEY" desc:"Access key for the S3 bucket."`
	SecretKey           string `yaml:"secret_key" env:"STORAGE_USERS_S3NG_SECRET_KEY" desc:"Secret key for the S3 bucket."`
	Endpoint            string `yaml:"endpoint" env:"STORAGE_USERS_S3NG_ENDPOINT" desc:"Endpoint for the S3 bucket."`
	Bucket              string `yaml:"bucket" env:"STORAGE_USERS_S3NG_BUCKET" desc:"Name of the S3 bucket."`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template" env:"STORAGE_USERS_S3NG_PERSONAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct personal space aliases."`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template" env:"STORAGE_USERS_S3NG_GENERAL_SPACE_ALIAS_TEMPLATE" desc:"Template string to construct general space aliases."`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder" env:"STORAGE_USERS_S3NG_SHARE_FOLDER" desc:"Name of the folder jailing all shares."`
}

type OwnCloudSQLDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root" env:"STORAGE_USERS_OWNCLOUDSQL_DATADIR" desc:"Path for the persistence directory."`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder           string `yaml:"share_folder" env:"STORAGE_USERS_OWNCLOUDSQL_SHARE_FOLDER" desc:"Name of the folder jailing all shares."`
	UserLayout            string `yaml:"user_layout" env:"STORAGE_USERS_OWNCLOUDSQL_LAYOUT" desc:"Path layout to use to navigate into a users folder in an owncloud data directory"`
	UploadInfoDir         string `yaml:"upload_info_dir" env:"STORAGE_USERS_OWNCLOUDSQL_UPLOADINFO_DIR" desc:"Path to a directory, where uploads will be stored temporarily."`
	DBUsername            string `yaml:"db_username" env:"STORAGE_USERS_OWNCLOUDSQL_DB_USERNAME" desc:"Username for the database."`
	DBPassword            string `yaml:"db_password" env:"STORAGE_USERS_OWNCLOUDSQL_DB_PASSWORD" desc:"Password for the database."`
	DBHost                string `yaml:"db_host" env:"STORAGE_USERS_OWNCLOUDSQL_DB_HOST" desc:"Hostname or IP of the database server."`
	DBPort                int    `yaml:"db_port" env:"STORAGE_USERS_OWNCLOUDSQL_DB_PORT" desc:"Port that the database server is listening on."`
	DBName                string `yaml:"db_name" env:"STORAGE_USERS_OWNCLOUDSQL_DB_NAME" desc:"Name of the database to be used."`
	UsersProviderEndpoint string `yaml:"users_provider_endpoint" env:"STORAGE_USERS_OWNCLOUDSQL_USERS_PROVIDER_ENDPOINT" desc:"Endpoint of the users provider."`
}

type Events struct {
	Addr      string `yaml:"endpoint" env:"STORAGE_USERS_EVENTS_ENDPOINT" desc:"The address of the streaming service"`
	ClusterID string `yaml:"cluster" env:"STORAGE_USERS_EVENTS_CLUSTER" desc:"The clusterID of the streaming service. Mandatory when using the NATS service."`
}
type S3Driver struct {
	// Root is the absolute path to the location of the data
	Root      string `yaml:"root"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
}
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
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	GRPCURI     string
	UserLayout  string
}

type LocalDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	UserLayout  string `yaml:"user_layout"`
}
