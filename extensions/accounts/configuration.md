---
title: "Configuration"
date: "2020-12-15T13:17:28+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/accounts/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-accounts reads `accounts.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### accounts remove

Removes an existing account

Usage: `accounts remove [command options] [arguments...]`

--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

### accounts add

Create a new account

Usage: `accounts add [command options] [arguments...]`

### accounts update

Make changes to an existing account

Usage: `accounts update [command options] [arguments...]`

### accounts inspect

Show detailed data on an existing account

Usage: `accounts inspect [command options] [arguments...]`

--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

### accounts version

Print the versions of the running instances

Usage: `accounts version [command options] [arguments...]`

--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

### accounts server

Start ocis accounts service

Usage: `accounts server [command options] [arguments...]`

--tracing-enabled | $ACCOUNTS_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $ACCOUNTS_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $ACCOUNTS_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $ACCOUNTS_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $ACCOUNTS_TRACING_SERVICE  
: Service name for tracing. Default: `accounts`.

--http-namespace | $ACCOUNTS_HTTP_NAMESPACE  
: Set the base namespace for the http namespace. Default: `com.owncloud.web`.

--http-addr | $ACCOUNTS_HTTP_ADDR  
: Address to bind http server. Default: `0.0.0.0:9181`.

--http-root | $ACCOUNTS_HTTP_ROOT  
: Root path of http server. Default: `/`.

--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

--grpc-addr | $ACCOUNTS_GRPC_ADDR  
: Address to bind grpc server. Default: `0.0.0.0:9180`.

--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

--asset-path | $ACCOUNTS_ASSET_PATH  
: Path to custom assets.

--jwt-secret | $ACCOUNTS_JWT_SECRET  
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

--storage-disk-path | $ACCOUNTS_STORAGE_DISK_PATH  
: Path on the local disk, e.g. /var/tmp/ocis/accounts.

--storage-cs3-provider-addr | $ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR  
: bind address for the metadata storage provider. Default: `localhost:9215`.

--storage-cs3-data-url | $ACCOUNTS_STORAGE_CS3_DATA_URL  
: http endpoint of the metadata storage. Default: `http://localhost:9216`.

--storage-cs3-data-prefix | $ACCOUNTS_STORAGE_CS3_DATA_PREFIX  
: path prefix for the http endpoint of the metadata storage, without leading slash. Default: `data`.

--storage-cs3-jwt-secret | $ACCOUNTS_STORAGE_CS3_JWT_SECRET  
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

--service-user-uuid | $ACCOUNTS_SERVICE_USER_UUID  
: uuid of the internal service user (required on EOS). Default: `95cb8724-03b2-11eb-a0a6-c33ef8ef53ad`.

--service-user-username | $ACCOUNTS_SERVICE_USER_USERNAME  
: username of the internal service user (required on EOS).

### accounts list

List existing accounts

Usage: `accounts list [command options] [arguments...]`

--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

### accounts rebuildIndex

Rebuilds the service's index, i.e. deleting and then re-adding all existing documents

Usage: `accounts rebuildIndex [command options] [arguments...]`

### accounts ocis-accounts

Provide accounts and groups for oCIS

Usage: `accounts ocis-accounts [command options] [arguments...]`

--log-level | $ACCOUNTS_LOG_LEVEL  
: Set logging level. Default: `info`.

--log-pretty | $ACCOUNTS_LOG_PRETTY  
: Enable pretty logging. Default: `true`.

--log-color | $ACCOUNTS_LOG_COLOR  
: Enable colored logging. Default: `true`.

