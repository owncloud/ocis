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
	Cache   *Cache   `yaml:"cache"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	API API `yaml:"api"`

	Reva          *shared.Reva          `yaml:"reva"`
	TokenManager  *TokenManager         `yaml:"token_manager"`
	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`

	Application Application `yaml:"application"`
	Spaces      Spaces      `yaml:"spaces"`
	Identity    Identity    `yaml:"identity"`
	Events      Events      `yaml:"events"`

	Keycloak       Keycloak       `yaml:"keycloak"`
	ServiceAccount ServiceAccount `yaml:"service_account"`

	FilesSharing FilesSharing `yaml:"files_sharing"`

	Context context.Context `yaml:"-"`
}

type Spaces struct {
	WebDavBase                      string `yaml:"webdav_base" env:"OCIS_URL;GRAPH_SPACES_WEBDAV_BASE" desc:"The public facing URL of WebDAV."`
	WebDavPath                      string `yaml:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH" desc:"The WebDAV subpath for spaces."`
	DefaultQuota                    string `yaml:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA" desc:"The default quota in bytes."`
	ExtendedSpacePropertiesCacheTTL int    `yaml:"extended_space_properties_cache_ttl" env:"GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL" desc:"Max TTL in seconds for the spaces property cache."`
	UsersCacheTTL                   int    `yaml:"users_cache_ttl" env:"GRAPH_SPACES_USERS_CACHE_TTL" desc:"Max TTL in seconds for the spaces users cache."`
	GroupsCacheTTL                  int    `yaml:"groups_cache_ttl" env:"GRAPH_SPACES_GROUPS_CACHE_TTL" desc:"Max TTL in seconds for the spaces groups cache."`
}

