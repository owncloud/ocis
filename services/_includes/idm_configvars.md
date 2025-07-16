## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>IDM_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>IDM_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>IDM_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>IDM_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>IDM_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>IDM_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>IDM_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>IDM_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| IDM_DEBUG_ADDR | string | 127.0.0.1:9239 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| IDM_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| IDM_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| IDM_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| IDM_LDAPS_ADDR | string | 127.0.0.1:9235 | Listen address for the LDAPS listener (ip-addr:port).|
| IDM_LDAPS_CERT | string | /var/lib/ocis/idm/ldap.crt | File name of the TLS server certificate for the LDAPS listener. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm.|
| IDM_LDAPS_KEY | string | /var/lib/ocis/idm/ldap.key | File name for the TLS certificate key for the server certificate. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm.|
| IDM_DATABASE_PATH | string | /var/lib/ocis/idm/ocis.boltdb | Full path to the IDM backend database. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm.|
| IDM_CREATE_DEMO_USERS | bool | false | Flag to enable or disable the creation of the demo users.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER | string | https://localhost:9200 | The OIDC issuer URL to assign to the demo users.|
| IDM_ADMIN_PASSWORD | string |  | Password to set for the oCIS 'admin' user. Either cleartext or an argon2id hash.|
| IDM_SVC_PASSWORD | string |  | Password to set for the 'idm' service user. Either cleartext or an argon2id hash.|
| IDM_REVASVC_PASSWORD | string |  | Password to set for the 'reva' service user. Either cleartext or an argon2id hash.|
| IDM_IDPSVC_PASSWORD | string |  | Password to set for the 'idp' service user. Either cleartext or an argon2id hash.|
| OCIS_ADMIN_USER_ID<br/>IDM_ADMIN_USER_ID | string |  | ID of the user that should receive admin privileges. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand.|