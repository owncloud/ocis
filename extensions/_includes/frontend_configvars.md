## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED;FRONTEND_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE;FRONTEND_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT;FRONTEND_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR;FRONTEND_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL;FRONTEND_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY;FRONTEND_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR;FRONTEND_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE;FRONTEND_LOG_FILE | string |  | The target log file.|
| FRONTEND_DEBUG_ADDR | string | 127.0.0.1:9141 | |
| FRONTEND_DEBUG_TOKEN | string |  | |
| FRONTEND_DEBUG_PPROF | bool | false | |
| FRONTEND_DEBUG_ZPAGES | bool | false | |
| FRONTEND_HTTP_ADDR | string | 127.0.0.1:9140 | The address of the http service.|
| FRONTEND_HTTP_PROTOCOL | string | tcp | The transport protocol of the http service.|
| STORAGE_TRANSFER_SECRET | string |  | |
| OCIS_JWT_SECRET;FRONTEND_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| OCIS_MACHINE_AUTH_API_KEY;FRONTEND_MACHINE_AUTH_API_KEY | string |  | |
| FRONTEND_ENABLE_PROJECT_SPACES | bool | true | Indicates to clients that project spaces are supposed to be made available.|
| FRONTEND_ENABLE_SHARE_JAIL | bool | true | Indicates to clients that the share jail is supposed to be used.|
| OCIS_URL;FRONTEND_PUBLIC_URL | string | https://localhost:9200 | |
| OCIS_INSECURE;FRONTEND_ARCHIVER_INSECURE | bool | false | |
| OCIS_INSECURE;FRONTEND_APPPROVIDER_INSECURE | bool | false | |