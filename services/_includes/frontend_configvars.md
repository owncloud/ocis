## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>FRONTEND_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>FRONTEND_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>FRONTEND_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>FRONTEND_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>FRONTEND_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>FRONTEND_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>FRONTEND_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>FRONTEND_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| FRONTEND_DEBUG_ADDR | string | 127.0.0.1:9141 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| FRONTEND_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| FRONTEND_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| FRONTEND_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| FRONTEND_HTTP_ADDR | string | 127.0.0.1:9140 | The bind address of the HTTP service.|
| FRONTEND_HTTP_PROTOCOL | string | tcp | The transport protocol of the HTTP service.|
| FRONTEND_HTTP_PREFIX | string |  | The Path prefix where the frontend can be accessed (defaults to /).|
| STORAGE_TRANSFER_SECRET | string |  | Transfer secret for signing file up- and download requests.|
| OCIS_JWT_SECRET<br/>FRONTEND_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| OCIS_MACHINE_AUTH_API_KEY<br/>FRONTEND_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used to validate internal requests necessary to access resources from other services.|
| FRONTEND_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the loading of user's group memberships from the reva access token.|
| FRONTEND_ENABLE_FAVORITES | bool | false | Enables the support for favorites in the frontend.|
| FRONTEND_ENABLE_PROJECT_SPACES | bool | true | Indicates to clients that project spaces are supposed to be made available.|
| FRONTEND_ENABLE_SHARE_JAIL | bool | true | Indicates to clients that the share jail is supposed to be used.|
| FRONTEND_UPLOAD_MAX_CHUNK_SIZE | int | 100000000 | Sets the max chunk sizes for uploads via the frontend.|
| FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE | string |  | Advise TUS to replace PATCH requests by POST requests.|
| FRONTEND_DEFAULT_UPLOAD_PROTOCOL | string | tus | The default upload protocol to use in the frontend (e.g. tus).|
| FRONTEND_ENABLE_RESHARING | bool | true | Enables the support for resharing in the frontend.|
| OCIS_URL<br/>FRONTEND_PUBLIC_URL | string | https://localhost:9200 | The public facing URL of the oCIS frontend.|
| OCIS_INSECURE<br/>FRONTEND_APP_HANDLER_INSECURE | bool | false | Allow insecure connections to the frontend.|
| FRONTEND_ARCHIVER_MAX_NUM_FILES | int64 | 10000 | Max number of files that can be packed into an archive.|
| FRONTEND_ARCHIVER_MAX_SIZE | int64 | 1073741824 | Max size of the zip archive the archiver can create.|
| OCIS_INSECURE<br/>FRONTEND_ARCHIVER_INSECURE | bool | false | Allow insecure connections to the archiver.|
| FRONTEND_DATA_GATEWAY_PREFIX | string | data | Path prefix for the data gateway.|
| FRONTEND_OCS_PREFIX | string | ocs | Path prefix for the OCS service|
| FRONTEND_OCS_SHARE_PREFIX | string | /Shares | Path prefix for shares.|
| FRONTEND_OCS_HOME_NAMESPACE | string | /users/{{.Id.OpaqueId}} | Homespace namespace identifier.|
| FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE | string | {{.Mail}} | Additional information attribute for the user like {{.Mail}}.|
| FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL | int | 0 | Max TTL for the resource info cache.|
| FRONTEND_CHECKSUMS_SUPPORTED_TYPES | []string | [sha1 md5 adler32] | Supported checksum types to be announced to the client. You can provide multiple types separated by blank or comma.|
| FRONTEND_CHECKSUMS_PREFERRED_UPLOAD_TYPES | string |  | Preferred checksum types to be announced to the client for uploads (e.g. md5)|