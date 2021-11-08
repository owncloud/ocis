---
title: "Configuration"
date: "2021-11-08T07:48:13+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/storage/templates
geekdocFilePath: CONFIGURATION.tmpl
---

{{< toc >}}

## Configuration

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-proxy reads `proxy.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/storage/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Storage service for oCIS

Usage: `storage [global options] command [command options] [arguments...]`



-config-file |  $STORAGE_CONFIG_FILE
: Path to config file.


-log-level |  $STORAGE_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $STORAGE_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $STORAGE_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





















































































































































































































































## Sub Commands

### storage auth-machine

Start authprovider for machine auth

Usage: `storage auth-machine [command options] [arguments...]`




























































































































































-debug-addr |  $STORAGE_AUTH_MACHINE_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9167"`.


-machine-auth-api-key |  $STORAGE_AUTH_MACHINE_AUTH_API_KEY , $OCIS_MACHINE_AUTH_API_KEY
: the API key to be used for the machine auth driver in reva. Default: `"change-me-please"`.


-network |  $STORAGE_AUTH_MACHINE_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_AUTH_MACHINE_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9166"`.


-service |  $STORAGE_AUTH_MACHINE_SERVICES
: --service authprovider [--service otherservice]. Default: `cli.NewStringSlice("authprovider")`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


























































































### storage sharing

Start sharing service

Usage: `storage sharing [command options] [arguments...]`







-debug-addr |  $STORAGE_SHARING_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9151"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-network |  $STORAGE_SHARING_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_SHARING_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9150"`.


-service |  $STORAGE_SHARING_SERVICES
: --service usershareprovider [--service publicshareprovider]. Default: `cli.NewStringSlice("usershareprovider", "publicshareprovider")`.


-user-driver |  $STORAGE_SHARING_USER_DRIVER
: driver to use for the UserShareProvider. Default: `"json"`.


-user-json-file |  $STORAGE_SHARING_USER_JSON_FILE
: file used to persist shares for the UserShareProvider. Default: `flags.OverrideDefaultString(cfg.Reva.Sharing.UserJSONFile, path.Join(defaults.BaseDataPath(), "storage", "shares.json"))`.


-public-driver |  $STORAGE_SHARING_PUBLIC_DRIVER
: driver to use for the PublicShareProvider. Default: `"json"`.


-public-json-file |  $STORAGE_SHARING_PUBLIC_JSON_FILE
: file used to persist shares for the PublicShareProvider. Default: `flags.OverrideDefaultString(cfg.Reva.Sharing.PublicJSONFile, path.Join(defaults.BaseDataPath(), "storage", "publicshares.json"))`.


-public-password-hash-cost |  $STORAGE_SHARING_PUBLIC_PASSWORD_HASH_COST
: the cost of hashing the public shares passwords. Default: `11`.


-public-enable-expired-shares-cleanup |  $STORAGE_SHARING_PUBLIC_ENABLE_EXPIRED_SHARES_CLEANUP
: whether to periodically delete expired public shares. Default: `true`.


-public-janitor-run-interval |  $STORAGE_SHARING_PUBLIC_JANITOR_RUN_INTERVAL
: the time period in seconds after which to start a janitor run. Default: `60`.









































































































































































































































### storage storage-home

Start storage-home service

Usage: `storage storage-home [command options] [arguments...]`



















-debug-addr |  $STORAGE_HOME_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9156"`.


-grpc-network |  $STORAGE_HOME_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_HOME_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9154"`.


-http-network |  $STORAGE_HOME_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_HOME_HTTP_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9155"`.


-driver |  $STORAGE_HOME_DRIVER
: storage driver for home mount: eg. local, eos, owncloud, ocis or s3. Default: `"ocis"`.


-read-only |  $STORAGE_HOME_READ_ONLY , $OCIS_STORAGE_READ_ONLY
: use storage driver in read-only mode. Default: `false`.


