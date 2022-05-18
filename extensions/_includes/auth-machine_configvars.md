## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>AUTH_MACHINE_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>AUTH_MACHINE_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>AUTH_MACHINE_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>AUTH_MACHINE_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>AUTH_MACHINE_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>AUTH_MACHINE_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>AUTH_MACHINE_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>AUTH_MACHINE_LOG_FILE | string |  | The target log file.|
| AUTH_MACHINE_DEBUG_ADDR | string | 127.0.0.1:9167 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| AUTH_MACHINE_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| AUTH_MACHINE_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| AUTH_MACHINE_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| AUTH_MACHINE_GRPC_ADDR | string | 127.0.0.1:9166 | The address of the grpc service.|
| AUTH_MACHINE_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>AUTH_MACHINE_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| AUTH_MACHINE_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| OCIS_MACHINE_AUTH_API_KEY<br/>AUTH_MACHINE_API_KEY | string |  | |