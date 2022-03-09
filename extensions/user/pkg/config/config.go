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
	URI              string   `env:"LDAP_URI;USERS_LDAP_URI"`
	CACert           string   `env:"LDAP_CACERT;USERS_LDAP_CACERT"`
	Insecure         bool     `env:"LDAP_INSECURE;USERS_LDAP_INSECURE"`
	BindDN           string   `env:"LDAP_BIND_DN;USERS_LDAP_BIND_DN"`
	BindPassword     string   `env:"LDAP_BIND_PASSWORD;USERS_LDAP_BIND_PASSWORD"`
	UserBaseDN       string   `env:"LDAP_USER_BASE_DN;USERS_LDAP_USER_BASE_DN"`
	GroupBaseDN      string   `env:"LDAP_GROUP_BASE_DN;USERS_LDAP_GROUP_BASE_DN"`
	UserScope        string   `env:"LDAP_USER_SCOPE;USERS_LDAP_USER_SCOPE"`
	GroupScope       string   `env:"LDAP_GROUP_SCOPE;USERS_LDAP_GROUP_SCOPE"`
	UserFilter       string   `env:"LDAP_USERFILTER;USERS_LDAP_USERFILTER"`
	GroupFilter      string   `env:"LDAP_GROUPFILTER;USERS_LDAP_USERFILTER"`
	UserObjectClass  string   `env:"LDAP_USER_OBJECTCLASS;USERS_LDAP_USER_OBJECTCLASS"`
	GroupObjectClass string   `env:"LDAP_GROUP_OBJECTCLASS;USERS_LDAP_GROUP_OBJECTCLASS"`
	LoginAttributes  []string `env:"LDAP_LOGIN_ATTRIBUTES;USERS_LDAP_LOGIN_ATTRIBUTES"`
	IDP              string   `env:"OCIS_URL;USERS_IDP_URL"` // TODO what is this for?
	GatewayEndpoint  string   // TODO do we need this here?
	UserSchema       LDAPUserSchema
	GroupSchema      LDAPGroupSchema
}

type LDAPUserSchema struct {
	ID              string `env:"LDAP_USER_SCHEMA_ID;USERS_LDAP_USER_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;USERS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_USER_SCHEMA_MAIL;USERS_LDAP_USER_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_USER_SCHEMA_DISPLAYNAME;USERS_LDAP_USER_SCHEMA_DISPLAYNAME"`
	Username        string `env:"LDAP_USER_SCHEMA_USERNAME;USERS_LDAP_USER_SCHEMA_USERNAME"`
}

type LDAPGroupSchema struct {
	ID              string `env:"LDAP_GROUP_SCHEMA_ID;USERS_LDAP_GROUP_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;USERS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_GROUP_SCHEMA_MAIL;USERS_LDAP_GROUP_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;USERS_LDAP_GROUP_SCHEMA_DISPLAYNAME"`
	Groupname       string `env:"LDAP_GROUP_SCHEMA_GROUPNAME;USERS_LDAP_GROUP_SCHEMA_GROUPNAME"`
	Member          string `env:"LDAP_GROUP_SCHEMA_MEMBER;USERS_LDAP_GROUP_SCHEMA_MEMBER"`
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
