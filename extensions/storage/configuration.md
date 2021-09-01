---
title: "Configuration"
date: "2021-09-01T16:20:09+0000"
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

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

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

### storage auth-basic

Start authprovider for basic auth

Usage: `storage auth-basic [command options] [arguments...]`













































































































































-debug-addr |  $STORAGE_AUTH_BASIC_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9147"`.


-auth-driver |  $STORAGE_AUTH_DRIVER
: auth driver: 'demo', 'json' or 'ldap'. Default: `"ldap"`.


-auth-json |  $STORAGE_AUTH_JSON
: Path to users.json file. Default: `""`.


-network |  $STORAGE_AUTH_BASIC_GRPC_NETWORK
: Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_AUTH_BASIC_GRPC_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9146"`.


-service |  $STORAGE_AUTH_BASIC_SERVICES
: --service authprovider [--service otherservice]. Default: `cli.NewStringSlice("authprovider")`.


-gateway-url |  $STORAGE_GATEWAY_ENDPOINT
: URL to use for the storage gateway service. Default: `"localhost:9142"`.






































































































### storage frontend

Start frontend service

Usage: `storage frontend [command options] [arguments...]`


























-debug-addr |  $STORAGE_FRONTEND_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9141"`.


-transfer-secret |  $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `"replace-me-with-a-transfer-secret"`.


-webdav-namespace |  $STORAGE_WEBDAV_NAMESPACE
: Namespace prefix for the /webdav endpoint. Default: `"/home/"`.


-dav-files-namespace |  $STORAGE_DAV_FILES_NAMESPACE
: Namespace prefix for the webdav /dav/files endpoint. Default: `"/users/"`.


-network |  $STORAGE_FRONTEND_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_FRONTEND_HTTP_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9140"`.


-public-url |  $STORAGE_FRONTEND_PUBLIC_URL , $OCIS_URL
: URL to use for the storage service. Default: `"https://localhost:9200"`.


-service |  $STORAGE_FRONTEND_SERVICES
: --service ocdav [--service ocs]. Default: `cli.NewStringSlice("datagateway", "ocdav", "ocs")`.


-datagateway-prefix |  $STORAGE_FRONTEND_DATAGATEWAY_PREFIX
: datagateway prefix. Default: `"data"`.


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


-gateway-url |  $STORAGE_GATEWAY_ENDPOINT
: URL to use for the storage gateway service. Default: `"localhost:9142"`.


-default-upload-protocol |  $STORAGE_FRONTEND_DEFAULT_UPLOAD_PROTOCOL
: Default upload chunking protocol to be used out of tus/v1/ng. Default: `"tus"`.


-upload-max-chunk-size |  $STORAGE_FRONTEND_UPLOAD_MAX_CHUNK_SIZE
: Max chunk size in bytes to advertise to clients through capabilities, or 0 for unlimited. Default: `0`.


-upload-http-method-override |  $STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE
: Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH. Default: `""`.


-checksum-suppored-type |  $STORAGE_FRONTEND_CHECKSUM_SUPPORTED_TYPES
: --checksum-suppored-type sha1 [--checksum-suppored-type adler32]. Default: `cli.NewStringSlice("sha1", "md5", "adler32")`.


-checksum-preferred-upload-type |  $STORAGE_FRONTEND_CHECKSUM_PREFERRED_UPLOAD_TYPE
: Specify the preferred checksum algorithm used for uploads. Default: `""`.


-user-agent-whitelist-lock-in |  $STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT
: --user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of comma separated [UserAgent:challenge] values, it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the storage-frontend it configures Reva..










































































































































































































### storage storage-home

Start storage-home service

Usage: `storage storage-home [command options] [arguments...]`























































































































































































-debug-addr |  $STORAGE_HOME_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9156"`.


-grpc-network |  $STORAGE_HOME_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_HOME_GRPC_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9154"`.


-http-network |  $STORAGE_HOME_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_HOME_HTTP_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9155"`.


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
: path to tmp folder. Default: `"/var/tmp/ocis/tmp/home"`.


-enable-home |  $STORAGE_HOME_ENABLE_HOME
: enable the creation of home directories. Default: `true`.


-gateway-endpoint |  $STORAGE_GATEWAY_ENDPOINT
: endpoint to use for the storage gateway service. Default: `"localhost:9142"`.


-users-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9144"`.



















































### storage users