type LDAP struct {
	URI                string `yaml:"uri" env:"OCIS_LDAP_URI;GRAPH_LDAP_URI" desc:"URI of the LDAP Server to connect to. Supported URI schemes are 'ldaps://' and 'ldap://'"`
	CACert             string `yaml:"cacert" env:"OCIS_LDAP_CACERT;GRAPH_LDAP_CACERT" desc:"Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the LDAP service. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/idm."`
	Insecure           bool   `yaml:"insecure" env:"OCIS_LDAP_INSECURE;GRAPH_LDAP_INSECURE" desc:"Disable TLS certificate validation for the LDAP connections. Do not set this in production environments."`
	BindDN             string `yaml:"bind_dn" env:"OCIS_LDAP_BIND_DN;GRAPH_LDAP_BIND_DN" desc:"LDAP DN to use for simple bind authentication with the target LDAP server."`
	BindPassword       string `yaml:"bind_password" env:"OCIS_LDAP_BIND_PASSWORD;LDAP_BIND_PASSWORD;GRAPH_LDAP_BIND_PASSWORD" desc:"Password to use for authenticating the 'bind_dn'." deprecationVersion:"4.0.2" removalVersion:"5.0.0" deprecationInfo:"LDAP_BIND_PASSWORD changing name for consistency" deprecationReplacement:"OCIS_LDAP_BIND_PASSWORD"`
	UseServerUUID      bool   `yaml:"use_server_uuid" env:"GRAPH_LDAP_SERVER_UUID" desc:"If set to true, rely on the LDAP Server to generate a unique ID for users and groups, like when using 'entryUUID' as the user ID attribute."`
	UsePasswordModExOp bool   `yaml:"use_password_modify_exop" env:"GRAPH_LDAP_SERVER_USE_PASSWORD_MODIFY_EXOP" desc:"Use the 'Password Modify Extended Operation' for updating user passwords."`
	WriteEnabled       bool   `yaml:"write_enabled" env:"OCIS_LDAP_SERVER_WRITE_ENABLED;GRAPH_LDAP_SERVER_WRITE_ENABLED" desc:"Allow creating, modifying and deleting LDAP users via the GRAPH API. This can only be set to 'true' when keeping default settings for the LDAP user and group attribute types (the 'OCIS_LDAP_USER_SCHEMA_* and 'OCIS_LDAP_GROUP_SCHEMA_* variables)."`
	RefintEnabled      bool   `yaml:"refint_enabled" env:"GRAPH_LDAP_REFINT_ENABLED" desc:"Signals that the server has the refint plugin enabled, which makes some actions not needed."`

	UserBaseDN               string `yaml:"user_base_dn" env:"OCIS_LDAP_USER_BASE_DN;GRAPH_LDAP_USER_BASE_DN" desc:"Search base DN for looking up LDAP users."`
	UserSearchScope          string `yaml:"user_search_scope" env:"OCIS_LDAP_USER_SCOPE;GRAPH_LDAP_USER_SCOPE" desc:"LDAP search scope to use when looking up users. Supported scopes are 'base', 'one' and 'sub'."`
	UserFilter               string `yaml:"user_filter" env:"OCIS_LDAP_USER_FILTER;GRAPH_LDAP_USER_FILTER" desc:"LDAP filter to add to the default filters for user search like '(objectclass=ownCloud)'."`
	UserObjectClass          string `yaml:"user_objectclass" env:"OCIS_LDAP_USER_OBJECTCLASS;GRAPH_LDAP_USER_OBJECTCLASS" desc:"The object class to use for users in the default user search filter ('inetOrgPerson')."`
	UserEmailAttribute       string `yaml:"user_mail_attribute" env:"OCIS_LDAP_USER_SCHEMA_MAIL;GRAPH_LDAP_USER_EMAIL_ATTRIBUTE" desc:"LDAP Attribute to use for the email address of users."`
	UserDisplayNameAttribute string `yaml:"user_displayname_attribute" env:"LDAP_USER_SCHEMA_DISPLAY_NAME;GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE" desc:"LDAP Attribute to use for the displayname of users."`
	UserNameAttribute        string `yaml:"user_name_attribute" env:"OCIS_LDAP_USER_SCHEMA_USERNAME;GRAPH_LDAP_USER_NAME_ATTRIBUTE" desc:"LDAP Attribute to use for username of users."`
	UserIDAttribute          string `yaml:"user_id_attribute" env:"OCIS_LDAP_USER_SCHEMA_ID;GRAPH_LDAP_USER_UID_ATTRIBUTE" desc:"LDAP Attribute to use as the unique ID for users. This should be a stable globally unique ID like a UUID."`
	UserIDIsOctetString      bool   `yaml:"user_id_is_octet_string" env:"OCIS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING;GRAPH_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'ID' attribute for users is of the 'OCTETSTRING' syntax. This is required when using the 'objectGUID' attribute of Active Directory for the user ID's."`
	UserTypeAttribute        string `yaml:"user_type_attribute" env:"OCIS_LDAP_USER_SCHEMA_USER_TYPE;GRAPH_LDAP_USER_TYPE_ATTRIBUTE" desc:"LDAP Attribute to distinguish between 'Member' and 'Guest' users. Default is 'ownCloudUserType'."`
	UserEnabledAttribute     string `yaml:"user_enabled_attribute" env:"OCIS_LDAP_USER_ENABLED_ATTRIBUTE;GRAPH_USER_ENABLED_ATTRIBUTE" desc:"LDAP Attribute to use as a flag telling if the user is enabled or disabled."`
	DisableUserMechanism     string `yaml:"disable_user_mechanism" env:"OCIS_LDAP_DISABLE_USER_MECHANISM;GRAPH_DISABLE_USER_MECHANISM" desc:"An option to control the behavior for disabling users. Supported options are 'none', 'attribute' and 'group'. If set to 'group', disabling a user via API will add the user to the configured group for disabled users, if set to 'attribute' this will be done in the ldap user entry, if set to 'none' the disable request is not processed. Default is 'attribute'."`
	LdapDisabledUsersGroupDN string `yaml:"ldap_disabled_users_group_dn" env:"OCIS_LDAP_DISABLED_USERS_GROUP_DN;GRAPH_DISABLED_USERS_GROUP_DN" desc:"The distinguished name of the group to which added users will be classified as disabled when 'disable_user_mechanism' is set to 'group'."`

	GroupBaseDN          string `yaml:"group_base_dn" env:"OCIS_LDAP_GROUP_BASE_DN;GRAPH_LDAP_GROUP_BASE_DN" desc:"Search base DN for looking up LDAP groups."`
	GroupCreateBaseDN    string `yaml:"group_create_base_dn" env:"GRAPH_LDAP_GROUP_CREATE_BASE_DN" desc:"Parent DN under which new groups are created. This DN needs to be subordinate to the 'GRAPH_LDAP_GROUP_BASE_DN'. This setting is only relevant when 'GRAPH_LDAP_SERVER_WRITE_ENABLED' is 'true'. It defaults to the value of 'GRAPH_LDAP_GROUP_BASE_DN'. All groups outside of this subtree are treated as readonly groups and cannot be updated."`
	GroupSearchScope     string `yaml:"group_search_scope" env:"OCIS_LDAP_GROUP_SCOPE;GRAPH_LDAP_GROUP_SEARCH_SCOPE" desc:"LDAP search scope to use when looking up groups. Supported scopes are 'base', 'one' and 'sub'."`
	GroupFilter          string `yaml:"group_filter" env:"OCIS_LDAP_GROUP_FILTER;GRAPH_LDAP_GROUP_FILTER" desc:"LDAP filter to add to the default filters for group searches."`
	GroupObjectClass     string `yaml:"group_objectclass" env:"OCIS_LDAP_GROUP_OBJECTCLASS;GRAPH_LDAP_GROUP_OBJECTCLASS" desc:"The object class to use for groups in the default group search filter ('groupOfNames')."`
	GroupNameAttribute   string `yaml:"group_name_attribute" env:"OCIS_LDAP_GROUP_SCHEMA_GROUPNAME;GRAPH_LDAP_GROUP_NAME_ATTRIBUTE" desc:"LDAP Attribute to use for the name of groups."`
	GroupMemberAttribute string `yaml:"group_member_attribute" env:"OCIS_LDAP_GROUP_SCHEMA_MEMBER;GRAPH_LDAP_GROUP_MEMBER_ATTRIBUTE" desc:"LDAP Attribute that is used for group members."`
	GroupIDAttribute     string `yaml:"group_id_attribute" env:"OCIS_LDAP_GROUP_SCHEMA_ID;GRAPH_LDAP_GROUP_ID_ATTRIBUTE" desc:"LDAP Attribute to use as the unique id for groups. This should be a stable globally unique ID like a UUID."`
	GroupIDIsOctetString bool   `yaml:"group_id_is_octet_string" env:"OCIS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING;GRAPH_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING" desc:"Set this to true if the defined 'ID' attribute for groups is of the 'OCTETSTRING' syntax. This is required when using the 'objectGUID' attribute of Active Directory for the group ID's."`

	EducationResourcesEnabled bool `yaml:"education_resources_enabled" env:"GRAPH_LDAP_EDUCATION_RESOURCES_ENABLED" desc:"Enable LDAP support for managing education related resources."`
	EducationConfig           LDAPEducationConfig
}

