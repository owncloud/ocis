---
title: "Configuration"
date: "2021-04-30T11:06:38+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/thumbnails/templates
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/thumbnails/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### thumbnails health

Check health status

Usage: `thumbnails health [command options] [arguments...]`

-debug-addr |  $THUMBNAILS_DEBUG_ADDR
: Address to debug endpoint. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9189")`.

### thumbnails ocis-thumbnails

Example usage

Usage: `thumbnails ocis-thumbnails [command options] [arguments...]`

### thumbnails server

Start integrated server

Usage: `thumbnails server [command options] [arguments...]`

-log-file |  $THUMBNAILS_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.

-log-level |  $THUMBNAILS_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.

-log-pretty |  $THUMBNAILS_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.

-log-color |  $THUMBNAILS_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.

-config-file |  $THUMBNAILS_CONFIG_FILE
: Path to config file.

-tracing-enabled |  $THUMBNAILS_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $THUMBNAILS_TRACING_TYPE
: Tracing backend type. Default: `flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger")`.

-tracing-endpoint |  $THUMBNAILS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `flags.OverrideDefaultString(cfg.Tracing.Endpoint, "")`.

-tracing-collector |  $THUMBNAILS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `flags.OverrideDefaultString(cfg.Tracing.Collector, "")`.

-tracing-service |  $THUMBNAILS_TRACING_SERVICE
: Service name for tracing. Default: `flags.OverrideDefaultString(cfg.Tracing.Service, "thumbnails")`.

-debug-addr |  $THUMBNAILS_DEBUG_ADDR
: Address to bind debug server. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9189")`.

-debug-token |  $THUMBNAILS_DEBUG_TOKEN
: Token to grant metrics access. Default: `flags.OverrideDefaultString(cfg.Debug.Token, "")`.

-debug-pprof |  $THUMBNAILS_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $THUMBNAILS_DEBUG_ZPAGES
: Enable zpages debugging.

-grpc-name |  $THUMBNAILS_GRPC_NAME
: Name of the service. Default: `flags.OverrideDefaultString(cfg.Server.Name, "thumbnails")`.

-grpc-addr |  $THUMBNAILS_GRPC_ADDR
: Address to bind grpc server. Default: `flags.OverrideDefaultString(cfg.Server.Address, "0.0.0.0:9185")`.

-grpc-namespace |  $THUMBNAILS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `flags.OverrideDefaultString(cfg.Server.Namespace, "com.owncloud.api")`.

-filesystemstorage-root |  $THUMBNAILS_FILESYSTEMSTORAGE_ROOT
: Root path of the filesystem storage directory. Default: `/var/tmp/ocis/thumbnails`.

-reva-gateway-addr |  $THUMBNAILS_REVA_GATEWAY , $PROXY_REVA_GATEWAY_ADDR
: Reva gateway address. Default: `flags.OverrideDefaultString(cfg.Thumbnail.RevaGateway, "127.0.0.1:9142")`.

-webdavsource-insecure |  $THUMBNAILS_WEBDAVSOURCE_INSECURE
: Whether to skip certificate checks. Default: `flags.OverrideDefaultBool(cfg.Thumbnail.WebdavAllowInsecure, true)`.

### thumbnails version

Print the versions of the running instances

Usage: `thumbnails version [command options] [arguments...]`

-grpc-name |  $THUMBNAILS_GRPC_NAME
: Name of the service. Default: `flags.OverrideDefaultString(cfg.Server.Name, "thumbnails")`.

-grpc-namespace |  $THUMBNAILS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `flags.OverrideDefaultString(cfg.Server.Namespace, "com.owncloud.api")`.