-mount-path |  $STORAGE_HOME_MOUNT_PATH
: mount path. Default: `"/home"`.


-mount-id |  $STORAGE_HOME_MOUNT_ID
: mount id. Default: `"1284d238-aa92-42ce-bdc4-0b0000009157"`.


-expose-data-server |  $STORAGE_HOME_EXPOSE_DATA_SERVER
: exposes a dedicated data server. Default: `false`.


-data-server-url |  $STORAGE_HOME_DATA_SERVER_URL
: data server url. Default: `"http://localhost:9155/data"`.


-http-prefix |  $STORAGE_HOME_HTTP_PREFIX
: prefix for the http endpoint, without leading slash. Default: `"data"`.


-tmp-folder |  $STORAGE_HOME_TMP_FOLDER
: path to tmp folder. Default: `flags.OverrideDefaultString(cfg.Reva.StorageHome.TempFolder, path.Join(defaults.BaseDataPath(), "tmp", "home"))`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-users-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9144"`.


























































































































































































































### storage users

Start users service

Usage: `storage users [command options] [arguments...]`







































-debug-addr |  $STORAGE_USERPROVIDER_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9145"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-network |  $STORAGE_USERPROVIDER_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_USERPROVIDER_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9144"`.


-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: URL to use for the storage service. Default: `"localhost:9144"`.


-service |  $STORAGE_USERPROVIDER_SERVICES
: --service userprovider [--service otherservice]. Default: `cli.NewStringSlice("userprovider")`.


-driver |  $STORAGE_USERPROVIDER_DRIVER
: user driver: 'demo', 'json', 'ldap', 'owncloudsql' or 'rest'. Default: `"ldap"`.


-json-config |  $STORAGE_USERPROVIDER_JSON
: Path to users.json file. Default: `""`.


-user-groups-cache-expiration |  $STORAGE_USER_CACHE_EXPIRATION
: Time in minutes for redis cache expiration.. Default: `5`.


-owncloudsql-dbhost |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_DBHOST
: hostname of the mysql db. Default: `"mysql"`.


-owncloudsql-dbport |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPORT
: port of the mysql db. Default: `3306`.


-owncloudsql-dbname |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_DBNAME
: database name of the owncloud db. Default: `"owncloud"`.


-owncloudsql-dbuser |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_DBUSER
: user name to use when connecting to the mysql owncloud db. Default: `"owncloud"`.


-owncloudsql-dbpass |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_DBPASS
: password to use when connecting to the mysql owncloud db. Default: `"secret"`.


-owncloudsql-idp |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_IDP , $OCIS_URL
: Identity provider to use for users. Default: `"https://localhost:9200"`.


-owncloudsql-nobody |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_NOBODY
: fallback user id to use when user has no id. Default: `99`.


-owncloudsql-join-username |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_USERNAME
: join the username from the oc_preferences table. Default: `false`.


-owncloudsql-join-ownclouduuid |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_JOIN_OWNCLOUDUUID
: join the ownclouduuid from the oc_preferences table. Default: `false`.


-owncloudsql-enable-medial-search |  $STORAGE_USERPROVIDER_OWNCLOUDSQL_ENABLE_MEDIAL_SEARCH
: enable medial search when finding users. Default: `false`.


































































































































































































### storage auth-bearer

Start authprovider for bearer auth

Usage: `storage auth-bearer [command options] [arguments...]`


















































































































































-debug-addr |  $STORAGE_AUTH_BEARER_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9149"`.


-oidc-issuer |  $STORAGE_OIDC_ISSUER , $OCIS_URL
: OIDC issuer. Default: `"https://localhost:9200"`.


-oidc-insecure |  $STORAGE_OIDC_INSECURE
: OIDC allow insecure communication. Default: `true`.


-oidc-id-claim |  $STORAGE_OIDC_ID_CLAIM
: OIDC id claim. Default: `"preferred_username"`.


