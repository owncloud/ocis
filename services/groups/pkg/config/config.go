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

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"GROUPS_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the loading of user's group memberships from the reva access token."`

	Driver  string  `yaml:"driver" desc:"The driver which should be used by the groups service like 'ldap'."`
	Drivers Drivers `yaml:"drivers"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;GROUPS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;GROUPS_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;GROUPS_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;GROUPS_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;GROUPS_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;GROUPS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;GROUPS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;GROUPS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"GROUPS_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"GROUPS_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"GROUPS_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"GROUPS_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string `yaml:"addr" env:"GROUPS_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"GROUPS_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service."`
}

type Drivers struct {
	LDAP        LDAPDriver        `yaml:"ldap"`
	OwnCloudSQL OwnCloudSQLDriver `yaml:"owncloudsql"`

	JSON JSONDriver   `yaml:"json,omitempty"` // not supported by the oCIS product, therefore not part of docs
	REST RESTProvider `yaml:"rest,omitempty"` // not supported by the oCIS product, therefore not part of docs
}

type LDAPDriver struct {
	URI              string          `yaml:"uri" env:"LDAP_URI;GROUPS_LDAP_URI" desc:"URI of the LDAP Server to connect to. Supported URI schemes are 'ldaps://' and 'ldap://'"`
	CACert           string          `yaml:"ca_cert" env:"LDAP_CACERT;GROUPS_LDAP_CACERT" desc:"Path to a CA certificate file for validating the LDAP server's TLS certificate. If empty the system default CA bundle will be used."`
	Insecure         bool            `yaml:"insecure" env:"LDAP_INSECURE;GROUPS_LDAP_INSECURE" desc:"Disable TLS certificate validation for the LDAP connections. Do not set this in production environments."`
	BindDN           string          `yaml:"bind_dn" env:"LDAP_BIND_DN;GROUPS_LDAP_BIND_DN" desc:"LDAP DN to use for simple bind authentication with the target LDAP server."`
	BindPassword     string          `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;GROUPS_LDAP_BIND_PASSWORD" desc:"Password to use for authenticating the 'bind_dn'."`
	UserBaseDN       string          `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;GROUPS_LDAP_USER_BASE_DN" desc:"Search base DN for looking up LDAP users."`
	GroupBaseDN      string          `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;GROUPS_LDAP_GROUP_BASE_DN" desc:"Search base DN for looking up LDAP groups."`
	UserScope        string          `yaml:"user_scope" env:"LDAP_USER_SCOPE;GROUPS_LDAP_USER_SCOPE" desc:"LDAP search scope to use when looking up users. Supported scopes are 'base', 'one' and 'sub'."`
	GroupScope       string          `yaml:"group_scope" env:"LDAP_GROUP_SCOPE;GROUPS_LDAP_GROUP_SCOPE" desc:"LDAP search scope to use when looking up groups. Supported scopes are 'base', 'one' and 'sub'."`
	UserFilter       string          `yaml:"user_filter" env:"LDAP_USER_FILTER;GROUPS_LDAP_USER_FILTER" desc:"LDAP filter to add to the default filters for user search like '(objectclass=ownCloud)'."`
	GroupFilter      string          `yaml:"group_filter" env:"LDAP_GROUP_FILTER;GROUPS_LDAP_GROUP_FILTER" desc:"LDAP filter to add to the default filters for group searches."`
	UserObjectClass  string          `yaml:"user_object_class" env:"LDAP_USER_OBJECTCLASS;GROUPS_LDAP_USER_OBJECTCLASS" desc:"The object class to use for users in the default user search filter ('inetOrgPerson')."`
	GroupObjectClass string          `yaml:"group_object_class" env:"LDAP_GROUP_OBJECTCLASS;GROUPS_LDAP_GROUP_OBJECTCLASS" desc:"The object class to use for groups in the default group search filter ('groupOfNames')."`
	IDP              string          `yaml:"idp" env:"OCIS_URL;OCIS_OIDC_ISSUER;GROUPS_IDP_URL" desc:"The identity provider value to set in the group IDs of the CS3 group objects for groups returned by this group provider."`
	UserSchema       LDAPUserSchema  `yaml:"user_schema"`
	GroupSchema      LDAPGroupSchema `yaml:"group_schema"`
}

