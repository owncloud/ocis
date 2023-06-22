## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>POSTPROCESSING_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>POSTPROCESSING_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>POSTPROCESSING_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>POSTPROCESSING_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>POSTPROCESSING_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>POSTPROCESSING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>POSTPROCESSING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>POSTPROCESSING_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| POSTPROCESSING_DEBUG_ADDR | string | 127.0.0.1:9255 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| POSTPROCESSING_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| POSTPROCESSING_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| POSTPROCESSING_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_PERSISTENT_STORE<br/>POSTPROCESSING_STORE | string | memory | The type of the store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details.|
| OCIS_PERSISTENT_STORE_NODES<br/>POSTPROCESSING_STORE_NODES | []string | [] | A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store.|
| POSTPROCESSING_STORE_DATABASE | string | postprocessing | The database name the configured store should use.|
| POSTPROCESSING_STORE_TABLE | string | postprocessing | The database table the store should use.|
| OCIS_PERSISTENT_STORE_TTL<br/>POSTPROCESSING_STORE_TTL | Duration | 0s | Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks).|
| OCIS_PERSISTENT_STORE_SIZE<br/>POSTPROCESSING_STORE_SIZE | int | 0 | The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512.|
| OCIS_EVENTS_ENDPOINT<br/>POSTPROCESSING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>POSTPROCESSING_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>POSTPROCESSING_EVENTS_TLS_INSECURE | bool | false | Whether the ocis server should skip the client certificate verification during the TLS handshake.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>POSTPROCESSING_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| POSTPROCESSING_STEPS | []string | [] | A comma separated list of postprocessing steps, processed in order of their appearance. Currently supported values by the system are: 'virusscan', 'policies' and 'delay'. Custom steps are allowed. See the documentation for instructions.|
| POSTPROCESSING_VIRUSSCAN | bool | false | After uploading a file but before making it available for download, virus scanning the file can be enabled. Needs as prerequisite the antivirus service to be enabled and configured.|
| POSTPROCESSING_DELAY | Duration | 0s | After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. The duration can be set as number followed by a unit identifier like s, m or h. If a duration is set but the keyword 'delay' is not explicitely added to 'POSTPROCESSING_STEPS', the delay step will be processed as last step. In such a case, a log entry will be written on service startup to remind the admin about that situation.|