// LDAPEducationConfig represents the LDAP configuration for education related resources
type LDAPEducationConfig struct {
	SchoolBaseDN      string `yaml:"school_base_dn" env:"GRAPH_LDAP_SCHOOL_BASE_DN" desc:"Search base DN for looking up LDAP schools."`
	SchoolSearchScope string `yaml:"school_search_scope" env:"GRAPH_LDAP_SCHOOL_SEARCH_SCOPE" desc:"LDAP search scope to use when looking up schools. Supported scopes are 'base', 'one' and 'sub'."`

	SchoolFilter      string `yaml:"school_filter" env:"GRAPH_LDAP_SCHOOL_FILTER" desc:"LDAP filter to add to the default filters for school searches."`
	SchoolObjectClass string `yaml:"school_objectclass" env:"GRAPH_LDAP_SCHOOL_OBJECTCLASS" desc:"The object class to use for schools in the default school search filter."`

	SchoolNameAttribute   string `yaml:"school_name_attribute" env:"GRAPH_LDAP_SCHOOL_NAME_ATTRIBUTE" desc:"LDAP Attribute to use for the name of a school."`
	SchoolNumberAttribute string `yaml:"school_number_attribute" env:"GRAPH_LDAP_SCHOOL_NUMBER_ATTRIBUTE" desc:"LDAP Attribute to use for the number of a school."`
	SchoolIDAttribute     string `yaml:"school_id_attribute" env:"GRAPH_LDAP_SCHOOL_ID_ATTRIBUTE" desc:"LDAP Attribute to use as the unique id for schools. This should be a stable globally unique ID like a UUID."`

	SchoolTerminationGraceDays int `yaml:"school_termination_min_grace_days" env:"GRAPH_LDAP_SCHOOL_TERMINATION_MIN_GRACE_DAYS" desc:"When setting a 'terminationDate' for a school, require the date to be at least this number of days in the future."`
}