Start users service

Usage: `storage users [command options] [arguments...]`


-debug-addr |  $STORAGE_USERPROVIDER_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9145"`.


-network |  $STORAGE_USERPROVIDER_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_USERPROVIDER_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9144"`.


-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: URL to use for the storage service. Default: `"localhost:9144"`.


-service |  $STORAGE_USERPROVIDER_SERVICES
: --service userprovider [--service otherservice]. Default: `cli.NewStringSlice("userprovider")`.


-driver |  $STORAGE_USERPROVIDER_DRIVER
: user driver: 'demo', 'json', 'ldap', or 'rest'. Default: `"ldap"`.


-json-config |  $STORAGE_USERPROVIDER_JSON
: Path to users.json file. Default: `""`.


-user-groups-cache-expiration |  $STORAGE_USER_CACHE_EXPIRATION
: Time in minutes for redis cache expiration.. Default: `5`.
















































































































































































































































### storage groups

Start groups service

Usage: `storage groups [command options] [arguments...]`










































































































































































-debug-addr |  $STORAGE_GROUPPROVIDER_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9161"`.


-network |  $STORAGE_GROUPPROVIDER_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_GROUPPROVIDER_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9160"`.


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






































































































































































































### storage sharing

Start sharing service

Usage: `storage sharing [command options] [arguments...]`





















































































































































-debug-addr |  $STORAGE_SHARING_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9151"`.


-network |  $STORAGE_SHARING_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_SHARING_GRPC_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9150"`.


-service |  $STORAGE_SHARING_SERVICES
: --service usershareprovider [--service publicshareprovider]. Default: `cli.NewStringSlice("usershareprovider", "publicshareprovider")`.


-user-driver |  $STORAGE_SHARING_USER_DRIVER
: driver to use for the UserShareProvider. Default: `"json"`.


-user-json-file |  $STORAGE_SHARING_USER_JSON_FILE
: file used to persist shares for the UserShareProvider. Default: `"/var/tmp/ocis/storage/shares.json"`.


-public-driver |  $STORAGE_SHARING_PUBLIC_DRIVER
: driver to use for the PublicShareProvider. Default: `"json"`.


-public-json-file |  $STORAGE_SHARING_PUBLIC_JSON_FILE
: file used to persist shares for the PublicShareProvider. Default: `"/var/tmp/ocis/storage/publicshares.json"`.


-public-password-hash-cost |  $STORAGE_SHARING_PUBLIC_PASSWORD_HASH_COST
: the cost of hashing the public shares passwords. Default: `11`.


-public-enable-expired-shares-cleanup |  $STORAGE_SHARING_PUBLIC_ENABLE_EXPIRED_SHARES_CLEANUP
: whether to periodically delete expired public shares. Default: `true`.


-public-janitor-run-interval |  $STORAGE_SHARING_PUBLIC_JANITOR_RUN_INTERVAL
: the time period in seconds after which to start a janitor run. Default: `60`.


























































































### storage health

Check health status

Usage: `storage health [command options] [arguments...]`


















































































































-debug-addr |  $STORAGE_DEBUG_ADDR
: Address to debug endpoint. Default: `"0.0.0.0:9109"`.







































































































































### storage auth-bearer

Start authprovider for bearer auth

Usage: `storage auth-bearer [command options] [arguments...]`








































































-debug-addr |  $STORAGE_AUTH_BEARER_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9149"`.


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
: Address to bind storage service. Default: `"0.0.0.0:9148"`.


-service |  $STORAGE_AUTH_BEARER_SERVICES
: --service authprovider [--service otherservice]. Default: `cli.NewStringSlice("authprovider")`.


-gateway-url |  $STORAGE_GATEWAY_ENDPOINT
: URL to use for the storage gateway service. Default: `"localhost:9142"`.








































































































































































### storage gateway

Start gateway

Usage: `storage gateway [command options] [arguments...]`

























































































































































































































-debug-addr |  $STORAGE_GATEWAY_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9143"`.


-transfer-secret |  $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `"replace-me-with-a-transfer-secret"`.


-transfer-expires |  $STORAGE_TRANSFER_EXPIRES
: Transfer token ttl in seconds. Default: `24 * 60 * 60`.


-network |  $STORAGE_GATEWAY_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_GATEWAY_GRPC_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9142"`.


-endpoint |  $STORAGE_GATEWAY_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9142"`.


