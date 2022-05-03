## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>AUTH_BASIC_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>AUTH_BASIC_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>AUTH_BASIC_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>AUTH_BASIC_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>AUTH_BASIC_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>AUTH_BASIC_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>AUTH_BASIC_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>AUTH_BASIC_LOG_FILE | string |  | The target log file.|
| AUTH_BASIC_DEBUG_ADDR | string | 127.0.0.1:9147 | |
| AUTH_BASIC_DEBUG_TOKEN | string |  | |
| AUTH_BASIC_DEBUG_PPROF | bool | false | |
| AUTH_BASIC_DEBUG_ZPAGES | bool | false | |
| AUTH_BASIC_GRPC_ADDR | string | 127.0.0.1:9146 | The address of the grpc service.|
| AUTH_BASIC_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>AUTH_BASIC_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| AUTH_BASIC_AUTH_PROVIDER | string | ldap | The auth provider which should be used by the service|
| AUTH_BASIC_JSON_PROVIDER_FILE | string |  | The file to which the json provider writes the data.|
| LDAP_URI<br/>AUTH_BASIC_LDAP_URI | string | ldaps://localhost:9235 | |
| LDAP_CACERT<br/>AUTH_BASIC_LDAP_CACERT | string | ~/.ocis/idm/ldap.crt | |
| LDAP_INSECURE<br/>AUTH_BASIC_LDAP_INSECURE | bool | false | |
| LDAP_BIND_DN<br/>AUTH_BASIC_LDAP_BIND_DN | string | uid=reva,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD<br/>AUTH_BASIC_LDAP_BIND_PASSWORD | string |  | |
| LDAP_USER_BASE_DN<br/>AUTH_BASIC_LDAP_USER_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_GROUP_BASE_DN<br/>AUTH_BASIC_LDAP_GROUP_BASE_DN | string | ou=groups,o=libregraph-idm | |
| LDAP_USER_SCOPE<br/>AUTH_BASIC_LDAP_USER_SCOPE | string | sub | |
| LDAP_GROUP_SCOPE<br/>AUTH_BASIC_LDAP_GROUP_SCOPE | string | sub | |
| LDAP_USERFILTER<br/>AUTH_BASIC_LDAP_USERFILTER | string |  | |
| LDAP_GROUPFILTER<br/>AUTH_BASIC_LDAP_USERFILTER | string |  | |
| LDAP_USER_OBJECTCLASS<br/>AUTH_BASIC_LDAP_USER_OBJECTCLASS | string | inetOrgPerson | |
| LDAP_GROUP_OBJECTCLASS<br/>AUTH_BASIC_LDAP_GROUP_OBJECTCLASS | string | groupOfNames | |
| LDAP_LOGIN_ATTRIBUTES<br/>AUTH_BASIC_LDAP_LOGIN_ATTRIBUTES |  | [uid mail] | |
| OCIS_URL<br/>AUTH_BASIC_IDP_URL | string | https://localhost:9200 | |
| LDAP_USER_SCHEMA_ID<br/>AUTH_BASIC_LDAP_USER_SCHEMA_ID | string | ownclouduuid | |
| LDAP_USER_SCHEMA_ID_IS_OCTETSTRING<br/>AUTH_BASIC_LDAP_USER_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_USER_SCHEMA_MAIL<br/>AUTH_BASIC_LDAP_USER_SCHEMA_MAIL | string | mail | |
| LDAP_USER_SCHEMA_DISPLAYNAME<br/>AUTH_BASIC_LDAP_USER_SCHEMA_DISPLAYNAME | string | displayname | |
| LDAP_USER_SCHEMA_USERNAME<br/>AUTH_BASIC_LDAP_USER_SCHEMA_USERNAME | string | uid | |
| LDAP_GROUP_SCHEMA_ID<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_ID | string | ownclouduuid | |
| LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_ID_IS_OCTETSTRING | bool | false | |
| LDAP_GROUP_SCHEMA_MAIL<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_MAIL | string | mail | |
| LDAP_GROUP_SCHEMA_DISPLAYNAME<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_DISPLAYNAME | string | cn | |
| LDAP_GROUP_SCHEMA_GROUPNAME<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_GROUPNAME | string | cn | |
| LDAP_GROUP_SCHEMA_MEMBER<br/>AUTH_BASIC_LDAP_GROUP_SCHEMA_MEMBER | string | member | |