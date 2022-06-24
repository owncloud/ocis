## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>GATEWAY_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>GATEWAY_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>GATEWAY_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>GATEWAY_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>GATEWAY_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>GATEWAY_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>GATEWAY_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>GATEWAY_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| GATEWAY_DEBUG_ADDR | string | 127.0.0.1:9143 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| GATEWAY_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| GATEWAY_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| GATEWAY_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| GATEWAY_GRPC_ADDR | string | 127.0.0.1:9142 | The address of the grpc service.|
| GATEWAY_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>GATEWAY_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| GATEWAY_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the encoding of the user's groupmember ships in the reva access token. To reduces token size, especially when users are members of a large number of groups.|
| GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT | bool | true | Commit shares to storage grants (default: true).|
| GATEWAY_COMMIT_SHARE_TO_STORAGE_REF | bool | true | Commit shares to storage (default: true)|
| GATEWAY_SHARE_FOLDER_NAME | string | Shares | Name of the gateway share folder|
| GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN | bool | true | Disable creation of the homespace on login|
| STORAGE_TRANSFER_SECRET | string |  | The storage transfer secret|
| GATEWAY_TRANSFER_EXPIRES | int | 86400 | Expiry for the gateway tokens|
| GATEWAY_ETAG_CACHE_TTL | int | 0 | Max TTL for the gatways ETAG cache.|
| OCIS_URL<br/>GATEWAY_FRONTEND_PUBLIC_URL | string | https://localhost:9200 | The public facing url of the ocis frontend.|
| GATEWAY_USERS_ENDPOINT | string | localhost:9144 | The users api endpoint.|
| GATEWAY_GROUPS_ENDPOINT | string | localhost:9160 | The groups api endpoint.|
| GATEWAY_PERMISSIONS_ENDPOINT | string | localhost:9191 | The permission api endpoint.|
| GATEWAY_SHARING_ENDPOINT | string | localhost:9150 | The share api endpoint.|
| GATEWAY_AUTH_BASIC_ENDPOINT | string | localhost:9146 | The auth basic api endpoint.|
| GATEWAY_AUTH_BEARER_ENDPOINT | string | localhost:9148 | The auth bearer api endpoint.|
| GATEWAY_AUTH_MACHINE_ENDPOINT | string | localhost:9166 | The auth machine api endpoint.|
| GATEWAY_STORAGE_PUBLIC_LINK_ENDPOINT | string | localhost:9178 | The storage puliclink api endpoint.|
| GATEWAY_STORAGE_USERS_ENDPOINT | string | localhost:9157 | The storage users api endpoint.|
| GATEWAY_STORAGE_SHARES_ENDPOINT | string | localhost:9154 | The storage shares api endpoint.|
| GATEWAY_APP_REGISTRY_ENDPOINT | string | localhost:9242 | The app registry api endpoint.|