-oidc-uid-claim |  $STORAGE_OIDC_UID_CLAIM
: OIDC uid claim. Default: `""`.


-oidc-gid-claim |  $STORAGE_OIDC_GID_CLAIM
: OIDC gid claim. Default: `""`.


-network |  $STORAGE_AUTH_BEARER_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_AUTH_BEARER_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9148"`.


-service |  $STORAGE_AUTH_BEARER_SERVICES
: --service authprovider [--service otherservice]. Default: `cli.NewStringSlice("authprovider")`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.
































































































### storage frontend

Start frontend service

Usage: `storage frontend [command options] [arguments...]`


































































































































































-debug-addr |  $STORAGE_FRONTEND_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9141"`.


-transfer-secret |  $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `"replace-me-with-a-transfer-secret"`.


-chunk-folder |  $STORAGE_CHUNK_FOLDER
: temp directory for chunked uploads. Default: `flags.OverrideDefaultString(cfg.Reva.OCDav.WebdavNamespace, path.Join(defaults.BaseDataPath(), "tmp", "chunks"))`.


-webdav-namespace |  $STORAGE_WEBDAV_NAMESPACE
: Namespace prefix for the /webdav endpoint. Default: `"/home/"`.


-dav-files-namespace |  $STORAGE_DAV_FILES_NAMESPACE
: Namespace prefix for the webdav /dav/files endpoint. Default: `"/users/"`.


-archiver-max-num-files |  $STORAGE_ARCHIVER_MAX_NUM_FILES
: Maximum number of files to be included in the archiver. Default: `10000`.


-archiver-max-size |  $STORAGE_ARCHIVER_MAX_SIZE
: Maximum size for the sum of the sizes of all the files included in the archive. Default: `1073741824`.


-network |  $STORAGE_FRONTEND_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_FRONTEND_HTTP_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9140"`.


-public-url |  $STORAGE_FRONTEND_PUBLIC_URL , $OCIS_URL
: URL to use for the storage service. Default: `"https://localhost:9200"`.


-service |  $STORAGE_FRONTEND_SERVICES
: --service ocdav [--service ocs]. Default: `cli.NewStringSlice("datagateway", "ocdav", "ocs", "appprovider")`.


-approvider-prefix |  $STORAGE_FRONTEND_APPPROVIDER_PREFIX
: approvider prefix. Default: `""`.


-archiver-prefix |  $STORAGE_FRONTEND_ARCHIVER_PREFIX
: archiver prefix. Default: `"archiver"`.


-datagateway-prefix |  $STORAGE_FRONTEND_DATAGATEWAY_PREFIX
: datagateway prefix. Default: `"data"`.


-favorites |  $STORAGE_FRONTEND_FAVORITES
: announces favorites support to clients. Default: `false`.


-ocdav-prefix |  $STORAGE_FRONTEND_OCDAV_PREFIX
: owncloud webdav endpoint prefix. Default: `""`.


-ocs-prefix |  $STORAGE_FRONTEND_OCS_PREFIX
: open collaboration services endpoint prefix. Default: `"ocs"`.


-ocs-share-prefix |  $STORAGE_FRONTEND_OCS_SHARE_PREFIX
: the prefix prepended to the path of shared files. Default: `"/Shares"`.


-ocs-home-namespace |  $STORAGE_FRONTEND_OCS_HOME_NAMESPACE
: the prefix prepended to the incoming requests in OCS. Default: `"/home"`.


-ocs-resource-info-cache-ttl |  $STORAGE_FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL
: the TTL for statted resources in the share cache. Default: `0`.


-ocs-cache-warmup-driver |  $STORAGE_FRONTEND_OCS_CACHE_WARMUP_DRIVER
: the driver to be used for warming up the share cache. Default: `""`.


