## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORAGE_SHARES_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORAGE_SHARES_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>STORAGE_SHARES_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>STORAGE_SHARES_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>STORAGE_SHARES_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>STORAGE_SHARES_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORAGE_SHARES_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORAGE_SHARES_LOG_FILE | string |  | The target log file.|
| STORAGE_SHARES_DEBUG_ADDR | string | 127.0.0.1:9156 | |
| STORAGE_SHARES_DEBUG_TOKEN | string |  | |
| STORAGE_SHARES_DEBUG_PPROF | bool | false | |
| STORAGE_SHARES_DEBUG_ZPAGES | bool | false | |
| STORAGE_SHARES_GRPC_ADDR | string | 127.0.0.1:9154 | The address of the grpc service.|
| STORAGE_SHARES_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>STORAGE_SHARES_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| STORAGE_SHARES_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| STORAGE_SHARES_READ_ONLY | bool | false | |
| STORAGE_SHARES_USER_SHARE_PROVIDER_ENDPOINT | string | localhost:9150 | |