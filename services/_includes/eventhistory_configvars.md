## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>EVENTHISTORY_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>EVENTHISTORY_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>EVENTHISTORY_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>EVENTHISTORY_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>EVENTHISTORY_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>EVENTHISTORY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>EVENTHISTORY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>EVENTHISTORY_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| EVENTHISTORY_DEBUG_ADDR | string | 127.0.0.1:9270 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| EVENTHISTORY_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| EVENTHISTORY_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| EVENTHISTORY_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| EVENTHISTORY_GRPC_ADDR | string | 127.0.0.1:9274 | The bind address of the GRPC service.|
| OCIS_EVENTS_ENDPOINT<br/>EVENTHISTORY_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>EVENTHISTORY_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>EVENTHISTORY_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>EVENTHISTORY_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. Will be seen as empty if NOTIFICATIONS_EVENTS_TLS_INSECURE is provided.|
| OCIS_EVENTS_ENABLE_TLS<br/>EVENTHISTORY_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>EVENTHISTORY_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>EVENTHISTORY_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_PERSISTENT_STORE<br/>EVENTHISTORY_STORE | string | nats-js-kv | The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>EVENTHISTORY_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| EVENTHISTORY_STORE_DATABASE | string | eventhistory | The database name the configured store should use.|
| EVENTHISTORY_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>EVENTHISTORY_STORE_TTL | Duration | 336h0m0s | Time to live for events in the store. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>EVENTHISTORY_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>EVENTHISTORY_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|