package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	*shared.Commons `yaml:"-"`
	Service         Service  `yaml:"-"`
	Tracing         *Tracing `yaml:"tracing"`
	Log             *Log     `yaml:"log"`
	Debug           Debug    `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"GROUPS_SKIP_USER_GROUPS_IN_TOKEN"`

	GroupMembersCacheExpiration int     `yaml:"group_members_cache_expiration"`
	Driver                      string  `yaml:"driver"`
	Drivers                     Drivers `yaml:"drivers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GROUPS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GROUPS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GROUPS_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GROUPS_TRACING_COLLECTOR"`
}

type Log struct {
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
	Addr      string `yaml:"addr" env:"GROUPS_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"GROUPS_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type Drivers struct {
	LDAP        LDAPDriver        `yaml:"ldap"`
	OwnCloudSQL OwnCloudSQLDriver `yaml:"owncloudsql"`

	JSON JSONDriver   `yaml:"json,omitempty"` // not supported by the oCIS product, therefore not part of docs
	REST RESTProvider `yaml:"rest,omitempty"` // not supported by the oCIS product, therefore not part of docs
}
type LDAPDriver struct {
	URI              string          `yaml:"uri" env:"LDAP_URI;GROUPS_LDAP_URI"`
	CACert           string          `yaml:"ca_cert" env:"LDAP_CACERT;GROUPS_LDAP_CACERT"`
	Insecure         bool            `yaml:"insecure" env:"LDAP_INSECURE;GROUPS_LDAP_INSECURE"`
	BindDN           string          `yaml:"bind_dn" env:"LDAP_BIND_DN;GROUPS_LDAP_BIND_DN"`
	BindPassword     string          `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;GROUPS_LDAP_BIND_PASSWORD"`
	UserBaseDN       string          `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;GROUPS_LDAP_USER_BASE_DN"`
	GroupBaseDN      string          `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;GROUPS_LDAP_GROUP_BASE_DN"`
	UserScope        string          `yaml:"user_scope" env:"LDAP_USER_SCOPE;GROUPS_LDAP_USER_SCOPE"`
	GroupScope       string          `yaml:"group_scope" env:"LDAP_GROUP_SCOPE;GROUPS_LDAP_GROUP_SCOPE"`
	UserFilter       string          `yaml:"user_filter" env:"LDAP_USERFILTER;GROUPS_LDAP_USERFILTER"`
	GroupFilter      string          `yaml:"group_filter" env:"LDAP_GROUPFILTER;GROUPS_LDAP_USERFILTER"`
	UserObjectClass  string          `yaml:"user_object_class" env:"LDAP_USER_OBJECTCLASS;GROUPS_LDAP_USER_OBJECTCLASS"`
	GroupObjectClass string          `yaml:"group_object_class" env:"LDAP_GROUP_OBJECTCLASS;GROUPS_LDAP_GROUP_OBJECTCLASS"`
	LoginAttributes  []string        `yaml:"login_attributes" env:"LDAP_LOGIN_ATTRIBUTES;GROUPS_LDAP_LOGIN_ATTRIBUTES"`
	IDP              string          `yaml:"idp" env:"OCIS_URL;OCIS_OIDC_ISSUER;GROUPS_IDP_URL"`
	UserSchema       LDAPUserSchema  `yaml:"user_schema"`
	GroupSchema      LDAPGroupSchema `yaml:"group_schema"`
}

type LDAPUserSchema struct {
	ID              string `yaml:"id" env:"LDAP_USER_SCHEMA_ID;GROUPS_LDAP_USER_SCHEMA_ID"`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `yaml:"mail" env:"LDAP_USER_SCHEMA_MAIL;GROUPS_LDAP_USER_SCHEMA_MAIL"`
	DisplayName     string `yaml:"display_name" env:"LDAP_USER_SCHEMA_DISPLAYNAME;GROUPS_LDAP_USER_SCHEMA_DISPLAYNAME"`
	Username        string `yaml:"user_name" env:"LDAP_USER_SCHEMA_USERNAME;GROUPS_LDAP_USER_SCHEMA_USERNAME"`
}

type LDAPGroupSchema struct {
	ID              string `yaml:"id" env:"LDAP_GROUP_SCHEMA_ID;GROUPS_LDAP_GROUP_SCHEMA_ID"`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `yaml:"mail" env:"LDAP_GROUP_SCHEMA_MAIL;GROUPS_LDAP_GROUP_SCHEMA_MAIL"`
	DisplayName     string `yaml:"display_name" env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;GROUPS_LDAP_GROUP_SCHEMA_DISPLAYNAME"`
	Groupname       string `yaml:"group_name" env:"LDAP_GROUP_SCHEMA_GROUPNAME;GROUPS_LDAP_GROUP_SCHEMA_GROUPNAME"`
	Member          string `yaml:"member" env:"LDAP_GROUP_SCHEMA_MEMBER;GROUPS_LDAP_GROUP_SCHEMA_MEMBER"`
}

type OwnCloudSQLDriver struct {
	DBUsername         string `yaml:"db_username" env:"GROUPS_OWNCLOUDSQL_DB_USERNAME"`
	DBPassword         string `yaml:"db_password" env:"GROUPS_OWNCLOUDSQL_DB_PASSWORD"`
	DBHost             string `yaml:"db_host" env:"GROUPS_OWNCLOUDSQL_DB_HOST"`
	DBPort             int    `yaml:"db_port" env:"GROUPS_OWNCLOUDSQL_DB_PORT"`
	DBName             string `yaml:"db_name" env:"GROUPS_OWNCLOUDSQL_DB_NAME"`
	IDP                string `yaml:"idp" env:"GROUPS_OWNCLOUDSQL_IDP"`
	Nobody             int64  `yaml:"nobody" env:"GROUPS_OWNCLOUDSQL_NOBODY"` // TODO what is this?
	JoinUsername       bool   `yaml:"join_username" env:"GROUPS_OWNCLOUDSQL_JOIN_USERNAME"`
	JoinOwnCloudUUID   bool   `yaml:"join_owncloud_uuid" env:"GROUPS_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID"`
	EnableMedialSearch bool   `yaml:"enable_medial_search" env:"GROUPS_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH"`
}

type JSONDriver struct {
	File string
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
