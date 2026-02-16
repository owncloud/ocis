## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>AUDIT_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>AUDIT_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>AUDIT_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>AUDIT_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>AUDIT_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>AUDIT_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>AUDIT_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>AUDIT_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| AUDIT_DEBUG_ADDR | string | 127.0.0.1:9229 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| AUDIT_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| AUDIT_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| AUDIT_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_EVENTS_ENDPOINT<br/>AUDIT_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>AUDIT_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>AUDIT_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>AUDIT_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided AUDIT_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>AUDIT_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>AUDIT_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>AUDIT_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| AUDIT_LOG_TO_CONSOLE | bool | true | Logs to stdout if set to 'true'. Independent of the LOG_TO_FILE option.|
| AUDIT_LOG_TO_FILE | bool | false | Logs to file if set to 'true'. Independent of the LOG_TO_CONSOLE option.|
| AUDIT_FILEPATH | string |  | Filepath of the logfile. Mandatory if LOG_TO_FILE is set to 'true'.|
| AUDIT_FORMAT | string | json | Log format. Supported values are '' (empty) and 'json'. Using 'json' is advised, '' (empty) renders the 'minimal' format. See the text description for more details.|