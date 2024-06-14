## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>CLIENTLOG_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>CLIENTLOG_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>CLIENTLOG_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>CLIENTLOG_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>CLIENTLOG_USERLOG_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>CLIENTLOG_USERLOG_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>CLIENTLOG_USERLOG_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>CLIENTLOG_USERLOG_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| CLIENTLOG_DEBUG_ADDR | string | 127.0.0.1:9197 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| CLIENTLOG_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| CLIENTLOG_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| CLIENTLOG_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_PERSISTENT_STORE<br/>ACTIVITYLOG_STORE | string | nats-js-kv | The type of the store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>ACTIVITYLOG_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| ACTIVITYLOG_STORE_DATABASE | string | activitylog | The database name the configured store should use.|
| ACTIVITYLOG_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>ACTIVITYLOG_STORE_TTL | Duration | 0s | Time to live for events in the store. See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_SIZE<br/>ACTIVITYLOG_STORE_SIZE | int | 0 | The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitly set as default.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>ACTIVITYLOG_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>ACTIVITYLOG_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | CS3 gateway used to look up user metadata|
| OCIS_SERVICE_ACCOUNT_ID<br/>ACTIVITYLOG_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>ACTIVITYOG_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|