type LDAPUserSchema struct {
	ID              string `yaml:"id" env:"LDAP_USER_SCHEMA_ID;GROUPS_LDAP_USER_SCHEMA_ID" desc:"LDAP Attribute to use as the unique id for users. This should be a stable globally unique id like a UUID."`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'id' attribute for users is of the 'OCTETSTRING' syntax. This is e.g. required when using the 'objectGUID' attribute of Active Directory for the user id's."`
	Mail            string `yaml:"mail" env:"LDAP_USER_SCHEMA_MAIL;GROUPS_LDAP_USER_SCHEMA_MAIL" desc:"LDAP Attribute to use for the email address of users."`
	DisplayName     string `yaml:"display_name" env:"LDAP_USER_SCHEMA_DISPLAYNAME;GROUPS_LDAP_USER_SCHEMA_DISPLAYNAME" desc:"LDAP Attribute to use for the displayname of users."`
	Username        string `yaml:"user_name" env:"LDAP_USER_SCHEMA_USERNAME;GROUPS_LDAP_USER_SCHEMA_USERNAME" desc:"LDAP Attribute to use for username of users."`
}

type LDAPGroupSchema struct {
	ID              string `yaml:"id" env:"LDAP_GROUP_SCHEMA_ID;GROUPS_LDAP_GROUP_SCHEMA_ID" desc:"LDAP Attribute to use as the unique id for groups. This should be a stable globally unique ID like a UUID."`
	IDIsOctetString bool   `yaml:"id_is_octet_string" env:"LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;GROUPS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'id' attribute for groups is of the 'OCTETSTRING' syntax. This is e.g. required when using the 'objectGUID' attribute of Active Directory for the group ID's."`
	Mail            string `yaml:"mail" env:"LDAP_GROUP_SCHEMA_MAIL;GROUPS_LDAP_GROUP_SCHEMA_MAIL" desc:"LDAP Attribute to use for the email address of groups (can be empty)."`
	DisplayName     string `yaml:"display_name" env:"LDAP_GROUP_SCHEMA_DISPLAYNAME;GROUPS_LDAP_GROUP_SCHEMA_DISPLAYNAME" desc:"LDAP Attribute to use for the displayname of groups (often the same as groupname attribute)."`
	Groupname       string `yaml:"group_name" env:"LDAP_GROUP_SCHEMA_GROUPNAME;GROUPS_LDAP_GROUP_SCHEMA_GROUPNAME" desc:"LDAP Attribute to use for the name of groups."`
	Member          string `yaml:"member" env:"LDAP_GROUP_SCHEMA_MEMBER;GROUPS_LDAP_GROUP_SCHEMA_MEMBER" desc:"LDAP Attribute that is used for group members."`
}

type OwnCloudSQLDriver struct {
	DBUsername         string `yaml:"db_username" env:"GROUPS_OWNCLOUDSQL_DB_USERNAME" desc:"Database user to use for authenticating with the owncloud database."`
	DBPassword         string `yaml:"db_password" env:"GROUPS_OWNCLOUDSQL_DB_PASSWORD" desc:"Password for the database user."`
	DBHost             string `yaml:"db_host" env:"GROUPS_OWNCLOUDSQL_DB_HOST" desc:"Hostname of the database server."`
	DBPort             int    `yaml:"db_port" env:"GROUPS_OWNCLOUDSQL_DB_PORT" desc:"Network port to use for the database connection."`
	DBName             string `yaml:"db_name" env:"GROUPS_OWNCLOUDSQL_DB_NAME" desc:"Name of the owncloud database."`
	IDP                string `yaml:"idp" env:"GROUPS_OWNCLOUDSQL_IDP" desc:"The identity provider value to set in the userids of the CS3 user objects for users returned by this user provider."`
	Nobody             int64  `yaml:"nobody" env:"GROUPS_OWNCLOUDSQL_NOBODY" desc:"Fallback number if no numeric UID and GID properties are provided."`
	JoinUsername       bool   `yaml:"join_username" env:"GROUPS_OWNCLOUDSQL_JOIN_USERNAME" desc:"Join the user properties table to read usernames."`
	JoinOwnCloudUUID   bool   `yaml:"join_owncloud_uuid" env:"GROUPS_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID" desc:"Join the user properties table to read user IDs."`
	EnableMedialSearch bool   `yaml:"enable_medial_search" env:"GROUPS_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH" desc:"Allow 'medial search' when searching for users instead of just doing a prefix search. This allows finding 'Alice' when searching for 'lic'."`
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
