## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_LOG_LEVEL<br/>AUDIT_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>AUDIT_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>AUDIT_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>AUDIT_LOG_FILE | string |  | |
| AUDIT_DEBUG_ADDR | string |  | |
| AUDIT_DEBUG_TOKEN | string |  | |
| AUDIT_DEBUG_PPROF | bool | false | |
| AUDIT_DEBUG_ZPAGES | bool | false | |
| AUDIT_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| AUDIT_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|
| AUDIT_EVENTS_GROUP | string | audit | the customergroup of the service. One group will only get one vopy of an event|
| AUDIT_LOG_TO_CONSOLE | bool | true | logs to Stdout if true|
| AUDIT_LOG_TO_FILE | bool | false | logs to file if true|
| AUDIT_FILEPATH | string |  | filepath to the logfile. Mandatory if LogToFile is true|
| AUDIT_FORMAT | string | json | log format. using json is advised|