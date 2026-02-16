## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>GRAPH_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>GRAPH_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>GRAPH_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>GRAPH_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>GRAPH_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>GRAPH_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>GRAPH_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>GRAPH_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| OCIS_CACHE_STORE<br/>GRAPH_CACHE_STORE | string | memory | The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details.|
| OCIS_CACHE_STORE_NODES<br/>GRAPH_CACHE_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| GRAPH_CACHE_STORE_DATABASE | string | cache-roles | The database name the configured store should use.|
| GRAPH_CACHE_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_CACHE_TTL<br/>GRAPH_CACHE_TTL | Duration | 336h0m0s | Time to live for cache records in the graph. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details.|
| OCIS_CACHE_DISABLE_PERSISTENCE<br/>GRAPH_CACHE_DISABLE_PERSISTENCE | bool | false | Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false.|
| OCIS_CACHE_AUTH_USERNAME<br/>GRAPH_CACHE_AUTH_USERNAME | string |  | The username to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_CACHE_AUTH_PASSWORD<br/>GRAPH_CACHE_AUTH_PASSWORD | string |  | The password to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured.|
| GRAPH_DEBUG_ADDR | string | 127.0.0.1:9124 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| GRAPH_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| GRAPH_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| GRAPH_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| GRAPH_HTTP_ADDR | string | 127.0.0.1:9120 | The bind address of the HTTP service.|
| GRAPH_HTTP_ROOT | string | /graph | Subdirectory that serves as the root for this HTTP service.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| GRAPH_HTTP_API_TOKEN | string |  | An optional API bearer token|
| OCIS_CORS_ALLOW_ORIGINS<br/>GRAPH_CORS_ALLOW_ORIGINS | []string | [*] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>GRAPH_CORS_ALLOW_METHODS | []string | [GET POST PUT PATCH DELETE OPTIONS] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>GRAPH_CORS_ALLOW_HEADERS | []string | [Authorization Origin Content-Type Accept X-Requested-With X-Request-Id Purge Restore] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>GRAPH_CORS_ALLOW_CREDENTIALS | bool | true | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| GRAPH_GROUP_MEMBERS_PATCH_LIMIT | int | 20 | The amount of group members allowed to be added with a single patch request.|
| GRAPH_USERNAME_MATCH | string | default | Apply restrictions to usernames. Supported values are 'default' and 'none'. When set to 'default', user names must not start with a number and are restricted to ASCII characters. When set to 'none', no restrictions are applied. The default value is 'default'.|
| GRAPH_ASSIGN_DEFAULT_USER_ROLE | bool | true | Whether to assign newly created users the default role 'User'. Set this to 'false' if you want to assign roles manually, or if the role assignment should happen at first login. Set this to 'true' (the default) to assign the role 'User' when creating a new user.|
| GRAPH_IDENTITY_SEARCH_MIN_LENGTH | int | 3 | The minimum length the search term needs to have for unprivileged users when searching for users or groups.|
| OCIS_USER_SEARCH_DISPLAYED_ATTRIBUTES | []string | [] | The attributes to display in the user search results.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | The CS3 gateway endpoint.|
| OCIS_GRPC_CLIENT_TLS_MODE | string |  | TLS mode for grpc connection to the go-micro based grpc services. Possible values are 'off', 'insecure' and 'on'. 'off': disables transport security for the clients. 'insecure' allows using transport security, but disables certificate verification (to be used with the autogenerated self-signed certificates). 'on' enables transport security, including server certificate verification.|
| OCIS_GRPC_CLIENT_TLS_CACERT | string |  | Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the go-micro based grpc services.|
| OCIS_JWT_SECRET<br/>GRAPH_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| GRAPH_APPLICATION_ID | string |  | The ocis application ID shown in the graph. All app roles are tied to this ID.|
| GRAPH_APPLICATION_DISPLAYNAME | string | ownCloud Infinite Scale | The ocis application name.|
| OCIS_URL<br/>GRAPH_SPACES_WEBDAV_BASE | string | https://localhost:9200 | The public facing URL of WebDAV.|
| GRAPH_SPACES_WEBDAV_PATH | string | /dav/spaces/ | The WebDAV sub-path for spaces.|
| GRAPH_SPACES_DEFAULT_QUOTA | string | 1000000000 | The default quota in bytes.|
| GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL | int | 60000000000 | Max TTL in seconds for the spaces property cache.|
| GRAPH_SPACES_USERS_CACHE_TTL | int | 60000000000 | Max TTL in seconds for the spaces users cache.|
| GRAPH_SPACES_GROUPS_CACHE_TTL | int | 60000000000 | Max TTL in seconds for the spaces groups cache.|
| GRAPH_SPACES_STORAGE_USERS_ADDRESS | string | com.owncloud.api.storage-users | The address of the storage-users service.|
| OCIS_DEFAULT_LANGUAGE | string |  | The default language used by services and the WebUI. If not defined, English will be used as default. See the documentation for more details.|
| OCIS_TRANSLATION_PATH<br/>GRAPH_TRANSLATION_PATH | string |  | (optional) Set this to a path with custom translations to overwrite the builtin translations. Note that file and folder naming rules apply, see the documentation for more details.|
| GRAPH_IDENTITY_BACKEND | string | ldap | The user identity backend to use. Supported backend types are 'ldap' and 'cs3'.|
| OCIS_LDAP_URI<br/>GRAPH_LDAP_URI | string | ldaps://localhost:9235 | URI of the LDAP Server to connect to. Supported URI schemes are 'ldaps://' and 'ldap://'|
| OCIS_LDAP_CACERT<br/>GRAPH_LDAP_CACERT | string | /var/lib/ocis/idm/ldap.crt | Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the LDAP service. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm.|
| OCIS_LDAP_INSECURE<br/>GRAPH_LDAP_INSECURE | bool | false | Disable TLS certificate validation for the LDAP connections. Do not set this in production environments.|
| OCIS_LDAP_BIND_DN<br/>GRAPH_LDAP_BIND_DN | string | uid=libregraph,ou=sysusers,o=libregraph-idm | LDAP DN to use for simple bind authentication with the target LDAP server.|
| OCIS_LDAP_BIND_PASSWORD<br/>GRAPH_LDAP_BIND_PASSWORD | string |  | Password to use for authenticating the 'bind_dn'.|
| GRAPH_LDAP_SERVER_UUID | bool | false | If set to true, rely on the LDAP Server to generate a unique ID for users and groups, like when using 'entryUUID' as the user ID attribute.|
| GRAPH_LDAP_SERVER_USE_PASSWORD_MODIFY_EXOP | bool | true | Use the 'Password Modify Extended Operation' for updating user passwords.|
| OCIS_LDAP_SERVER_WRITE_ENABLED<br/>GRAPH_LDAP_SERVER_WRITE_ENABLED | bool | true | Allow creating, modifying and deleting LDAP users via the GRAPH API. This can only be set to 'true' when keeping default settings for the LDAP user and group attribute types (the 'OCIS_LDAP_USER_SCHEMA_* and 'OCIS_LDAP_GROUP_SCHEMA_* variables).|
| GRAPH_LDAP_REFINT_ENABLED | bool | false | Signals that the server has the refint plugin enabled, which makes some actions not needed.|
| OCIS_LDAP_USER_BASE_DN<br/>GRAPH_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | Search base DN for looking up LDAP users.|
| OCIS_LDAP_USER_SCOPE<br/>GRAPH_LDAP_USER_SCOPE | string | sub | LDAP search scope to use when looking up users. Supported scopes are 'base', 'one' and 'sub'.|
| OCIS_LDAP_USER_FILTER<br/>GRAPH_LDAP_USER_FILTER | string |  | LDAP filter to add to the default filters for user search like '(objectclass=ownCloud)'.|
| OCIS_LDAP_USER_OBJECTCLASS<br/>GRAPH_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | The object class to use for users in the default user search filter ('inetOrgPerson').|
| OCIS_LDAP_USER_SCHEMA_MAIL<br/>GRAPH_LDAP_USER_EMAIL_ATTRIBUTE | string | mail | LDAP Attribute to use for the email address of users.|
| OCIS_LDAP_USER_SCHEMA_DISPLAYNAME<br/>GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE | string | displayName | LDAP Attribute to use for the display name of users.|
| OCIS_LDAP_USER_SCHEMA_USERNAME<br/>GRAPH_LDAP_USER_NAME_ATTRIBUTE | string | uid | LDAP Attribute to use for username of users.|
| OCIS_LDAP_USER_SCHEMA_ID<br/>GRAPH_LDAP_USER_UID_ATTRIBUTE | string | owncloudUUID | LDAP Attribute to use as the unique ID for users. This should be a stable globally unique ID like a UUID.|
| OCIS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING<br/>GRAPH_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING | bool | false | Set this to true if the defined 'ID' attribute for users is of the 'OCTETSTRING' syntax. This is required when using the 'objectGUID' attribute of Active Directory for the user ID's.|
| OCIS_LDAP_USER_SCHEMA_USER_TYPE<br/>GRAPH_LDAP_USER_TYPE_ATTRIBUTE | string | ownCloudUserType | LDAP Attribute to distinguish between 'Member' and 'Guest' users. Default is 'ownCloudUserType'.|
| OCIS_LDAP_USER_ENABLED_ATTRIBUTE<br/>GRAPH_USER_ENABLED_ATTRIBUTE | string | ownCloudUserEnabled | LDAP Attribute to use as a flag telling if the user is enabled or disabled.|
| OCIS_LDAP_USER_SCHEMA_EXTERNAL_ID<br/>GRAPH_LDAP_EXTERNAL_ID_ATTRIBUTE | string | owncloudExternalID | LDAP attribute that references the external ID of users during the provisioning process. The final ID is provided by an external identity provider. If it is not set, a default attribute will be used instead.|
| OCIS_LDAP_DISABLE_USER_MECHANISM<br/>GRAPH_DISABLE_USER_MECHANISM | string | attribute | An option to control the behavior for disabling users. Supported options are 'none', 'attribute' and 'group'. If set to 'group', disabling a user via API will add the user to the configured group for disabled users, if set to 'attribute' this will be done in the ldap user entry, if set to 'none' the disable request is not processed. Default is 'attribute'.|
| OCIS_LDAP_DISABLED_USERS_GROUP_DN<br/>GRAPH_DISABLED_USERS_GROUP_DN | string | cn=DisabledUsersGroup,ou=groups,o=libregraph-idm | The distinguished name of the group to which added users will be classified as disabled when 'disable_user_mechanism' is set to 'group'.|
| OCIS_LDAP_GROUP_BASE_DN<br/>GRAPH_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | Search base DN for looking up LDAP groups.|
| GRAPH_LDAP_GROUP_CREATE_BASE_DN | string | ou=groups,o=libregraph-idm | Parent DN under which new groups are created. This DN needs to be subordinate to the 'GRAPH_LDAP_GROUP_BASE_DN'. This setting is only relevant when 'GRAPH_LDAP_SERVER_WRITE_ENABLED' is 'true'. It defaults to the value of 'GRAPH_LDAP_GROUP_BASE_DN'. All groups outside of this subtree are treated as readonly groups and cannot be updated.|
| OCIS_LDAP_GROUP_SCOPE<br/>GRAPH_LDAP_GROUP_SEARCH_SCOPE | string | sub | LDAP search scope to use when looking up groups. Supported scopes are 'base', 'one' and 'sub'.|
| OCIS_LDAP_GROUP_FILTER<br/>GRAPH_LDAP_GROUP_FILTER | string |  | LDAP filter to add to the default filters for group searches.|
| OCIS_LDAP_GROUP_OBJECTCLASS<br/>GRAPH_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | The object class to use for groups in the default group search filter ('groupOfNames').|
| OCIS_LDAP_GROUP_SCHEMA_GROUPNAME<br/>GRAPH_LDAP_GROUP_NAME_ATTRIBUTE | string | cn | LDAP Attribute to use for the name of groups.|
| OCIS_LDAP_GROUP_SCHEMA_MEMBER<br/>GRAPH_LDAP_GROUP_MEMBER_ATTRIBUTE | string | member | LDAP Attribute that is used for group members.|
| OCIS_LDAP_GROUP_SCHEMA_ID<br/>GRAPH_LDAP_GROUP_ID_ATTRIBUTE | string | owncloudUUID | LDAP Attribute to use as the unique id for groups. This should be a stable globally unique ID like a UUID.|
| OCIS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING<br/>GRAPH_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING | bool | false | Set this to true if the defined 'ID' attribute for groups is of the 'OCTETSTRING' syntax. This is required when using the 'objectGUID' attribute of Active Directory for the group ID's.|
| GRAPH_LDAP_EDUCATION_RESOURCES_ENABLED | bool | false | Enable LDAP support for managing education related resources.|
| GRAPH_LDAP_SCHOOL_BASE_DN | string |  | Search base DN for looking up LDAP schools.|
| GRAPH_LDAP_SCHOOL_SEARCH_SCOPE | string |  | LDAP search scope to use when looking up schools. Supported scopes are 'base', 'one' and 'sub'.|
| GRAPH_LDAP_SCHOOL_FILTER | string |  | LDAP filter to add to the default filters for school searches.|
| GRAPH_LDAP_SCHOOL_OBJECTCLASS | string |  | The object class to use for schools in the default school search filter.|
| GRAPH_LDAP_SCHOOL_NAME_ATTRIBUTE | string |  | LDAP Attribute to use for the name of a school.|
| GRAPH_LDAP_SCHOOL_NUMBER_ATTRIBUTE | string |  | LDAP Attribute to use for the number of a school.|
| GRAPH_LDAP_SCHOOL_ID_ATTRIBUTE | string |  | LDAP Attribute to use as the unique id for schools. This should be a stable globally unique ID like a UUID.|
| GRAPH_LDAP_SCHOOL_TERMINATION_MIN_GRACE_DAYS | int | 0 | When setting a 'terminationDate' for a school, require the date to be at least this number of days in the future.|
| GRAPH_LDAP_REQUIRE_EXTERNAL_ID | bool | false | If enabled, the 'OCIS_LDAP_USER_SCHEMA_EXTERNAL_ID' is used as primary identifier for the provisioning API.|
| OCIS_LDAP_USER_MEMBER_ATTRIBUTE | string |  | LDAP Attribute to signal the user is member of an instance. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_USER_GUEST_ATTRIBUTE | string |  | LDAP Attribute to signal the user is guest of an instance. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_GROUP_AFFILIATION_ATTRIBUTE | string |  | LDAP Attribute to signal which instance the group is belonging to. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_PRECISE_SEARCH_ATTRIBUTE | string |  | LDAP Attribute to be used for searching users on other instances. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_INSTANCE_MAPPER_ENABLED | bool | false | The InstanceMapper allows mapping instance names (user readable) to instance IDs (machine readable) based on an LDAP query. See other _INSTANCE_MAPPER_ env vars. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_INSTANCE_MAPPER_BASE_DN | string |  | BaseDN of the 'instancename to instanceid' mapper in LDAP. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_INSTANCE_MAPPER_NAME_ATTRIBUTE | string |  | LDAP Attribute of the instance name. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_INSTANCE_MAPPER_ID_ATTRIBUTE | string |  | LDAP Attribute of the instance ID. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_CROSS_INSTANCE_REFERENCE_TEMPLATE | string |  | Template for the users unique reference across oCIS instances. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_LDAP_INSTANCE_URL_TEMPLATE | string |  | Template for the instance URL. Requires OCIS_MULTI_INSTANCE_ENABLED.|
| OCIS_ENABLE_OCM<br/>GRAPH_INCLUDE_OCM_SHAREES | bool | false | Include OCM sharees when listing users.|
| OCIS_EVENTS_ENDPOINT<br/>GRAPH_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Set to a empty string to disable emitting events.|
| OCIS_EVENTS_CLUSTER<br/>GRAPH_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_INSECURE<br/>GRAPH_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>GRAPH_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided GRAPH_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>GRAPH_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>GRAPH_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>GRAPH_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| GRAPH_AVAILABLE_ROLES | []string | [b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 a8d5fe5e-96e3-418d-825b-534dbdf22b99 fb6c3e19-e378-47e5-b277-9732f9de6e21 58c63c02-1d89-4572-916a-870abc5a1b7d 2d00ce52-1fc2-4dbc-8b95-a73b73395f5a 1c996275-f1c9-4e71-abdf-a42f6495e960 312c0871-5ef7-4b3a-85b6-0e4074c64049] | A comma separated list of roles that are available for assignment.|
| OCIS_MAX_CONCURRENCY<br/>GRAPH_MAX_CONCURRENCY | int | 20 | The maximum number of concurrent requests the service will handle.|
| OCIS_KEYCLOAK_BASE_PATH<br/>GRAPH_KEYCLOAK_BASE_PATH | string |  | The URL to access keycloak.|
| OCIS_KEYCLOAK_CLIENT_ID<br/>GRAPH_KEYCLOAK_CLIENT_ID | string |  | The client id to authenticate with keycloak.|
| OCIS_KEYCLOAK_CLIENT_SECRET<br/>GRAPH_KEYCLOAK_CLIENT_SECRET | string |  | The client secret to use in authentication.|
| OCIS_KEYCLOAK_CLIENT_REALM<br/>GRAPH_KEYCLOAK_CLIENT_REALM | string |  | The realm the client is defined in.|
| OCIS_KEYCLOAK_USER_REALM<br/>GRAPH_KEYCLOAK_USER_REALM | string |  | The realm users are defined.|
| OCIS_KEYCLOAK_INSECURE_SKIP_VERIFY<br/>GRAPH_KEYCLOAK_INSECURE_SKIP_VERIFY | bool | false | Disable TLS certificate validation for Keycloak connections. Do not set this in production environments.|
| OCIS_SERVICE_ACCOUNT_ID<br/>GRAPH_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>GRAPH_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|
| OCIS_MULTI_INSTANCE_ENABLED | bool | false | Enable multiple instances of Infinite Scale.|
| OCIS_MULTI_INSTANCE_INSTANCEID | string |  | The unique ID of this instance.|
| OCIS_MULTI_INSTANCE_QUERY_TEMPLATE | string |  | The regular expression extracting username and instancename from a user provided search.|
| OCIS_MAX_TAG_LENGTH | int | 100 | Define the maximum tag length. Defaults to 100 if not set. Set to 0 to not limit the tag length. Changes only impact the validation of new tags.|