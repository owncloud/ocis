## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>SHARING_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>SHARING_TRACING_TYPE | string |  | |
| OCIS_TRACING_ENDPOINT<br/>SHARING_TRACING_ENDPOINT | string |  | The endpoint to the tracing collector.|
| OCIS_TRACING_COLLECTOR<br/>SHARING_TRACING_COLLECTOR | string |  | |
| OCIS_LOG_LEVEL<br/>SHARING_LOG_LEVEL | string |  | The log level.|
| OCIS_LOG_PRETTY<br/>SHARING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>SHARING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>SHARING_LOG_FILE | string |  | The target log file.|
| SHARING_DEBUG_ADDR | string | 127.0.0.1:9151 | |
| SHARING_DEBUG_TOKEN | string |  | |
| SHARING_DEBUG_PPROF | bool | false | |
| SHARING_DEBUG_ZPAGES | bool | false | |
| SHARING_GRPC_ADDR | string | 127.0.0.1:9150 | The address of the grpc service.|
| SHARING_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>SHARING_JWT_SECRET | string |  | |
| REVA_GATEWAY | string | 127.0.0.1:9142 | |
| SHARING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| SHARING_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|
| SHARING_SKIP_USER_GROUPS_IN_TOKEN | bool | false | |
| SHARING_USER_DRIVER | string | json | |
| SHARING_USER_JSON_FILE | string | ~/.ocis/storage/shares.json | |
| SHARING_USER_CS3_PROVIDER_ADDR | string | 127.0.0.1:9215 | |
| SHARING_USER_CS3_SERVICE_USER_ID | string |  | |
| OCIS_URL<br/>SHARING_USER_CS3_SERVICE_USER_IDP | string | internal | |
| OCIS_MACHINE_AUTH_API_KEY<br/>SHARING_USER_CS3_MACHINE_AUTH_API_KEY | string |  | |
| SHARING_USER_OWNCLOUDSQL_DB_USERNAME | string |  | |
| SHARING_USER_OWNCLOUDSQL_DB_PASSWORD | string |  | |
| SHARING_USER_OWNCLOUDSQL_DB_HOST | string |  | |
| SHARING_USER_OWNCLOUDSQL_DB_PORT | int | 0 | |
| SHARING_USER_OWNCLOUDSQL_DB_NAME | string |  | |
| SHARING_USER_OWNCLOUDSQL_USER_STORAGE_MOUNT_ID | string |  | |
| SHARING_PUBLIC_DRIVER | string | json | |
| SHARING_PUBLIC_JSON_FILE | string | ~/.ocis/storage/publicshares.json | |
| SHARING_PUBLIC_CS3_PROVIDER_ADDR | string | 127.0.0.1:9215 | |
| SHARING_PUBLIC_CS3_SERVICE_USER_ID | string |  | |
| OCIS_URL<br/>SHARING_PUBLIC_CS3_SERVICE_USER_IDP | string | internal | |
| OCIS_MACHINE_AUTH_API_KEY<br/>SHARING_PUBLIC_CS3_MACHINE_AUTH_API_KEY | string |  | |