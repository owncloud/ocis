## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| STORAGE_USERS_DEBUG_ADDR | string | 127.0.0.1:9159 | |
| STORAGE_USERS_DEBUG_TOKEN | string |  | |
| STORAGE_USERS_DEBUG_PPROF | bool | false | |
| STORAGE_USERS_DEBUG_ZPAGES | bool | false | |
| STORAGE_USERS_GRPC_ADDR | string | 127.0.0.1:9157 | The address of the grpc service.|
| STORAGE_USERS_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| STORAGE_USERS_GRPC_ADDR | string | 127.0.0.1:9158 | The address of the grpc service.|
| STORAGE_USERS_GRPC_PROTOCOL | string | tcp | The transport protocol of the grpc service.|
| STORAGE_USERS_DRIVER | string | ocis | The storage driver which should be used by the service|
| STORAGE_USERS_LOCAL_ROOT | string | ~/.ocis/storage/local/users | |
| STORAGE_USERS_OCIS_ROOT | string | ~/.ocis/storage/users | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DATADIR | string | ~/.ocis/storage/owncloud | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_SHARE_FOLDER | string | /Shares | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_LAYOUT | string | {{.Username}} | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR | string | ~/.ocis/storage/uploadinfo | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBUSERNAME | string | owncloud | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPASSWORD | string | owncloud | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBHOST | string |  | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPORT | int | 3306 | |
| STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBNAME | string | owncloud | |
| OCIS_INSECURE;STORAGE_USERS_DATAPROVIDER_INSECURE | bool | false | |