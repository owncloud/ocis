## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>WEBDAV_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>WEBDAV_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>WEBDAV_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>WEBDAV_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>WEBDAV_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>WEBDAV_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>WEBDAV_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>WEBDAV_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| WEBDAV_DEBUG_ADDR | string | 127.0.0.1:9119 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| WEBDAV_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| WEBDAV_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| WEBDAV_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| WEBDAV_HTTP_ADDR | string | 127.0.0.1:9115 | The HTTP API address.|
| WEBDAV_HTTP_ROOT | string | / | The HTTP API root path.|
| OCIS_URL<br/>OCIS_PUBLIC_URL | string | https://127.0.0.1:9200 | |
| WEBDAV_WEBDAV_NAMESPACE | string | /users/{{.Id.OpaqueId}} | CS3 path layout to use when forwarding /webdav requests|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint|