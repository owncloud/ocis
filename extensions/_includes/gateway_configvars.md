## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED;GATEWAY_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE;GATEWAY_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT;GATEWAY_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR;GATEWAY_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL;GATEWAY_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY;GATEWAY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR;GATEWAY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE;GATEWAY_LOG_FILE | string |  | The target log file.|
| GATEWAY_DEBUG_ADDR | string | 127.0.0.1:9143 | |
| GATEWAY_DEBUG_TOKEN | string |  | |
| GATEWAY_DEBUG_PPROF | bool | false | |
| GATEWAY_DEBUG_ZPAGES | bool | false | |
| GATEWAY_GRPC_ADDR | string | 127.0.0.1:9142 | The address of the grpc service.|
| GATEWAY_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET;GATEWAY_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| STORAGE_TRANSFER_SECRET | string |  | |
| OCIS_URL;GATEWAY_FRONTEND_PUBLIC_URL | string | https://localhost:9200 | |