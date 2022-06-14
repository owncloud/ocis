## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>FRONTEND_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>FRONTEND_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>FRONTEND_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>FRONTEND_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>FRONTEND_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>FRONTEND_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>FRONTEND_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>FRONTEND_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| FRONTEND_DEBUG_ADDR | string | 127.0.0.1:9141 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| FRONTEND_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| FRONTEND_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| FRONTEND_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| FRONTEND_HTTP_ADDR | string | 127.0.0.1:9140 | The address of the http service.|
| FRONTEND_HTTP_PROTOCOL | string | tcp | The transport protocol of the http service.|
| FRONTEND_HTTP_PREFIX | string |  | |
| STORAGE_TRANSFER_SECRET | string |  | |
| OCIS_JWT_SECRET<br/>FRONTEND_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| OCIS_MACHINE_AUTH_API_KEY<br/>FRONTEND_MACHINE_AUTH_API_KEY | string |  | |
| FRONTEND_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| FRONTEND_ENABLE_FAVORITES | bool | false | |
| FRONTEND_ENABLE_PROJECT_SPACES | bool | true | Indicates to clients that project spaces are supposed to be made available.|
| FRONTEND_ENABLE_SHARE_JAIL | bool | true | Indicates to clients that the share jail is supposed to be used.|
| FRONTEND_UPLOAD_MAX_CHUNK_SIZE | int | 100000000 | |
| FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE | string |  | |
| FRONTEND_DEFAULT_UPLOAD_PROTOCOL | string | tus | |
| OCIS_URL<br/>FRONTEND_PUBLIC_URL | string | https://localhost:9200 | |
| OCIS_INSECURE<br/>FRONTEND_APP_HANDLER_INSECURE | bool | false | |
| FRONTEND_ARCHIVER_MAX_NUM_FILES | int64 | 10000 | |
| FRONTEND_ARCHIVER_MAX_SIZE | int64 | 1073741824 | |
| OCIS_INSECURE<br/>FRONTEND_ARCHIVER_INSECURE | bool | false | |
| FRONTEND_DATA_GATEWAY_PREFIX | string | data | |
| FRONTEND_OCS_PREFIX | string | ocs | |
| FRONTEND_OCS_SHARE_PREFIX | string | /Shares | |
| FRONTEND_OCS_HOME_NAMESPACE | string | /users/{{.Id.OpaqueId}} | |
| FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE | string | {{.Mail}} | |
| FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL | int | 0 | |
| FRONTEND_CHECKSUMS_SUPPORTED_TYPES |  | [sha1 md5 adler32] | |
| FRONTEND_CHECKSUMS_PREFERRED_UPLOAD_TYPES | string |  | |