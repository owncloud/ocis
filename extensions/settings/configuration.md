---
title: "Configuration"
date: "2021-02-26T04:43:22+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/settings/templates
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### settings version

Print the versions of the running instances

Usage: `settings version [command options] [arguments...]`

-grpc-namespace |  $SETTINGS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

-name |  $SETTINGS_NAME
: service name. Default: `settings`.

### settings health

Check health status

Usage: `settings health [command options] [arguments...]`

-debug-addr |  $SETTINGS_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9194`.

### settings ocis-settings

Provide settings and permissions for oCIS

Usage: `settings ocis-settings [command options] [arguments...]`

-log-level |  $SETTINGS_LOG_LEVEL
: Set logging level. Default: `info`.

-log-pretty |  $SETTINGS_LOG_PRETTY
: Enable pretty logging. Default: `true`.

-log-color |  $SETTINGS_LOG_COLOR
: Enable colored logging. Default: `true`.

### settings server

Start integrated server

Usage: `settings server [command options] [arguments...]`

-config-file |  $SETTINGS_CONFIG_FILE
: Path to config file.

-tracing-enabled |  $SETTINGS_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $SETTINGS_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.

-tracing-endpoint |  $SETTINGS_TRACING_ENDPOINT
: Endpoint for the agent.

-tracing-collector |  $SETTINGS_TRACING_COLLECTOR
: Endpoint for the collector.

-tracing-service |  $SETTINGS_TRACING_SERVICE
: Service name for tracing. Default: `settings`.

-debug-addr |  $SETTINGS_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9194`.

-debug-token |  $SETTINGS_DEBUG_TOKEN
: Token to grant metrics access.

-debug-pprof |  $SETTINGS_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $SETTINGS_DEBUG_ZPAGES
: Enable zpages debugging.

-http-addr |  $SETTINGS_HTTP_ADDR
: Address to bind http server. Default: `0.0.0.0:9190`.

-http-namespace |  $SETTINGS_HTTP_NAMESPACE
: Set the base namespace for the http namespace. Default: `com.owncloud.web`.

-http-root |  $SETTINGS_HTTP_ROOT
: Root path of http server. Default: `/`.

-grpc-addr |  $SETTINGS_GRPC_ADDR
: Address to bind grpc server. Default: `0.0.0.0:9191`.

-asset-path |  $SETTINGS_ASSET_PATH
: Path to custom assets.

-grpc-namespace |  $SETTINGS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

-name |  $SETTINGS_NAME
: service name. Default: `settings`.

-data-path |  $SETTINGS_DATA_PATH
: Mount path for the storage. Default: `/var/tmp/ocis/settings`.

-jwt-secret |  $SETTINGS_JWT_SECRET , $OCIS_JWT_SECRET
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

