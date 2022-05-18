## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>OCS_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>OCS_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>OCS_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>OCS_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>OCS_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>OCS_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>OCS_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>OCS_LOG_FILE | string |  | |
| OCS_DEBUG_ADDR | string | 127.0.0.1:9114 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| OCS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| OCS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| OCS_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| OCS_HTTP_ADDR | string | 127.0.0.1:9110 | |
| OCS_HTTP_ROOT | string | /ocs | |
| OCIS_JWT_SECRET<br/>OCS_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>OCS_IDM_ADDRESS | string | https://localhost:9200 | |
| OCS_ACCOUNT_BACKEND_TYPE | string | cs3 | |
| STORAGE_USERS_DRIVER<br/>OCS_STORAGE_USERS_DRIVER | string | ocis | |
| OCIS_MACHINE_AUTH_API_KEY<br/>OCS_MACHINE_AUTH_API_KEY | string |  | |