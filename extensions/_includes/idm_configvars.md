## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>IDM_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>IDM_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>IDM_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>IDM_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>IDM_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>IDM_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>IDM_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>IDM_LOG_FILE | string |  | |
| IDM_DEBUG_ADDR | string | 127.0.0.1:9239 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| IDM_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| IDM_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| IDM_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| IDM_LDAPS_ADDR | string | 127.0.0.1:9235 | Listen address for the ldaps listener (ip-addr:port)|
| IDM_LDAPS_CERT | string | ~/.ocis/idm/ldap.crt | File name of the TLS server certificate for the ldaps listener|
| IDM_LDAPS_KEY | string | ~/.ocis/idm/ldap.key | File name for the TLS certificate key for the server certificate|
| IDM_DATABASE_PATH | string | ~/.ocis/idm/ocis.boltdb | Full path to the idm backend database|
| IDM_CREATE_DEMO_USERS<br/>ACCOUNTS_DEMO_USERS_AND_GROUPS | bool | false | Flag to enabe/disable the creation of the demo users|
| IDM_ADMIN_PASSWORD | string |  | Password to set for the ocis "admin" user. Either cleartext or an argon2id hash|
| IDM_SVC_PASSWORD | string |  | Password to set for the "idm" service user. Either cleartext or an argon2id hash|
| IDM_REVASVC_PASSWORD | string |  | Password to set for the "reva" service user. Either cleartext or an argon2id hash|
| IDM_IDPSVC_PASSWORD | string |  | Password to set for the "idp" service user. Either cleartext or an argon2id hash|
| OCIS_ADMIN_USER_ID<br/>IDM_ADMIN_USER_ID | string |  | |