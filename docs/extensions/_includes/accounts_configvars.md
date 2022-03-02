## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| ACCOUNTS_DEBUG_ADDR | string | 127.0.0.1:9182 | |
| ACCOUNTS_DEBUG_TOKEN | string |  | |
| ACCOUNTS_DEBUG_PPROF | bool | false | |
| ACCOUNTS_DEBUG_ZPAGES | bool | false | |
| ACCOUNTS_HTTP_ADDR | string | 127.0.0.1:9181 | The address of the http service.|
| ACCOUNTS_HTTP_ROOT | string | / | The root path of the http service.|
| ACCOUNTS_CACHE_TTL | int | 604800 | The cache time for the static assets.|
| ACCOUNTS_GRPC_ADDR | string | 127.0.0.1:9180 | The address of the grpc service.|
| OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET | string | Pive-Fumkiu4 | The secret to mint jwt tokens.|
| ACCOUNTS_ASSET_PATH | string |  | The path to the ui assets.|
| ACCOUNTS_STORAGE_BACKEND | string | CS3 | Defines which storage implementation is to be used|
| ACCOUNTS_STORAGE_DISK_PATH | string | ~/.ocis/accounts | The path where the accounts data is stored.|
| ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR | string | localhost:9215 | The address to the storage provider.|
| ACCOUNTS_UID_INDEX_LOWER_BOUND | int64 | 0 | The lowest possible uid value for the indexer.|
| ACCOUNTS_UID_INDEX_UPPER_BOUND | int64 | 1000 | The highest possible uid value for the indexer.|
| ACCOUNTS_GID_INDEX_LOWER_BOUND | int64 | 0 | The lowest possible gid value for the indexer.|
| ACCOUNTS_GID_INDEX_UPPER_BOUND | int64 | 1000 | The highest possible gid value for the indexer.|
| ACCOUNTS_SERVICE_USER_UUID | string | 95cb8724-03b2-11eb-a0a6-c33ef8ef53ad | The id of the accounts service user.|
| ACCOUNTS_SERVICE_USER_USERNAME | string |  | The username of the accounts service user.|
| ACCOUNTS_SERVICE_USER_UID | int64 | 0 | The uid of the accounts service user.|
| ACCOUNTS_SERVICE_USER_GID | int64 | 0 | The gid of the accounts service user.|
| ACCOUNTS_HASH_DIFFICULTY | int | 11 | The hash difficulty makes sure that validating a password takes at least a certain amount of time.|
| ACCOUNTS_DEMO_USERS_AND_GROUPS | bool | true | If this flag is set the service will setup the demo users and groups.|