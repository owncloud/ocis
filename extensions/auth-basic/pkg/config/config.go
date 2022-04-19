package config

import "github.com/owncloud/ocis/ocis-pkg/shared"

type Config struct {
	*shared.Commons `yaml:"-"`
	Tracing         *TracingConfig `yaml:"tracing"`
	Logging         *LoggingConfig `yaml:"log"`
	Service         ServiceConfig
	DebugService    DebugServiceConfig `yaml:"debug"`
	Supervised      bool
}
type TracingConfig struct {
	Enabled     bool
	Endpoint    string
	Collector   string
	ServiceName string
	Type        string
}

type LoggingConfig struct {
	Level  string
	Pretty bool
	Color  bool
	File   string
}

type ServiceConfig struct {
	JWTSecret             string
	GatewayEndpoint       string
	SkipUserGroupsInToken bool
	Network               string // TODO: name transport or protocol?
	Address               string
	AuthManager           string
	AuthManagers          AuthManagers
}

type DebugServiceConfig struct {
	Address string
	Pprof   bool
	Zpages  bool
	Token   string
}

type AuthManagers struct {
	JSON        JSONManager
	LDAP        LDAPManager
	OwnCloudSQL OwnCloudSQLManager
}

type JSONManager struct {
	Users string // TODO is there a better name?
}

type LDAPManager struct {
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

type OwnCloudSQLManager struct {
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
