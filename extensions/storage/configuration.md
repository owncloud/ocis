---
title: "Configuration"
date: "2020-10-09T04:02:19+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: configuration.md
---

{{< toc >}}

## Configuration

oCIS Single Binary is not responsible for configuring extensions. Instead, each extension could either be configured by environment variables, cli flags or config files.

Each extension has its dedicated documentation page (e.g. https://owncloud.github.io/extensions/ocis_proxy/configuration) which lists all possible configurations. Config files and environment variables are picked up if you use the `./bin/ocis server` command within the oCIS single binary. Command line flags must be set explicitly on the extensions subcommands.

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-proxy reads `proxy.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Storage service for oCIS

Usage: `storage [global options] command [command options] [arguments...]`

--config-file | $STORAGE_CONFIG_FILE
: Path to config file.

--log-level | $STORAGE_LOG_LEVEL
: Set logging level. Default: `info`.

--log-pretty | $STORAGE_LOG_PRETTY
: Enable pretty logging.

--log-color | $STORAGE_LOG_COLOR
: Enable colored logging.

## Sub Commands

### storage storage-eos-data

Start storage-eos-data service

Usage: `storage storage-eos-data [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_OC_DATA_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9161`.

--network | $STORAGE_STORAGE_EOS_DATA_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_EOS_DATA_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `http`.

--addr | $STORAGE_STORAGE_EOS_DATA_ADDR
: Address to bind storage service. Default: `0.0.0.0:9160`.

--url | $STORAGE_STORAGE_EOS_DATA_URL
: URL to use for the storage service. Default: `localhost:9160`.

--driver | $STORAGE_STORAGE_EOS_DATA_DRIVER
: storage driver for eos data mount: eg. local, eos, owncloud, ocis or s3. Default: `eos`.

--prefix | $STORAGE_STORAGE_EOS_DATA_PREFIX
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $STORAGE_STORAGE_EOS_DATA_TEMP_FOLDER
: temp folder. Default: `/var/tmp/`.

--gateway-url | $STORAGE_GATEWAY_URL
: URL to use for the storage gateway service. Default: `localhost:9142`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

### storage frontend

Start frontend service

Usage: `storage frontend [command options] [arguments...]`

--debug-addr | $STORAGE_FRONTEND_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9141`.

--transfer-secret | $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `replace-me-with-a-transfer-secret`.

--webdav-namespace | $WEBDAV_NAMESPACE
: Namespace prefix for the /webdav endpoint. Default: `/home/`.

--dav-files-namespace | $DAV_FILES_NAMESPACE
: Namespace prefix for the webdav /dav/files endpoint. Default: `/oc/`.

--network | $STORAGE_FRONTEND_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_FRONTEND_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `http`.

--addr | $STORAGE_FRONTEND_ADDR
: Address to bind storage service. Default: `0.0.0.0:9140`.

--url | $STORAGE_FRONTEND_URL
: URL to use for the storage service. Default: `https://localhost:9200`.

--datagateway-prefix | $STORAGE_FRONTEND_DATAGATEWAY_PREFIX
: datagateway prefix. Default: `data`.

--ocdav-prefix | $STORAGE_FRONTEND_OCDAV_PREFIX
: owncloud webdav endpoint prefix.

--ocs-prefix | $STORAGE_FRONTEND_OCS_PREFIX
: open collaboration services endpoint prefix. Default: `ocs`.

--gateway-url | $STORAGE_GATEWAY_URL
: URL to use for the storage gateway service. Default: `localhost:9142`.

--upload-disable-tus | $STORAGE_FRONTEND_UPLOAD_DISABLE_TUS
: Disables TUS upload mechanism. Default: `false`.

--upload-http-method-override | $STORAGE_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE
: Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH.

### storage storage-home-data

Start storage-home-data service

Usage: `storage storage-home-data [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_HOME_DATA_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9157`.

--network | $STORAGE_STORAGE_HOME_DATA_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_HOME_DATA_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `http`.

--addr | $STORAGE_STORAGE_HOME_DATA_ADDR
: Address to bind storage service. Default: `0.0.0.0:9156`.

--url | $STORAGE_STORAGE_HOME_DATA_URL
: URL to use for the storage service. Default: `localhost:9156`.

--driver | $STORAGE_STORAGE_HOME_DATA_DRIVER
: storage driver for home data mount: eg. local, eos, owncloud, ocis or s3. Default: `owncloud`.

--prefix | $STORAGE_STORAGE_HOME_DATA_PREFIX
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $STORAGE_STORAGE_HOME_DATA_TEMP_FOLDER
: temp folder. Default: `/var/tmp/`.

--enable-home | $STORAGE_STORAGE_HOME_ENABLE_HOME
: enable the creation of home directories. Default: `true`.

--gateway-url | $STORAGE_GATEWAY_URL
: URL to use for the storage gateway service. Default: `localhost:9142`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

### storage gateway

Start gateway

Usage: `storage gateway [command options] [arguments...]`

--debug-addr | $STORAGE_GATEWAY_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9143`.

--transfer-secret | $STORAGE_TRANSFER_SECRET
: Transfer secret for datagateway. Default: `replace-me-with-a-transfer-secret`.

--network | $STORAGE_GATEWAY_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_GATEWAY_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_GATEWAY_ADDR
: Address to bind storage service. Default: `0.0.0.0:9142`.

--url | $STORAGE_GATEWAY_URL
: URL to use for the storage service. Default: `localhost:9142`.

--commit-share-to-storage-grant | $STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT
: Commit shares to the share manager. Default: `true`.

--commit-share-to-storage-ref | $STORAGE_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF
: Commit shares to the storage. Default: `true`.

--share-folder | $STORAGE_GATEWAY_SHARE_FOLDER
: mount shares in this folder of the home storage provider. Default: `Shares`.

--disable-home-creation-on-login | $STORAGE_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN
: Disable creation of home folder on login.

--storage-registry-driver | $STORAGE_STORAGE_REGISTRY_DRIVER
: driver of the storage registry. Default: `static`.

--storage-home-provider | $STORAGE_STORAGE_HOME_PROVIDER
: mount point of the storage provider for user homes in the global namespace. Default: `/home`.

--frontend-url | $STORAGE_FRONTEND_URL
: URL to use for the storage service. Default: `https://localhost:9200`.

--datagateway-url | $STORAGE_DATAGATEWAY_URL
: URL to use for the storage datagateway. Default: `https://localhost:9200/data`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

--auth-basic-url | $STORAGE_AUTH_BASIC_URL
: URL to use for the storage service. Default: `localhost:9146`.

--auth-bearer-url | $STORAGE_AUTH_BEARER_URL
: URL to use for the storage service. Default: `localhost:9148`.

--sharing-url | $STORAGE_SHARING_URL
: URL to use for the storage service. Default: `localhost:9150`.

--storage-root-url | $STORAGE_STORAGE_ROOT_URL
: URL to use for the storage service. Default: `localhost:9152`.

--storage-root-mount-path | $STORAGE_STORAGE_ROOT_MOUNT_PATH
: mount path. Default: `/`.

--storage-root-mount-id | $STORAGE_STORAGE_ROOT_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009152`.

--storage-home-url | $STORAGE_STORAGE_HOME_URL
: URL to use for the storage service. Default: `localhost:9154`.

--storage-home-mount-path | $STORAGE_STORAGE_HOME_MOUNT_PATH
: mount path. Default: `/home`.

--storage-home-mount-id | $STORAGE_STORAGE_HOME_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009154`.

--storage-eos-url | $STORAGE_STORAGE_EOS_URL
: URL to use for the storage service. Default: `localhost:9158`.

--storage-eos-mount-path | $STORAGE_STORAGE_EOS_MOUNT_PATH
: mount path. Default: `/eos`.

--storage-eos-mount-id | $STORAGE_STORAGE_EOS_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009158`.

--storage-oc-url | $STORAGE_STORAGE_OC_URL
: URL to use for the storage service. Default: `localhost:9162`.

--storage-oc-mount-path | $STORAGE_STORAGE_OC_MOUNT_PATH
: mount path. Default: `/oc`.

--storage-oc-mount-id | $STORAGE_STORAGE_OC_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--public-link-url | $STORAGE_STORAGE_PUBLIC_LINK_URL
: URL to use for the public links service. Default: `localhost:9178`.

--storage-public-link-mount-path | $STORAGE_STORAGE_PUBLIC_LINK_MOUNT_PATH
: mount path. Default: `/public/`.

### storage storage-metadata

Start storage-metadata service

Usage: `storage storage-metadata [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_METADATA_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9217`.

--network | $STORAGE_STORAGE_METADATA_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--provider-addr | $STORAGE_STORAGE_METADATA_PROVIDER_ADDR
: Address to bind storage service. Default: `0.0.0.0:9215`.

--data-server-url | $STORAGE_STORAGE_METADATA_DATA_SERVER_URL
: URL of the data-server the storage-provider uses. Default: `http://localhost:9216`.

--data-server-addr | $STORAGE_STORAGE_METADATA_DATA_SERVER_ADDR
: Address to bind the metadata data-server to. Default: `0.0.0.0:9216`.

--storage-provider-driver | $STORAGE_STORAGE_METADATA_PROVIDER_DRIVER
: storage driver for metadata mount: eg. local, eos, owncloud, ocis or s3. Default: `local`.

--data-provider-driver | $STORAGE_STORAGE_METADATA_DATA_PROVIDER_DRIVER
: storage driver for data-provider mount: eg. local, eos, owncloud, ocis or s3. Default: `local`.

--storage-root | $STORAGE_STORAGE_METADATA_ROOT
: the path to the metadata storage root. Default: `/var/tmp/ocis/metadata`.

### storage auth-basic

Start authprovider for basic auth

Usage: `storage auth-basic [command options] [arguments...]`

--debug-addr | $STORAGE_AUTH_BASIC_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9147`.

--auth-driver | $STORAGE_AUTH_DRIVER
: auth driver: 'demo', 'json' or 'ldap'. Default: `ldap`.

--auth-json | $STORAGE_AUTH_JSON
: Path to users.json file.

--network | $STORAGE_AUTH_BASIC_NETWORK
: Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_AUTH_BASIC_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_AUTH_BASIC_ADDR
: Address to bind storage service. Default: `0.0.0.0:9146`.

--url | $STORAGE_AUTH_BASIC_URL
: URL to use for the storage service. Default: `localhost:9146`.

### storage storage-root

Start storage-root service

Usage: `storage storage-root [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_ROOT_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9153`.

--network | $STORAGE_STORAGE_ROOT_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_ROOT_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_STORAGE_ROOT_ADDR
: Address to bind storage service. Default: `0.0.0.0:9152`.

--url | $STORAGE_STORAGE_ROOT_URL
: URL to use for the storage service. Default: `localhost:9152`.

--driver | $STORAGE_STORAGE_ROOT_DRIVER
: storage driver for root mount: eg. local, eos, owncloud, ocis or s3. Default: `local`.

--mount-path | $STORAGE_STORAGE_ROOT_MOUNT_PATH
: mount path. Default: `/`.

--mount-id | $STORAGE_STORAGE_ROOT_MOUNT_ID
: mount id. Default: `123e4567-e89b-12d3-a456-426655440001`.

--expose-data-server | $STORAGE_STORAGE_ROOT_EXPOSE_DATA_SERVER
: exposes a dedicated data server.

--data-server-url | $STORAGE_STORAGE_ROOT_DATA_SERVER_URL
: data server url.

### storage storage-public-link

Start storage-public-link service

Usage: `storage storage-public-link [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_PUBLIC_LINK_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9179`.

--network | $STORAGE_STORAGE_PUBLIC_LINK_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_PUBLIC_LINK_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_STORAGE_PUBLIC_LINK_ADDR
: Address to bind storage service. Default: `0.0.0.0:9178`.

--url | $STORAGE_STORAGE_PUBLIC_LINK_URL
: Address to bind storage service. Default: `localhost:9178`.

--mount-path | $STORAGE_STORAGE_PUBLIC_LINK_MOUNT_PATH
: mount path. Default: `/public/`.

--gateway-url | $STORAGE_GATEWAY_URL
: URL to use for the storage gateway service. Default: `localhost:9142`.

### storage storage-home

Start storage-home service

Usage: `storage storage-home [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_HOME_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9155`.

--network | $STORAGE_STORAGE_HOME_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_HOME_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_STORAGE_HOME_ADDR
: Address to bind storage service. Default: `0.0.0.0:9154`.

--url | $STORAGE_STORAGE_HOME_URL
: URL to use for the storage service. Default: `localhost:9154`.

--driver | $STORAGE_STORAGE_HOME_DRIVER
: storage driver for home mount: eg. local, eos, owncloud, ocis or s3. Default: `owncloud`.

--mount-path | $STORAGE_STORAGE_HOME_MOUNT_PATH
: mount path. Default: `/home`.

--mount-id | $STORAGE_STORAGE_HOME_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--expose-data-server | $STORAGE_STORAGE_HOME_EXPOSE_DATA_SERVER
: exposes a dedicated data server. Default: `false`.

--data-server-url | $STORAGE_STORAGE_HOME_DATA_SERVER_URL
: data server url. Default: `http://localhost:9156/data`.

--enable-home | $STORAGE_STORAGE_HOME_ENABLE_HOME
: enable the creation of home directories. Default: `true`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

### storage sharing

Start sharing service

Usage: `storage sharing [command options] [arguments...]`

--debug-addr | $STORAGE_SHARING_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9151`.

--network | $STORAGE_SHARING_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_SHARING_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_SHARING_ADDR
: Address to bind storage service. Default: `0.0.0.0:9150`.

--url | $STORAGE_SHARING_URL
: URL to use for the storage service. Default: `localhost:9150`.

--user-driver | $STORAGE_SHARING_USER_DRIVER
: driver to use for the UserShareProvider. Default: `json`.

--user-json-file | $STORAGE_SHARING_USER_JSON_FILE
: file used to persist shares for the UserShareProvider. Default: `/var/tmp/ocis/shares.json`.

--public-driver | $STORAGE_SHARING_PUBLIC_DRIVER
: driver to use for the PublicShareProvider. Default: `json`.

### storage storage-eos

Start storage-eos service

Usage: `storage storage-eos [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_EOS_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9159`.

--network | $STORAGE_STORAGE_EOS_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_EOS_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_STORAGE_EOS_ADDR
: Address to bind storage service. Default: `0.0.0.0:9158`.

--url | $STORAGE_STORAGE_EOS_URL
: URL to use for the storage service. Default: `localhost:9158`.

--driver | $STORAGE_STORAGE_EOS_DRIVER
: storage driver for eos mount: eg. local, eos, owncloud, ocis or s3. Default: `eos`.

--mount-path | $STORAGE_STORAGE_EOS_MOUNT_PATH
: mount path. Default: `/eos`.

--mount-id | $STORAGE_STORAGE_EOS_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009158`.

--expose-data-server | $STORAGE_STORAGE_EOS_EXPOSE_DATA_SERVER
: exposes a dedicated data server. Default: `false`.

--data-server-url | $STORAGE_STORAGE_EOS_DATA_SERVER_URL
: data server url. Default: `http://localhost:9160/data`.

### storage storage

Storage service for oCIS

Usage: `storage storage [command options] [arguments...]`

--config-file | $STORAGE_CONFIG_FILE
: Path to config file.

--log-level | $STORAGE_LOG_LEVEL
: Set logging level. Default: `info`.

--log-pretty | $STORAGE_LOG_PRETTY
: Enable pretty logging.

--log-color | $STORAGE_LOG_COLOR
: Enable colored logging.

### storage auth-bearer

Start authprovider for bearer auth

Usage: `storage auth-bearer [command options] [arguments...]`

--debug-addr | $STORAGE_AUTH_BEARER_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9149`.

--oidc-issuer | $STORAGE_OIDC_ISSUER
: OIDC issuer. Default: `https://localhost:9200`.

--oidc-insecure | $STORAGE_OIDC_INSECURE
: OIDC allow insecure communication. Default: `true`.

--oidc-id-claim | $STORAGE_OIDC_ID_CLAIM
: OIDC id claim. Default: `preferred_username`.

--oidc-uid-claim | $STORAGE_OIDC_UID_CLAIM
: OIDC uid claim.

--oidc-gid-claim | $STORAGE_OIDC_GID_CLAIM
: OIDC gid claim.

--network | $STORAGE_AUTH_BEARER_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_AUTH_BEARER_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_AUTH_BEARER_ADDR
: Address to bind storage service. Default: `0.0.0.0:9148`.

--url | $STORAGE_AUTH_BEARER_URL
: URL to use for the storage service. Default: `localhost:9148`.

### storage storage-oc-data

Start storage-oc-data service

Usage: `storage storage-oc-data [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_OC_DATA_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9165`.

--network | $STORAGE_STORAGE_OC_DATA_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_OC_DATA_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `http`.

--addr | $STORAGE_STORAGE_OC_DATA_ADDR
: Address to bind storage service. Default: `0.0.0.0:9164`.

--url | $STORAGE_STORAGE_OC_DATA_URL
: URL to use for the storage service. Default: `localhost:9164`.

--driver | $STORAGE_STORAGE_OC_DATA_DRIVER
: storage driver for oc data mount: eg. local, eos, owncloud, ocis or s3. Default: `owncloud`.

--prefix | $STORAGE_STORAGE_OC_DATA_PREFIX
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $STORAGE_STORAGE_OC_DATA_TEMP_FOLDER
: temp folder. Default: `/var/tmp/`.

--gateway-url | $STORAGE_GATEWAY_URL
: URL to use for the storage gateway service. Default: `localhost:9142`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

### storage users

Start users service

Usage: `storage users [command options] [arguments...]`

--debug-addr | $STORAGE_SHARING_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9145`.

--network | $STORAGE_USERS_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_USERS_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_USERS_ADDR
: Address to bind storage service. Default: `0.0.0.0:9144`.

--url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

--driver | $STORAGE_USERS_DRIVER
: user driver: 'demo', 'json', 'ldap', or 'rest'. Default: `ldap`.

--json-config | $STORAGE_USERS_JSON
: Path to users.json file.

--rest-client-id | $STORAGE_REST_CLIENT_ID
: User rest driver Client ID.

--rest-client-secret | $STORAGE_REST_CLIENT_SECRET
: User rest driver Client Secret.

--rest-redis-address | $STORAGE_REST_REDIS_ADDRESS
: Address for redis server. Default: `localhost:6379`.

--rest-redis-username | $STORAGE_REST_REDIS_USERNAME
: Username for redis server.

--rest-redis-password | $STORAGE_REST_REDIS_PASSWORD
: Password for redis server.

--rest-id-provider | $STORAGE_REST_ID_PROVIDER
: The OIDC Provider.

--rest-api-base-url | $STORAGE_REST_API_BASE_URL
: Base API Endpoint.

--rest-oidc-token-endpoint | $STORAGE_REST_OIDC_TOKEN_ENDPOINT
: Endpoint to generate token to access the API.

--rest-target-api | $STORAGE_REST_TARGET_API
: The target application.

### storage health

Check health status

Usage: `storage health [command options] [arguments...]`

--debug-addr | $STORAGE_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9109`.

### storage storage-oc

Start storage-oc service

Usage: `storage storage-oc [command options] [arguments...]`

--debug-addr | $STORAGE_STORAGE_OC_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9163`.

--network | $STORAGE_STORAGE_OC_NETWORK
: Network to use for the storage service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $STORAGE_STORAGE_OC_PROTOCOL
: protocol for storage service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $STORAGE_STORAGE_OC_ADDR
: Address to bind storage service. Default: `0.0.0.0:9162`.

--url | $STORAGE_STORAGE_OC_URL
: URL to use for the storage service. Default: `localhost:9162`.

--driver | $STORAGE_STORAGE_OC_DRIVER
: storage driver for oc mount: eg. local, eos, owncloud, ocis or s3. Default: `owncloud`.

--mount-path | $STORAGE_STORAGE_OC_MOUNT_PATH
: mount path. Default: `/oc`.

--mount-id | $STORAGE_STORAGE_OC_MOUNT_ID
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--expose-data-server | $STORAGE_STORAGE_OC_EXPOSE_DATA_SERVER
: exposes a dedicated data server. Default: `false`.

--data-server-url | $STORAGE_STORAGE_OC_DATA_SERVER_URL
: data server url. Default: `http://localhost:9164/data`.

--users-url | $STORAGE_USERS_URL
: URL to use for the storage service. Default: `localhost:9144`.

