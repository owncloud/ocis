## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>POSTPROCESSING_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>POSTPROCESSING_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>POSTPROCESSING_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>POSTPROCESSING_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>POSTPROCESSING_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>POSTPROCESSING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>POSTPROCESSING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>POSTPROCESSING_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| POSTPROCESSING_DEBUG_ADDR | string | 127.0.0.1:9255 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| POSTPROCESSING_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| POSTPROCESSING_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| POSTPROCESSING_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_PERSISTENT_STORE<br/>POSTPROCESSING_STORE | string | nats-js-kv | The type of the store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>POSTPROCESSING_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| POSTPROCESSING_STORE_DATABASE | string | postprocessing | The database name the configured store should use.|
| POSTPROCESSING_STORE_TABLE | string |  | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>POSTPROCESSING_STORE_TTL | Duration | 0s | Time to live for events in the store. See the Environment Variable Types description for more details.|
| OCIS_PERSISTENT_STORE_AUTH_USERNAME<br/>POSTPROCESSING_STORE_AUTH_USERNAME | string |  | The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_PERSISTENT_STORE_AUTH_PASSWORD<br/>POSTPROCESSING_STORE_AUTH_PASSWORD | string |  | The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_EVENTS_ENDPOINT<br/>POSTPROCESSING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>POSTPROCESSING_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>POSTPROCESSING_EVENTS_TLS_INSECURE | bool | false | Whether the ocis server should skip the client certificate verification during the TLS handshake.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>POSTPROCESSING_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>POSTPROCESSING_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>POSTPROCESSING_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| POSTPROCESSING_WORKERS | int | 3 | The number of concurrent go routines that fetch events from the event queue.|
| POSTPROCESSING_STEPS | []string | [] | A list of postprocessing steps processed in order of their appearance. Currently supported values by the system are: 'virusscan', 'policies' and 'delay'. Custom steps are allowed. See the documentation for instructions. See the Environment Variable Types description for more details.|
| POSTPROCESSING_DELAY | Duration | 0s | After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. If a duration is set but the keyword 'delay' is not explicitely added to 'POSTPROCESSING_STEPS', the delay step will be processed as last step. In such a case, a log entry will be written on service startup to remind the admin about that situation. See the Environment Variable Types description for more details.|
| POSTPROCESSING_RETRY_BACKOFF_DURATION | Duration | 5s | The base for the exponential backoff duration before retrying a failed postprocessing step. See the Environment Variable Types description for more details.|
| POSTPROCESSING_MAX_RETRIES | int | 14 | The maximum number of retries for a failed postprocessing step.|