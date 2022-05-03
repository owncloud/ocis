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
| SEARCH_DEBUG_ADDR | string | 127.0.0.1:9224 | |
| SEARCH_DEBUG_TOKEN | string |  | |
| SEARCH_DEBUG_PPROF | bool | false | |
| SEARCH_DEBUG_ZPAGES | bool | false | |
| ACCOUNTS_GRPC_ADDR | string | 127.0.0.1:9220 | The address of the grpc service.|
| SEARCH_DATA_PATH | string | ~/.ocis/search | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| SEARCH_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| SEARCH_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|
| SEARCH_EVENTS_GROUP | string | search | the customergroup of the service. One group will only get one copy of an event|
| OCIS_MACHINE_AUTH_API_KEY<br/>SEARCH_MACHINE_AUTH_API_KEY | string | change-me-please | |