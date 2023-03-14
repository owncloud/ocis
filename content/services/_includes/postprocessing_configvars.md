## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>POSTPROCESSING_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>POSTPROCESSING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>POSTPROCESSING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>POSTPROCESSING_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| POSTPROCESSING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| POSTPROCESSING_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>POSTPROCESSING_EVENTS_TLS_INSECURE | bool | false | Whether the ocis server should skip the client certificate verification during the TLS handshake.|
| POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>POSTPROCESSING_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| POSTPROCESSING_STEPS | []string | [] | A comma separated list of postprocessing steps, processed in order of their appearance. Currently supported values by the system are: 'virusscan', 'policies' and 'delay'. Custom steps are allowed. See the documentation for instructions.|
| POSTPROCESSING_VIRUSSCAN | bool | false | After uploading a file but before making it available for download, virus scanning the file can be enabled. Needs as prerequisite the antivirus service to be enabled and configured.|
| POSTPROCESSING_DELAY | Duration | 0s | After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. The duration can be set as number followed by a unit identifier like s, m or h. If a duration is set but the keyword 'delay' is not explicitely added to 'POSTPROCESSING_STEPS', the delay step will be processed as last step. In such a case, a log entry will be written on service startup to remind the admin about that situation.|