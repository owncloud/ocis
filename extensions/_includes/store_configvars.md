## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORE_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>STORE_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>STORE_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>STORE_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>STORE_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>STORE_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>STORE_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>STORE_LOG_FILE | string |  | |
| STORE_DEBUG_ADDR | string | 127.0.0.1:9464 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| STORE_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| STORE_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| STORE_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| STORE_GRPC_ADDR | string | 127.0.0.1:9460 | |
| STORE_DATA_PATH | string | ~/.ocis/store | |