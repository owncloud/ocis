## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>USERLOG_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>USERLOG_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>USERLOG_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>USERLOG_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>USERLOG_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>USERLOG_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>USERLOG_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>USERLOG_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| USERLOG_DEBUG_ADDR | string | 127.0.0.1:9214 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| USERLOG_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| USERLOG_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| USERLOG_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| USERLOG_HTTP_ADDR | string | 127.0.0.1:9210 | The bind address of the HTTP service.|
| USERLOG_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>USERLOG_CORS_ALLOW_ORIGINS | []string | [*] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>USERLOG_CORS_ALLOW_METHODS | []string | [GET] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>USERLOG_CORS_ALLOW_HEADERS | []string | [Authorization Origin Content-Type Accept X-Requested-With X-Request-Id Ocs-Apirequest] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>USERLOG_CORS_ALLOW_CREDENTIALS | bool | true | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| OCIS_JWT_SECRET<br/>USERLOG_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | CS3 gateway used to look up user metadata|
| OCIS_TRANSLATION_PATH<br/>USERLOG_TRANSLATION_PATH | string |  | (optional) Set this to a path with custom translations to overwrite the builtin translations. Note that file and folder naming rules apply, see the documentation for more details.|
| OCIS_DEFAULT_LANGUAGE | string |  | The default language used by services and the WebUI. If not defined, English will be used as default. See the documentation for more details.|
| OCIS_EVENTS_ENDPOINT<br/>USERLOG_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>USERLOG_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>USERLOG_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>USERLOG_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>USERLOG_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>USERLOG_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>USERLOG_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_MAX_CONCURRENCY<br/>USERLOG_MAX_CONCURRENCY | int | 1 | Maximum number of concurrent go-routines. Higher values can potentially get work done faster but will also cause more load on the system. Values of 0 or below will be ignored and the default value will be used.|
| OCIS_PERSISTENT_STORE<br/>USERLOG_STORE | string | memory | The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>USERLOG_STORE_NODES | []string | [] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| USERLOG_STORE_DATABASE | string | userlog | The database name the configured store should use.|
| USERLOG_STORE_TABLE | string | events | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>USERLOG_STORE_TTL | Duration | 336h0m0s | Time to live for events in the store. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>USERLOG_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>USERLOG_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_DISABLE_SSE,USERLOG_DISABLE_SSE | bool | false | Disables server-sent events (sse). When disabled, clients will no longer receive sse notifications.|
| USERLOG_GLOBAL_NOTIFICATIONS_SECRET | string |  | The secret to secure the global notifications endpoint. Only system admins and users knowing that secret can call the global notifications POST/DELETE endpoints.|
| OCIS_SERVICE_ACCOUNT_ID<br/>USERLOG_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>USERLOG_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|