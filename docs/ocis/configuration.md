---
title: "Configuration"
date: "2020-10-05T21:19:39+0200"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
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

ownCloud Infinite Scale Stack

Usage: `ocis [global options] command [command options] [arguments...]`

--config-file | $OCIS_CONFIG_FILE
: Path to config file.

--log-level | $OCIS_LOG_LEVEL
: Set logging level. Default: `info`.

--log-pretty | $OCIS_LOG_PRETTY
: Enable pretty logging. Default: `true`.

--log-color | $OCIS_LOG_COLOR
: Enable colored logging. Default: `true`.

## Sub Commands

### ocis kill

Kill an extension by name

Usage: `ocis kill [command options] [arguments...]`

### ocis server

Start fullstack server

Usage: `ocis server [command options] [arguments...]`

--tracing-enabled | $OCIS_TRACING_ENABLED
: Enable sending traces.

--tracing-type | $OCIS_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `localhost:6831`.

--tracing-collector | $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `http://localhost:14268/api/traces`.

--tracing-service | $OCIS_TRACING_SERVICE
: Service name for tracing. Default: `ocis`.

--debug-addr | $OCIS_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9010`.

--debug-token | $OCIS_DEBUG_TOKEN
: Token to grant metrics access.

--debug-pprof | $OCIS_DEBUG_PPROF
: Enable pprof debugging.

--debug-zpages | $OCIS_DEBUG_ZPAGES
: Enable zpages debugging.

--http-addr | $OCIS_HTTP_ADDR
: Address to bind http server. Default: `0.0.0.0:9000`.

--http-root | $OCIS_HTTP_ROOT
: Root path of http server. Default: `/`.

--grpc-addr | $OCIS_GRPC_ADDR
: Address to bind grpc server. Default: `0.0.0.0:9001`.

### ocis run

Runs an extension

Usage: `ocis run [command options] [arguments...]`

### ocis health

Check health status

Usage: `ocis health [command options] [arguments...]`

--debug-addr | $OCIS_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9010`.

### ocis list

Lists running ocis extensions

Usage: `ocis list [command options] [arguments...]`

### List of available Extension subcommands

There are more subcommands to start the individual extensions. Please check the documentation about their usage and options in the dedicated section of the documentation.

#### ocis konnectd

Start konnectd server

#### ocis storage-frontend

Start storage frontend

#### ocis accounts

Start accounts server

#### ocis storage-storage-root

Start storage root storage

#### ocis storage-gateway

Start storage gateway

#### ocis storage-storage-home

Start storage storage service for home mount

#### ocis storage-storage-public-link

Start storage public link storage

#### ocis store

Start a go-micro store

#### ocis glauth

Start glauth server

#### ocis storage-auth-bearer

Start storage auth-bearer service

#### ocis storage-sharing

Start storage sharing service

#### ocis webdav

Start webdav server

#### ocis storage-storage-oc-data

Start storage storage data provider for oc mount

#### ocis thumbnails

Start thumbnails server

#### ocis proxy

Start proxy server

#### ocis settings

Start settings server

#### ocis storage-auth-basic

Start storage auth-basic service

#### ocis storage-storage-metadata

Start storage storage service for metadata mount

#### ocis ocs

Start ocs server

#### ocis storage-storage-eos

Start storage storage service for eos mount

#### ocis storage-storage-eos-data

Start storage storage data provider for eos mount

#### ocis storage-storage-home-data

Start storage storage data provider for home mount

#### ocis storage-storage-oc

Start storage storage service for oc mount

#### ocis storage-users

Start storage users service

#### ocis phoenix

Start phoenix server

