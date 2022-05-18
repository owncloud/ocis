## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>APP_REGISTRY_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>APP_REGISTRY_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>APP_REGISTRY_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>APP_REGISTRY_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>APP_REGISTRY_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>APP_REGISTRY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>APP_REGISTRY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>APP_REGISTRY_LOG_FILE | string |  | The target log file.|
| APP_REGISTRY_DEBUG_ADDR | string | 127.0.0.1:9243 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| APP_REGISTRY_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| APP_REGISTRY_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| APP_REGISTRY_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| APP_REGISTRY_GRPC_ADDR | string | 127.0.0.1:9242 | The address of the grpc service.|
| APP_REGISTRY_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>APP_REGISTRY_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |