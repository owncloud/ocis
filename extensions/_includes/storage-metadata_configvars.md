## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED;STORAGE_METADATA_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE;STORAGE_METADATA_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT;STORAGE_METADATA_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR;STORAGE_METADATA_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL;STORAGE_METADATA_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY;STORAGE_METADATA_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR;STORAGE_METADATA_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE;STORAGE_METADATA_LOG_FILE | string |  | The target log file.|
| STORAGE_METADATA_DEBUG_ADDR | string | 127.0.0.1:9217 | |
| STORAGE_METADATA_DEBUG_TOKEN | string |  | |
| STORAGE_METADATA_DEBUG_PPROF | bool | false | |
| STORAGE_METADATA_DEBUG_ZPAGES | bool | false | |
| STORAGE_METADATA_GRPC_ADDR | string | 127.0.0.1:9215 | The address of the grpc service.|
| STORAGE_METADATA_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| STORAGE_METADATA_HTTP_ADDR | string | 127.0.0.1:9216 | The address of the http service.|
| STORAGE_METADATA_HTTP_PROTOCOL | string | tcp | The transport protocol of the http service.|
| OCIS_JWT_SECRET;STORAGE_METADATA_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| STORAGE_METADATA_MACHINE_AUTH_API_KEY | string |  | |
| STORAGE_METADATA_DRIVER | string | ocis | The driver which should be used by the service|
| STORAGE_METADATA_DRIVER_OCIS_ROOT | string | ~/.ocis/storage/metadata | |
| OCIS_INSECURE;STORAGE_METADATA_DATAPROVIDER_INSECURE | bool | false | |