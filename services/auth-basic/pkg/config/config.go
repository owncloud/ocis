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

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *Reva         `yaml:"reva"`

	SkipUserGroupsInToken bool          `yaml:"skip_user_groups_in_token" env:"AUTH_BASIC_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the encoding of the user's group memberships in the reva access token. This reduces the token size, especially when users are members of a large number of groups."`
	AuthProvider          string        `yaml:"auth_provider" env:"AUTH_BASIC_AUTH_PROVIDER" desc:"The auth provider which should be used by the service like 'ldap'."`
	AuthProviders         AuthProviders `yaml:"auth_providers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_BASIC_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_BASIC_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_BASIC_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_BASIC_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_BASIC_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_BASIC_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_BASIC_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_BASIC_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_BASIC_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"AUTH_BASIC_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"AUTH_BASIC_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"AUTH_BASIC_DEBUG_ZPAGES" desc:"Enables zpages, which can  be used for collecting and viewing traces in-memory."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"AUTH_BASIC_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"AUTH_BASIC_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service."`
}

type AuthProviders struct {
	LDAP        LDAPProvider        `yaml:"ldap"`
	OwnCloudSQL OwnCloudSQLProvider `yaml:"owncloudsql"`
	JSON        JSONProvider        `yaml:"json,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type JSONProvider struct {
	File string `yaml:"file,omitempty"`
}

type LDAPProvider struct {
	URI              string          `yaml:"uri" env:"LDAP_URI;AUTH_BASIC_LDAP_URI" desc:"URI of the LDAP Server to connect to. Supported URI schemes are 'ldaps://' and 'ldap://'"`
	CACert           string          `yaml:"ca_cert" env:"LDAP_CACERT;AUTH_BASIC_LDAP_CACERT" desc:"Path to a CA certificate file for validating the LDAP server's TLS certificate. If empty the system default CA bundle will be used."`
	Insecure         bool            `yaml:"insecure" env:"LDAP_INSECURE;AUTH_BASIC_LDAP_INSECURE" desc:"Disable TLS certificate validation for the LDAP connections. Do not set this in production environments."`
	BindDN           string          `yaml:"bind_dn" env:"LDAP_BIND_DN;AUTH_BASIC_LDAP_BIND_DN" desc:"LDAP DN to use for simple bind authentication with the target LDAP server."`
	BindPassword     string          `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;AUTH_BASIC_LDAP_BIND_PASSWORD" desc:"Password to use for authenticating the 'bind_dn'."`
	UserBaseDN       string          `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;AUTH_BASIC_LDAP_USER_BASE_DN" desc:"Search base DN for looking up LDAP users."`
	GroupBaseDN      string          `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;AUTH_BASIC_LDAP_GROUP_BASE_DN" desc:"Search base DN for looking up LDAP groups."`
	UserScope        string          `yaml:"user_scope" env:"LDAP_USER_SCOPE;AUTH_BASIC_LDAP_USER_SCOPE" desc:"LDAP search scope to use when looking up users. Supported values are 'base', 'one' and 'sub'."`
	GroupScope       string          `yaml:"group_scope" env:"LDAP_GROUP_SCOPE;AUTH_BASIC_LDAP_GROUP_SCOPE" desc:"LDAP search scope to use when looking up groups. Supported values are 'base', 'one' and 'sub'."`
	UserFilter       string          `yaml:"user_filter" env:"LDAP_USER_FILTER;AUTH_BASIC_LDAP_USER_FILTER" desc:"LDAP filter to add to the default filters for user search like '(objectclass=ownCloud)'."`
	GroupFilter      string          `yaml:"group_filter" env:"LDAP_GROUP_FILTER;AUTH_BASIC_LDAP_GROUP_FILTER" desc:"LDAP filter to add to the default filters for group searches."`
	UserObjectClass  string          `yaml:"user_object_class" env:"LDAP_USER_OBJECTCLASS;AUTH_BASIC_LDAP_USER_OBJECTCLASS" desc:"The object class to use for users in the default user search filter ('inetOrgPerson')."`
	GroupObjectClass string          `yaml:"group_object_class" env:"LDAP_GROUP_OBJECTCLASS;AUTH_BASIC_LDAP_GROUP_OBJECTCLASS" desc:"The object class to use for groups in the default group search filter ('groupOfNames'). "`
	LoginAttributes  []string        `yaml:"login_attributes" env:"LDAP_LOGIN_ATTRIBUTES;AUTH_BASIC_LDAP_LOGIN_ATTRIBUTES" desc:"The user object attributes, that can be used for login."`
	IDP              string          `yaml:"idp" env:"OCIS_URL;OCIS_OIDC_ISSUER;AUTH_BASIC_IDP_URL" desc:"The identity provider value to set in the userids of the CS3 user objects for users returned by this user provider."`
	UserSchema       LDAPUserSchema  `yaml:"user_schema"`
	GroupSchema      LDAPGroupSchema `yaml:"group_schema"`
}

