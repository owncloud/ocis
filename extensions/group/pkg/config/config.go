package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing,omitempty"`
	Logging         *Logging `yaml:"log,omitempty"`
	Debug           Debug    `yaml:"debug,omitempty"`
	Supervised      bool     `yaml:"supervised,omitempty"`

	GRPC GRPCConfig `yaml:"grpc,omitempty"`

	TokenManager *TokenManager `yaml:"token_manager,omitempty"`
	Reva         *Reva         `yaml:"reva,omitempty"`

	SkipUserGroupsInToken       bool    `yaml:"skip_user_groups_in_token,omitempty"`
	GroupMembersCacheExpiration int     `yaml:"group_members_cache_expiration,omitempty"`
	Driver                      string  `yaml:"driver,omitempty"`
	Drivers                     Drivers `yaml:"drivers,omitempty"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GROUPS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GROUPS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GROUPS_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GROUPS_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GROUPS_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;GROUPS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;GROUPS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;GROUPS_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"GROUPS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"GROUPS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"GROUPS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"GROUPS_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"GROUPS_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"GROUPS_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
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
	URI              string   `env:"LDAP_URI;GROUPS_LDAP_URI"`
	CACert           string   `env:"LDAP_CACERT;GROUPS_LDAP_CACERT"`
	Insecure         bool     `env:"LDAP_INSECURE;GROUPS_LDAP_INSECURE"`
	BindDN           string   `env:"LDAP_BIND_DN;GROUPS_LDAP_BIND_DN"`
	BindPassword     string   `env:"LDAP_BIND_PASSWORD;GROUPS_LDAP_BIND_PASSWORD"`
	UserBaseDN       string   `env:"LDAP_USER_BASE_DN;GROUPS_LDAP_USER_BASE_DN"`
	GroupBaseDN      string   `env:"LDAP_GROUP_BASE_DN;GROUPS_LDAP_GROUP_BASE_DN"`
	UserFilter       string   `env:"LDAP_USERFILTER;GROUPS_LDAP_USERFILTER"`
	GroupFilter      string   `env:"LDAP_GROUPFILTER;GROUPS_LDAP_USERFILTER"`
	UserObjectClass  string   `env:"LDAP_USER_OBJECTCLASS;GROUPS_LDAP_USER_OBJECTCLASS"`
	GroupObjectClass string   `env:"LDAP_GROUP_OBJECTCLASS;GROUPS_LDAP_GROUP_OBJECTCLASS"`
	LoginAttributes  []string `env:"LDAP_LOGIN_ATTRIBUTES;GROUPS_LDAP_LOGIN_ATTRIBUTES"`
	IDP              string   `env:"OCIS_URL;GROUPS_IDP_URL"` // TODO what is this for?
	GatewayEndpoint  string   // TODO do we need this here?
	UserSchema       LDAPUserSchema
	GroupSchema      LDAPGroupSchema
}

type LDAPUserSchema struct {
	ID              string `env:"LDAP_USER_SCHEMA_ID;GROUPS_LDAP_USER_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_USER_SCHEMA_MAIL;GROUPS_LDAP_USER_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_USER_SCHEMA_DISPLAYNAME;GROUPS_LDAP_USER_SCHEMA_DISPLAYNAME"`
	Username        string `env:"LDAP_USER_SCHEMA_USERNAME;GROUPS_LDAP_USER_SCHEMA_USERNAME"`
}

type LDAPGroupSchema struct {
	ID              string `env:"LDAP_GROUP_SCHEMA_ID;GROUPS_LDAP_GROUP_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_GROUP_SCHEMA_MAIL;GROUPS_LDAP_GROUP_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;GROUPS_LDAP_GROUP_SCHEMA_DISPLAYNAME"`
	Groupname       string `env:"LDAP_GROUP_SCHEMA_GROUPNAME;GROUPS_LDAP_GROUP_SCHEMA_GROUPNAME"`
	Member          string `env:"LDAP_GROUP_SCHEMA_MEMBER;GROUPS_LDAP_GROUP_SCHEMA_MEMBER"`
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