-service |  $STORAGE_GATEWAY_SERVICES
: --service gateway [--service authregistry]. Default: `cli.NewStringSlice("gateway", "authregistry", "storageregistry")`.


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


-storage-registry-driver |  $STORAGE_STORAGE_REGISTRY_DRIVER
: driver of the storage registry. Default: `"static"`.


-storage-registry-rule |  $STORAGE_STORAGE_REGISTRY_RULES
: `Replaces the generated storage registry rules with this set: --storage-registry-rule "/eos=localhost:9158" [--storage-registry-rule "1284d238-aa92-42ce-bdc4-0b0000009162=localhost:9162"]`. Default: `cli.NewStringSlice()`.


-storage-home-provider |  $STORAGE_STORAGE_REGISTRY_HOME_PROVIDER
: mount point of the storage provider for user homes in the global namespace. Default: `"/home"`.


-storage-registry-json |  $STORAGE_STORAGE_REGISTRY_JSON
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

### storage storage-metadata

Start storage-metadata service

Usage: `storage storage-metadata [command options] [arguments...]`










-debug-addr |  $STORAGE_METADATA_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9217"`.


-grpc-network |  $STORAGE_METADATA_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_METADATA_GRPC_PROVIDER_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9215"`.


-data-server-url |  $STORAGE_METADATA_DATA_SERVER_URL
: URL of the data-provider the storage-provider uses. Default: `"http://localhost:9216"`.


-http-network |  $STORAGE_METADATA_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_METADATA_HTTP_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9216"`.


-tmp-folder |  $STORAGE_METADATA_TMP_FOLDER
: path to tmp folder. Default: `"/var/tmp/ocis/tmp/metadata"`.


-driver |  $STORAGE_METADATA_DRIVER
: storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3. Default: `"ocis"`.


-gateway-endpoint |  $STORAGE_GATEWAY_ENDPOINT
: endpoint to use for the gateway service. Default: `"localhost:9142"`.


-userprovider-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the userprovider service. Default: `"localhost:9144"`.


-storage-root |  $STORAGE_METADATA_ROOT
: the path to the metadata storage root. Default: `"/var/tmp/ocis/storage/metadata"`.





























-config-file |  $STORAGE_CONFIG_FILE
: Path to config file.


-log-level |  $STORAGE_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $STORAGE_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $STORAGE_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.






































































































































































































### storage storage-public-link

Start storage-public-link service

Usage: `storage storage-public-link [command options] [arguments...]`





















-debug-addr |  $STORAGE_PUBLIC_LINK_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9179"`.


-network |  $STORAGE_PUBLIC_LINK_GRPC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-addr |  $STORAGE_PUBLIC_LINK_GRPC_ADDR
: Address to bind storage service. Default: `"0.0.0.0:9178"`.


-mount-path |  $STORAGE_PUBLIC_LINK_MOUNT_PATH
: mount path. Default: `"/public"`.


-gateway-endpoint |  $STORAGE_GATEWAY_ENDPOINT
: endpoint to use for the storage gateway service. Default: `"localhost:9142"`.
































































































































































































































### storage storage-users

Start storage-users service

Usage: `storage storage-users [command options] [arguments...]`




















































-debug-addr |  $STORAGE_USERS_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9159"`.


-grpc-network |  $STORAGE_USERS_GRPC_NETWORK
: Network to use for the users storage, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-grpc-addr |  $STORAGE_USERS_GRPC_ADDR
: GRPC Address to bind users storage. Default: `"0.0.0.0:9157"`.


-http-network |  $STORAGE_USERS_HTTP_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `"tcp"`.


-http-addr |  $STORAGE_USERS_HTTP_ADDR
: HTTP Address to bind users storage. Default: `"0.0.0.0:9158"`.


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
: path to tmp folder. Default: `"/var/tmp/ocis/tmp/users"`.


-gateway-endpoint |  $STORAGE_GATEWAY_ENDPOINT
: endpoint to use for the storage gateway service. Default: `"localhost:9142"`.


-users-endpoint |  $STORAGE_USERPROVIDER_ENDPOINT
: endpoint to use for the storage service. Default: `"localhost:9144"`.























































































































































































## Config for the different Storage Drivers

You can set different storage drivers for the Storage Providers. Please check the storage provider configuration.

Example: Set the home and users Storage Provider to `ocis`

