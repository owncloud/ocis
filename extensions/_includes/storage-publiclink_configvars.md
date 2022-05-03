## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORAGE_PUBLICLINK_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORAGE_PUBLICLINK_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>STORAGE_PUBLICLINK_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>STORAGE_PUBLICLINK_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>STORAGE_PUBLICLINK_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>STORAGE_PUBLICLINK_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORAGE_PUBLICLINK_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORAGE_PUBLICLINK_LOG_FILE | string |  | The target log file.|
| STORAGE_PUBLICLINK_DEBUG_ADDR | string | 127.0.0.1:9179 | |
| STORAGE_PUBLICLINK_DEBUG_TOKEN | string |  | |
| STORAGE_PUBLICLINK_DEBUG_PPROF | bool | false | |
| STORAGE_PUBLICLINK_DEBUG_ZPAGES | bool | false | |
| STORAGE_PUBLICLINK_GRPC_ADDR | string | 127.0.0.1:9178 | The address of the grpc service.|
| STORAGE_PUBLICLINK_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>STORAGE_PUBLICLINK_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |