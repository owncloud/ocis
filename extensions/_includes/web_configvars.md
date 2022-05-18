## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>WEB_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>WEB_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>WEB_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>WEB_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>WEB_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>WEB_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>WEB_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>WEB_LOG_FILE | string |  | |
| WEB_DEBUG_ADDR | string | 127.0.0.1:9104 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| WEB_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| WEB_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| WEB_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| WEB_HTTP_ADDR | string | 127.0.0.1:9100 | |
| WEB_HTTP_ROOT | string | / | |
| WEB_CACHE_TTL | int | 604800 | |
| WEB_ASSET_PATH | string |  | |
| WEB_UI_CONFIG | string |  | |
| WEB_UI_PATH | string |  | |
| OCIS_URL<br/>WEB_UI_THEME_SERVER | string | https://localhost:9200 | |
| WEB_UI_THEME_PATH | string | /themes/owncloud/theme.json | |
| OCIS_URL<br/>WEB_UI_CONFIG_SERVER | string | https://localhost:9200 | |
|  | string |  | |
| WEB_UI_CONFIG_VERSION | string | 0.1.0 | |
| WEB_OIDC_METADATA_URL | string | https://localhost:9200/.well-known/openid-configuration | |
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>WEB_OIDC_AUTHORITY | string | https://localhost:9200 | |
| WEB_OIDC_CLIENT_ID | string | web | |
| WEB_OIDC_RESPONSE_TYPE | string | code | |
| WEB_OIDC_SCOPE | string | openid profile email | |