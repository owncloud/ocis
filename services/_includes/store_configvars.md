## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORE_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORE_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>STORE_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>STORE_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>STORE_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>STORE_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORE_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORE_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| STORE_DEBUG_ADDR | string | 127.0.0.1:9464 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| STORE_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| STORE_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| STORE_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| STORE_GRPC_ADDR | string | 127.0.0.1:9460 | The bind address of the GRPC service.|
| STORE_DATA_PATH | string | ~/.ocis/store | Path for the store persistence directory.|