`STORAGE_HOME_DRIVER=ocis`
`STORAGE_USERS_DRIVER=ocis`

### Local Driver

-storage-local-root |  $STORAGE_DRIVER_LOCAL_ROOT
: the path to the local storage root. Default: `"/var/tmp/ocis/storage/local"`.

### Eos Driver

-storage-eos-namespace |  $STORAGE_DRIVER_EOS_NAMESPACE
: Namespace for metadata operations. Default: `"/eos/dockertest/reva"`.

-storage-eos-shadow-namespace |  $STORAGE_DRIVER_EOS_SHADOW_NAMESPACE
: Shadow namespace where share references are stored.

-storage-eos-share-folder |  $STORAGE_DRIVER_EOS_SHARE_FOLDER
: name of the share folder. Default: `"/Shares"`.

-storage-eos-binary |  $STORAGE_DRIVER_EOS_BINARY
: Location of the eos binary. Default: `"/usr/bin/eos"`.

-storage-eos-xrdcopy-binary |  $STORAGE_DRIVER_EOS_XRDCOPY_BINARY
: Location of the xrdcopy binary. Default: `"/usr/bin/xrdcopy"`.

-storage-eos-master-url |  $STORAGE_DRIVER_EOS_MASTER_URL
: URL of the Master EOS MGM. Default: `"root://eos-mgm1.eoscluster.cern.ch:1094"`.

-storage-eos-slave-url |  $STORAGE_DRIVER_EOS_SLAVE_URL
: URL of the Slave EOS MGM. Default: `"root://eos-mgm1.eoscluster.cern.ch:1094"`.

-storage-eos-cache-directory |  $STORAGE_DRIVER_EOS_CACHE_DIRECTORY
: Location on the local fs where to store reads. Default: `flags.OverrideDefaultString(cfg.Reva.Storages.EOS.CacheDirectory, os.TempDir())`.

-storage-eos-enable-logging |  $STORAGE_DRIVER_EOS_ENABLE_LOGGING
: Enables logging of the commands executed.

-storage-eos-show-hidden-sysfiles |  $STORAGE_DRIVER_EOS_SHOW_HIDDEN_SYSFILES
: show internal EOS files like .sys.v# and .sys.a# files..

-storage-eos-force-singleuser-mode |  $STORAGE_DRIVER_EOS_FORCE_SINGLEUSER_MODE
: force connections to EOS to use SingleUsername.

-storage-eos-use-keytab |  $STORAGE_DRIVER_EOS_USE_KEYTAB
: authenticate requests by using an EOS keytab.

-storage-eos-enable-home |  $STORAGE_DRIVER_EOS_ENABLE_HOME
: enable the creation of home directories.

-storage-eos-sec-protocol |  $STORAGE_DRIVER_EOS_SEC_PROTOCOL
: the xrootd security protocol to use between the server and EOS.

-storage-eos-keytab |  $STORAGE_DRIVER_EOS_KEYTAB
: the location of the keytab to use to authenticate to EOS.

-storage-eos-single-username |  $STORAGE_DRIVER_EOS_SINGLE_USERNAME
: the username to use when SingleUserMode is enabled.

-storage-eos-layout |  $STORAGE_DRIVER_EOS_LAYOUT
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `"{{substr 0 1 .Username}}/{{.Username}}"`.

-storage-eos-gatewaysvc |  $STORAGE_DRIVER_EOS_GATEWAYSVC
: URL to use for the storage gateway service. Default: `"localhost:9142"`.

### owCloud Driver

-storage-owncloud-datadir |  $STORAGE_DRIVER_OWNCLOUD_DATADIR
: the path to the owncloud data directory. Default: `"/var/tmp/ocis/storage/owncloud"`.

-storage-owncloud-uploadinfo-dir |  $STORAGE_DRIVER_OWNCLOUD_UPLOADINFO_DIR
: the path to the tus upload info directory. Default: `"/var/tmp/ocis/storage/uploadinfo"`.

-storage-owncloud-share-folder |  $STORAGE_DRIVER_OWNCLOUD_SHARE_FOLDER
: name of the shares folder. Default: `"/Shares"`.

-storage-owncloud-scan |  $STORAGE_DRIVER_OWNCLOUD_SCAN
: scan files on startup to add fileids. Default: `true`.

