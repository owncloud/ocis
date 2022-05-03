## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED;APP_PROVIDER_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE;APP_PROVIDER_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT;APP_PROVIDER_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR;APP_PROVIDER_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL;APP_PROVIDER_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY;APP_PROVIDER_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR;APP_PROVIDER_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE;APP_PROVIDER_LOG_FILE | string |  | The target log file.|
| APP_PROVIDER_DEBUG_ADDR | string | 127.0.0.1:9165 | |
| APP_PROVIDER_DEBUG_TOKEN | string |  | |
| APP_PROVIDER_DEBUG_PPROF | bool | false | |
| APP_PROVIDER_DEBUG_ZPAGES | bool | false | |
| APP_PROVIDER_GRPC_ADDR | string | 127.0.0.1:9164 | The address of the grpc service.|
| APP_PROVIDER_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET;APP_PROVIDER_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |