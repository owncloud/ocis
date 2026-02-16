## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>INVITATIONS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>INVITATIONS_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>INVITATIONS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>INVITATIONS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>INVITATIONS_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>INVITATIONS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>INVITATIONS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>INVITATIONS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| INVITATIONS_DEBUG_ADDR | string | 127.0.0.1:9269 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| INVITATIONS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| INVITATIONS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| INVITATIONS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| INVITATIONS_HTTP_ADDR | string | 127.0.0.1:9265 | The bind address of the HTTP service.|
| INVITATIONS_HTTP_ROOT | string | /graph/v1.0 | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>INVITATIONS_CORS_ALLOW_ORIGINS | []string | [https://localhost:9200] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>INVITATIONS_CORS_ALLOW_METHODS | []string | [] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>INVITATIONS_CORS_ALLOW_HEADERS | []string | [] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>INVITATIONS_CORS_ALLOW_CREDENTIALS | bool | false | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| OCIS_KEYCLOAK_BASE_PATH<br/>INVITATIONS_KEYCLOAK_BASE_PATH | string |  | The URL to access keycloak.|
| OCIS_KEYCLOAK_CLIENT_ID<br/>INVITATIONS_KEYCLOAK_CLIENT_ID | string |  | The client ID to authenticate with keycloak.|
| OCIS_KEYCLOAK_CLIENT_SECRET<br/>INVITATIONS_KEYCLOAK_CLIENT_SECRET | string |  | The client secret to use in authentication.|
| OCIS_KEYCLOAK_CLIENT_REALM<br/>INVITATIONS_KEYCLOAK_CLIENT_REALM | string |  | The realm the client is defined in.|
| OCIS_KEYCLOAK_USER_REALM<br/>INVITATIONS_KEYCLOAK_USER_REALM | string |  | The realm users are defined.|
| OCIS_KEYCLOAK_INSECURE_SKIP_VERIFY<br/>INVITATIONS_KEYCLOAK_INSECURE_SKIP_VERIFY | bool | false | Disable TLS certificate validation for Keycloak connections. Do not set this in production environments.|
| OCIS_JWT_SECRET<br/>INVITATIONS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|