type Identity struct {
	Backend string `yaml:"backend" env:"GRAPH_IDENTITY_BACKEND" desc:"The user identity backend to use. Supported backend types are 'ldap' and 'cs3'."`
	LDAP    LDAP   `yaml:"ldap"`
}

// API represents API configuration parameters.
type API struct {
	GroupMembersPatchLimit  int    `yaml:"group_members_patch_limit" env:"GRAPH_GROUP_MEMBERS_PATCH_LIMIT" desc:"The amount of group members allowed to be added with a single patch request."`
	UsernameMatch           string `yaml:"graph_username_match" env:"GRAPH_USERNAME_MATCH" desc:"Apply restrictions to usernames. Supported values are 'default' and 'none'. When set to 'default', user names must not start with a number and are restricted to ASCII characters. When set to 'none', no restrictions are applied. The default value is 'default'."`
	AssignDefaultUserRole   bool   `yaml:"graph_assign_default_user_role" env:"GRAPH_ASSIGN_DEFAULT_USER_ROLE" desc:"Whether to assign newly created users the default role 'User'. Set this to 'false' if you want to assign roles manually, or if the role assignment should happen at first login. Set this to 'true' (the default) to assign the role 'User' when creating a new user."`
	IdentitySearchMinLength int    `yaml:"graph_identity_search_min_length" env:"GRAPH_IDENTITY_SEARCH_MIN_LENGTH" desc:"The minimum length the search term needs to have for unprivileged users when searching for users or groups."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;GRAPH_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Set to a empty string to disable emitting events."`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;GRAPH_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;GRAPH_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;GRAPH_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided GRAPH_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;GRAPH_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;GRAPH_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details."`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;GRAPH_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details."`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;GRAPH_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details."`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;GRAPH_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials."`
}

// Keycloak configuration
type Keycloak struct {
	BasePath           string `yaml:"base_path" env:"OCIS_KEYCLOAK_BASE_PATH;GRAPH_KEYCLOAK_BASE_PATH" desc:"The URL to access keycloak."`
	ClientID           string `yaml:"client_id" env:"OCIS_KEYCLOAK_CLIENT_ID;GRAPH_KEYCLOAK_CLIENT_ID" desc:"The client id to authenticate with keycloak."`
	ClientSecret       string `yaml:"client_secret" env:"OCIS_KEYCLOAK_CLIENT_SECRET;GRAPH_KEYCLOAK_CLIENT_SECRET" desc:"The client secret to use in authentication."`
	ClientRealm        string `yaml:"client_realm" env:"OCIS_KEYCLOAK_CLIENT_REALM;GRAPH_KEYCLOAK_CLIENT_REALM" desc:"The realm the client is defined in."`
	UserRealm          string `yaml:"user_realm" env:"OCIS_KEYCLOAK_USER_REALM;GRAPH_KEYCLOAK_USER_REALM" desc:"The realm users are defined."`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" env:"OCIS_KEYCLOAK_INSECURE_SKIP_VERIFY;GRAPH_KEYCLOAK_INSECURE_SKIP_VERIFY" desc:"Disable TLS certificate validation for Keycloak connections. Do not set this in production environments."`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ServiceAccountID     string `yaml:"service_account_id" env:"OCIS_SERVICE_ACCOUNT_ID;GRAPH_SERVICE_ACCOUNT_ID" desc:"The ID of the service account the service should use. See the 'auth-service' service description for more details."`
	ServiceAccountSecret string `yaml:"service_account_secret" env:"OCIS_SERVICE_ACCOUNT_SECRET;GRAPH_SERVICE_ACCOUNT_SECRET" desc:"The service account secret."`
}

// FilesSharing is the configuration for the files sharing
type FilesSharing struct {
	EnableResharing bool `yaml:"enable_resharing" env:"OCIS_ENABLE_RESHARING;GRAPH_ENABLE_RESHARING" desc:"Changing this value is NOT supported. Enables the support for resharing."`
}
