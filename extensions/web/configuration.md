---
title: "Configuration"
date: "2021-11-11T07:43:20+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/web/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-web reads `web.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/web/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Serve ownCloud Web for oCIS

Usage: `web [global options] command [command options] [arguments...]`

































## Sub Commands

### web server

Start integrated server

Usage: `web server [command options] [arguments...]`


-log-level |  $WEB_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $WEB_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $WEB_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.



-log-file |  $WEB_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-config-file |  $WEB_CONFIG_FILE
: Path to config file.


-tracing-enabled |  $WEB_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $WEB_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $WEB_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $WEB_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $WEB_TRACING_SERVICE
: Service name for tracing. Default: `"web"`.


-debug-addr |  $WEB_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9104"`.


-debug-token |  $WEB_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $WEB_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $WEB_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $WEB_HTTP_ADDR
: Address to bind http server. Default: `"127.0.0.1:9100"`.


-http-root |  $WEB_HTTP_ROOT
: Root path of http server. Default: `"/"`.


-http-namespace |  $WEB_NAMESPACE
: Set the base namespace for the http namespace. Default: `"com.owncloud.web"`.


-http-cache-ttl |  $WEB_CACHE_TTL
: Set the static assets caching duration in seconds. Default: `604800`.


-asset-path |  $WEB_ASSET_PATH
: Path to custom assets. Default: `""`.


-web-config |  $WEB_UI_CONFIG
: Path to web config. Default: `""`.


-web-config-server |  $WEB_UI_CONFIG_SERVER , $OCIS_URL
: Configuration server URL. Default: `"https://localhost:9200"`.


-web-theme-server |  $WEB_UI_THEME_SERVER , $OCIS_URL
: Theme server URL. Default: `"https://localhost:9200"`.


-web-config-theme |  $WEB_UI_THEME_PATH
: Theme path on the theme server. Default: `"/themes/owncloud/theme.json"`.


-web-config-version |  $WEB_UI_CONFIG_VERSION
: Version. Default: `"0.1.0"`.


-web-config-app |  $WEB_UI_CONFIG_APPS
: `--web-config-app files [--web-config-app draw-io]`. Default: `cli.NewStringSlice("files", "search", "media-viewer", "external")`.


-oidc-metadata-url |  $WEB_OIDC_METADATA_URL
: OpenID Connect metadata URL, defaults to <WEB_OIDC_AUTHORITY>/.well-known/openid-configuration. Default: `""`.


-oidc-authority |  $WEB_OIDC_AUTHORITY , $OCIS_URL
: OpenID Connect authority. Default: `"https://localhost:9200"`.


-oidc-client-id |  $WEB_OIDC_CLIENT_ID
: OpenID Connect client ID. Default: `"web"`.


-oidc-response-type |  $WEB_OIDC_RESPONSE_TYPE
: OpenID Connect response type. Default: `"code"`.


-oidc-scope |  $WEB_OIDC_SCOPE
: OpenID Connect scope. Default: `"openid profile email"`.


-extensions | 
: Run specific extensions during supervised mode. This flag is set by the runtime.

### web health

Check health status

Usage: `web health [command options] [arguments...]`





-debug-addr |  $WEB_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9104"`.





























