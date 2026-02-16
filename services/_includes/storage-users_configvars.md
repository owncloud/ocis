## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| STORAGE_USERS_SERVICE_NAME | string | storage-users | Service name to use. Change this when starting an additional storage provider with a custom configuration to prevent it from colliding with the default 'storage-users' service.|
| OCIS_TRACING_ENABLED<br/>STORAGE_USERS_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>STORAGE_USERS_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger', 'otlp' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>STORAGE_USERS_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>STORAGE_USERS_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>STORAGE_USERS_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>STORAGE_USERS_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>STORAGE_USERS_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>STORAGE_USERS_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| STORAGE_USERS_DEBUG_ADDR | string | 127.0.0.1:9159 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| STORAGE_USERS_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| STORAGE_USERS_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| STORAGE_USERS_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| STORAGE_USERS_GRPC_ADDR | string | 127.0.0.1:9157 | The bind address of the GRPC service.|
| OCIS_GRPC_PROTOCOL<br/>STORAGE_USERS_GRPC_PROTOCOL | string | tcp | The transport protocol of the GPRC service.|
| STORAGE_USERS_HTTP_ADDR | string | 127.0.0.1:9158 | The bind address of the HTTP service.|
| STORAGE_USERS_HTTP_PROTOCOL | string | tcp | The transport protocol of the HTTP service.|
| OCIS_CORS_ALLOW_ORIGINS<br/>STORAGE_USERS_CORS_ALLOW_ORIGINS | []string | [https://localhost:9200] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>STORAGE_USERS_CORS_ALLOW_METHODS | []string | [POST HEAD PATCH OPTIONS GET DELETE] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>STORAGE_USERS_CORS_ALLOW_HEADERS | []string | [Authorization Origin X-Requested-With X-Request-Id X-HTTP-Method-Override Content-Type Upload-Length Upload-Offset Tus-Resumable Upload-Metadata Upload-Defer-Length Upload-Concat Upload-Incomplete Upload-Draft-Interop-Version] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>STORAGE_USERS_CORS_ALLOW_CREDENTIALS | bool | false | Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| OCIS_CORS_EXPOSE_HEADERS<br/>STORAGE_USERS_CORS_EXPOSE_HEADERS | []string | [Upload-Offset Location Upload-Length Tus-Version Tus-Resumable Tus-Max-Size Tus-Extension Upload-Metadata Upload-Defer-Length Upload-Concat Upload-Incomplete Upload-Draft-Interop-Version] | A list of exposed CORS headers. See following chapter for more details: *Access-Control-Expose-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_MAX_AGE<br/>STORAGE_USERS_CORS_MAX_AGE | uint | 86400 | The max cache duration of preflight headers. See following chapter for more details: *Access-Control-Max-Age* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age. See the Environment Variable Types description for more details.|
| OCIS_JWT_SECRET<br/>STORAGE_USERS_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| OCIS_REVA_GATEWAY | string | com.owncloud.api.gateway | The CS3 gateway endpoint.|
| OCIS_GRPC_CLIENT_TLS_MODE | string |  | TLS mode for grpc connection to the go-micro based grpc services. Possible values are 'off', 'insecure' and 'on'. 'off': disables transport security for the clients. 'insecure' allows using transport security, but disables certificate verification (to be used with the autogenerated self-signed certificates). 'on' enables transport security, including server certificate verification.|
| OCIS_GRPC_CLIENT_TLS_CACERT | string |  | Path/File name for the root CA certificate (in PEM format) used to validate TLS server certificates of the go-micro based grpc services.|
| STORAGE_USERS_SKIP_USER_GROUPS_IN_TOKEN | bool | false | Disables the loading of user's group memberships from the reva access token.|
| STORAGE_USERS_GRACEFUL_SHUTDOWN_TIMEOUT | int | 30 | The number of seconds to wait for the 'storage-users' service to shutdown cleanly before exiting with an error that gets logged. Note: This setting is only applicable when running the 'storage-users' service as a standalone service. See the text description for more details.|
| STORAGE_USERS_DRIVER | string | ocis | The storage driver which should be used by the service. Defaults to 'ocis', Supported values are: 'ocis', 's3ng' and 'owncloudsql'. The 'ocis' driver stores all data (blob and meta data) in an POSIX compliant volume. The 's3ng' driver stores metadata in a POSIX compliant volume and uploads blobs to the s3 bucket.|
| OCIS_DECOMPOSEDFS_PROPAGATOR<br/>STORAGE_USERS_OCIS_PROPAGATOR | string | sync | The propagator used for decomposedfs. At the moment, only 'sync' is fully supported, 'async' is available as an experimental option.|
| STORAGE_USERS_ASYNC_PROPAGATOR_PROPAGATION_DELAY | Duration | 0s | The delay between a change made to a tree and the propagation start on treesize and treetime. Multiple propagations are computed to a single one. See the Environment Variable Types description for more details.|
| STORAGE_USERS_OCIS_ROOT | string | /var/lib/ocis/storage/users | The directory where the filesystem storage will store blobs and metadata. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users.|
| STORAGE_USERS_OCIS_USER_LAYOUT | string | {{.Id.OpaqueId}} | Template string for the user storage layout in the user directory.|
| STORAGE_USERS_PERMISSION_ENDPOINT<br/>STORAGE_USERS_OCIS_PERMISSIONS_ENDPOINT | string | com.owncloud.api.settings | Endpoint of the permissions service. The endpoints can differ for 'ocis' and 's3ng'.|
| STORAGE_USERS_OCIS_PERSONAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.User.Username \| lower}} | Template string to construct personal space aliases.|
| STORAGE_USERS_OCIS_PERSONAL_SPACE_PATH_TEMPLATE | string |  | Template string to construct the paths of the personal space roots.|
| STORAGE_USERS_OCIS_GENERAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.SpaceName \| replace &#34; &#34; &#34;-&#34; \| lower}} | Template string to construct general space aliases.|
| STORAGE_USERS_OCIS_GENERAL_SPACE_PATH_TEMPLATE | string |  | Template string to construct the paths of the projects space roots.|
| STORAGE_USERS_OCIS_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_OCIS_MAX_ACQUIRE_LOCK_CYCLES | int | 20 | When trying to lock files, ocis will try this amount of times to acquire the lock before failing. After each try it will wait for an increasing amount of time. Values of 0 or below will be ignored and the default value will be used.|
| STORAGE_USERS_OCIS_LOCK_CYCLE_DURATION_FACTOR | int | 30 | When trying to lock files, ocis will multiply the cycle with this factor and use it as a millisecond timeout. Values of 0 or below will be ignored and the default value will be used.|
| OCIS_MAX_CONCURRENCY<br/>STORAGE_USERS_OCIS_MAX_CONCURRENCY | int | 5 | Maximum number of concurrent go-routines. Higher values can potentially get work done faster but will also cause more load on the system. Values of 0 or below will be ignored and the default value will be used.|
| OCIS_ASYNC_UPLOADS | bool | true | Enable asynchronous file uploads.|
| OCIS_SPACES_MAX_QUOTA<br/>STORAGE_USERS_OCIS_MAX_QUOTA | uint64 | 0 | Set a global max quota for spaces in bytes. A value of 0 equals unlimited. If not using the global OCIS_SPACES_MAX_QUOTA, you must define the FRONTEND_MAX_QUOTA in the frontend service.|
| OCIS_DISABLE_VERSIONING | bool | false | Disables versioning of files. When set to true, new uploads with the same filename will overwrite existing files instead of creating a new version.|
| OCIS_DECOMPOSEDFS_PROPAGATOR<br/>STORAGE_USERS_S3NG_PROPAGATOR | string | sync | The propagator used for decomposedfs. At the moment, only 'sync' is fully supported, 'async' is available as an experimental option.|
| STORAGE_USERS_ASYNC_PROPAGATOR_PROPAGATION_DELAY | Duration | 0s | The delay between a change made to a tree and the propagation start on treesize and treetime. Multiple propagations are computed to a single one. See the Environment Variable Types description for more details.|
| STORAGE_USERS_S3NG_ROOT | string | /var/lib/ocis/storage/users | The directory where the filesystem storage will store metadata for blobs. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users.|
| STORAGE_USERS_S3NG_USER_LAYOUT | string | {{.Id.OpaqueId}} | Template string for the user storage layout in the user directory.|
| STORAGE_USERS_PERMISSION_ENDPOINT<br/>STORAGE_USERS_S3NG_PERMISSIONS_ENDPOINT | string | com.owncloud.api.settings | Endpoint of the permissions service. The endpoints can differ for 'ocis' and 's3ng'.|
| STORAGE_USERS_S3NG_REGION | string | default | Region of the S3 bucket.|
| STORAGE_USERS_S3NG_ACCESS_KEY | string |  | Access key for the S3 bucket.|
| STORAGE_USERS_S3NG_SECRET_KEY | string |  | Secret key for the S3 bucket.|
| STORAGE_USERS_S3NG_ENDPOINT | string |  | Endpoint for the S3 bucket.|
| STORAGE_USERS_S3NG_BUCKET | string |  | Name of the S3 bucket.|
| STORAGE_USERS_S3NG_PUT_OBJECT_DISABLE_CONTENT_SHA256 | bool | false | Disable sending content sha256 when copying objects to S3.|
| STORAGE_USERS_S3NG_PUT_OBJECT_DISABLE_MULTIPART | bool | true | Disable multipart uploads when copying objects to S3|
| STORAGE_USERS_S3NG_PUT_OBJECT_SEND_CONTENT_MD5 | bool | true | Send a Content-MD5 header when copying objects to S3.|
| STORAGE_USERS_S3NG_PUT_OBJECT_CONCURRENT_STREAM_PARTS | bool | true | Always precreate parts when copying objects to S3.|
| STORAGE_USERS_S3NG_PUT_OBJECT_NUM_THREADS | uint | 4 | Number of concurrent uploads to use when copying objects to S3.|
| STORAGE_USERS_S3NG_PUT_OBJECT_PART_SIZE | uint64 | 0 | Part size for concurrent uploads to S3. If no value or 0 is set, the library's default value of 16MB is used. The value range is min 5MB and max 5GB.|
| STORAGE_USERS_S3NG_PERSONAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.User.Username \| lower}} | Template string to construct personal space aliases.|
| STORAGE_USERS_S3NG_PERSONAL_SPACE_PATH_TEMPLATE | string |  | Template string to construct the paths of the personal space roots.|
| STORAGE_USERS_S3NG_GENERAL_SPACE_ALIAS_TEMPLATE | string | {{.SpaceType}}/{{.SpaceName \| replace &#34; &#34; &#34;-&#34; \| lower}} | Template string to construct general space aliases.|
| STORAGE_USERS_S3NG_GENERAL_SPACE_PATH_TEMPLATE | string |  | Template string to construct the paths of the projects space roots.|
| STORAGE_USERS_S3NG_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_S3NG_MAX_ACQUIRE_LOCK_CYCLES | int | 20 | When trying to lock files, ocis will try this amount of times to acquire the lock before failing. After each try it will wait for an increasing amount of time. Values of 0 or below will be ignored and the default value of 20 will be used.|
| STORAGE_USERS_S3NG_LOCK_CYCLE_DURATION_FACTOR | int | 30 | When trying to lock files, ocis will multiply the cycle with this factor and use it as a millisecond timeout. Values of 0 or below will be ignored and the default value of 30 will be used.|
| OCIS_MAX_CONCURRENCY<br/>STORAGE_USERS_S3NG_MAX_CONCURRENCY | int | 5 | Maximum number of concurrent go-routines. Higher values can potentially get work done faster but will also cause more load on the system. Values of 0 or below will be ignored and the default value of 100 will be used.|
| OCIS_DISABLE_VERSIONING | bool | false | Disables versioning of files. When set to true, new uploads with the same filename will overwrite existing files instead of creating a new version.|
| STORAGE_USERS_OWNCLOUDSQL_DATADIR | string | /var/lib/ocis/storage/owncloud | The directory where the filesystem storage will store SQL migration data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/owncloud.|
| STORAGE_USERS_OWNCLOUDSQL_SHARE_FOLDER | string | /Shares | Name of the folder jailing all shares.|
| STORAGE_USERS_OWNCLOUDSQL_LAYOUT | string | {{.Username}} | Path layout to use to navigate into a users folder in an owncloud data directory|
| STORAGE_USERS_OWNCLOUDSQL_UPLOADINFO_DIR | string | /var/lib/ocis/storage/uploadinfo | The directory where the filesystem will store uploads temporarily. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/uploadinfo.|
| STORAGE_USERS_OWNCLOUDSQL_DB_USERNAME | string | owncloud | Username for the database.|
| STORAGE_USERS_OWNCLOUDSQL_DB_PASSWORD | string | owncloud | Password for the database.|
| STORAGE_USERS_OWNCLOUDSQL_DB_HOST | string |  | Hostname or IP of the database server.|
| STORAGE_USERS_OWNCLOUDSQL_DB_PORT | int | 3306 | Port that the database server is listening on.|
| STORAGE_USERS_OWNCLOUDSQL_DB_NAME | string | owncloud | Name of the database to be used.|
| STORAGE_USERS_OWNCLOUDSQL_USERS_PROVIDER_ENDPOINT | string | com.owncloud.api.users | Endpoint of the users provider.|
| STORAGE_USERS_POSIX_ROOT | string | /var/lib/ocis/storage/users | The directory where the filesystem storage will store its data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/storage/users.|
| STORAGE_USERS_POSIX_PERSONAL_SPACE_PATH_TEMPLATE | string | users/{{.User.Username}} | Template string to construct the paths of the personal space roots.|
| STORAGE_USERS_POSIX_GENERAL_SPACE_PATH_TEMPLATE | string | projects/{{.SpaceId}} | Template string to construct the paths of the projects space roots.|
| STORAGE_USERS_PERMISSION_ENDPOINT<br/>STORAGE_USERS_POSIX_PERMISSIONS_ENDPOINT | string | com.owncloud.api.settings | Endpoint of the permissions service. The endpoints can differ for 'ocis', 'posix' and 's3ng'.|
| OCIS_ASYNC_UPLOADS | bool | true | Enable asynchronous file uploads.|
| STORAGE_USERS_POSIX_SCAN_DEBOUNCE_DELAY | Duration | 1s | The time in milliseconds to wait before scanning the filesystem for changes after a change has been detected.|
| STORAGE_USERS_POSIX_USE_SPACE_GROUPS | bool | false | Use space groups to manage permissions on spaces.|
| STORAGE_USERS_POSIX_WATCH_TYPE | string |  | Type of the watcher to use for getting notified about changes to the filesystem. Currently available options are 'inotifywait' (default), 'gpfswatchfolder' and 'gpfsfileauditlogging'.|
| STORAGE_USERS_POSIX_WATCH_PATH | string |  | Path to the watch directory/file. Only applies to the 'gpfsfileauditlogging' and 'inotifywait' watcher, in which case it is the path of the file audit log file/base directory to watch.|
| STORAGE_USERS_POSIX_WATCH_FOLDER_KAFKA_BROKERS | string |  | Comma-separated list of kafka brokers to read the watchfolder events from.|
| STORAGE_USERS_DATA_SERVER_URL | string | http://localhost:9158/data | URL of the data server, needs to be reachable by the data gateway provided by the frontend service or the user if directly exposed.|
| STORAGE_USERS_DATA_GATEWAY_URL | string | https://localhost:9200/data | URL of the data gateway server|
| STORAGE_USERS_TRANSFER_EXPIRES | int64 | 86400 | The time after which the token for upload postprocessing expires|
| OCIS_EVENTS_ENDPOINT<br/>STORAGE_USERS_EVENTS_ENDPOINT | string | 127.0.0.1:9233 | The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture.|
| OCIS_EVENTS_CLUSTER<br/>STORAGE_USERS_EVENTS_CLUSTER | string | ocis-cluster | The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system.|
| OCIS_INSECURE<br/>STORAGE_USERS_EVENTS_TLS_INSECURE | bool | false | Whether to verify the server TLS certificates.|
| OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE<br/>STORAGE_USERS_EVENTS_TLS_ROOT_CA_CERTIFICATE | string |  | The root CA certificate used to validate the server's TLS certificate. If provided STORAGE_USERS_EVENTS_TLS_INSECURE will be seen as false.|
| OCIS_EVENTS_ENABLE_TLS<br/>STORAGE_USERS_EVENTS_ENABLE_TLS | bool | false | Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| STORAGE_USERS_EVENTS_NUM_CONSUMERS | int | 0 | The amount of concurrent event consumers to start. Event consumers are used for post-processing files. Multiple consumers increase parallelisation, but will also increase CPU and memory demands. The setting has no effect when the OCIS_ASYNC_UPLOADS is set to false. The default and minimum value is 1.|
| OCIS_EVENTS_AUTH_USERNAME<br/>STORAGE_USERS_EVENTS_AUTH_USERNAME | string |  | The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_EVENTS_AUTH_PASSWORD<br/>STORAGE_USERS_EVENTS_AUTH_PASSWORD | string |  | The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services.|
| OCIS_CACHE_STORE<br/>STORAGE_USERS_FILEMETADATA_CACHE_STORE | string | memory | The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details.|
| OCIS_CACHE_STORE_NODES<br/>STORAGE_USERS_FILEMETADATA_CACHE_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| OCIS_CACHE_DATABASE | string | storage-users | The database name the configured store should use.|
| OCIS_CACHE_TTL<br/>STORAGE_USERS_FILEMETADATA_CACHE_TTL | Duration | 24m0s | Default time to live for user info in the user info cache. Only applied when access tokens has no expiration. See the Environment Variable Types description for more details.|
| OCIS_CACHE_DISABLE_PERSISTENCE<br/>STORAGE_USERS_FILEMETADATA_CACHE_DISABLE_PERSISTENCE | bool | false | Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false.|
| OCIS_CACHE_AUTH_USERNAME<br/>STORAGE_USERS_FILEMETADATA_CACHE_AUTH_USERNAME | string |  | The username to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_CACHE_AUTH_PASSWORD<br/>STORAGE_USERS_FILEMETADATA_CACHE_AUTH_PASSWORD | string |  | The password to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_CACHE_STORE<br/>STORAGE_USERS_ID_CACHE_STORE | string | memory | The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details.|
| OCIS_CACHE_STORE_NODES<br/>STORAGE_USERS_ID_CACHE_STORE_NODES | []string | [127.0.0.1:9233] | A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details.|
| OCIS_CACHE_DATABASE | string | ids-storage-users | The database name the configured store should use.|
| OCIS_CACHE_TTL<br/>STORAGE_USERS_ID_CACHE_TTL | Duration | 24m0s | Default time to live for user info in the user info cache. Only applied when access tokens have no expiration. Defaults to 300s which is derived from the underlaying package though not explicitly set as default. See the Environment Variable Types description for more details.|
| OCIS_CACHE_DISABLE_PERSISTENCE<br/>STORAGE_USERS_ID_CACHE_DISABLE_PERSISTENCE | bool | false | Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false.|
| OCIS_CACHE_AUTH_USERNAME<br/>STORAGE_USERS_ID_CACHE_AUTH_USERNAME | string |  | The username to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured.|
| OCIS_CACHE_AUTH_PASSWORD<br/>STORAGE_USERS_ID_CACHE_AUTH_PASSWORD | string |  | The password to authenticate with the cache store. Only applies when store type 'nats-js-kv' is configured.|
| STORAGE_USERS_MOUNT_ID | string |  | Mount ID of this storage.|
| STORAGE_USERS_EXPOSE_DATA_SERVER | bool | false | Exposes the data server directly to users and bypasses the data gateway. Ensure that the data server address is reachable by users.|
| STORAGE_USERS_READ_ONLY | bool | false | Set this storage to be read-only.|
| STORAGE_USERS_UPLOAD_EXPIRATION | int64 | 86400 | Duration in seconds after which uploads will expire. Note that when setting this to a low number, uploads could be cancelled before they are finished and return a 403 to the user.|
| OCIS_ADMIN_USER_ID<br/>STORAGE_USERS_PURGE_TRASH_BIN_USER_ID | string |  | ID of the user who collects all necessary information for deletion. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand.|
| STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE | Duration | 720h0m0s | Specifies the period of time in which items that have been in the personal trash-bin for longer than this value should be deleted. A value of 0 means no automatic deletion. See the Environment Variable Types description for more details.|
| STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE | Duration | 720h0m0s | Specifies the period of time in which items that have been in the project trash-bin for longer than this value should be deleted. A value of 0 means no automatic deletion. See the Environment Variable Types description for more details.|
| OCIS_SERVICE_ACCOUNT_ID<br/>STORAGE_USERS_SERVICE_ACCOUNT_ID | string |  | The ID of the service account the service should use. See the 'auth-service' service description for more details.|
| OCIS_SERVICE_ACCOUNT_SECRET<br/>STORAGE_USERS_SERVICE_ACCOUNT_SECRET | string |  | The service account secret.|
| OCIS_GATEWAY_GRPC_ADDR<br/>STORAGE_USERS_GATEWAY_GRPC_ADDR | string | 127.0.0.1:9142 | The bind address of the gateway GRPC address.|
| OCIS_MACHINE_AUTH_API_KEY<br/>STORAGE_USERS_MACHINE_AUTH_API_KEY | string |  | Machine auth API key used to validate internal requests necessary for the access to resources from other services.|
| STORAGE_USERS_CLI_MAX_ATTEMPTS_RENAME_FILE | int | 0 | The maximum number of attempts to rename a file when a user restores a file to an existing destination with the same name. The minimum value is 100.|