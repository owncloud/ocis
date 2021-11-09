---
title: "Configuration"
date: "2021-11-09T08:53:31+0000"
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

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### thumbnails health

Check health status

Usage: `thumbnails health [command options] [arguments...]`


-debug-addr |  $THUMBNAILS_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9189"`.


























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


-tracing-enabled |  $THUMBNAILS_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $THUMBNAILS_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $THUMBNAILS_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $THUMBNAILS_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $THUMBNAILS_TRACING_SERVICE
: Service name for tracing. Default: `"thumbnails"`.


-debug-addr |  $THUMBNAILS_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9189"`.


-debug-token |  $THUMBNAILS_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $THUMBNAILS_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $THUMBNAILS_DEBUG_ZPAGES
: Enable zpages debugging.


-grpc-name |  $THUMBNAILS_GRPC_NAME
: Name of the service. Default: `"thumbnails"`.


-grpc-addr |  $THUMBNAILS_GRPC_ADDR
: Address to bind grpc server. Default: `"127.0.0.1:9185"`.


-grpc-namespace |  $THUMBNAILS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `"com.owncloud.api"`.


-filesystemstorage-root |  $THUMBNAILS_FILESYSTEMSTORAGE_ROOT
: Root path of the filesystem storage directory. Default: `flags.OverrideDefaultString(cfg.Thumbnail.FileSystemStorage.RootDirectory, path.Join(defaults.BaseDataPath(), "thumbnails"))`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-webdavsource-insecure |  $THUMBNAILS_WEBDAVSOURCE_INSECURE
: Whether to skip certificate checks. Default: `true`.


-thumbnail-resolution |  $THUMBNAILS_RESOLUTIONS
: --thumbnail-resolution 16x16 [--thumbnail-resolution 32x32]. Default: `cli.NewStringSlice("16x16", "32x32", "64x64", "128x128", "1920x1080", "3840x2160", "7680x4320")`.


-webdav-namespace |  $STORAGE_WEBDAV_NAMESPACE
: Namespace prefix for the webdav endpoint. Default: `"/home"`.


-extensions | 
: Run specific extensions during supervised mode.



### thumbnails version

Print the versions of the running instances

Usage: `thumbnails version [command options] [arguments...]`


























-grpc-name |  $THUMBNAILS_GRPC_NAME
: Name of the service. Default: `"thumbnails"`.


-grpc-namespace |  $THUMBNAILS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace. Default: `"com.owncloud.api"`.