-storage-owncloud-redis |  $STORAGE_DRIVER_OWNCLOUD_REDIS_ADDR
: the address of the redis server. Default: `":6379"`.

-storage-owncloud-enable-home |  $STORAGE_DRIVER_OWNCLOUD_ENABLE_HOME
: enable the creation of home storages. Default: `false`.

-storage-owncloud-layout |  $STORAGE_DRIVER_OWNCLOUD_LAYOUT
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `. Default: `"{{.Id.OpaqueId}}"`.

### ownCloudSQL Driver

-storage-owncloudsql-datadir |  $STORAGE_DRIVER_OWNCLOUDSQL_DATADIR
: the path to the owncloudsql data directory. Default: `"/var/tmp/ocis/storage/owncloud"`.

-storage-owncloudsql-uploadinfo-dir |  $STORAGE_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR
: the path to the tus upload info directory. Default: `"/var/tmp/ocis/storage/uploadinfo"`.

-storage-owncloudsql-share-folder |  $STORAGE_DRIVER_OWNCLOUDSQL_SHARE_FOLDER
: name of the shares folder. Default: `"/Shares"`.

-storage-owncloudsql-enable-home |  $STORAGE_DRIVER_OWNCLOUDSQL_ENABLE_HOME
: enable the creation of home storages. Default: `false`.

-storage-owncloudsql-layout |  $STORAGE_DRIVER_OWNCLOUDSQL_LAYOUT
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `. Default: `"{{.Username}}"`.

-storage-owncloudsql-dbusername |  $STORAGE_DRIVER_OWNCLOUDSQL_DBUSERNAME
: `"username for accessing the database" `. Default: `"owncloud"`.

-storage-owncloudsql-dbpassword |  $STORAGE_DRIVER_OWNCLOUDSQL_DBPASSWORD
: `"password for accessing the database" `. Default: `"owncloud"`.

-storage-owncloudsql-dbhost |  $STORAGE_DRIVER_OWNCLOUDSQL_DBHOST
: `"the database hostname or IP address" `. Default: `cfg.Reva.Storages.OwnCloudSQL.DBHost`.

-storage-owncloudsql-dbport |  $STORAGE_DRIVER_OWNCLOUDSQL_DBPORT
: `"port the database listens on" `. Default: `3306`.

-storage-owncloudsql-dbname |  $STORAGE_DRIVER_OWNCLOUDSQL_DBNAME
: `"the database name" `. Default: `"owncloud"`.

### Ocis Driver

-storage-ocis-root |  $STORAGE_DRIVER_OCIS_ROOT
: the path to the local storage root. Default: `"/var/tmp/ocis/storage/users"`.

-storage-ocis-enable-home |  $STORAGE_DRIVER_OCIS_ENABLE_HOME
: enable the creation of home storages. Default: `false`.

-storage-ocis-layout |  $STORAGE_DRIVER_OCIS_LAYOUT
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `. Default: `"{{.Id.OpaqueId}}"`.

-service-user-uuid |  $STORAGE_DRIVER_OCIS_SERVICE_USER_UUID
: uuid of the internal service user. Default: `95cb8724-03b2-11eb-a0a6-c33ef8ef53ad`.

### S3ng Driver

-storage-s3ng-root |  $STORAGE_DRIVER_S3NG_ROOT
: the path to the local storage root. Default: `"/var/tmp/ocis/storage/users"`.

-storage-s3ng-enable-home |  $STORAGE_DRIVER_S3NG_ENABLE_HOME
: enable the creation of home storages. Default: `false`.

-storage-s3ng-layout |  $STORAGE_DRIVER_S3NG_LAYOUT
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `. Default: `"{{.Id.OpaqueId}}"`.

-storage-s3ng-region |  $STORAGE_DRIVER_S3NG_REGION
: `"the s3 region" `. Default: `default`.

-storage-s3ng-accesskey |  $STORAGE_DRIVER_S3NG_ACCESS_KEY
: `"the s3 access key" `.

-storage-s3ng-secretkey |  $STORAGE_DRIVER_S3NG_SECRET_KEY
: `"the secret s3 api key" `.

-storage-s3ng-endpoint |  $STORAGE_DRIVER_S3NG_ENDPOINT
: `"s3 compatible API endpoint" `.

-storage-s3ng-bucket |  $STORAGE_DRIVER_S3NG_BUCKET
: `"bucket where the data will be stored in`.

