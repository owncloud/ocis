## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>SHARING_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>SHARING_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>SHARING_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>SHARING_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>SHARING_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>SHARING_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>SHARING_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>SHARING_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| SHARING_DEBUG_ADDR | string | 127.0.0.1:9151 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| SHARING_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint|
| SHARING_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling|
| SHARING_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| SHARING_GRPC_ADDR | string | 127.0.0.1:9150 | The bind address of the GRPC service.|
| SHARING_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| OCIS_JWT_SECRET<br/>SHARING_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| SHARING_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | the address of the streaming service|
| SHARING_EVENTS_CLUSTER | string | ocis-cluster | the clusterID of the streaming service. Mandatory when using nats|
| SHARING_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the loading of user's group memberships from the reva access token.|
| SHARING_USER_DRIVER | string | cs3 | Driver to be used to persist shares. Possible values: "json", "cs3", "owncloudsql"|
| SHARING_USER_JSON_FILE | string | ~/.ocis/storage/shares.json | Path to the json file where shares will be persisted.|
| SHARING_USER_CS3_PROVIDER_ADDR | string | 127.0.0.1:9215 | GRPC address of the storage-system extension.|
| OCIS_SYSTEM_USER_ID<br/>SHARING_USER_CS3_SYSTEM_USER_ID | string |  | ID of the oCIS storage-system system user. Admins need to set the ID for the storage-system system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format.|
| OCIS_SYSTEM_USER_IDP<br/>SHARING_USER_CS3_SYSTEM_USER_IDP | string | internal | IDP of the oCIS storage-system system user.|
| OCIS_SYSTEM_USER_API_KEY<br/>SHARING_USER_CS3_SYSTEM_USER_API_KEY | string |  | API key for the storage-system system user.|
| SHARING_USER_OWNCLOUDSQL_DB_USERNAME | string | owncloud | Username for the database.|
| SHARING_USER_OWNCLOUDSQL_DB_PASSWORD | string |  | Password for the database.|
| SHARING_USER_OWNCLOUDSQL_DB_HOST | string | mysql | Hostname or IP of the database server.|
| SHARING_USER_OWNCLOUDSQL_DB_PORT | int | 3306 | Port that the database server is listening on.|
| SHARING_USER_OWNCLOUDSQL_DB_NAME | string | owncloud | Name of the database to be used.|
| SHARING_USER_OWNCLOUDSQL_USER_STORAGE_MOUNT_ID | string |  | Mount ID of the ownCloudSQL users storage for mapping ownCloud 10 shares.|
| SHARING_PUBLIC_DRIVER | string | cs3 | Driver to be used to persist public shares. Possible values: "json", "cs3"|
| SHARING_PUBLIC_JSON_FILE | string | ~/.ocis/storage/publicshares.json | Path to the JSON file where public share meta-data will be stored. This JSON file contains the information about public shares that have been created.|
| SHARING_PUBLIC_CS3_PROVIDER_ADDR | string | 127.0.0.1:9215 | GRPC address of the storage-system extension.|
| OCIS_SYSTEM_USER_ID<br/>SHARING_PUBLIC_CS3_SYSTEM_USER_ID | string |  | ID of the oCIS storage-system system user. Admins need to set the ID for the storage-system system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format.|
| OCIS_SYSTEM_USER_IDP<br/>SHARING_PUBLIC_CS3_SYSTEM_USER_IDP | string | internal | IDP of the oCIS storage-system system user.|
| OCIS_SYSTEM_USER_API_KEY<br/>SHARING_USER_CS3_SYSTEM_USER_API_KEY | string |  | API key for the storage-system system user.|