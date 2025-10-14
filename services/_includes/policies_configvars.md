## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| POLICIES_GRPC_ADDR | string | 127.0.0.1:9125 | The bind address of the GRPC service.|
| POLICIES_DEBUG_ADDR | string | 127.0.0.1:9129 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| POLICIES_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| POLICIES_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| POLICIES_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| OCIS_EVENTS_ENDPOINT<br/>POLICIES_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>POLICIES_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>POLICIES_EVENTS_TLS_INSECURE | bool | false | Whether the server should skip the client certificate verification during the TLS handshake.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>POLICIES_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided POLICIES_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>POLICIES_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_USERNAME<br/>POLICIES_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>POLICIES_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_LOG_LEVEL<br/>POLICIES_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>POLICIES_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>POLICIES_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>POLICIES_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| POLICIES_ENGINE_TIMEOUT | Duration | 10s | Sets the timeout the rego expression evaluation can take. Rules default to deny if the timeout was reached. See the Environment Variable Types description for more details.|
| POLICIES_ENGINE_MIMES | string |  | Sets the mimes file path which maps mimetypes to associated file extensions. See the text description for details.|
| POLICIES_POSTPROCESSING_QUERY | string |  | Defines the 'Complete Rules' variable defined in the rego rule set this step uses for its evaluation. Defaults to deny if the variable was not found.|
| OCIS_TRACING_ENABLED<br/>POLICIES_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>POLICIES_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>POLICIES_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>POLICIES_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|