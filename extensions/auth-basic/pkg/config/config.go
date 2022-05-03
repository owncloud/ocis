package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
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

	SkipUserGroupsInToken bool          `yaml:"-"`
	AuthProvider          string        `yaml:"auth_provider" env:"AUTH_BASIC_AUTH_PROVIDER" desc:"The auth provider which should be used by the service"`
	AuthProviders         AuthProviders `yaml:"auth_providers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_BASIC_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_BASIC_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_BASIC_TRACING_ENDPOINT" desc:"The endpoint to the tracing collector."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_BASIC_TRACING_COLLECTOR"`
}

type Log struct {
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
	Addr      string `yaml:"addr" env:"AUTH_BASIC_GRPC_ADDR" desc:"The address of the grpc service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"AUTH_BASIC_GRPC_PROTOCOL" desc:"The transport protocol of the grpc service."`
}

type AuthProviders struct {
	LDAP        LDAPProvider        `yaml:"ldap"`
	OwnCloudSQL OwnCloudSQLProvider `yaml:"owncloud_sql"`
	JSON        JSONProvider        `yaml:"json,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type JSONProvider struct {
	File string `yaml:"file,omitempty"`
}

type LDAPProvider struct {
	URI              string          `yaml:"uri" env:"LDAP_URI;AUTH_BASIC_LDAP_URI"`
	CACert           string          `yaml:"ca_cert" env:"LDAP_CACERT;AUTH_BASIC_LDAP_CACERT"`
	Insecure         bool            `yaml:"insecure" env:"LDAP_INSECURE;AUTH_BASIC_LDAP_INSECURE"`
	BindDN           string          `yaml:"bind_dn" env:"LDAP_BIND_DN;AUTH_BASIC_LDAP_BIND_DN"`
	BindPassword     string          `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;AUTH_BASIC_LDAP_BIND_PASSWORD"`
	UserBaseDN       string          `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;AUTH_BASIC_LDAP_USER_BASE_DN"`
	GroupBaseDN      string          `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;AUTH_BASIC_LDAP_GROUP_BASE_DN"`
	UserScope        string          `yaml:"user_scope" env:"LDAP_USER_SCOPE;AUTH_BASIC_LDAP_USER_SCOPE"`
	GroupScope       string          `yaml:"group_scope" env:"LDAP_GROUP_SCOPE;AUTH_BASIC_LDAP_GROUP_SCOPE"`
	UserFilter       string          `yaml:"user_filter" env:"LDAP_USERFILTER;AUTH_BASIC_LDAP_USERFILTER"`
	GroupFilter      string          `yaml:"group_filter" env:"LDAP_GROUPFILTER;AUTH_BASIC_LDAP_USERFILTER"`
	UserObjectClass  string          `yaml:"user_object_filter" env:"LDAP_USER_OBJECTCLASS;AUTH_BASIC_LDAP_USER_OBJECTCLASS"`
	GroupObjectClass string          `yaml:"group_object_class" env:"LDAP_GROUP_OBJECTCLASS;AUTH_BASIC_LDAP_GROUP_OBJECTCLASS"`
	LoginAttributes  []string        `yaml:"login_attributes" env:"LDAP_LOGIN_ATTRIBUTES;AUTH_BASIC_LDAP_LOGIN_ATTRIBUTES"`
	IDP              string          `yaml:"idp" env:"OCIS_URL;AUTH_BASIC_IDP_URL"`
	UserSchema       LDAPUserSchema  `yaml:"user_schema"`
	GroupSchema      LDAPGroupSchema `yaml:"group_schema"`
}

type LDAPUserSchema struct {
	ID              string `yaml:"id" env:"LDAP_USER_SCHEMA_ID;AUTH_BASIC_LDAP_USER_SCHEMA_ID"`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `yaml:"mail" env:"LDAP_USER_SCHEMA_MAIL;AUTH_BASIC_LDAP_USER_SCHEMA_MAIL"`
	DisplayName     string `yaml:"display_name" env:"LDAP_USER_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_USER_SCHEMA_DISPLAYNAME"`
	Username        string `yaml:"user_name" env:"LDAP_USER_SCHEMA_USERNAME;AUTH_BASIC_LDAP_USER_SCHEMA_USERNAME"`
}

type LDAPGroupSchema struct {
	ID              string `yaml:"id" env:"LDAP_GROUP_SCHEMA_ID;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID"`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING"`
	Mail            string `yaml:"mail" env:"LDAP_GROUP_SCHEMA_MAIL;AUTH_BASIC_LDAP_GROUP_SCHEMA_MAIL"`
	DisplayName     string `yaml:"display_name" env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_DISPLAYNAME"`
	Groupname       string `yaml:"group_name" env:"LDAP_GROUP_SCHEMA_GROUPNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_GROUPNAME"`
	Member          string `yaml:"member" env:"LDAP_GROUP_SCHEMA_MEMBER;AUTH_BASIC_LDAP_GROUP_SCHEMA_MEMBER"`
}

type OwnCloudSQLProvider struct {
	DBUsername       string `yaml:"db_username" env:"AUTH_BASIC_OWNCLOUDSQL_DB_USERNAME"`
	DBPassword       string `yaml:"db_password" env:"AUTH_BASIC_OWNCLOUDSQL_DB_PASSWORD"`
	DBHost           string `yaml:"db_host" env:"AUTH_BASIC_OWNCLOUDSQL_DB_HOST"`
	DBPort           int    `yaml:"db_port" env:"AUTH_BASIC_OWNCLOUDSQL_DB_PORT"`
	DBName           string `yaml:"db_name" env:"AUTH_BASIC_OWNCLOUDSQL_DB_NAME"`
	IDP              string `yaml:"idp" env:"AUTH_BASIC_OWNCLOUDSQL_IDP"`
	Nobody           int64  `yaml:"nobody" env:"AUTH_BASIC_OWNCLOUDSQL_NOBODY"` // TODO what is this?
	JoinUsername     bool   `yaml:"join_username" env:"AUTH_BASIC_OWNCLOUDSQL_JOIN_USERNAME"`
	JoinOwnCloudUUID bool   `yaml:"join_owncloud_uuid" env:"AUTH_BASIC_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID"`
}
