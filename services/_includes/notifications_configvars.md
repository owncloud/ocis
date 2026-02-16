## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>NOTIFICATIONS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>NOTIFICATIONS_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>NOTIFICATIONS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>NOTIFICATIONS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>NOTIFICATIONS_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>NOTIFICATIONS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>NOTIFICATIONS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>NOTIFICATIONS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| NOTIFICATIONS_DEBUG_ADDR | string | 127.0.0.1:9174 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| NOTIFICATIONS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| NOTIFICATIONS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| NOTIFICATIONS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_URL<br/>NOTIFICATIONS_WEB_UI_URL | string | https://localhost:9200 | The public facing URL of the oCIS Web UI, used e.g. when sending notification eMails|
| NOTIFICATIONS_SMTP_HOST | string |  | SMTP host to connect to.|
| NOTIFICATIONS_SMTP_PORT | int | 0 | Port of the SMTP host to connect to.|
| NOTIFICATIONS_SMTP_SENDER | string |  | Sender address of emails that will be sent (e.g. 'ownCloud <noreply@example.com>'.|
| NOTIFICATIONS_SMTP_USERNAME | string |  | Username for the SMTP host to connect to.|
| NOTIFICATIONS_SMTP_PASSWORD | string |  | Password for the SMTP host to connect to.|
| NOTIFICATIONS_SMTP_INSECURE | bool | false | Allow insecure connections to the SMTP server.|
| NOTIFICATIONS_SMTP_AUTHENTICATION | string |  | Authentication method for the SMTP communication. Possible values are 'login', 'plain', 'crammd5', 'none' or 'auto'. If set to 'auto' or unset, the authentication method is automatically negotiated with the server.|
| NOTIFICATIONS_SMTP_ENCRYPTION | string | none | Encryption method for the SMTP communication. Possible values are 'starttls', 'ssltls' and 'none'.|
| OCIS_EVENTS_ENDPOINT<br/>NOTIFICATIONS_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>NOTIFICATIONS_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>NOTIFICATIONS_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>NOTIFICATIONS_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>NOTIFICATIONS_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>NOTIFICATIONS_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>NOTIFICATIONS_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EMAIL_TEMPLATE_PATH<br/>NOTIFICATIONS_EMAIL_TEMPLATE_PATH | string |  | Path to Email notification templates overriding embedded ones.|
| OCIS_TRANSLATION_PATH<br/>NOTIFICATIONS_TRANSLATION_PATH | string |  | (optional) Set this to a path with custom translations to overwrite the builtin translations. Note that file and folder naming rules apply, see the documentation for more details.|
| OCIS_DEFAULT_LANGUAGE | string |  | The default language used by services and the WebUI. If not defined, English will be used as default. See the documentation for more details.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | CS3 gateway used to look up user metadata|
| OCIS_GRPC_CLIENT_TLS_MODE | string |  | TLS mode for grpc connection to the go-micro based grpc services. Possible values are 'off', 'insecure' and 'on'. 'off': disables transport security for the clients. 'insecure' allows using transport security, but disables certificate verification (to be used with the autogenerated self-signed certificates). 'on' enables transport security, including server certificate verification.|
| OCIS_GRPC_CLIENT_TLS_CACERT | string |  | Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the go-micro based grpc services.|
| OCIS_SERVICE_ACCOUNT_ID<br/>NOTIFICATIONS_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>NOTIFICATIONS_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|
| OCIS_PERSISTENT_STORE<br/>NOTIFICATIONS_STORE | string | nats-js-kv | The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>NOTIFICATIONS_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| NOTIFICATIONS_STORE_DATABASE | string | notifications | The database name the configured store should use.|
| NOTIFICATIONS_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>NOTIFICATIONS_STORE_TTL | Duration | 336h0m0s | Time to live for notifications in the store. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>NOTIFICATIONS_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>NOTIFICATIONS_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|