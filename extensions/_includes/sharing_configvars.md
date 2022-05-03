## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED;SHARING_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE;SHARING_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT;SHARING_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR;SHARING_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL;SHARING_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY;SHARING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR;SHARING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE;SHARING_LOG_FILE | string |  | The target log file.|
| SHARING_DEBUG_ADDR | string | 127.0.0.1:9151 | |
| SHARING_DEBUG_TOKEN | string |  | |
| SHARING_DEBUG_PPROF | bool | false | |
| SHARING_DEBUG_ZPAGES | bool | false | |
| SHARING_GRPC_ADDR | string | 127.0.0.1:9150 | The address of the grpc service.|
| SHARING_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET;SHARING_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| SHARING_USER_JSON_FILE | string | ~/.ocis/storage/shares.json | |
| SHARING_USER_SQL_USERNAME | string |  | |
| SHARING_USER_SQL_PASSWORD | string |  | |
| SHARING_USER_SQL_HOST | string |  | |
| SHARING_USER_SQL_PORT | int | 1433 | |
| SHARING_USER_SQL_NAME | string |  | |
| OCIS_URL;SHARING_CS3_SERVICE_USER_IDP | string | internal | |
| OCIS_MACHINE_AUTH_API_KEY | string |  | |
| OCIS_MACHINE_AUTH_API_KEY | string |  | |