-ocs-additional-info-attribute |  $STORAGE_FRONTEND_OCS_ADDITIONAL_INFO_ATTRIBUTE
: the additional info to be returned when searching for users. Default: `"{{.Mail}}"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-default-upload-protocol |  $STORAGE_FRONTEND_DEFAULT_UPLOAD_PROTOCOL
: Default upload chunking protocol to be used out of tus/v1/ng. Default: `"tus"`.


-upload-max-chunk-size |  $STORAGE_FRONTEND_UPLOAD_MAX_CHUNK_SIZE
: Max chunk size in bytes to advertise to clients through capabilities, or 0 for unlimited. Default: `1e+8`.


-upload-http-method-override |  $STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE
: Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH. Default: `""`.


-checksum-supported-type |  $STORAGE_FRONTEND_CHECKSUM_SUPPORTED_TYPES
: --checksum-supported-type sha1 [--checksum-supported-type adler32]. Default: `cli.NewStringSlice("sha1", "md5", "adler32")`.


-checksum-preferred-upload-type |  $STORAGE_FRONTEND_CHECKSUM_PREFERRED_UPLOAD_TYPE
: Specify the preferred checksum algorithm used for uploads. Default: `""`.


-archiver-url |  $STORAGE_FRONTEND_ARCHIVER_URL
: URL where the archiver is reachable. Default: `"/archiver"`.


-appprovider-apps-url |  $STORAGE_FRONTEND_APP_PROVIDER_APPS_URL
: URL where the app listing of the app provider is reachable. Default: `"/app/list"`.


-appprovider-open-url |  $STORAGE_FRONTEND_APP_PROVIDER_OPEN_URL
: URL where files can be handed over to an application from the app provider. Default: `"/app/open"`.


-user-agent-whitelist-lock-in |  $STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT
: --user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of comma separated [UserAgent:challenge] values, it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the storage-frontend it configures Reva..


























































### storage gateway

Start gateway

Usage: `storage gateway [command options] [arguments...]`


































































































































































































-debug-addr |  $STORAGE_GATEWAY_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9143"`.


-transfer-secret |  $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `"replace-me-with-a-transfer-secret"`.


-transfer-expires |  $STORAGE_TRANSFER_EXPIRES
: Transfer token ttl in seconds. Default: `24 * 60 * 60`.


-network |  $STORAGE_GATEWAY_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_GATEWAY_GRPC_ADDR
: Address to bind REVA service. Default: `"127.0.0.1:9142"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-service |  $STORAGE_GATEWAY_SERVICES
: --service gateway [--service authregistry]. Default: `cli.NewStringSlice("gateway", "authregistry", "storageregistry", "appregistry")`.


-commit-share-to-storage-grant |  $STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT
: Commit shares to the share manager. Default: `true`.


-commit-share-to-storage-ref |  $STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF
: Commit shares to the storage. Default: `true`.


-share-folder |  $STORAGE_GATEWAY_SHARE_FOLDER
: mount shares in this folder of the home storage provider. Default: `"Shares"`.


-disable-home-creation-on-login |  $STORAGE_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN
: Disable creation of home folder on login.


-storage-home-mapping |  $STORAGE_GATEWAY_HOME_MAPPING
: mapping template for user home paths to user-specific mount points, e.g. /home/{{substr 0 1 .Username}}. Default: `""`.


-etag-cache-ttl |  $STORAGE_GATEWAY_ETAG_CACHE_TTL
: TTL for the home and shares directory etags cache. Default: `0`.


-auth-basic-endpoint |  $STORAGE_AUTH_BASIC_ENDPOINT
: endpoint to use for the basic auth provider. Default: `"localhost:9146"`.


-auth-bearer-endpoint |  $STORAGE_AUTH_BEARER_ENDPOINT
: endpoint to use for the bearer auth provider. Default: `"localhost:9148"`.


-auth-machine-endpoint |  $STORAGE_AUTH_MACHINE_ENDPOINT
: endpoint to use for the machine auth provider. Default: `"localhost:9166"`.


-storage-registry-driver |  $STORAGE_STORAGE_REGISTRY_DRIVER
: driver of the storage registry. Default: `"static"`.


-storage-registry-rule |  $STORAGE_STORAGE_REGISTRY_RULES
: `Replaces the generated storage registry rules with this set: --storage-registry-rule "/eos=localhost:9158" [--storage-registry-rule "1284d238-aa92-42ce-bdc4-0b0000009162=localhost:9162"]`. Default: `cli.NewStringSlice()`.


-storage-home-provider |  $STORAGE_STORAGE_REGISTRY_HOME_PROVIDER
: mount point of the storage provider for user homes in the global namespace. Default: `"/home"`.


-storage-registry-json |  $STORAGE_STORAGE_REGISTRY_JSON
: JSON file containing the storage registry rules. Default: `""`.


-app-registry-driver |  $STORAGE_APP_REGISTRY_DRIVER
: driver of the app registry. Default: `"static"`.


-app-registry-mimetypes-json |  $STORAGE_APP_REGISTRY_MIMETYPES_JSON
: JSON file containing the storage registry rules. Default: `""`.


-public-url |  $STORAGE_FRONTEND_PUBLIC_URL , $OCIS_URL
: URL to use for the storage service. Default: `"https://localhost:9200"`.


-datagateway-url |  $STORAGE_DATAGATEWAY_PUBLIC_URL
: URL to use for the storage datagateway, defaults to <STORAGE_FRONTEND_PUBLIC_URL>/data. Default: `""`.


-userprovider-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the userprovider. Default: `"localhost:9144"`.


-groupprovider-endpoint |  $STORAGE_GROUPPROVIDER_ENDPOINT
: endpoint to use for the groupprovider. Default: `"localhost:9160"`.


-sharing-endpoint |  $STORAGE_SHARING_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9150"`.


