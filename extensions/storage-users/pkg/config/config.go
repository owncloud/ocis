package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing,omitempty"`
	Logging         *Logging `yaml:"log,omitempty"`
	Debug           Debug    `yaml:"debug,omitempty"`
	Supervised      bool     `yaml:"supervised,omitempty"`

	GRPC GRPCConfig `yaml:"grpc,omitempty"`
	HTTP HTTPConfig `yaml:"http,omitempty"`

	TokenManager *TokenManager `yaml:"token_manager,omitempty"`
	Reva         *Reva         `yaml:"reva,omitempty"`

	Context context.Context `yaml:"context,omitempty"`

	SkipUserGroupsInToken bool    `yaml:"skip_user_groups_in_token,omitempty"`
	Driver                string  `yaml:"driver,omitempty" env:"STORAGE_USERS_DRIVER" desc:"The storage driver which should be used by the service"`
	Drivers               Drivers `yaml:"drivers,omitempty"`
	DataServerURL         string  `yaml:"data_server_url,omitempty"`
	TempFolder            string  `yaml:"temp_folder,omitempty"`
	DataProviderInsecure  bool    `yaml:"data_provider_insecure,omitempty" env:"OCIS_INSECURE;STORAGE_USERS_DATAPROVIDER_INSECURE"`
	Events                Events  `yaml:"events,omitempty"`
	MountID               string  `yaml:"mount_id,omitempty"`
	ExposeDataServer      bool    `yaml:"expose_data_server,omitempty"`
	ReadOnly              bool    `yaml:"readonly,omitempty"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_USERS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_USERS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_USERS_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_USERS_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_USERS_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_USERS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_USERS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_USERS_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_USERS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORAGE_USERS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_USERS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_USERS_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"STORAGE_USERS_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"STORAGE_USERS_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type HTTPConfig struct {
	Addr     string `yaml:"addr" env:"STORAGE_USERS_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"STORAGE_USERS_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
	Prefix   string
}

type Drivers struct {
	EOS         EOSDriver
	Local       LocalDriver
	OCIS        OCISDriver
	S3          S3Driver
	S3NG        S3NGDriver
	OwnCloudSQL OwnCloudSQLDriver
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
	Root string `yaml:"root" env:"STORAGE_USERS_LOCAL_ROOT"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
	UserLayout  string
}

type OCISDriver struct {
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root" env:"STORAGE_USERS_OCIS_ROOT"`
	UserLayout          string
	PermissionsEndpoint string
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template"`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
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

type S3NGDriver struct {
	// Root is the absolute path to the location of the data
	Root                string `yaml:"root"`
	UserLayout          string
	PermissionsEndpoint string
	Region              string `yaml:"region"`
	AccessKey           string `yaml:"access_key"`
	SecretKey           string `yaml:"secret_key"`
	Endpoint            string `yaml:"endpoint"`
	Bucket              string `yaml:"bucket"`
	// PersonalSpaceAliasTemplate  contains the template used to construct
	// the personal space alias, eg: `"{{.SpaceType}}/{{.User.Username | lower}}"`
	PersonalSpaceAliasTemplate string `yaml:"personalspacealias_template"`
	// GeneralSpaceAliasTemplate contains the template used to construct
	// the general space alias, eg: `{{.SpaceType}}/{{.SpaceName | replace " " "-" | lower}}`
	GeneralSpaceAliasTemplate string `yaml:"generalspacealias_template"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder string `yaml:"share_folder"`
}

type OwnCloudSQLDriver struct {
	// Root is the absolute path to the location of the data
	Root string `yaml:"root" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DATADIR"`
	//ShareFolder defines the name of the folder jailing all shares
	ShareFolder           string `yaml:"share_folder" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_SHARE_FOLDER"`
	UserLayout            string `env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_LAYOUT"`
	UploadInfoDir         string `yaml:"upload_info_dir" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR"`
	DBUsername            string `yaml:"db_username" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBUSERNAME"`
	DBPassword            string `yaml:"db_password" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPASSWORD"`
	DBHost                string `yaml:"db_host" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBHOST"`
	DBPort                int    `yaml:"db_port" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPORT"`
	DBName                string `yaml:"db_name" env:"STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBNAME"`
	UsersProviderEndpoint string
}

type Events struct {
	Addr      string
	ClusterID string
}
