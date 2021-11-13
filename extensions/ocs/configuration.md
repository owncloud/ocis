---
title: "Configuration"
date: "2021-11-13T00:09:27+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/ocs/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-ocs reads `ocs.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/ocs/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### ocs health

Check health status

Usage: `ocs health [command options] [arguments...]`


-debug-addr |  $OCS_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9114"`.
































### ocs ocis-ocs

Serve OCS API for oCIS

Usage: `ocs ocis-ocs [command options] [arguments...]`

































### ocs server

Start integrated server

Usage: `ocs server [command options] [arguments...]`



-log-file |  $OCS_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-log-level |  $OCS_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $OCS_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $OCS_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.


-config-file |  $OCS_CONFIG_FILE
: Path to config file.


-tracing-enabled |  $OCS_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces. Default: `false`.


-tracing-type |  $OCS_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $OCS_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $OCS_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $OCS_TRACING_SERVICE
: Service name for tracing. Default: `"ocs"`.


-debug-addr |  $OCS_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9114"`.


-debug-token |  $OCS_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $OCS_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $OCS_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $OCS_HTTP_ADDR
: Address to bind http server. Default: `"127.0.0.1:9110"`.


-http-namespace |  $OCS_NAMESPACE
: Set the base namespace for the http namespace. Default: `"com.owncloud.web"`.


-cors-allowed-origins |  $OCS_CORS_ALLOW_ORIGINS , $OCIS_CORS_ALLOW_ORIGINS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("*")`.


-cors-allowed-methods |  $OCS_CORS_ALLOW_METHODS , $OCIS_CORS_ALLOW_METHODS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")`.


-cors-allowed-headers |  $OCS_CORS_ALLOW_HEADERS , $OCIS_CORS_ALLOW_HEADERS
: Set the allowed CORS origins. Default: `cli.NewStringSlice("Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With")`.


-cors-allow-credentials |  $OCS_CORS_ALLOW_CREDENTIALS , $OCIS_CORS_ALLOW_CREDENTIALS
: Allow credentials for CORS. Default: `true`.


-name |  $OCS_NAME
: Service name. Default: `"ocs"`.


-http-root |  $OCS_HTTP_ROOT
: Root path of http server. Default: `"/ocs"`.


-jwt-secret |  $OCS_JWT_SECRET , $OCIS_JWT_SECRET
: Used to dismantle the access token, should equal reva's jwt-secret. Default: `"Pive-Fumkiu4"`.


-account-backend-type |  $OCS_ACCOUNT_BACKEND_TYPE
: account-backend-type. Default: `"accounts"`.


-reva-gateway-addr |  $REVA_GATEWAY
: Address of REVA gateway endpoint. Default: `"127.0.0.1:9142"`.


-machine-auth-api-key |  $OCS_MACHINE_AUTH_API_KEY , $OCIS_MACHINE_AUTH_API_KEY
: the API key to be used for the machine auth driver in reva. Default: `"change-me-please"`.


-idm-address |  $OCS_IDM_ADDRESS , $OCIS_URL
: keeps track of the IDM Address. Needed because of Reva requisite of uniqueness for users. Default: `"https://localhost:9200"`.


-users-driver |  $OCS_STORAGE_USERS_DRIVER , $STORAGE_USERS_DRIVER
: storage driver for users mount: eg. local, eos, owncloud, ocis or s3. Default: `"ocis"`.


-extensions | 
: Run specific extensions during supervised mode.



### ocs version

Print the versions of the running instances

Usage: `ocs version [command options] [arguments...]`
































-http-namespace |  $OCS_NAMESPACE
: Set the base namespace for the http namespace. Default: `"com.owncloud.web"`.


-name |  $OCS_NAME
: Service name. Default: `"ocs"`.