-appprovider-endpoint |  $STORAGE_APPPROVIDER_ENDPOINT
: endpoint to use for the app provider. Default: `"localhost:9164"`.


-storage-home-endpoint |  $STORAGE_HOME_ENDPOINT
: endpoint to use for the home storage. Default: `"localhost:9154"`.


-storage-home-mount-path |  $STORAGE_HOME_MOUNT_PATH
: mount path. Default: `"/home"`.


-storage-home-mount-id |  $STORAGE_HOME_MOUNT_ID
: mount id. Default: `"1284d238-aa92-42ce-bdc4-0b0000009154"`.


-storage-users-endpoint |  $STORAGE_USERS_ENDPOINT
: endpoint to use for the users storage. Default: `"localhost:9157"`.


-storage-users-mount-path |  $STORAGE_USERS_MOUNT_PATH
: mount path. Default: `"/users"`.


-storage-users-mount-id |  $STORAGE_USERS_MOUNT_ID
: mount id. Default: `"1284d238-aa92-42ce-bdc4-0b0000009157"`.


-public-link-endpoint |  $STORAGE_PUBLIC_LINK_ENDPOINT
: endpoint to use for the public links service. Default: `"localhost:9178"`.


-storage-public-link-mount-path |  $STORAGE_PUBLIC_LINK_MOUNT_PATH
: mount path. Default: `"/public"`.






















### storage storage-public-link

Start storage-public-link service

Usage: `storage storage-public-link [command options] [arguments...]`


































-debug-addr |  $STORAGE_PUBLIC_LINK_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9179"`.


-network |  $STORAGE_PUBLIC_LINK_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_PUBLIC_LINK_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9178"`.


-mount-path |  $STORAGE_PUBLIC_LINK_MOUNT_PATH
: mount path. Default: `"/public"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.





















































































































































































































### storage auth-basic

Start authprovider for basic auth

Usage: `storage auth-basic [command options] [arguments...]`











































































































































-debug-addr |  $STORAGE_AUTH_BASIC_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9147"`.


-auth-driver |  $STORAGE_AUTH_DRIVER
: auth driver: 'demo', 'json' or 'ldap'. Default: `"ldap"`.


