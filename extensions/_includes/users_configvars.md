## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>USERS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>USERS_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>USERS_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>USERS_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>USERS_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>USERS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>USERS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>USERS_LOG_FILE | string |  | The target log file.|
| USERS_DEBUG_ADDR | string | 127.0.0.1:9145 | |
| USERS_DEBUG_TOKEN | string |  | |
| USERS_DEBUG_PPROF | bool | false | |
| USERS_DEBUG_ZPAGES | bool | false | |
| USERS_GRPC_ADDR | string | 127.0.0.1:9144 | The address of the grpc service.|
| USERS_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>USERS_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| USERS_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| LDAP_URI<br/>USERS_LDAP_URI | string | ldaps://localhost:9235 | |
| LDAP_CACERT<br/>USERS_LDAP_CACERT | string | ~/.ocis/idm/ldap.crt | |
| LDAP_INSECURE<br/>USERS_LDAP_INSECURE | bool | false | |
| LDAP_BIND_DN<br/>USERS_LDAP_BIND_DN | string | uid=reva,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD<br/>USERS_LDAP_BIND_PASSWORD | string |  | |
| LDAP_USER_BASE_DN<br/>USERS_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_GROUP_BASE_DN<br/>USERS_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | |
| LDAP_USER_SCOPE<br/>USERS_LDAP_USER_SCOPE | string | sub | |
| LDAP_GROUP_SCOPE<br/>USERS_LDAP_GROUP_SCOPE | string | sub | |
| LDAP_USERFILTER<br/>USERS_LDAP_USERFILTER | string |  | |
| LDAP_GROUPFILTER<br/>USERS_LDAP_USERFILTER | string |  | |
| LDAP_USER_OBJECTCLASS<br/>USERS_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | |
| LDAP_GROUP_OBJECTCLASS<br/>USERS_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | |
| LDAP_LOGIN_ATTRIBUTES<br/>USERS_LDAP_LOGIN_ATTRIBUTES |  | [uid mail] | |
| OCIS_URL<br/>USERS_IDP_URL | string | https://localhost:9200 | |
| LDAP_USER_SCHEMA_ID<br/>USERS_LDAP_USER_SCHEMA_ID | string | ownclouduuid | |
| LDAP_USER_SCHEMA_ID_IS_OCTETSTRING<br/>USERS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_USER_SCHEMA_MAIL<br/>USERS_LDAP_USER_SCHEMA_MAIL | string | mail | |
| LDAP_USER_SCHEMA_DISPLAYNAME<br/>USERS_LDAP_USER_SCHEMA_DISPLAYNAME | string | displayname | |
| LDAP_USER_SCHEMA_USERNAME<br/>USERS_LDAP_USER_SCHEMA_USERNAME | string | uid | |
| LDAP_GROUP_SCHEMA_ID<br/>USERS_LDAP_GROUP_SCHEMA_ID | string | ownclouduuid | |
| LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING<br/>USERS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_GROUP_SCHEMA_MAIL<br/>USERS_LDAP_GROUP_SCHEMA_MAIL | string | mail | |
| LDAP_GROUP_SCHEMA_DISPLAYNAME<br/>USERS_LDAP_GROUP_SCHEMA_DISPLAYNAME | string | cn | |
| LDAP_GROUP_SCHEMA_GROUPNAME<br/>USERS_LDAP_GROUP_SCHEMA_GROUPNAME | string | cn | |
| LDAP_GROUP_SCHEMA_MEMBER<br/>USERS_LDAP_GROUP_SCHEMA_MEMBER | string | member | |
| USERS_OWNCLOUDSQL_DB_USERNAME | string | owncloud | |
| USERS_OWNCLOUDSQL_DB_PASSWORD | string | secret | |
| USERS_OWNCLOUDSQL_DB_HOST | string | mysql | |
| USERS_OWNCLOUDSQL_DB_PORT | int | 3306 | |
| USERS_OWNCLOUDSQL_DB_NAME | string | owncloud | |
| USERS_OWNCLOUDSQL_IDP | string | https://localhost:9200 | |
| USERS_OWNCLOUDSQL_NOBODY | int64 | 90 | |
| USERS_OWNCLOUDSQL_JOIN_USERNAME | bool | false | |
| USERS_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID | bool | false | |
| USERS_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH | bool | false | |