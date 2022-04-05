package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Reva         Reva         `yaml:"reva"`
	TokenManager TokenManager `yaml:"token_manager"`

	Spaces   Spaces   `yaml:"spaces"`
	Identity Identity `yaml:"identity"`

	Context context.Context `yaml:"-"`
}

type Spaces struct {
	WebDavBase                      string `yaml:"webdav_base" env:"OCIS_URL;GRAPH_SPACES_WEBDAV_BASE"`
	WebDavPath                      string `yaml:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH"`
	DefaultQuota                    string `yaml:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA"`
	Insecure                        bool   `yaml:"insecure" env:"OCIS_INSECURE;GRAPH_SPACES_INSECURE"`
	ExtendedSpacePropertiesCacheTTL int    `yaml:"extended_space_properties_cache_ttl" env:"GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL"`
}

type LDAP struct {
	URI           string `yaml:"uri" env:"GRAPH_LDAP_URI"`
	Insecure      bool   `yaml:"insecure" env:"OCIS_INSECURE;GRAPH_LDAP_INSECURE"`
	BindDN        string `yaml:"bind_dn" env:"GRAPH_LDAP_BIND_DN"`
	BindPassword  string `yaml:"bind_password" env:"GRAPH_LDAP_BIND_PASSWORD"`
	UseServerUUID bool   `yaml:"use_server_uuid" env:"GRAPH_LDAP_SERVER_UUID"`
	WriteEnabled  bool   `yaml:"write_enabled" env:"GRAPH_LDAP_SERVER_WRITE_ENABLED"`

	UserBaseDN               string `yaml:"user_base_dn" env:"GRAPH_LDAP_USER_BASE_DN"`
	UserSearchScope          string `yaml:"user_search_scope" env:"GRAPH_LDAP_USER_SCOPE"`
	UserFilter               string `yaml:"user_filter" env:"GRAPH_LDAP_USER_FILTER"`
	UserEmailAttribute       string `yaml:"user_mail_attribute" env:"GRAPH_LDAP_USER_EMAIL_ATTRIBUTE"`
	UserDisplayNameAttribute string `yaml:"user_displayname_attribute" env:"GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE"`
	UserNameAttribute        string `yaml:"user_name_attribute" env:"GRAPH_LDAP_USER_NAME_ATTRIBUTE"`
	UserIDAttribute          string `yaml:"user_id_attribute" env:"GRAPH_LDAP_USER_UID_ATTRIBUTE"`

	GroupBaseDN        string `yaml:"group_base_dn" env:"GRAPH_LDAP_GROUP_BASE_DN"`
	GroupSearchScope   string `yaml:"group_search_scope" env:"GRAPH_LDAP_GROUP_SEARCH_SCOPE"`
	GroupFilter        string `yaml:"group_filter" env:"GRAPH_LDAP_GROUP_FILTER"`
	GroupNameAttribute string `yaml:"group_name_attribute" env:"GRAPH_LDAP_GROUP_NAME_ATTRIBUTE"`
	GroupIDAttribute   string `yaml:"group_id_attribute" env:"GRAPH_LDAP_GROUP_ID_ATTRIBUTE"`
}

type Identity struct {
	Backend string `yaml:"backend" env:"GRAPH_IDENTITY_BACKEND"`
	LDAP    LDAP   `yaml:"ldap"`
}
