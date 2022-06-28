## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>SETTINGS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>SETTINGS_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>SETTINGS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>SETTINGS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>SETTINGS_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>SETTINGS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>SETTINGS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>SETTINGS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| SETTINGS_DEBUG_ADDR | string | 127.0.0.1:9194 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| SETTINGS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| SETTINGS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| SETTINGS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| SETTINGS_HTTP_ADDR | string | 127.0.0.1:9190 | |
| SETTINGS_HTTP_ROOT | string | / | |
| SETTINGS_CACHE_TTL | int | 604800 | |
| SETTINGS_GRPC_ADDR | string | 127.0.0.1:9191 | |
| SETTINGS_STORE_TYPE | string | metadata | |
| SETTINGS_DATA_PATH | string | ~/.ocis/settings | |
| STORAGE_GATEWAY_GRPC_ADDR | string | 127.0.0.1:9215 | |
| STORAGE_GRPC_ADDR | string | 127.0.0.1:9215 | |
| OCIS_SYSTEM_USER_ID<br/>SETTINGS_SYSTEM_USER_ID | string |  | |
| OCIS_SYSTEM_USER_IDP<br/>SETTINGS_SYSTEM_USER_IDP | string | internal | |
| OCIS_SYSTEM_USER_API_KEY | string |  | |
| OCIS_ADMIN_USER_ID<br/>SETTINGS_ADMIN_USER_ID | string |  | ID of a user, that should receive admin privileges.|
| SETTINGS_ASSET_PATH | string |  | |
| OCIS_JWT_SECRET<br/>SETTINGS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| SETTINGS_SETUP_DEFAULT_ASSIGNMENTS<br/>ACCOUNTS_DEMO_USERS_AND_GROUPS | bool | false | If the default role assignments for the demo users should be setup.|