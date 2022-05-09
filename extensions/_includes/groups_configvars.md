## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>GROUPS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>GROUPS_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>GROUPS_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>GROUPS_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>GROUPS_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>GROUPS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>GROUPS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>GROUPS_LOG_FILE | string |  | The target log file.|
| GROUPS_DEBUG_ADDR | string | 127.0.0.1:9161 | |
| GROUPS_DEBUG_TOKEN | string |  | |
| GROUPS_DEBUG_PPROF | bool | false | |
| GROUPS_DEBUG_ZPAGES | bool | false | |
| GROUPS_GRPC_ADDR | string | 127.0.0.1:9160 | The address of the grpc service.|
| GROUPS_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>GROUPS_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| GROUPS_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| LDAP_URI<br/>GROUPS_LDAP_URI | string | ldaps://localhost:9235 | |
| LDAP_CACERT<br/>GROUPS_LDAP_CACERT | string | ~/.ocis/idm/ldap.crt | |
| LDAP_INSECURE<br/>GROUPS_LDAP_INSECURE | bool | false | |
| LDAP_BIND_DN<br/>GROUPS_LDAP_BIND_DN | string | uid=reva,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD<br/>GROUPS_LDAP_BIND_PASSWORD | string |  | |
| LDAP_USER_BASE_DN<br/>GROUPS_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_GROUP_BASE_DN<br/>GROUPS_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | |
| LDAP_USER_SCOPE<br/>GROUPS_LDAP_USER_SCOPE | string | sub | |
| LDAP_GROUP_SCOPE<br/>GROUPS_LDAP_GROUP_SCOPE | string | sub | |
| LDAP_USERFILTER<br/>GROUPS_LDAP_USERFILTER | string |  | |
| LDAP_GROUPFILTER<br/>GROUPS_LDAP_USERFILTER | string |  | |
| LDAP_USER_OBJECTCLASS<br/>GROUPS_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | |
| LDAP_GROUP_OBJECTCLASS<br/>GROUPS_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | |
| LDAP_LOGIN_ATTRIBUTES<br/>GROUPS_LDAP_LOGIN_ATTRIBUTES |  | [uid mail] | |
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>GROUPS_IDP_URL | string | https://localhost:9200 | |
| LDAP_USER_SCHEMA_ID<br/>GROUPS_LDAP_USER_SCHEMA_ID | string | ownclouduuid | |
| LDAP_USER_SCHEMA_ID_IS_OCTETSTRING<br/>GROUPS_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_USER_SCHEMA_MAIL<br/>GROUPS_LDAP_USER_SCHEMA_MAIL | string | mail | |
| LDAP_USER_SCHEMA_DISPLAYNAME<br/>GROUPS_LDAP_USER_SCHEMA_DISPLAYNAME | string | displayname | |
| LDAP_USER_SCHEMA_USERNAME<br/>GROUPS_LDAP_USER_SCHEMA_USERNAME | string | uid | |
| LDAP_GROUP_SCHEMA_ID<br/>GROUPS_LDAP_GROUP_SCHEMA_ID | string | ownclouduuid | |
| LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING<br/>GROUPS_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_GROUP_SCHEMA_MAIL<br/>GROUPS_LDAP_GROUP_SCHEMA_MAIL | string | mail | |
| LDAP_GROUP_SCHEMA_DISPLAYNAME<br/>GROUPS_LDAP_GROUP_SCHEMA_DISPLAYNAME | string | cn | |
| LDAP_GROUP_SCHEMA_GROUPNAME<br/>GROUPS_LDAP_GROUP_SCHEMA_GROUPNAME | string | cn | |
| LDAP_GROUP_SCHEMA_MEMBER<br/>GROUPS_LDAP_GROUP_SCHEMA_MEMBER | string | member | |
| GROUPS_OWNCLOUDSQL_DB_USERNAME | string | owncloud | |
| GROUPS_OWNCLOUDSQL_DB_PASSWORD | string |  | |
| GROUPS_OWNCLOUDSQL_DB_HOST | string | mysql | |
| GROUPS_OWNCLOUDSQL_DB_PORT | int | 3306 | |
| GROUPS_OWNCLOUDSQL_DB_NAME | string | owncloud | |
| GROUPS_OWNCLOUDSQL_IDP | string | https://localhost:9200 | |
| GROUPS_OWNCLOUDSQL_NOBODY | int64 | 90 | |
| GROUPS_OWNCLOUDSQL_JOIN_USERNAME | bool | false | |
| GROUPS_OWNCLOUDSQL_JOIN_OWNCLOUD_UUID | bool | false | |
| GROUPS_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH | bool | false | |