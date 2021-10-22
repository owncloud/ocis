---
title: "Configuration"
date: "2021-10-22T14:37:44+0000"
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/webdav/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Serve WebDAV API for oCIS

Usage: `webdav [global options] command [command options] [arguments...]`





























## Sub Commands

### webdav health

Check health status

Usage: `webdav health [command options] [arguments...]`


-debug-addr |  $WEBDAV_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9119"`.




























### webdav server

Start integrated server

Usage: `webdav server [command options] [arguments...]`



-log-file |  $WEBDAV_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-log-level |  $WEBDAV_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $WEBDAV_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $WEBDAV_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.


-config-file |  $WEBDAV_CONFIG_FILE
: Path to config file.


-tracing-enabled |  $WEBDAV_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $WEBDAV_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $WEBDAV_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $WEBDAV_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $WEBDAV_TRACING_SERVICE
: Service name for tracing. Default: `"webdav"`.


-debug-addr |  $WEBDAV_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9119"`.


-debug-token |  $WEBDAV_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $WEBDAV_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $WEBDAV_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $WEBDAV_HTTP_ADDR
: Address to bind http server. Default: `"127.0.0.1:9115"`.


-http-namespace |  $WEBDAV_HTTP_NAMESPACE
: Set the base namespace for service discovery. Default: `"com.owncloud.web"`.


-cors-allowed-origins |  $WEBDAV_CORS_ALLOW_ORIGINS , $OCIS_CORS_ALLOW_ORIGINS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("*")`.


-cors-allowed-methods |  $WEBDAV_CORS_ALLOW_METHODS , $OCIS_CORS_ALLOW_METHODS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")`.


-cors-allowed-headers |  $WEBDAV_CORS_ALLOW_HEADERS , $OCIS_CORS_ALLOW_HEADERS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With")`.


-cors-allow-credentials |  $WEBDAV_CORS_ALLOW_CREDENTIALS , $OCIS_CORS_ALLOW_CREDENTIALS
: Allow credentials for CORS. Default: `true`.


-service-name |  $WEBDAV_SERVICE_NAME
: Service name. Default: `"webdav"`.


-http-root |  $WEBDAV_HTTP_ROOT
: Root path of http server. Default: `"/"`.


-ocis-public-url |  $OCIS_PUBLIC_URL , $OCIS_URL
: The domain under which oCIS is reachable. Default: `"https://127.0.0.1:9200"`.


-webdav-namespace |  $STORAGE_WEBDAV_NAMESPACE
: Namespace prefix for the /webdav endpoint. Default: `"/home"`.


-extensions | 
: Run specific extensions during supervised mode. This flag is set by the runtime.



### webdav version

Print the versions of the running instances

Usage: `webdav version [command options] [arguments...]`




























-http-namespace |  $WEBDAV_HTTP_NAMESPACE
: Set the base namespace for service discovery. Default: `"com.owncloud.web"`.


-service-name |  $WEBDAV_SERVICE_NAME
: Service name. Default: `"webdav"`.

