---
title: "Configuration"
date: "2020-06-24T22:26:40+0200"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis-reva
geekdocEditPath: edit/master/docs
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

Example service for Reva/oCIS

Usage: `ocis-reva [global options] command [command options] [arguments...]`

--config-file | $REVA_CONFIG_FILE  
: Path to config file.

--log-level | $REVA_LOG_LEVEL  
: Set logging level. Default: `info`.

--log-pretty | $REVA_LOG_PRETTY  
: Enable pretty logging.

--log-color | $REVA_LOG_COLOR  
: Enable colored logging.

## Sub Commands

### ocis-reva users

Start reva users service

Usage: `ocis-reva users [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_SHARING_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9145`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--ldap-hostname | $REVA_LDAP_HOSTNAME  
: LDAP hostname. Default: `localhost`.

--ldap-base-dn | $REVA_LDAP_BASE_DN  
: LDAP basedn. Default: `dc=example,dc=org`.

--ldap-userfilter | $REVA_LDAP_USERFILTER  
: LDAP userfilter. Default: `(&(objectclass=posixAccount)(cn=%s*))`.

--ldap-groupfilter | $REVA_LDAP_GROUPFILTER  
: LDAP groupfilter. Default: `(&(objectclass=posixGroup)(cn=%s*))`.

--ldap-bind-dn | $REVA_LDAP_BIND_DN  
: LDAP bind dn. Default: `cn=reva,ou=sysusers,dc=example,dc=org`.

--ldap-bind-password | $REVA_LDAP_BIND_PASSWORD  
: LDAP bind password. Default: `reva`.

--ldap-idp | $REVA_LDAP_IDP  
: Identity provider to use for users. Default: `https://localhost:9200`.

--ldap-schema-uid | $REVA_LDAP_SCHEMA_UID  
: LDAP schema uid. Default: `uid`.

--ldap-schema-mail | $REVA_LDAP_SCHEMA_MAIL  
: LDAP schema mail. Default: `mail`.

--ldap-schema-displayName | $REVA_LDAP_SCHEMA_DISPLAYNAME  
: LDAP schema displayName. Default: `sn`.

--ldap-schema-cn | $REVA_LDAP_SCHEMA_CN  
: LDAP schema cn. Default: `cn`.

