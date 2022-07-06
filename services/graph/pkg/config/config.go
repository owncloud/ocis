package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Reva         *Reva         `yaml:"reva"`
	TokenManager *TokenManager `yaml:"token_manager"`

	Spaces   Spaces   `yaml:"spaces"`
	Identity Identity `yaml:"identity"`
	Events   Events   `yaml:"events"`

	Context context.Context `yaml:"-"`
}

type Spaces struct {
	WebDavBase                      string `yaml:"webdav_base" env:"OCIS_URL;GRAPH_SPACES_WEBDAV_BASE" desc:"The public facing URL of WebDAV."`
	WebDavPath                      string `yaml:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH" desc:"The WebDAV subpath for spaces."`
	DefaultQuota                    string `yaml:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA" desc:"The default quota in bytes."`
	ExtendedSpacePropertiesCacheTTL int    `yaml:"extended_space_properties_cache_ttl" env:"GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL" desc:"Max TTL in seconds for the spaces property cache."`
}

type LDAP struct {
	URI                string `yaml:"uri" env:"LDAP_URI;GRAPH_LDAP_URI" desc:"URI of the LDAP Server to connect to. Supported URI schemes are 'ldaps://' and 'ldap://'"`
	CACert             string `yaml:"cacert" env:"LDAP_CACERT;GRAPH_LDAP_CACERT" desc:"The certificate to verify TLS connections."`
	Insecure           bool   `yaml:"insecure" env:"LDAP_INSECURE;GRAPH_LDAP_INSECURE" desc:"Disable TLS certificate validation for the LDAP connections. Do not set this in production environments."`
	BindDN             string `yaml:"bind_dn" env:"LDAP_BIND_DN;GRAPH_LDAP_BIND_DN" desc:"LDAP DN to use for simple bind authentication with the target LDAP server."`
	BindPassword       string `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;GRAPH_LDAP_BIND_PASSWORD" desc:"Password to use for authenticating the 'bind_dn'."`
	UseServerUUID      bool   `yaml:"use_server_uuid" env:"GRAPH_LDAP_SERVER_UUID" desc:"If set to true, rely on the LDAP Server to generate a unique ID for users and groups, like when using 'entryUUID' as the user ID attribute."`
	UsePasswordModExOp bool   `yaml:"use_password_modify_exop" env:"GRAPH_LDAP_SERVER_USE_PASSWORD_MODIFY_EXOP" desc:"User the Password Modify Extended Operation for updating user passwords."`
	WriteEnabled       bool   `yaml:"write_enabled" env:"GRAPH_LDAP_SERVER_WRITE_ENABLED" desc:"Allow to create, modify and delete LDAP users via GRAPH API. This is only works when the default Schema is used."`

	UserBaseDN               string `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;GRAPH_LDAP_USER_BASE_DN" desc:"Search base DN for looking up LDAP users."`
	UserSearchScope          string `yaml:"user_search_scope" env:"LDAP_USER_SCOPE;GRAPH_LDAP_USER_SCOPE" desc:"LDAP search scope to use when looking up users. Supported scopes are 'base', 'one' and 'sub'."`
	UserFilter               string `yaml:"user_filter" env:"LDAP_USER_FILTER;GRAPH_LDAP_USER_FILTER" desc:"LDAP filter to add to the default filters for user search like '(objectclass=ownCloud)'."`
	UserObjectClass          string `yaml:"user_objectclass" env:"LDAP_USER_OBJECTCLASS;GRAPH_LDAP_USER_OBJECTCLASS" desc:"The object class to use for users in the default user search filter ('inetOrgPerson')."`
	UserEmailAttribute       string `yaml:"user_mail_attribute" env:"LDAP_USER_SCHEMA_MAIL;GRAPH_LDAP_USER_EMAIL_ATTRIBUTE" desc:"LDAP Attribute to use for the email address of users."`
	UserDisplayNameAttribute string `yaml:"user_displayname_attribute" env:"LDAP_USER_SCHEMA_DISPLAY_NAME;GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE" desc:"LDAP Attribute to use for the displayname of users."`
	UserNameAttribute        string `yaml:"user_name_attribute" env:"LDAP_USER_SCHEMA_USERNAME;GRAPH_LDAP_USER_NAME_ATTRIBUTE" desc:"LDAP Attribute to use for username of users."`
	UserIDAttribute          string `yaml:"user_id_attribute" env:"LDAP_USER_SCHEMA_ID;GRAPH_LDAP_USER_UID_ATTRIBUTE" desc:"LDAP Attribute to use as the unique ID for users. This should be a stable globally unique ID like a UUID."`

	GroupBaseDN        string `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;GRAPH_LDAP_GROUP_BASE_DN" desc:"Search base DN for looking up LDAP groups."`
	GroupSearchScope   string `yaml:"group_search_scope" env:"LDAP_GROUP_SCOPE;GRAPH_LDAP_GROUP_SEARCH_SCOPE" desc:"LDAP search scope to use when looking up groups. Supported scopes are 'base', 'one' and 'sub'."`
	GroupFilter        string `yaml:"group_filter" env:"LDAP_GROUP_FILTER;GRAPH_LDAP_GROUP_FILTER" desc:"LDAP filter to add to the default filters for group searches."`
	GroupObjectClass   string `yaml:"group_objectclass" env:"LDAP_GROUP_OBJECTCLASS;GRAPH_LDAP_GROUP_OBJECTCLASS" desc:"The object class to use for groups in the default group search filter ('groupOfNames'). "`
	GroupNameAttribute string `yaml:"group_name_attribute" env:"LDAP_GROUP_SCHEMA_GROUPNAME;GRAPH_LDAP_GROUP_NAME_ATTRIBUTE" desc:"LDAP Attribute to use for the name of groups."`
	GroupIDAttribute   string `yaml:"group_id_attribute" env:"LDAP_GROUP_SCHEMA_ID;GRAPH_LDAP_GROUP_ID_ATTRIBUTE" desc:"LDAP Attribute to use as the unique id for groups. This should be a stable globally unique ID like a UUID."`
}

type Identity struct {
	Backend string `yaml:"backend" env:"GRAPH_IDENTITY_BACKEND" desc:"The user identity backend to use. Supported backend types are 'ldap' and 'cs3'."`
	LDAP    LDAP   `yaml:"ldap"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"GRAPH_EVENTS_ENDPOINT" desc:"The address of the streaming service."`
	Cluster  string `yaml:"cluster" env:"GRAPH_EVENTS_CLUSTER" desc:"The clusterID of the streaming service. Mandatory when using the NATS service."`
}
