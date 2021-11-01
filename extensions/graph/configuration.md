---
title: "Configuration"
date: "2021-11-01T12:24:18+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/graph/templates
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/graph/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### graph health

Check health status

Usage: `graph health [command options] [arguments...]`






-debug-addr |  $GRAPH_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9124"`.




















### graph ocis-graph

Serve Graph API for oCIS

Usage: `graph ocis-graph [command options] [arguments...]`


-config-file |  $GRAPH_CONFIG_FILE
: Path to config file. Default: `""`.


-log-level |  $GRAPH_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $GRAPH_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $GRAPH_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





















### graph server

Start integrated server

Usage: `graph server [command options] [arguments...]`







-log-file |  $GRAPH_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-tracing-enabled |  $GRAPH_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $GRAPH_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $GRAPH_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $GRAPH_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $GRAPH_TRACING_SERVICE
: Service name for tracing. Default: `"graph"`.


-debug-addr |  $GRAPH_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9124"`.


-debug-token |  $GRAPH_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $GRAPH_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $GRAPH_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $GRAPH_HTTP_ADDR
: Address to bind http server. Default: `"127.0.0.1:9120"`.


-http-root |  $GRAPH_HTTP_ROOT
: Root path of http server. Default: `"/graph"`.


-http-namespace |  $GRAPH_HTTP_NAMESPACE
: Set the base namespace for the http service for service discovery. Default: `"com.owncloud.web"`.


-spaces-webdav-base |  $GRAPH_SPACES_WEBDAV_BASE , $OCIS_URL
: spaces webdav base URL to use when rendering drive WabDAV URLs. Default: `"https://localhost:9200"`.


-spaces-webdav-path |  $GRAPH_SPACES_WEBDAV_PATH
: spaces webdav path to use when rendering drive WabDAV URLs. Default: `"/dav/spaces/"`.


-default-space-quota |  $GRAPH_SPACES_DEFAULT_QUOTA
: default quota used for all spaces if no custom quota was given. Default: `"1000000000"`.


-jwt-secret |  $GRAPH_JWT_SECRET , $OCIS_JWT_SECRET
: Used to validate the reva access JWT, should equal reva's jwt-secret. Default: `"Pive-Fumkiu4"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-extensions | 
: Run specific extensions during supervised mode. This flag is set by the runtime.

