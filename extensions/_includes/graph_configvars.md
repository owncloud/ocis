## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>GRAPH_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>GRAPH_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>GRAPH_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>GRAPH_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>GRAPH_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>GRAPH_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>GRAPH_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>GRAPH_LOG_FILE | string |  | |
| GRAPH_DEBUG_ADDR | string | 127.0.0.1:9124 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| GRAPH_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| GRAPH_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| GRAPH_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| GRAPH_HTTP_ADDR | string | 127.0.0.1:9120 | |
| GRAPH_HTTP_ROOT | string | /graph | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| OCIS_JWT_SECRET<br/>GRAPH_JWT_SECRET | string |  | |
| OCIS_URL<br/>GRAPH_SPACES_WEBDAV_BASE | string | https://localhost:9200 | |
| GRAPH_SPACES_WEBDAV_PATH | string | /dav/spaces/ | |
| GRAPH_SPACES_DEFAULT_QUOTA | string | 1000000000 | |
| OCIS_INSECURE<br/>GRAPH_SPACES_INSECURE | bool | false | |
| GRAPH_SPACES_EXTENDED_SPACE_PROPERTIES_CACHE_TTL | int | 0 | |
| GRAPH_IDENTITY_BACKEND | string | ldap | |
| LDAP_URI<br/>GRAPH_LDAP_URI | string | ldaps://localhost:9235 | |
| OCIS_INSECURE<br/>GRAPH_LDAP_INSECURE | bool | true | |
| LDAP_BIND_DN<br/>GRAPH_LDAP_BIND_DN | string | uid=libregraph,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD<br/>GRAPH_LDAP_BIND_PASSWORD | string |  | |
| GRAPH_LDAP_SERVER_UUID | bool | false | |
| GRAPH_LDAP_SERVER_WRITE_ENABLED | bool | true | |
| LDAP_USER_BASE_DN<br/>GRAPH_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_USER_SCOPE<br/>GRAPH_LDAP_USER_SCOPE | string | sub | |
| LDAP_USER_FILTER<br/>GRAPH_LDAP_USER_FILTER | string |  | |
| LDAP_USER_OBJECTCLASS<br/>GRAPH_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | |
| LDAP_USER_SCHEMA_MAIL<br/>GRAPH_LDAP_USER_EMAIL_ATTRIBUTE | string | mail | |
| LDAP_USER_SCHEMA_DISPLAY_NAME<br/>GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE | string | displayName | |
| LDAP_USER_SCHEMA_USERNAME<br/>GRAPH_LDAP_USER_NAME_ATTRIBUTE | string | uid | |
| LDAP_USER_SCHEMA_ID<br/>GRAPH_LDAP_USER_UID_ATTRIBUTE | string | owncloudUUID | |
| LDAP_GROUP_BASE_DN<br/>GRAPH_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | |
| LDAP_GROUP_SCOPE<br/>GRAPH_LDAP_GROUP_SEARCH_SCOPE | string | sub | |
| LDAP_GROUP_FILTER<br/>GRAPH_LDAP_GROUP_FILTER | string |  | |
| LDAP_GROUP_OBJECTCLASS<br/>GRAPH_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | |
| LDAP_GROUP_SCHEMA_GROUPNAME<br/>GRAPH_LDAP_GROUP_NAME_ATTRIBUTE | string | cn | |
| LDAP_GROUP_SCHEMA_ID<br/>GRAPH_LDAP_GROUP_ID_ATTRIBUTE | string | owncloudUUID | |
| GRAPH_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| GRAPH_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|