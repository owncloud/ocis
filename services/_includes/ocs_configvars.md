## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>OCS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>OCS_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>OCS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>OCS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>OCS_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>OCS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>OCS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>OCS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| OCS_DEBUG_ADDR | string | 127.0.0.1:9114 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| OCS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| OCS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| OCS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCS_HTTP_ADDR | string | 127.0.0.1:9110 | The bind address of the HTTP service.|
| OCS_HTTP_ROOT | string | /ocs | Subdirectory that serves as the root for this HTTP service.|
| OCIS_JWT_SECRET<br/>OCS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>OCS_IDM_ADDRESS | string | https://localhost:9200 | URL of the OIDC issuer. It defaults to URL of the builtin IDP.|
| OCIS_MACHINE_AUTH_API_KEY<br/>OCS_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used for accessing the 'auth-machine' service to impersonate users.|