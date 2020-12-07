---
title: "Configuration"
date: "2020-12-07T10:59:47+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/phoenix/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-phoenix reads `phoenix.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### phoenix ocis-phoenix

Serve Phoenix for oCIS

Usage: `phoenix ocis-phoenix [command options] [arguments...]`

--log-level | $PHOENIX_LOG_LEVEL  
: Set logging level. Default: `info`.

--log-pretty | $PHOENIX_LOG_PRETTY  
: Enable pretty logging. Default: `true`.

--log-color | $PHOENIX_LOG_COLOR  
: Enable colored logging. Default: `true`.

### phoenix server

Start integrated server

Usage: `phoenix server [command options] [arguments...]`

--config-file | $PHOENIX_CONFIG_FILE  
: Path to config file.

--tracing-enabled | $PHOENIX_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $PHOENIX_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $PHOENIX_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $PHOENIX_TRACING_COLLECTOR  
: Endpoint for the collector.

--tracing-service | $PHOENIX_TRACING_SERVICE  
: Service name for tracing. Default: `phoenix`.

--debug-addr | $PHOENIX_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9104`.

--debug-token | $PHOENIX_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $PHOENIX_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $PHOENIX_DEBUG_ZPAGES  
: Enable zpages debugging.

--http-addr | $PHOENIX_HTTP_ADDR  
: Address to bind http server. Default: `0.0.0.0:9100`.

--http-root | $PHOENIX_HTTP_ROOT  
: Root path of http server. Default: `/`.

--http-namespace | $PHOENIX_NAMESPACE  
: Set the base namespace for the http namespace. Default: `com.owncloud.web`.

--asset-path | $PHOENIX_ASSET_PATH  
: Path to custom assets.

--web-config | $PHOENIX_WEB_CONFIG  
: Path to phoenix config.

--web-config-server | $PHOENIX_WEB_CONFIG_SERVER  
: Server URL. Default: `https://localhost:9200`.

--web-config-theme | $PHOENIX_WEB_CONFIG_THEME  
: Theme. Default: `owncloud`.

--web-config-version | $PHOENIX_WEB_CONFIG_VERSION  
: Version. Default: `0.1.0`.

--oidc-metadata-url | $PHOENIX_OIDC_METADATA_URL  
: OpenID Connect metadata URL. Default: `https://localhost:9200/.well-known/openid-configuration`.

--oidc-authority | $PHOENIX_OIDC_AUTHORITY  
: OpenID Connect authority. Default: `https://localhost:9200`.

--oidc-client-id | $PHOENIX_OIDC_CLIENT_ID  
: OpenID Connect client ID. Default: `phoenix`.

--oidc-response-type | $PHOENIX_OIDC_RESPONSE_TYPE  
: OpenID Connect response type. Default: `code`.

--oidc-scope | $PHOENIX_OIDC_SCOPE  
: OpenID Connect scope. Default: `openid profile email`.

### phoenix health

Check health status

Usage: `phoenix health [command options] [arguments...]`

--debug-addr | $PHOENIX_DEBUG_ADDR  
: Address to debug endpoint. Default: `0.0.0.0:9104`.

