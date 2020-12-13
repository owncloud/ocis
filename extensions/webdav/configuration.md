---
title: "Configuration"
date: "2020-12-13T11:06:20+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/webdav/templates
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

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Serve WebDAV API for oCIS

Usage: `webdav [global options] command [command options] [arguments...]`

--log-level | $WEBDAV_LOG_LEVEL  
: Set logging level. Default: `info`.

--log-pretty | $WEBDAV_LOG_PRETTY  
: Enable pretty logging. Default: `true`.

--log-color | $WEBDAV_LOG_COLOR  
: Enable colored logging. Default: `true`.

## Sub Commands

### webdav server

Start integrated server

Usage: `webdav server [command options] [arguments...]`

--config-file | $WEBDAV_CONFIG_FILE  
: Path to config file.

--tracing-enabled | $WEBDAV_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $WEBDAV_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $WEBDAV_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $WEBDAV_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $WEBDAV_TRACING_SERVICE  
: Service name for tracing. Default: `webdav`.

--debug-addr | $WEBDAV_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9119`.

--debug-token | $WEBDAV_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $WEBDAV_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $WEBDAV_DEBUG_ZPAGES  
: Enable zpages debugging.

--http-addr | $WEBDAV_HTTP_ADDR  
: Address to bind http server. Default: `0.0.0.0:9115`.

--http-namespace | $WEBDAV_HTTP_NAMESPACE  
: Set the base namespace for service discovery. Default: `com.owncloud.web`.

--service-name | $WEBDAV_SERVICE_NAME  
: Service name. Default: `webdav`.

--http-root | $WEBDAV_HTTP_ROOT  
: Root path of http server. Default: `/`.

### webdav version

Print the versions of the running instances

Usage: `webdav version [command options] [arguments...]`

--http-namespace | $WEBDAV_HTTP_NAMESPACE  
: Set the base namespace for service discovery. Default: `com.owncloud.web`.

--service-name | $WEBDAV_SERVICE_NAME  
: Service name. Default: `webdav`.

### webdav health

Check health status

Usage: `webdav health [command options] [arguments...]`

--debug-addr | $WEBDAV_DEBUG_ADDR  
: Address to debug endpoint. Default: `0.0.0.0:9119`.

