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
	WebDavBase                      string `yaml:"webdav_base" env:"OCIS_URL;GRAPH_SPACES_WEBDAV_BASE" desc:"URL, where oCIS is reachable for users."`
	WebDavPath                      string `yaml:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH"`
	DefaultQuota                    string `yaml:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA"`
	Insecure                        bool   `yaml:"insecure" env:"OCIS_INSECURE;GRAPH_SPACES_INSECURE"`
	ExtendedSpacePropertiesCacheTTL int    `yaml:"extended_space_properties_cache_ttl" env:"GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL"`
}

type LDAP struct {
	URI           string `yaml:"uri" env:"LDAP_URI;GRAPH_LDAP_URI"`
	CACert        string `yaml:"cacert" env:"LDAP_CACERT;GRAPH_LDAP_CACERT" desc:"The certificate to verify TLS connections"`
	Insecure      bool   `yaml:"insecure" env:"LDAP_INSECURE;GRAPH_LDAP_INSECURE"`
	BindDN        string `yaml:"bind_dn" env:"LDAP_BIND_DN;GRAPH_LDAP_BIND_DN"`
	BindPassword  string `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;GRAPH_LDAP_BIND_PASSWORD"`
	UseServerUUID bool   `yaml:"use_server_uuid" env:"GRAPH_LDAP_SERVER_UUID"`
	WriteEnabled  bool   `yaml:"write_enabled" env:"GRAPH_LDAP_SERVER_WRITE_ENABLED"`

	UserBaseDN               string `yaml:"user_base_dn" env:"LDAP_USER_BASE_DN;GRAPH_LDAP_USER_BASE_DN"`
	UserSearchScope          string `yaml:"user_search_scope" env:"LDAP_USER_SCOPE;GRAPH_LDAP_USER_SCOPE"`
	UserFilter               string `yaml:"user_filter" env:"LDAP_USER_FILTER;GRAPH_LDAP_USER_FILTER"`
	UserObjectClass          string `yaml:"user_objectclass" env:"LDAP_USER_OBJECTCLASS;GRAPH_LDAP_USER_OBJECTCLASS"`
	UserEmailAttribute       string `yaml:"user_mail_attribute" env:"LDAP_USER_SCHEMA_MAIL;GRAPH_LDAP_USER_EMAIL_ATTRIBUTE"`
	UserDisplayNameAttribute string `yaml:"user_displayname_attribute" env:"LDAP_USER_SCHEMA_DISPLAY_NAME;GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE"`
	UserNameAttribute        string `yaml:"user_name_attribute" env:"LDAP_USER_SCHEMA_USERNAME;GRAPH_LDAP_USER_NAME_ATTRIBUTE"`
	UserIDAttribute          string `yaml:"user_id_attribute" env:"LDAP_USER_SCHEMA_ID;GRAPH_LDAP_USER_UID_ATTRIBUTE"`

	GroupBaseDN        string `yaml:"group_base_dn" env:"LDAP_GROUP_BASE_DN;GRAPH_LDAP_GROUP_BASE_DN"`
	GroupSearchScope   string `yaml:"group_search_scope" env:"LDAP_GROUP_SCOPE;GRAPH_LDAP_GROUP_SEARCH_SCOPE"`
	GroupFilter        string `yaml:"group_filter" env:"LDAP_GROUP_FILTER;GRAPH_LDAP_GROUP_FILTER"`
	GroupObjectClass   string `yaml:"group_objectclass" env:"LDAP_GROUP_OBJECTCLASS;GRAPH_LDAP_GROUP_OBJECTCLASS"`
	GroupNameAttribute string `yaml:"group_name_attribute" env:"LDAP_GROUP_SCHEMA_GROUPNAME;GRAPH_LDAP_GROUP_NAME_ATTRIBUTE"`
	GroupIDAttribute   string `yaml:"group_id_attribute" env:"LDAP_GROUP_SCHEMA_ID;GRAPH_LDAP_GROUP_ID_ATTRIBUTE"`
}

type Identity struct {
	Backend string `yaml:"backend" env:"GRAPH_IDENTITY_BACKEND" desc:"The user identity backend to use, defaults to 'ldap', can be 'cs3'."`
	LDAP    LDAP   `yaml:"ldap"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"GRAPH_EVENTS_ENDPOINT" desc:"the address of the streaming service"`
	Cluster  string `yaml:"cluster" env:"GRAPH_EVENTS_CLUSTER" desc:"the clusterID of the streaming service. Mandatory when using nats"`
}
