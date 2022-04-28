package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Logging         *Logging `yaml:"log"`
	Debug           Debug    `yaml:"debug"`
	Supervised      bool     `yaml:"supervised"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool          `yaml:"skip_user_groups_in_token"`
	AuthProvider          string        `yaml:"auth_provider" env:"AUTH_BASIC_AUTH_PROVIDER" desc:"The auth provider which should be used by the service"`
	AuthProviders         AuthProviders `yaml:"auth_providers"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_BASIC_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_BASIC_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_BASIC_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_BASIC_TRACING_COLLECTOR"`
}

type Logging struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_BASIC_LOG_LEVEL" desc:"The log level."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_BASIC_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_BASIC_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_BASIC_LOG_FILE" desc:"The target log file."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_BASIC_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"AUTH_BASIC_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"AUTH_BASIC_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"AUTH_BASIC_DEBUG_ZPAGES"`
}

type GRPCConfig struct {
	Addr     string `yaml:"addr" env:"AUTH_BASIC_GRPC_ADDR" desc:"The address of the grpc service."`
	Protocol string `yaml:"protocol" env:"AUTH_BASIC_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type AuthProviders struct {
	JSON        JSONProvider        `yaml:"json"`
	LDAP        LDAPProvider        `yaml:"ldap"`
	OwnCloudSQL OwnCloudSQLProvider `yaml:"owncloud_sql"`
}

type JSONProvider struct {
	File string `yaml:"file" env:"AUTH_BASIC_JSON_PROVIDER_FILE" desc:"The file to which the json provider writes the data."`
}

type LDAPProvider struct {
	URI              string          `yaml:"" env:"LDAP_URI;AUTH_BASIC_LDAP_URI"`
	CACert           string          `yaml:"" env:"LDAP_CACERT;AUTH_BASIC_LDAP_CACERT"`
	Insecure         bool            `yaml:"" env:"LDAP_INSECURE;AUTH_BASIC_LDAP_INSECURE"`
	BindDN           string          `yaml:"" env:"LDAP_BIND_DN;AUTH_BASIC_LDAP_BIND_DN"`
	BindPassword     string          `yaml:"" env:"LDAP_BIND_PASSWORD;AUTH_BASIC_LDAP_BIND_PASSWORD"`
	UserBaseDN       string          `yaml:"" env:"LDAP_USER_BASE_DN;AUTH_BASIC_LDAP_USER_BASE_DN"`
	GroupBaseDN      string          `yaml:"" env:"LDAP_GROUP_BASE_DN;AUTH_BASIC_LDAP_GROUP_BASE_DN"`
	UserScope        string          `yaml:"" env:"LDAP_USER_SCOPE;AUTH_BASIC_LDAP_USER_SCOPE"`
	GroupScope       string          `yaml:"" env:"LDAP_GROUP_SCOPE;AUTH_BASIC_LDAP_GROUP_SCOPE"`
	UserFilter       string          `yaml:"" env:"LDAP_USERFILTER;AUTH_BASIC_LDAP_USERFILTER"`
	GroupFilter      string          `yaml:"" env:"LDAP_GROUPFILTER;AUTH_BASIC_LDAP_USERFILTER"`
	UserObjectClass  string          `yaml:"" env:"LDAP_USER_OBJECTCLASS;AUTH_BASIC_LDAP_USER_OBJECTCLASS"`
	GroupObjectClass string          `yaml:"" env:"LDAP_GROUP_OBJECTCLASS;AUTH_BASIC_LDAP_GROUP_OBJECTCLASS"`
	LoginAttributes  []string        `yaml:"" env:"LDAP_LOGIN_ATTRIBUTES;AUTH_BASIC_LDAP_LOGIN_ATTRIBUTES"`
	IDP              string          `yaml:"" env:"OCIS_URL;AUTH_BASIC_IDP_URL"` // TODO what is this for?
	GatewayEndpoint  string          `yaml:""`                                   // TODO do we need this here?
	UserSchema       LDAPUserSchema  `yaml:""`
	GroupSchema      LDAPGroupSchema `yaml:""`
}

type LDAPUserSchema struct {
	ID              string `env:"LDAP_USER_SCHEMA_ID;AUTH_BASIC_LDAP_USER_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_USER_SCHEMA_MAIL;AUTH_BASIC_LDAP_USER_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_USER_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_USER_SCHEMA_DISPLAYNAME"`
	Username        string `env:"LDAP_USER_SCHEMA_USERNAME;AUTH_BASIC_LDAP_USER_SCHEMA_USERNAME"`
}

type LDAPGroupSchema struct {
	ID              string `env:"LDAP_GROUP_SCHEMA_ID;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID"`
	IDIsOctetString bool   `env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `env:"LDAP_GROUP_SCHEMA_MAIL;AUTH_BASIC_LDAP_GROUP_SCHEMA_MAIL"`
	DisplayName     string `env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_DISPLAYNAME"`
	Groupname       string `env:"LDAP_GROUP_SCHEMA_GROUPNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_GROUPNAME"`
	Member          string `env:"LDAP_GROUP_SCHEMA_MEMBER;AUTH_BASIC_LDAP_GROUP_SCHEMA_MEMBER"`
}

type OwnCloudSQLProvider struct {
	DBUsername       string
	DBPassword       string
	DBHost           string
	DBPort           int
	DBName           string
	IDP              string // TODO do we need this?
	Nobody           int64  // TODO what is this?
	JoinUsername     bool
	JoinOwnCloudUUID bool
}
