## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>SEARCH_TRACING_ENABLED | bool | false | |
| OCIS_TRACING_TYPE<br/>SEARCH_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>SEARCH_TRACING_ENDPOINT | string |  | |
| OCIS_TRACING_COLLECTOR<br/>SEARCH_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>SEARCH_LOG_LEVEL | string |  | |
| OCIS_LOG_PRETTY<br/>SEARCH_LOG_PRETTY | bool | false | |
| OCIS_LOG_COLOR<br/>SEARCH_LOG_COLOR | bool | false | |
| OCIS_LOG_FILE<br/>SEARCH_LOG_FILE | string |  | |
| SEARCH_DEBUG_ADDR | string | 127.0.0.1:9224 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| SEARCH_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| SEARCH_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| SEARCH_DEBUG_ZPAGES | bool | false | Enables zpages, which can  be used for collecting and viewing traces in-me|
| SEARCH_GRPC_ADDR | string | 127.0.0.1:9220 | The address of the grpc service.|
| SEARCH_DATA_PATH | string | ~/.ocis/search | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| SEARCH_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| SEARCH_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|
| SEARCH_EVENTS_GROUP | string | search | the customergroup of the service. One group will only get one copy of an event|
| OCIS_MACHINE_AUTH_API_KEY<br/>SEARCH_MACHINE_AUTH_API_KEY | string |  | |