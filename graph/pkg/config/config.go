package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	Reva         Reva         `ocisConfig:"reva"`
	TokenManager TokenManager `ocisConfig:"token_manager"`

	Spaces   Spaces   `ocisConfig:"spaces"`
	Identity Identity `ocisConfig:"identity"`

	Context context.Context
}

type Spaces struct {
	WebDavBase   string `ocisConfig:"webdav_base" env:"OCIS_URL;GRAPH_SPACES_WEBDAV_BASE"`
	WebDavPath   string `ocisConfig:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH"`
	DefaultQuota string `ocisConfig:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA"`
}

type LDAP struct {
	URI          string `ocisConfig:"uri" env:"GRAPH_LDAP_URI"`
	BindDN       string `ocisConfig:"bind_dn" env:"GRAPH_LDAP_BIND_DN"`
	BindPassword string `ocisConfig:"bind_password" env:"GRAPH_LDAP_BIND_PASSWORD"`

	UserBaseDN               string `ocisConfig:"user_base_dn" env:"GRAPH_LDAP_USER_BASE_DN"`
	UserSearchScope          string `ocisConfig:"user_search_scope" env:"GRAPH_LDAP_USER_SCOPE"`
	UserFilter               string `ocisConfig:"user_filter" env:"GRAPH_LDAP_USER_FILTER"`
	UserEmailAttribute       string `ocisConfig:"user_mail_attribute" env:"GRAPH_LDAP_USER_EMAIL_ATTRIBUTE"`
	UserDisplayNameAttribute string `ocisConfig:"user_displayname_attribute" env:"GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE"`
	UserNameAttribute        string `ocisConfig:"user_name_attribute" env:"GRAPH_LDAP_USER_NAME_ATTRIBUTE"`
	UserIDAttribute          string `ocisConfig:"user_id_attribute" env:"GRAPH_LDAP_USER_UID_ATTRIBUTE"`

	GroupBaseDN        string `ocisConfig:"group_base_dn" env:"GRAPH_LDAP_GROUP_BASE_DN"`
	GroupSearchScope   string `ocisConfig:"group_search_scope" env:"GRAPH_LDAP_GROUP_SEARCH_SCOPE"`
	GroupFilter        string `ocisConfig:"group_filter" env:"GRAPH_LDAP_GROUP_FILTER"`
	GroupNameAttribute string `ocisConfig:"group_name_attribute" env:"GRAPH_LDAP_GROUP_NAME_ATTRIBUTE"`
	GroupIDAttribute   string `ocisConfig:"group_id_attribute" env:"GRAPH_LDAP_GROUP_ID_ATTRIBUTE"`
}

type Identity struct {
	Backend string `ocisConfig:"backend" env:"GRAPH_IDENTITY_BACKEND"`
	LDAP    LDAP   `ocisConfig:"ldap"`
}