-auth-json |  $STORAGE_AUTH_JSON
: Path to users.json file. Default: `""`.


-network |  $STORAGE_AUTH_BASIC_GRPC_NETWORK
: Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_AUTH_BASIC_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9146"`.


-service |  $STORAGE_AUTH_BASIC_SERVICES
: --service authprovider [--service otherservice]. Default: `cli.NewStringSlice("authprovider")`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.










































































































### storage groups

Start groups service

Usage: `storage groups [command options] [arguments...]`

















































































































-debug-addr |  $STORAGE_GROUPPROVIDER_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9161"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-network |  $STORAGE_GROUPPROVIDER_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_GROUPPROVIDER_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9160"`.


-endpoint |  $STORAGE_GROUPPROVIDER_ENDPOINT
: URL to use for the storage service. Default: `"localhost:9160"`.


-service |  $STORAGE_GROUPPROVIDER_SERVICES
: --service groupprovider [--service otherservice]. Default: `cli.NewStringSlice("groupprovider")`.


-driver |  $STORAGE_GROUPPROVIDER_DRIVER
: group driver: 'json', 'ldap', or 'rest'. Default: `"ldap"`.


-json-config |  $STORAGE_GROUPPROVIDER_JSON
: Path to groups.json file. Default: `""`.


-group-members-cache-expiration |  $STORAGE_GROUP_CACHE_EXPIRATION
: Time in minutes for redis cache expiration.. Default: `5`.


































































































































### storage storage-metadata

Start storage-metadata service

Usage: `storage storage-metadata [command options] [arguments...]`



-config-file |  $STORAGE_CONFIG_FILE
: Path to config file.


-log-level |  $STORAGE_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $STORAGE_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $STORAGE_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





























































































-debug-addr |  $STORAGE_METADATA_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9217"`.


-grpc-network |  $STORAGE_METADATA_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_METADATA_GRPC_PROVIDER_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9215"`.


-data-server-url |  $STORAGE_METADATA_DATA_SERVER_URL
: URL of the data-provider the storage-provider uses. Default: `"http://localhost:9216"`.


-http-network |  $STORAGE_METADATA_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_METADATA_HTTP_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9216"`.


-tmp-folder |  $STORAGE_METADATA_TMP_FOLDER
: path to tmp folder. Default: `flags.OverrideDefaultString(cfg.Reva.StorageMetadata.TempFolder, path.Join(defaults.BaseDataPath(), "tmp", "metadata"))`.


-driver |  $STORAGE_METADATA_DRIVER
: storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3. Default: `"ocis"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-userprovider-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the userprovider service. Default: `"localhost:9144"`.
















































































































































### storage storage-users

Start storage-users service

Usage: `storage storage-users [command options] [arguments...]`




























































































































-debug-addr |  $STORAGE_USERS_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9159"`.


-grpc-network |  $STORAGE_USERS_GRPC_NETWORK
: Network to use for the users storage, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_USERS_GRPC_ADDR
: GRPC Address to bind users storage. Default: `"127.0.0.1:9157"`.


-http-network |  $STORAGE_USERS_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_USERS_HTTP_ADDR
: HTTP Address to bind users storage. Default: `"127.0.0.1:9158"`.


-driver |  $STORAGE_USERS_DRIVER
: storage driver for users mount: eg. local, eos, owncloud, ocis or s3. Default: `"ocis"`.


-read-only |  $STORAGE_USERS_READ_ONLY , $OCIS_STORAGE_READ_ONLY
: use storage driver in read-only mode. Default: `false`.


-mount-path |  $STORAGE_USERS_MOUNT_PATH
: mount path. Default: `"/users"`.


-mount-id |  $STORAGE_USERS_MOUNT_ID
: mount id. Default: `"1284d238-aa92-42ce-bdc4-0b0000009157"`.


-expose-data-server |  $STORAGE_USERS_EXPOSE_DATA_SERVER
: exposes a dedicated data server. Default: `false`.


