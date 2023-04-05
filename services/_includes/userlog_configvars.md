## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>USERLOG_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>USERLOG_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>USERLOG_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>USERLOG_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| USERLOG_DEBUG_ADDR | string |  | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| USERLOG_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| USERLOG_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| USERLOG_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| USERLOG_HTTP_ADDR | string | 127.0.0.1:0 | The bind address of the HTTP service.|
| USERLOG_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>USERLOG_CORS_ALLOW_ORIGINS | []string | [*] | A comma-separated list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin|
| OCIS_CORS_ALLOW_METHODS<br/>USERLOG_CORS_ALLOW_METHODS | []string | [GET] | A comma-separated list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method|
| OCIS_CORS_ALLOW_HEADERS<br/>USERLOG_CORS_ALLOW_HEADERS | []string | [Authorization Origin Content-Type Accept X-Requested-With] | A comma-separated list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>USERLOG_CORS_ALLOW_CREDENTIALS | bool | true | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| OCIS_JWT_SECRET<br/>USERLOG_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| OCIS_MACHINE_AUTH_API_KEY<br/>USERLOG_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used to validate internal requests necessary to access resources from other services.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | CS3 gateway used to look up user metadata|
| USERLOG_TRANSLATION_PATH | string |  | (optional) Set this to a path with custom translations to overwrite the builtin translations. See the documentation for more details.|
| OCIS_EVENTS_ENDPOINT<br/>USERLOG_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>USERLOG_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>USERLOG_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| USERLOG_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>USERLOG_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services..|
| OCIS_PERSISTENT_STORE<br/>USERLOG_STORE<br/>USERLOG_STORE_TYPE | string | memory | The type of the userlog store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>USERLOG_STORE_NODES<br/>USERLOG_STORE_ADDRESSES | []string | [] | A comma separated list of nodes to access the configured store. This has no effect when 'in-memory' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store.|
| USERLOG_STORE_DATABASE | string | userlog | The database name the configured store should use.|
| USERLOG_STORE_TABLE | string | events | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>USERLOG_STORE_TTL | Duration | 336h0m0s | Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks).|
| OCIS_PERSISTENT_STORE_SIZE<br/>USERLOG_STORE_SIZE | int | 0 | The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512.|