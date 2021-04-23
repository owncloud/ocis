---
title: "Configuration"
date: "2021-04-23T05:52:53+0000"
weight: 2
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/ocis/templates
geekdocFilePath: CONFIGURATION.tmpl
---

{{< toc >}}

## Configuration

oCIS Single Binary is not responsible for configuring extensions. Instead, each extension could either be configured by environment variables, cli flags or config files.

Each extension has its dedicated documentation page (e.g. [proxy configuration]({{< relref "../extensions/accounts/configuration.md" >}})) which lists all possible configurations. Config files and environment variables are picked up if you use the `./bin/ocis server` command within the oCIS single binary. Command line flags must be set explicitly on the extensions subcommands.

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-proxy reads `proxy.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/ocis/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

ownCloud Infinite Scale Stack

Usage: `ocis [global options] command [command options] [arguments...]`

-config-file |  $OCIS_CONFIG_FILE
: Load config file from a non standard location..

-ocis-log-level |  $OCIS_LOG_LEVEL
: Set logging level. Default: `info`.

-ocis-log-pretty |  $OCIS_LOG_PRETTY
: Enable pretty logging. Default: `false`.

-ocis-log-color |  $OCIS_LOG_COLOR
: Enable colored logging. Default: `true`.

-ocis-log-file |  $OCIS_LOG_FILE
: Enable log to file.

-tracing-enabled |  $OCIS_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $OCIS_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.

-tracing-endpoint |  $OCIS_TRACING_ENDPOINT
: Endpoint for the agent.

-tracing-collector |  $OCIS_TRACING_COLLECTOR
: Endpoint for the collector.

-tracing-service |  $OCIS_TRACING_SERVICE
: Service name for tracing. Default: `ocis`.

-jwt-secret |  $OCIS_JWT_SECRET
: Used to dismantle the access token, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

-runtime-port |  $OCIS_RUNTIME_PORT
: Configures which port the runtime starts. Default: `9250`.

-runtime-host |  $OCIS_RUNTIME_HOST
: Configures the host where the runtime process is running. Default: `localhost`.

## Sub Commands

### ocis health

Check health status

Usage: `ocis health [command options] [arguments...]`

-debug-addr |  $OCIS_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9010`.

### ocis list

Lists running ocis extensions

Usage: `ocis list [command options] [arguments...]`

### ocis kill

Kill an extension by name

Usage: `ocis kill [command options] [arguments...]`

### ocis server

Start fullstack server

Usage: `ocis server [command options] [arguments...]`

-debug-addr |  $OCIS_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9010`.

-debug-token |  $OCIS_DEBUG_TOKEN
: Token to grant metrics access.

-debug-pprof |  $OCIS_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $OCIS_DEBUG_ZPAGES
: Enable zpages debugging.

-http-addr |  $OCIS_HTTP_ADDR
: Address to bind http server. Default: `0.0.0.0:9000`.

-http-root |  $OCIS_HTTP_ROOT
: Root path of http server. Default: `/`.

-grpc-addr |  $OCIS_GRPC_ADDR
: Address to bind grpc server. Default: `0.0.0.0:9001`.

### ocis run

Runs an extension

Usage: `ocis run [command options] [arguments...]`

### List of available Extension subcommands

There are more subcommands to start the individual extensions. Please check the documentation about their usage and options in the dedicated section of the documentation.

#### ocis idp

Start idp server

#### ocis storage-auth-basic

Start storage auth-basic service

#### ocis storage-auth-bearer

Start storage auth-bearer service

#### ocis storage-public-link

Start storage public link storage

#### ocis ocs

Start ocs server

#### ocis proxy

Start proxy server

#### ocis accounts

Start accounts server

#### ocis version

Lists running services with version

#### ocis webdav

Start webdav server

#### ocis storage-users

Start storage and data provider for /users mount

#### ocis thumbnails

Start thumbnails server

#### ocis web

Start web server

#### ocis glauth

Start glauth server

#### ocis storage-metadata

Start storage and data service for metadata

#### ocis storage-userprovider

Start storage userprovider service

#### ocis store

Start a go-micro store

#### ocis storage-groupprovider

Start storage groupprovider service

#### ocis onlyoffice

Start onlyoffice server

#### ocis storage-frontend

Start storage frontend

#### ocis storage-home

Start storage and data provider for /home mount

#### ocis storage-gateway

Start storage gateway

#### ocis storage-sharing

Start storage sharing service

#### ocis settings

Start settings server

