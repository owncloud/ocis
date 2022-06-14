## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| IDP_PASSWORD_RESET_URI | string |  | The URI where a user can reset their password.|
| OCIS_TRACING_ENABLED<br/>IDP_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>IDP_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>IDP_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>IDP_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>IDP_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>IDP_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>IDP_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>IDP_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| IDP_DEBUG_ADDR | string | 127.0.0.1:9134 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| IDP_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| IDP_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| IDP_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| IDP_HTTP_ADDR | string | 127.0.0.1:9130 | |
| IDP_HTTP_ROOT | string | / | |
| IDP_TRANSPORT_TLS_CERT | string | ~/.ocis/idp/server.crt | |
| IDP_TRANSPORT_TLS_KEY | string | ~/.ocis/idp/server.key | |
| IDP_TLS | bool | false | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | CS3 gateway used to authenticate and look up users|
| OCIS_MACHINE_AUTH_API_KEY<br/>IDP_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used for accessing the 'auth-machine' service to impersonate users when looking up their userinfo via the 'cs3' backend.|
| IDP_ASSET_PATH | string |  | |
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>IDP_ISS | string | https://localhost:9200 | The OIDC issuer URL to use.|
| IDP_IDENTITY_MANAGER | string | ldap | The identity manager implementation to use, defaults to 'ldap', can be changed to 'cs3', 'kc', 'libregraph', 'cookie' or 'guest'.|
| IDP_URI_BASE_PATH | string |  | |
| IDP_SIGN_IN_URI | string |  | |
| IDP_SIGN_OUT_URI | string |  | |
| IDP_ENDPOINT_URI | string |  | |
| IDP_ENDSESSION_ENDPOINT_URI | string |  | |
| LDAP_INSECURE<br/>IDP_INSECURE | bool | false | Allow insecure connections to the user backend (eg. LDAP, CS3 api, ...).|
| IDP_ALLOW_CLIENT_GUESTS | bool | false | |
| IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION | bool | false | |
| IDP_ENCRYPTION_SECRET_FILE | string |  | |
| IDP_DISABLE_IDENTIFIER_WEBAPP | bool | true | |
| IDP_IDENTIFIER_SCOPES_CONF | string |  | |
| IDP_SIGNING_KID | string |  | |
| IDP_SIGNING_METHOD | string | PS256 | |
| IDP_SIGNING_PRIVATE_KEY_FILES |  | [] | |
| IDP_VALIDATION_KEYS_PATH | string |  | |
| IDP_ACCESS_TOKEN_EXPIRATION | uint64 | 86400 | |
| IDP_ID_TOKEN_EXPIRATION | uint64 | 3600 | |
| IDP_REFRESH_TOKEN_EXPIRATION | uint64 | 94608000 | |
|  | uint64 | 0 | |
| LDAP_URI<br/>IDP_LDAP_URI | string | ldaps://localhost:9235 | |
| LDAP_CACERT<br/>IDP_LDAP_TLS_CACERT | string | ~/.ocis/idm/ldap.crt | |
| LDAP_BIND_DN<br/>IDP_LDAP_BIND_DN | string | uid=idp,ou=sysusers,o=libregraph-idm | |
| LDAP_BIND_PASSWORD<br/>IDP_LDAP_BIND_PASSWORD | string |  | |
| LDAP_USER_BASE_DN<br/>IDP_LDAP_BASE_DN | string | ou=users,o=libregraph-idm | |
| LDAP_USER_SCOPE<br/>IDP_LDAP_SCOPE | string | sub | |
| IDP_LDAP_LOGIN_ATTRIBUTE | string | uid | |
| LDAP_USER_SCHEMA_MAIL<br/>IDP_LDAP_EMAIL_ATTRIBUTE | string | mail | |
| LDAP_USER_SCHEMA_USERNAME<br/>IDP_LDAP_NAME_ATTRIBUTE | string | displayName | |
| LDAP_USER_SCHEMA_ID<br/>IDP_LDAP_UUID_ATTRIBUTE | string | uid | |
| IDP_LDAP_UUID_ATTRIBUTE_TYPE | string | text | |
| LDAP_USER_FILTER<br/>IDP_LDAP_FILTER | string |  | |
| LDAP_USER_OBJECTCLASS<br/>IDP_LDAP_OBJECTCLASS | string | inetOrgPerson | |