## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>POSTPROCESSING_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>POSTPROCESSING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>POSTPROCESSING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>POSTPROCESSING_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| POSTPROCESSING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | Endpoint of the event system.|
| POSTPROCESSING_EVENTS_CLUSTER | string | ocis-cluster | Cluster ID of the event system.|
| OCIS_INSECURE<br/>SEARCH_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| SEARCH_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided SEARCH_EVENTS_TLS_INSECURE will be seen as false.|
| POSTPROCESSING_VIRUSSCAN | bool | false | should the system do a virusscan? Needs antivirus service|
| POSTPROCESSING_DELAY | Duration | 0s | the sytem sleeps for this time while postprocessing|