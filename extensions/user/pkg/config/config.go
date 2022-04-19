package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Logging         *Logging `yaml:"log"`
	Debug           Debug    `yaml:"debug"`
	Supervised      bool

	GRPC GRPCConfig `yaml:"grpc"`

	JWTSecret             string
	GatewayEndpoint       string
	SkipUserGroupsInToken bool
	UsersCacheExpiration  int
	Driver                string
	Drivers               Drivers
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;USERS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;USERS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;USERS_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;USERS_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;USERS_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;USERS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;USERS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;USERS_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"USERS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"USERS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"USERS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"USERS_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"USERS_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"USERS_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type Drivers struct {
	JSON        JSONDriver
	LDAP        LDAPDriver
	OwnCloudSQL OwnCloudSQLDriver
	REST        RESTProvider
}

type JSONDriver struct {
	File string
}

type LDAPDriver struct {
	URI              string
	CACert           string
	Insecure         bool
	BindDN           string
	BindPassword     string
	UserBaseDN       string
	GroupBaseDN      string
	UserFilter       string
	GroupFilter      string
	UserObjectClass  string
	GroupObjectClass string
	LoginAttributes  []string
	IDP              string // TODO what is this for?
	GatewayEndpoint  string // TODO do we need this here?
	UserSchema       LDAPUserSchema
	GroupSchema      LDAPGroupSchema
}

type LDAPUserSchema struct {
	ID              string
	IDIsOctetString bool
	Mail            string
	DisplayName     string
	Username        string
}

type LDAPGroupSchema struct {
	ID              string
	IDIsOctetString bool
	Mail            string
	DisplayName     string
	Groupname       string
	Member          string
}

type OwnCloudSQLDriver struct {
	DBUsername         string
	DBPassword         string
	DBHost             string
	DBPort             int
	DBName             string
	IDP                string // TODO do we need this?
	Nobody             int64  // TODO what is this?
	JoinUsername       bool
	JoinOwnCloudUUID   bool
	EnableMedialSearch bool
}

type RESTProvider struct {
	ClientID          string
	ClientSecret      string
	RedisAddr         string
	RedisUsername     string
	RedisPassword     string
	IDProvider        string
	APIBaseURL        string
	OIDCTokenEndpoint string
	TargetAPI         string
}
