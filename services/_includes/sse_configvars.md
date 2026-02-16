## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>SSE_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>SSE_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>SSE_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>SSE_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| SSE_DEBUG_ADDR | string | 127.0.0.1:9139 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| SSE_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| SSE_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| SSE_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_TRACING_ENABLED<br/>SSE_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>SSE_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>SSE_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>SSE_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| SSE_KEEPALIVE_INTERVAL | Duration | 0s | To prevent intermediate proxies from closing the SSE connection, send periodic SSE comments to keep it open.|
| OCIS_EVENTS_ENDPOINT<br/>SSE_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>SSE_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>SSE_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>SSE_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided SSE_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>SSE_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>SSE_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>SSE_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| SSE_HTTP_ADDR | string | 127.0.0.1:9135 | The bind address of the HTTP service.|
| SSE_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>SSE_CORS_ALLOW_ORIGINS | []string | [*] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>SSE_CORS_ALLOW_METHODS | []string | [GET] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>SSE_CORS_ALLOW_HEADERS | []string | [Authorization Origin Content-Type Accept X-Requested-With X-Request-Id Ocs-Apirequest] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>SSE_CORS_ALLOW_CREDENTIALS | bool | true | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| OCIS_JWT_SECRET<br/>SSE_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|