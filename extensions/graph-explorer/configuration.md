---
title: "Configuration"
date: "2021-05-03T12:28:09+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/graph-explorer/templates
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/graph-explorer/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Serve Graph-Explorer for oCIS

Usage: `graph-explorer [global options] command [command options] [arguments...]`

-log-level |  $GRAPH_EXPLORER_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.

-log-pretty |  $GRAPH_EXPLORER_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.

-log-color |  $GRAPH_EXPLORER_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.

## Sub Commands

### graph-explorer health

Check health status

Usage: `graph-explorer health [command options] [arguments...]`

-debug-addr |  $GRAPH_EXPLORER_DEBUG_ADDR
: Address to debug endpoint. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9136")`.

### graph-explorer server

Start integrated server

Usage: `graph-explorer server [command options] [arguments...]`

-log-file |  $GRAPH_EXPLORER_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.

-tracing-enabled |  $GRAPH_EXPLORER_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $GRAPH_EXPLORER_TRACING_TYPE
: Tracing backend type. Default: `flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger")`.

-tracing-endpoint |  $GRAPH_EXPLORER_TRACING_ENDPOINT
: Endpoint for the agent. Default: `flags.OverrideDefaultString(cfg.Tracing.Endpoint, "")`.

-tracing-collector |  $GRAPH_EXPLORER_TRACING_COLLECTOR
: Endpoint for the collector. Default: `flags.OverrideDefaultString(cfg.Tracing.Collector, "")`.

-tracing-service |  $GRAPH_EXPLORER_TRACING_SERVICE
: Service name for tracing. Default: `flags.OverrideDefaultString(cfg.Tracing.Service, "graph-explorer")`.

-debug-addr |  $GRAPH_EXPLORER_DEBUG_ADDR
: Address to bind debug server. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9136")`.

-debug-token |  $GRAPH_EXPLORER_DEBUG_TOKEN
: Token to grant metrics access. Default: `flags.OverrideDefaultString(cfg.Debug.Token, "")`.

-debug-pprof |  $GRAPH_EXPLORER_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $GRAPH_EXPLORER_DEBUG_ZPAGES
: Enable zpages debugging.

-http-addr |  $GRAPH_EXPLORER_HTTP_ADDR
: Address to bind http server. Default: `flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9135")`.

-http-root |  $GRAPH_EXPLORER_HTTP_ROOT
: Root path of http server. Default: `flags.OverrideDefaultString(cfg.HTTP.Root, "/graph-explorer")`.

-http-namespace |  $GRAPH_EXPLORER_NAMESPACE
: Set the base namespace for the http namespace. Default: `flags.OverrideDefaultString(cfg.HTTP.Namespace, "com.owncloud.web")`.

-issuer |  $GRAPH_EXPLORER_ISSUER , $OCIS_URL
: Set the OpenID Connect Provider. Default: `flags.OverrideDefaultString(cfg.GraphExplorer.Issuer, "https://localhost:9200")`.

-client-id |  $GRAPH_EXPLORER_CLIENT_ID
: Set the OpenID Client ID to send to the issuer. Default: `flags.OverrideDefaultString(cfg.GraphExplorer.ClientID, "ocis-explorer.js")`.

-graph-url |  $GRAPH_EXPLORER_GRAPH_URL
: Set the url to the graph api service. Default: `flags.OverrideDefaultString(cfg.GraphExplorer.GraphURL, "https://localhost:9200/graph")`.

