## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>STORAGE_USERS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORAGE_USERS_TRACING_TYPE | string |  | The type of tracing. Defaults to "", which is the same as "jaeger". Allowed tracing types are "jaeger" and "" as of now.|
| OCIS_TRACING_ENDPOINT<br/>STORAGE_USERS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>STORAGE_USERS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>STORAGE_USERS_LOG_LEVEL | string |  | The log level. Valid values are: "panic", "fatal", "error", "warn", "info", "debug", "trace".|
| OCIS_LOG_PRETTY<br/>STORAGE_USERS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORAGE_USERS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORAGE_USERS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| STORAGE_USERS_DEBUG_ADDR | string | 127.0.0.1:9159 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| STORAGE_USERS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| STORAGE_USERS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| STORAGE_USERS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| STORAGE_USERS_GRPC_ADDR | string | 127.0.0.1:9157 | The bind address of the GRPC service.|
| STORAGE_USERS_GRPC_PROTOCOL | string | tcp | The transport protocol of the GPRC service.|
| STORAGE_USERS_HTTP_ADDR | string | 127.0.0.1:9158 | The bind address of the HTTP service.|
| STORAGE_USERS_HTTP_PROTOCOL | string | tcp | The transport protocol of the HTTP service.|
| OCIS_JWT_SECRET<br/>STORAGE_USERS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| REVA_GATEWAY | string | 127.0.0.1:9142 | The CS3 gateway endpoint.|
| STORAGE_USERS_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the loading of user's group memberships from the reva access token.|
| STORAGE_USERS_DRIVER | string | ocis | The storage driver which should be used by the service|
| STORAGE_USERS_OCIS_ROOT | string | ~/.ocis/storage/users | Path for the persistence directory.|
| STORAGE_USERS_OCIS_USER_LAYOUT | string | {{.Id.OpaqueId}} | Template string for the user storage layout in the persistence directory.|
| STORAGE_USERS_PERMISSION_ENDPOINT,STORAGE_USERS_OCIS_PERMISSIONS_ENDPOINT | string | 127.0.0.1:9191 | Endpoint of the permissions service.|
| STORAGE_USERS_OCIS_PERSONAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.User.Username \| lower}} | Template string to construct personal space aliases.|
| STORAGE_USERS_OCIS_GENERAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.SpaceName \| replace &#34; &#34; &#34;-&#34; \| lower}} | Template string to construct general space aliases.|
| STORAGE_USERS_OCIS_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_S3NG_ROOT | string | ~/.ocis/storage/users | Path for the persistence directory.|
| STORAGE_USERS_S3NG_USER_LAYOUT | string | {{.Id.OpaqueId}} | Template string for the user storage layout in the persistence directory.|
| STORAGE_USERS_PERMISSION_ENDPOINT<br/>STORAGE_USERS_S3NG_PERMISSIONS_ENDPOINT | string | 127.0.0.1:9191 | Endpoint of the permissions service.|
| STORAGE_USERS_S3NG_REGION | string | default | Region of the S3 bucket.|
| STORAGE_USERS_S3NG_ACCESS_KEY | string |  | Access key for the S3 bucket.|
| STORAGE_USERS_S3NG_SECRET_KEY | string |  | Secret key for the S3 bucket.|
| STORAGE_USERS_S3NG_ENDPOINT | string |  | Endpoint for the S3 bucket.|
| STORAGE_USERS_S3NG_BUCKET | string |  | Name of the S3 bucket.|
| STORAGE_USERS_S3NG_PERSONAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.User.Username \| lower}} | Template string to construct personal space aliases.|
| STORAGE_USERS_S3NG_GENERAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.SpaceName \| replace &#34; &#34; &#34;-&#34; \| lower}} | Template string to construct general space aliases.|
| STORAGE_USERS_S3NG_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_OWNCLOUDSQL_DATADIR | string | ~/.ocis/storage/owncloud | Path for the persistence directory.|
| STORAGE_USERS_OWNCLOUDSQL_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_OWNCLOUDSQL_LAYOUT | string | {{.Username}} | Path layout to use to navigate into a users folder in an owncloud data directory|
| STORAGE_USERS_OWNCLOUDSQL_UPLOADINFO_DIR | string | ~/.ocis/storage/uploadinfo | Path to a directory, where uploads will be stored temporarily.|
| STORAGE_USERS_OWNCLOUDSQL_DB_USERNAME | string | owncloud | Username for the database.|
| STORAGE_USERS_OWNCLOUDSQL_DB_PASSWORD | string | owncloud | Password for the database.|
| STORAGE_USERS_OWNCLOUDSQL_DB_HOST | string |  | Hostname or IP of the database server.|
| STORAGE_USERS_OWNCLOUDSQL_DB_PORT | int | 3306 | Port that the database server is listening on.|
| STORAGE_USERS_OWNCLOUDSQL_DB_NAME | string | owncloud | Name of the database to be used.|
| STORAGE_USERS_OWNCLOUDSQL_USERS_PROVIDER_ENDPOINT | string | localhost:9144 | Endpoint of the users provider.|
| STORAGE_USERS_DATA_SERVER_URL | string | http://localhost:9158/data | URL of the data server, needs to be reachable by the data gateway provided by the frontend service or the user if directly exposed.|
| STORAGE_USERS_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | |
| STORAGE_USERS_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| STORAGE_USERS_MOUNT_ID | string | 1284d238-aa92-42ce-bdc4-0b0000009157 | Mount ID of this storage.|
| STORAGE_USERS_EXPOSE_DATA_SERVER | bool | false | Exposes the data server directly to users and bypasses the data gateway. Ensure that the data server address is reachable by users.|
| STORAGE_USERS_READ_ONLY | bool | false | Set this storage to be read-only.|