type LDAPUserSchema struct {
	ID              string `yaml:"id" env:"LDAP_USER_SCHEMA_ID;AUTH_BASIC_LDAP_USER_SCHEMA_ID" desc:"LDAP Attribute to use as the unique id for users. This should be a stable globally unique ID like a UUID."`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'id' attribute for users is of the 'OCTETSTRING' syntax. This is e.g. required when using the 'objectGUID' attribute of Active Directory for the user IDs."`
	Mail            string `yaml:"mail" env:"LDAP_USER_SCHEMA_MAIL;AUTH_BASIC_LDAP_USER_SCHEMA_MAIL" desc:"LDAP Attribute to use for the email address of users."`
	DisplayName     string `yaml:"display_name" env:"LDAP_USER_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_USER_SCHEMA_DISPLAYNAME" desc:"LDAP Attribute to use for the displayname of users."`
	Username        string `yaml:"user_name" env:"LDAP_USER_SCHEMA_USERNAME;AUTH_BASIC_LDAP_USER_SCHEMA_USERNAME" desc:"LDAP Attribute to use for username of users."`
}

type LDAPGroupSchema struct {
	ID              string `yaml:"id" env:"LDAP_GROUP_SCHEMA_ID;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID" desc:"LDAP Attribute to use as the unique id for groups. This should be a stable globally unique id (e.g. a UUID)."`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;AUTH_BASIC_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'id' attribute for groups is of the 'OCTETSTRING' syntax. This is e.g. required when using the 'objectGUID' attribute of Active Directory for the group IDs."`
	Mail            string `yaml:"mail" env:"LDAP_GROUP_SCHEMA_MAIL;AUTH_BASIC_LDAP_GROUP_SCHEMA_MAIL" desc:"LDAP Attribute to use for the email address of groups (can be empty)."`
	DisplayName     string `yaml:"display_name" env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_DISPLAYNAME" desc:"LDAP Attribute to use for the displayname of groups (often the same as groupname attribute)."`
	Groupname       string `yaml:"group_name" env:"LDAP_GROUP_SCHEMA_GROUPNAME;AUTH_BASIC_LDAP_GROUP_SCHEMA_GROUPNAME" desc:"LDAP Attribute to use for the name of groups."`
	Member          string `yaml:"member" env:"LDAP_GROUP_SCHEMA_MEMBER;AUTH_BASIC_LDAP_GROUP_SCHEMA_MEMBER" desc:"LDAP Attribute that is used for group members."`
}

type OwnCloudSQLProvider struct {
	DBUsername       string `yaml:"db_username" env:"AUTH_BASIC_OWNCLOUDSQL_DB_USERNAME" desc:"Database user to use for authenticating with the owncloud database."`
	DBPassword       string `yaml:"db_password" env:"AUTH_BASIC_OWNCLOUDSQL_DB_PASSWORD" desc:"Password for the database user."`
	DBHost           string `yaml:"db_host" env:"AUTH_BASIC_OWNCLOUDSQL_DB_HOST" desc:"Hostname of the database server."`
	DBPort           int    `yaml:"db_port" env:"AUTH_BASIC_OWNCLOUDSQL_DB_PORT" desc:"Network port to use for the database connection."`
	DBName           string `yaml:"db_name" env:"AUTH_BASIC_OWNCLOUDSQL_DB_NAME" desc:"Name of the owncloud database."`
	IDP              string `yaml:"idp" env:"AUTH_BASIC_OWNCLOUDSQL_IDP" desc:"The identity provider value to set in the userids of the CS3 user objects for users returned by this user provider."`
	Nobody           int64  `yaml:"nobody" env:"AUTH_BASIC_OWNCLOUDSQL_NOBODY" desc:"Fallback number if no numeric UID and GID properties are provided."`
	JoinUsername     bool   `yaml:"join_username" env:"AUTH_BASIC_OWNCLOUDSQL_JOIN_USERNAME" desc:"Join the user properties table to read usernames"`
	JoinOwnCloudUUID bool   `yaml:"join_owncloud_uuid" env:"AUTH_BASIC_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID" desc:"Join the user properties table to read user ID's."`
}
