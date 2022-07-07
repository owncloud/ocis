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
| SETTINGS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| SETTINGS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| SETTINGS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| SETTINGS_HTTP_ADDR | string | 127.0.0.1:9190 | The bind address of the HTTP service.|
| SETTINGS_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| SETTINGS_CACHE_TTL | int | 604800 | Browser cache control max-age value in seconds for settings Web UI assets.|
| SETTINGS_GRPC_ADDR | string | 127.0.0.1:9191 | The bind address of the GRPC service.|
| SETTINGS_STORE_TYPE | string | metadata | Store type configures the persistency driver. Supported values are "metadata" and "filesystem".|
| SETTINGS_DATA_PATH | string | ~/.ocis/settings | Path for the persistence directory.|
| STORAGE_GATEWAY_GRPC_ADDR | string | 127.0.0.1:9215 | GRPC address of the STORAGE-SYSTEM service.|
| STORAGE_GRPC_ADDR | string | 127.0.0.1:9215 | GRPC address of the STORAGE-SYSTEM service.|
| OCIS_SYSTEM_USER_ID<br/>SETTINGS_SYSTEM_USER_ID | string |  | ID of the oCIS STORAGE-SYSTEM system user. Admins need to set the ID for the STORAGE-SYSTEM system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format.|
| OCIS_SYSTEM_USER_IDP<br/>SETTINGS_SYSTEM_USER_IDP | string | internal | IDP of the oCIS STORAGE-SYSTEM system user.|
| OCIS_SYSTEM_USER_API_KEY | string |  | API key for the STORAGE-SYSTEM system user.|
| OCIS_ADMIN_USER_ID<br/>SETTINGS_ADMIN_USER_ID | string |  | ID of the user that should receive admin privileges.|
| SETTINGS_ASSET_PATH | string |  | Serve settings Web UI assets from a path on the filesystem instead of the builtin assets. Can be used for development and customization.|
| OCIS_JWT_SECRET<br/>SETTINGS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| SETTINGS_SETUP_DEFAULT_ASSIGNMENTS<br/>ACCOUNTS_DEMO_USERS_AND_GROUPS | bool | false | The default role assignments the demo users should be setup.|