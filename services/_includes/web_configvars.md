## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>WEB_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>WEB_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>WEB_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>WEB_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>WEB_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>WEB_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>WEB_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>WEB_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| WEB_DEBUG_ADDR | string | 127.0.0.1:9104 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| WEB_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| WEB_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| WEB_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| WEB_HTTP_ADDR | string | 127.0.0.1:9100 | The bind address of the HTTP service.|
| WEB_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| WEB_CACHE_TTL | int | 604800 | Cache policy in seconds for ownCloud Web assets.|
| WEB_ASSET_PATH | string |  | Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets.|
| WEB_UI_CONFIG | string |  | Read the ownCloud Web configuration from this file.|
| WEB_UI_PATH | string |  | Read the ownCloud Web configuration from this file path.|
| OCIS_URL<br/>WEB_UI_THEME_SERVER | string | https://localhost:9200 | URL to load themes from. Will be prepended to the theme path.|
| WEB_UI_THEME_PATH | string | /themes/owncloud/theme.json | URL path to load themes from. The theme server will be prepended.|
| OCIS_URL<br/>WEB_UI_CONFIG_SERVER | string | https://localhost:9200 | URL, where the oCIS APIs are reachable for ownCloud Web.|
| WEB_OIDC_METADATA_URL | string | https://localhost:9200/.well-known/openid-configuration | URL for the OIDC well-known configuration endpoint. Defaults to the oCIS API URL + "/.well-known/openid-configuration".|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>WEB_OIDC_AUTHORITY | string | https://localhost:9200 | URL of the OIDC issuer. It defaults to URL of the builtin IDP.|
| WEB_OIDC_CLIENT_ID | string | web | OIDC client ID, which ownCloud Web uses. This client needs to be set up in your IDP.|
| WEB_OIDC_RESPONSE_TYPE | string | code | OIDC response type to use for authentication.|
| WEB_OIDC_SCOPE | string | openid profile email | OIDC scopes to request during authentication.|