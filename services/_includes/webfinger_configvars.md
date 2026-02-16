## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>WEBFINGER_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>WEBFINGER_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>WEBFINGER_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>WEBFINGER_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>WEBFINGER_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>WEBFINGER_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>WEBFINGER_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>WEBFINGER_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| WEBFINGER_DEBUG_ADDR | string | 127.0.0.1:9279 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| WEBFINGER_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| WEBFINGER_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| WEBFINGER_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| WEBFINGER_HTTP_ADDR | string | 127.0.0.1:9275 | The bind address of the HTTP service.|
| WEBFINGER_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>WEBFINGER_CORS_ALLOW_ORIGINS | []string | [https://localhost:9200] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>WEBFINGER_CORS_ALLOW_METHODS | []string | [] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>WEBFINGER_CORS_ALLOW_HEADERS | []string | [] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>WEBFINGER_CORS_ALLOW_CREDENTIALS | bool | false | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| WEBFINGER_RELATIONS | []string | [http://openid.net/specs/connect/1.0/issuer http://webfinger.owncloud/rel/server-instance] | A list of relation URIs or registered relation types to add to webfinger responses. See the Environment Variable Types description for more details.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>WEBFINGER_OIDC_ISSUER | string | https://localhost:9200 | The identity provider href for the openid-discovery relation.|
| OCIS_URL<br/>WEBFINGER_OWNCLOUD_SERVER_INSTANCE_URL | string | https://localhost:9200 | The URL for the legacy ownCloud server instance relation (not to be confused with the product ownCloud Server). It defaults to the OCIS_URL but can be overridden to support some reverse proxy corner cases. To shard the deployment, multiple instances can be configured in the configuration file.|
| OCIS_INSECURE<br/>WEBFINGER_INSECURE | bool | false | Allow insecure connections to the WEBFINGER service.|