-data-server-url |  $STORAGE_USERS_DATA_SERVER_URL
: data server url. Default: `"http://localhost:9158/data"`.


-http-prefix |  $STORAGE_USERS_HTTP_PREFIX
: prefix for the http endpoint, without leading slash. Default: `"data"`.


-tmp-folder |  $STORAGE_USERS_TMP_FOLDER
: path to tmp folder. Default: `flags.OverrideDefaultString(cfg.Reva.StorageUsers.TempFolder, path.Join(defaults.BaseDataPath(), "tmp", "users"))`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-users-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9144"`.

















































































































### storage app-provider

Start appprovider for providing apps

Usage: `storage app-provider [command options] [arguments...]`











































































































































































































































-debug-addr |  $APP_PROVIDER_BASIC_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9165"`.


-network |  $APP_PROVIDER_BASIC_GRPC_NETWORK
: Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $APP_PROVIDER_BASIC_GRPC_ADDR
: Address to bind storage service. Default: `"127.0.0.1:9164"`.


-external-addr |  $APP_PROVIDER_BASIC_EXTERNAL_ADDR
: Address to connect to the storage service for other services. Default: `"127.0.0.1:9164"`.


-service |  $APP_PROVIDER_BASIC_SERVICES
: --service appprovider [--service otherservice]. Default: `cli.NewStringSlice("appprovider")`.


-driver |  $APP_PROVIDER_DRIVER
: Driver to use for app provider. Default: `""`.


-wopi-driver-app-apikey |  $APP_PROVIDER_WOPI_DRIVER_APP_API_KEY
: The API key used by the app, if applicable.. Default: `""`.


-wopi-driver-app-desktop-only |  $APP_PROVIDER_WOPI_DRIVER_APP_DESKTOP_ONLY
: Whether the app can be opened only on desktop. Default: `false`.


-wopi-driver-app-icon-uri |  $APP_PROVIDER_WOPI_DRIVER_APP_ICON_URI
: IOP Secret (Shared with WOPI server). Default: `""`.


-wopi-driver-app-internal-url |  $APP_PROVIDER_WOPI_DRIVER_APP_INTERNAL_URL
: The internal app URL in case of dockerized deployments. Defaults to AppURL. Default: `""`.


-wopi-driver-app-name |  $APP_PROVIDER_WOPI_DRIVER_APP_NAME
: The App user-friendly name.. Default: `""`.


-wopi-driver-app-url |  $APP_PROVIDER_WOPI_DRIVER_APP_URL
: App server URL. Default: `""`.


-wopi-driver-insecure |  $APP_PROVIDER_WOPI_DRIVER_INSECURE
: Disable SSL certificate verification of WOPI server and WOPI bridge. Default: `false`.


-wopi-driver-iopsecret |  $APP_PROVIDER_WOPI_DRIVER_IOP_SECRET
: IOP Secret (Shared with WOPI server). Default: `""`.


-wopi-driver-wopiurl |  $APP_PROVIDER_WOPI_DRIVER_WOPI_URL
: WOPI server URL. Default: `""`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.

### storage health

Check health status

Usage: `storage health [command options] [arguments...]`


-debug-addr |  $STORAGE_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9109"`.

























































































































































































































































### storage storage

Storage service for oCIS

Usage: `storage storage [command options] [arguments...]`



-config-file |  $STORAGE_CONFIG_FILE
: Path to config file.


-log-level |  $STORAGE_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $STORAGE_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $STORAGE_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





















































































































































































































































## Config for the different Storage Drivers

You can set different storage drivers for the Storage Providers. Please check the storage provider configuration.

Example: Set the home and users Storage Provider to `ocis`

`STORAGE_HOME_DRIVER=ocis`
`STORAGE_USERS_DRIVER=ocis`

### Local Driver

### Eos Driver

### owCloud Driver

### ownCloudSQL Driver

### Ocis Driver

### S3ng Driver

