## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>AUDIT_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>AUDIT_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>AUDIT_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>AUDIT_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| AUDIT_DEBUG_ADDR | string | 127.0.0.1:9234 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| AUDIT_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| AUDIT_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| AUDIT_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| AUDIT_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the streaming service.|
| AUDIT_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the streaming service. Mandatory when using nats.|
| AUDIT_EVENTS_GROUP | string | audit | The consumergroup of the service. One group will only get one copy of an event.|
| AUDIT_LOG_TO_CONSOLE | bool | true | Logs to Stdout if true. Independent of the log to file option.|
| AUDIT_LOG_TO_FILE | bool | false | Logs to file if true. Independent of the log to Stdout file option.|
| AUDIT_FILEPATH | string |  | Filepath to the logfile. Mandatory if LogToFile is true.|
| AUDIT_FORMAT | string | json | Log format. Using json is advised.|