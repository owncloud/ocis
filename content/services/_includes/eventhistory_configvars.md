## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>EVENTHISTORY_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>EVENTHISTORY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>EVENTHISTORY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>EVENTHISTORY_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| EVENTHISTORY_DEBUG_ADDR | string |  | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| EVENTHISTORY_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| EVENTHISTORY_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| EVENTHISTORY_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| EVENTHISTORY_GRPC_ADDR | string | 127.0.0.1:0 | The bind address of the GRPC service.|
| EVENTHISTORY_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| EVENTHISTORY_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>EVENTHISTORY_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| EVENTHISTORY_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>EVENTHISTORY_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services..|
| EVENTHISTORY_STORE_TYPE | string | mem | The type of the eventhistory store. Supported values are: 'mem', 'ocmem', 'etcd', 'redis', 'nats-js', 'noop'. See the text description for details.|
| EVENTHISTORY_STORE_ADDRESSES | string |  | A comma separated list of addresses to access the configured store. This has no effect when 'in-memory' stores are configured. Note that the behaviour how addresses are used is dependent on the library of the configured store.|
| EVENTHISTORY_STORE_DATABASE | string |  | (optional) The database name the configured store should use. This has no effect when 'in-memory' stores are configured.|
| EVENTHISTORY_STORE_TABLE | string |  | (optional) The database table the store should use. This has no effect when 'in-memory' stores are configured.|
| EVENTHISTORY_RECORD_EXPIRY | Duration | 336h0m0s | Time to life for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks).|
| EVENTHISTORY_STORE_SIZE | int | 0 | The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512.|