## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>GRAPH_EXPLORER_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>GRAPH_EXPLORER_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>GRAPH_EXPLORER_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>GRAPH_EXPLORER_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>GRAPH_EXPLORER_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>GRAPH_EXPLORER_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>GRAPH_EXPLORER_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>GRAPH_EXPLORER_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| GRAPH_EXPLORER_DEBUG_ADDR | string | 127.0.0.1:9136 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| GRAPH_EXPLORER_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| GRAPH_EXPLORER_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| GRAPH_EXPLORER_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| GRAPH_EXPLORER_HTTP_ADDR | string | 127.0.0.1:9135 | The bind address of the HTTP service.|
| GRAPH_EXPLORER_HTTP_ROOT | string | /graph-explorer | Subdirectory that serves as the root for this HTTP service.|
| GRAPH_EXPLORER_CLIENT_ID | string | ocis-explorer.js | OIDC client ID the graph explorer uses. This client needs to be set up in your IDP.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>GRAPH_EXPLORER_ISSUER | string | https://localhost:9200 | URL of the OIDC issuer. It defaults to URL of the builtin IDP.|
| OCIS_URL<br/>GRAPH_EXPLORER_GRAPH_URL_BASE | string | https://localhost:9200 | Base URL where the graph explorer is reachable for users.|
| GRAPH_EXPLORER_GRAPH_URL_PATH | string | /graph | URL path where the graph explorer is reachable for users.|