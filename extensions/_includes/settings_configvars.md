## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>SETTINGS_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>SETTINGS_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>SETTINGS_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>SETTINGS_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>SETTINGS_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>SETTINGS_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>SETTINGS_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>SETTINGS_LOG_FILE | string |  | |
| SETTINGS_DEBUG_ADDR | string | 127.0.0.1:9194 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| SETTINGS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| SETTINGS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| SETTINGS_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
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
| OCIS_ADMIN_USER_ID<br/>SETTINGS_ADMIN_USER_ID | string |  | |
| SETTINGS_ASSET_PATH | string |  | |
| OCIS_JWT_SECRET<br/>SETTINGS_JWT_SECRET | string |  | |