--network | $REVA_USERS_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_USERS_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_USERS_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9144`.

--url | $REVA_USERS_URL  
: URL to use for the reva service. Default: `localhost:9144`.

--driver | $REVA_USERS_DRIVER  
: user driver: 'demo', 'json' or 'ldap'. Default: `ldap`.

--json-config | $REVA_USERS_JSON  
: Path to users.json file.

### ocis-reva auth-basic

Start reva authprovider for basic auth

Usage: `ocis-reva auth-basic [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_AUTH_BASIC_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9147`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--users-driver | $REVA_USERS_DRIVER  
: user driver: 'demo', 'json' or 'ldap'. Default: `ldap`.

--users-json | $REVA_USERS_JSON  
: Path to users.json file.

--ldap-hostname | $REVA_LDAP_HOSTNAME  
: LDAP hostname. Default: `localhost`.

--ldap-base-dn | $REVA_LDAP_BASE_DN  
: LDAP basedn. Default: `dc=example,dc=org`.

--ldap-userfilter | $REVA_LDAP_USERFILTER  
: LDAP userfilter. Default: `(&(objectclass=posixAccount)(cn=%s))`.

--ldap-groupfilter | $REVA_LDAP_GROUPFILTER  
: LDAP groupfilter. Default: `(&(objectclass=posixGroup)(cn=%s))`.

--ldap-bind-dn | $REVA_LDAP_BIND_DN  
: LDAP bind dn. Default: `cn=reva,ou=sysusers,dc=example,dc=org`.

--ldap-bind-password | $REVA_LDAP_BIND_PASSWORD  
: LDAP bind password. Default: `reva`.

--ldap-idp | $REVA_LDAP_IDP  
: Identity provider to use for users. Default: `https://localhost:9200`.

--ldap-schema-uid | $REVA_LDAP_SCHEMA_UID  
: LDAP schema uid. Default: `uid`.

--ldap-schema-mail | $REVA_LDAP_SCHEMA_MAIL  
: LDAP schema mail. Default: `mail`.

--ldap-schema-displayName | $REVA_LDAP_SCHEMA_DISPLAYNAME  
: LDAP schema displayName. Default: `sn`.

--ldap-schema-cn | $REVA_LDAP_SCHEMA_CN  
: LDAP schema cn. Default: `cn`.

--network | $REVA_AUTH_BASIC_NETWORK  
: Network to use for the reva auth-basic service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_AUTH_BASIC_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_AUTH_BASIC_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9146`.

--url | $REVA_AUTH_BASIC_URL  
: URL to use for the reva service. Default: `localhost:9146`.

### ocis-reva storage-oc-data

Start reva storage-oc-data service

Usage: `ocis-reva storage-oc-data [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_OC_DATA_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9165`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_OC_DATA_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_OC_DATA_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `http`.

--addr | $REVA_STORAGE_OC_DATA_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9164`.

--url | $REVA_STORAGE_OC_DATA_URL  
: URL to use for the reva service. Default: `localhost:9164`.

--driver | $REVA_STORAGE_OC_DATA_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `owncloud`.

--prefix | $REVA_STORAGE_OC_DATA_PREFIX  
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $REVA_STORAGE_OC_DATA_TEMP_FOLDER  
: temp folder. Default: `/var/tmp/`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-example.org`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-example.org`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `false`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

--gateway-url | $REVA_GATEWAY_URL  
: URL to use for the reva gateway service. Default: `localhost:9142`.

### ocis-reva frontend

Start reva frontend service

Usage: `ocis-reva frontend [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_FRONTEND_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9141`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--transfer-secret | $REVA_TRANSFER_SECRET  
: Transfer secret for datagateway. Default: `replace-me-with-a-transfer-secret`.

--webdav-namespace | $WEBDAV_NAMESPACE  
: Namespace prefix for the /webdav endpoint. Default: `/home/`.

--dav-files-namespace | $DAV_FILES_NAMESPACE  
: Namespace prefix for the webdav /dav/files endpoint. Default: `/oc/`.

--network | $REVA_FRONTEND_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_FRONTEND_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `http`.

--addr | $REVA_FRONTEND_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9140`.

--url | $REVA_FRONTEND_URL  
: URL to use for the reva service. Default: `https://localhost:9200`.

--gateway-url | $REVA_GATEWAY_URL  
: URL to use for the reva gateway service. Default: `localhost:9142`.

--upload-disable-tus | $REVA_FRONTEND_UPLOAD_DISABLE_TUS  
: Disables TUS upload mechanism. Default: `false`.

--upload-http-method-override | $REVA_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE  
: Specify an HTTP method (ex: POST) that clients should to use when uploading instead of PATCH.

### ocis-reva storage-root

Start reva storage-root service

Usage: `ocis-reva storage-root [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_ROOT_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9153`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_ROOT_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_ROOT_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_STORAGE_ROOT_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9152`.

--url | $REVA_STORAGE_ROOT_URL  
: URL to use for the reva service. Default: `localhost:9152`.

--driver | $REVA_STORAGE_ROOT_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `local`.

--mount-path | $REVA_STORAGE_ROOT_MOUNT_PATH  
: mount path. Default: `/`.

--mount-id | $REVA_STORAGE_ROOT_MOUNT_ID  
: mount id. Default: `123e4567-e89b-12d3-a456-426655440001`.

--expose-data-server | $REVA_STORAGE_ROOT_EXPOSE_DATA_SERVER  
: exposes a dedicated data server.

--data-server-url | $REVA_STORAGE_ROOT_DATA_SERVER_URL  
: data server url.

--enable-home-creation | $REVA_STORAGE_HOME_ENABLE_HOME_CREATION  
: if enabled home dirs will be automatically created.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-example.org`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-example.org`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `false`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

### ocis-reva reva-storage-public-link

Start reva storage-public-link service

Usage: `ocis-reva reva-storage-public-link [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_PUBLIC_LINK_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:10053`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_PUBLIC_LINK_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_PUBLIC_LINK_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_STORAGE_PUBLIC_LINK_ADDR  
: Address to bind reva service. Default: `localhost:10054`.

--public_share_provider_addr | $REVA_STORAGE_PUBLICLINK_PUBLIC_SHARE_PROVIDER_ADDR  
: public share provider service address. Default: `localhost:9150`.

--user_provider_addr | $REVA_STORAGE_PUBLICLINK_USER_PROVIDER_ADDR  
: user provider service address. Default: `localhost:9144`.

--storage_provider_addr | $REVA_STORAGE_PUBLICLINK_STORAGE_PROVIDER_ADDR  
: storage provider service address. Default: `localhost:9154`.

--driver | $REVA_STORAGE_PUBLIC_LINK_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `owncloud`.

--mount-id | $REVA_STORAGE_PUBLIC_LINK_MOUNT_ID  
: mount id. Default: `e1a73ede-549b-4226-abdf-40e69ca8230d`.

--expose-data-server | $REVA_STORAGE_PUBLIC_LINK_EXPOSE_DATA_SERVER  
: exposes a dedicated data server. Default: `true`.

--data-server-url | $REVA_STORAGE_PUBLIC_LINK_DATA_SERVER_URL  
: data server url. Default: `http://localhost:9156/data`.

--enable-home-creation | $REVA_STORAGE_PUBLIC_LINK_ENABLE_HOME_CREATION  
: if enabled home dirs will be automatically created. Default: `true`.

--mount-path | $REVA_STORAGE_PUBLIC_LINK_MOUNT_PATH  
: mount path. Default: `/public/`.

--gateway-url | $REVA_GATEWAY_URL  
: URL to use for the reva gateway service. Default: `localhost:9142`.

### ocis-reva health

Check health status

Usage: `ocis-reva health [command options] [arguments...]`

--debug-addr | $REVA_DEBUG_ADDR  
: Address to debug endpoint. Default: `0.0.0.0:9109`.

### ocis-reva auth-bearer

Start reva authprovider for bearer auth

Usage: `ocis-reva auth-bearer [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_AUTH_BEARER_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9149`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--oidc-issuer | $REVA_OIDC_ISSUER  
: OIDC issuer. Default: `https://localhost:9200`.

--oidc-insecure | $REVA_OIDC_INSECURE  
: OIDC allow insecure communication. Default: `true`.

--oidc-id-claim | $REVA_OIDC_ID_CLAIM  
: OIDC id claim. Default: `preferred_username`.

--network | $REVA_AUTH_BEARER_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_AUTH_BEARER_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_AUTH_BEARER_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9148`.

--url | $REVA_AUTH_BEARER_URL  
: URL to use for the reva service. Default: `localhost:9148`.

### ocis-reva storage-home-data

Start reva storage-home-data service

Usage: `ocis-reva storage-home-data [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_HOME_DATA_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9157`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_HOME_DATA_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_HOME_DATA_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `http`.

--addr | $REVA_STORAGE_HOME_DATA_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9156`.

--url | $REVA_STORAGE_HOME_DATA_URL  
: URL to use for the reva service. Default: `localhost:9156`.

--driver | $REVA_STORAGE_HOME_DATA_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `owncloud`.

--prefix | $REVA_STORAGE_HOME_DATA_PREFIX  
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $REVA_STORAGE_HOME_DATA_TEMP_FOLDER  
: temp folder. Default: `/var/tmp/`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations. Default: `/eos/dockertest/reva/users`.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder. Default: `/Shares`.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories. Default: `true`.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "E/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `true`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

--gateway-url | $REVA_GATEWAY_URL  
: URL to use for the reva gateway service. Default: `localhost:9142`.

### ocis-reva gateway

Start reva gateway

Usage: `ocis-reva gateway [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_GATEWAY_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9143`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--transfer-secret | $REVA_TRANSFER_SECRET  
: Transfer secret for datagateway. Default: `replace-me-with-a-transfer-secret`.

--network | $REVA_GATEWAY_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_GATEWAY_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_GATEWAY_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9142`.

--url | $REVA_GATEWAY_URL  
: URL to use for the reva service. Default: `localhost:9142`.

--commit-share-to-storage-grant | $REVA_GATEWAY_COMMIT_SHARE_TO_STORAGE_GRANT  
: Commit shares to the share manager. Default: `true`.

--commit-share-to-storage-ref | $REVA_GATEWAY_COMMIT_SHARE_TO_STORAGE_REF  
: Commit shares to the storage. Default: `true`.

--share-folder | $REVA_GATEWAY_SHARE_FOLDER  
: mount shares in this folder of the home storage provider. Default: `Shares`.

--link_grants_file | $REVA_GATEWAY_LINK_GRANTS_FILE  
: when using a json manager, file to use as json serialized database. Default: `/var/tmp/reva/link_grants.json`.

--disable-home-creation-on-login | $REVA_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN  
: Disable creation of home folder on login.

--storage-home-provider | $REVA_STORAGE_HOME_PROVIDER  
: mount point of the storage provider for user homes in the global namespace. Default: `/home`.

--frontend-url | $REVA_FRONTEND_URL  
: URL to use for the reva service. Default: `https://localhost:9200`.

--datagateway-url | $REVA_DATAGATEWAY_URL  
: URL to use for the reva datagateway. Default: `https://localhost:9200/data`.

--users-url | $REVA_USERS_URL  
: URL to use for the reva service. Default: `localhost:9144`.

--auth-basic-url | $REVA_AUTH_BASIC_URL  
: URL to use for the reva service. Default: `localhost:9146`.

--auth-bearer-url | $REVA_AUTH_BEARER_URL  
: URL to use for the reva service. Default: `localhost:9148`.

--sharing-url | $REVA_SHARING_URL  
: URL to use for the reva service. Default: `localhost:9150`.

--storage-root-url | $REVA_STORAGE_ROOT_URL  
: URL to use for the reva service. Default: `localhost:9152`.

--storage-root-mount-path | $REVA_STORAGE_ROOT_MOUNT_PATH  
: mount path. Default: `/`.

--storage-root-mount-id | $REVA_STORAGE_ROOT_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009152`.

--storage-home-url | $REVA_STORAGE_HOME_URL  
: URL to use for the reva service. Default: `localhost:9154`.

--storage-home-mount-path | $REVA_STORAGE_HOME_MOUNT_PATH  
: mount path. Default: `/home`.

--storage-eos-url | $REVA_STORAGE_EOS_URL  
: URL to use for the reva service. Default: `localhost:9158`.

--storage-eos-mount-path | $REVA_STORAGE_EOS_MOUNT_PATH  
: mount path. Default: `/eos`.

--storage-eos-mount-id | $REVA_STORAGE_EOS_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009158`.

--storage-oc-url | $REVA_STORAGE_OC_URL  
: URL to use for the reva service. Default: `localhost:9162`.

--storage-oc-mount-path | $REVA_STORAGE_OC_MOUNT_PATH  
: mount path. Default: `/oc`.

--storage-oc-mount-id | $REVA_STORAGE_OC_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--public-links-url | $REVA_STORAGE_PUBLIC_LINK_URL  
: URL to use for the public links service. Default: `localhost:10054`.

### ocis-reva storage-home

Start reva storage-home service

Usage: `ocis-reva storage-home [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_HOME_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9155`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_HOME_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_HOME_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_STORAGE_HOME_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9154`.

--url | $REVA_STORAGE_HOME_URL  
: URL to use for the reva service. Default: `localhost:9154`.

--driver | $REVA_STORAGE_HOME_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `owncloud`.

--mount-path | $REVA_STORAGE_HOME_MOUNT_PATH  
: mount path. Default: `/home`.

--mount-id | $REVA_STORAGE_HOME_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--expose-data-server | $REVA_STORAGE_HOME_EXPOSE_DATA_SERVER  
: exposes a dedicated data server. Default: `false`.

--data-server-url | $REVA_STORAGE_HOME_DATA_SERVER_URL  
: data server url. Default: `http://localhost:9156/data`.

--enable-home-creation | $REVA_STORAGE_HOME_ENABLE_HOME_CREATION  
: if enabled home dirs will be automatically created. Default: `true`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations. Default: `/eos/dockertest/reva/users`.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder. Default: `/Shares`.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories. Default: `true`.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "E/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `true`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

### ocis-reva storage-eos-data

Start reva storage-eos-data service

Usage: `ocis-reva storage-eos-data [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_OC_DATA_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9161`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_EOS_DATA_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_EOS_DATA_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `http`.

--addr | $REVA_STORAGE_EOS_DATA_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9160`.

--url | $REVA_STORAGE_EOS_DATA_URL  
: URL to use for the reva service. Default: `localhost:9160`.

--driver | $REVA_STORAGE_EOS_DATA_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `eos`.

--prefix | $REVA_STORAGE_EOS_DATA_PREFIX  
: prefix for the http endpoint, without leading slash. Default: `data`.

--temp-folder | $REVA_STORAGE_EOS_DATA_TEMP_FOLDER  
: temp folder. Default: `/var/tmp/`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations. Default: `/eos/dockertest/reva`.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder. Default: `/Shares`.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `false`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

--gateway-url | $REVA_GATEWAY_URL  
: URL to use for the reva gateway service. Default: `localhost:9142`.

### ocis-reva sharing

Start reva sharing service

Usage: `ocis-reva sharing [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_SHARING_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9151`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_SHARING_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_SHARING_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_SHARING_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9150`.

--url | $REVA_SHARING_URL  
: URL to use for the reva service. Default: `localhost:9150`.

--user-driver | $REVA_SHARING_USER_DRIVER  
: driver to use for the UserShareProvider. Default: `json`.

--user-json-file | $REVA_SHARING_USER_JSON_FILE  
: file used to persist shares for the UserShareProvider. Default: `/var/tmp/reva/shares.json`.

--public-driver | $REVA_SHARING_PUBLIC_DRIVER  
: driver to use for the PublicShareProvider. Default: `json`.

### ocis-reva storage-eos

Start reva storage-eos service

Usage: `ocis-reva storage-eos [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_EOS_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9159`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_EOS_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_EOS_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_STORAGE_EOS_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9158`.

--url | $REVA_STORAGE_EOS_URL  
: URL to use for the reva service. Default: `localhost:9158`.

--driver | $REVA_STORAGE_EOS_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `eos`.

--mount-path | $REVA_STORAGE_EOS_MOUNT_PATH  
: mount path. Default: `/eos`.

--mount-id | $REVA_STORAGE_EOS_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009158`.

--expose-data-server | $REVA_STORAGE_EOS_EXPOSE_DATA_SERVER  
: exposes a dedicated data server. Default: `false`.

--data-server-url | $REVA_STORAGE_EOS_DATA_SERVER_URL  
: data server url. Default: `http://localhost:9160/data`.

--enable-home-creation | $REVA_STORAGE_EOS_ENABLE_HOME_CREATION  
: if enabled home dirs will be automatically created. Default: `false`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations. Default: `/eos/dockertest/reva`.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder. Default: `/Shares`.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-mgm1.eoscluster.cern.ch:1094`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `false`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

### ocis-reva storage-oc

Start reva storage-oc service

Usage: `ocis-reva storage-oc [command options] [arguments...]`

--tracing-enabled | $REVA_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $REVA_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $REVA_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $REVA_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $REVA_TRACING_SERVICE  
: Service name for tracing. Default: `reva`.

--debug-addr | $REVA_STORAGE_OC_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9163`.

--debug-token | $REVA_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $REVA_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $REVA_DEBUG_ZPAGES  
: Enable zpages debugging.

--jwt-secret | $REVA_JWT_SECRET  
: Shared jwt secret for reva service communication. Default: `Pive-Fumkiu4`.

--network | $REVA_STORAGE_OC_NETWORK  
: Network to use for the reva service, can be 'tcp', 'udp' or 'unix'. Default: `tcp`.

--protocol | $REVA_STORAGE_OC_PROTOCOL  
: protocol for reva service, can be 'http' or 'grpc'. Default: `grpc`.

--addr | $REVA_STORAGE_OC_ADDR  
: Address to bind reva service. Default: `0.0.0.0:9162`.

--url | $REVA_STORAGE_OC_URL  
: URL to use for the reva service. Default: `localhost:9162`.

--driver | $REVA_STORAGE_OC_DRIVER  
: storage driver, eg. local, eos, owncloud or s3. Default: `owncloud`.

--mount-path | $REVA_STORAGE_OC_MOUNT_PATH  
: mount path. Default: `/oc`.

--mount-id | $REVA_STORAGE_OC_MOUNT_ID  
: mount id. Default: `1284d238-aa92-42ce-bdc4-0b0000009162`.

--expose-data-server | $REVA_STORAGE_OC_EXPOSE_DATA_SERVER  
: exposes a dedicated data server. Default: `false`.

--data-server-url | $REVA_STORAGE_OC_DATA_SERVER_URL  
: data server url. Default: `http://localhost:9164/data`.

--enable-home-creation | $REVA_STORAGE_OC_ENABLE_HOME_CREATION  
: if enabled home dirs will be automatically created. Default: `false`.

--storage-eos-namespace | $REVA_STORAGE_EOS_NAMESPACE  
: Namespace for metadata operations.

--storage-eos-shadow-namespace | $REVA_STORAGE_EOS_SHADOW_NAMESPACE  
: Shadow namespace where share references are stored.

--storage-eos-share-folder | $REVA_STORAGE_EOS_SHARE_FOLDER  
: name of the share folder.

--storage-eos-binary | $REVA_STORAGE_EOS_BINARY  
: Location of the eos binary. Default: `/usr/bin/eos`.

--storage-eos-xrdcopy-binary | $REVA_STORAGE_EOS_XRDCOPY_BINARY  
: Location of the xrdcopy binary. Default: `/usr/bin/xrdcopy`.

--storage-eos-master-url | $REVA_STORAGE_EOS_MASTER_URL  
: URL of the Master EOS MGM. Default: `root://eos-example.org`.

--storage-eos-slave-url | $REVA_STORAGE_EOS_SLAVE_URL  
: URL of the Slave EOS MGM. Default: `root://eos-example.org`.

--storage-eos-cache-directory | $REVA_STORAGE_EOS_CACHE_DIRECTORY  
: Location on the local fs where to store reads. Default: `os.TempDir()`.

--storage-eos-enable-logging | $REVA_STORAGE_EOS_ENABLE_LOGGING  
: Enables logging of the commands executed.

--storage-eos-show-hidden-sysfiles | $REVA_STORAGE_EOS_SHOW_HIDDEN_SYSFILES  
: show internal EOS files like .sys.v# and .sys.a# files..

--storage-eos-force-singleuser-mode | $REVA_STORAGE_EOS_FORCE_SINGLEUSER_MODE  
: force connections to EOS to use SingleUsername.

--storage-eos-use-keytab | $REVA_STORAGE_EOS_USE_KEYTAB  
: authenticate requests by using an EOS keytab.

--storage-eos-enable-home | $REVA_STORAGE_EOS_ENABLE_HOME  
: enable the creation of home directories.

--storage-eos-sec-protocol | $REVA_STORAGE_EOS_SEC_PROTOCOL  
: the xrootd security protocol to use between the server and EOS.

--storage-eos-keytab | $REVA_STORAGE_EOS_KEYTAB  
: the location of the keytab to use to authenticate to EOS.

--storage-eos-single-username | $REVA_STORAGE_EOS_SINGLE_USERNAME  
: the username to use when SingleUserMode is enabled.

--storage-eos-layout | $REVA_STORAGE_EOS_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{substr 0 1 .Username}}/{{.Username}}`.

--storage-local-root | $REVA_STORAGE_LOCAL_ROOT  
: the path to the local storage root. Default: `/var/tmp/reva/root`.

--storage-owncloud-datadir | $REVA_STORAGE_OWNCLOUD_DATADIR  
: the path to the owncloud data directory. Default: `/var/tmp/reva/data`.

--storage-owncloud-scan | $REVA_STORAGE_OWNCLOUD_SCAN  
: scan files on startup to add fileids. Default: `true`.

--storage-owncloud-redis | $REVA_STORAGE_OWNCLOUD_REDIS_ADDR  
: the address of the redis server. Default: `:6379`.

--storage-owncloud-enable-home | $REVA_STORAGE_OWNCLOUD_ENABLE_HOME  
: enable the creation of home storages. Default: `false`.

--storage-owncloud-layout | $REVA_STORAGE_OWNCLOUD_LAYOUT  
: `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `. Default: `{{.Username}}`.

