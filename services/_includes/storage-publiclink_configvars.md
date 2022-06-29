## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORAGE_PUBLICLINK_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORAGE_PUBLICLINK_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>STORAGE_PUBLICLINK_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>STORAGE_PUBLICLINK_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>STORAGE_PUBLICLINK_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>STORAGE_PUBLICLINK_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORAGE_PUBLICLINK_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORAGE_PUBLICLINK_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| STORAGE_PUBLICLINK_DEBUG_ADDR | string | 127.0.0.1:9179 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| STORAGE_PUBLICLINK_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| STORAGE_PUBLICLINK_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| STORAGE_PUBLICLINK_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| STORAGE_PUBLICLINK_GRPC_ADDR | string | 127.0.0.1:9178 | The bind address of the GRPC service.|
| STORAGE_PUBLICLINK_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>STORAGE_PUBLICLINK_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| STORAGE_PUBLICLINK_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the loading of user's group memberships from the reva access token.|
| STORAGE_PUBLICLINK_STORAGE_PROVIDER_MOUNT_ID | string | 7993447f-687f-490d-875c-ac95e89a62a4 | Mount ID of this storage.|