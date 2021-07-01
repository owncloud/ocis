---
title: "Configuration"
date: "2021-07-01T14:17:59+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/proxy/templates
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

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/proxy/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### proxy health

Check health status

Usage: `proxy health [command options] [arguments...]`






-debug-addr |  $PROXY_DEBUG_ADDR
: Address to debug endpoint. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9109")`.































### proxy ocis-proxy

proxy for oCIS

Usage: `proxy ocis-proxy [command options] [arguments...]`


-log-level |  $PROXY_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $PROXY_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $PROXY_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.


-extensions | 
: Run specific extensions during supervised mode.
































### proxy server

Start integrated server

Usage: `proxy server [command options] [arguments...]`


-log-level |  $PROXY_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $PROXY_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $PROXY_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.


-extensions | 
: Run specific extensions during supervised mode.



-log-file |  $PROXY_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-config-file |  $PROXY_CONFIG_FILE
: Path to config file.


-tracing-enabled |  $PROXY_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $PROXY_TRACING_TYPE
: Tracing backend type. Default: `flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger")`.


-tracing-endpoint |  $PROXY_TRACING_ENDPOINT
: Endpoint for the agent.


-tracing-collector |  $PROXY_TRACING_COLLECTOR
: Endpoint for the collector.


-tracing-service |  $PROXY_TRACING_SERVICE
: Service name for tracing. Default: `flags.OverrideDefaultString(cfg.Tracing.Service, "proxy")`.


-debug-addr |  $PROXY_DEBUG_ADDR
: Address to bind debug server. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9205")`.


-debug-token |  $PROXY_DEBUG_TOKEN
: Token to grant metrics access.


-debug-pprof |  $PROXY_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $PROXY_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $PROXY_HTTP_ADDR
: Address to bind http server. Default: `flags.OverrideDefaultString(cfg.HTTP.Addr, "0.0.0.0:9200")`.


-http-root |  $PROXY_HTTP_ROOT
: Root path of http server. Default: `flags.OverrideDefaultString(cfg.HTTP.Root, "/")`.


-asset-path |  $PROXY_ASSET_PATH
: Path to custom assets. Default: `flags.OverrideDefaultString(cfg.Asset.Path, "")`.


-service-namespace |  $PROXY_SERVICE_NAMESPACE
: Set the base namespace for the service namespace. Default: `flags.OverrideDefaultString(cfg.Service.Namespace, "com.owncloud.web")`.


-service-name |  $PROXY_SERVICE_NAME
: Service name. Default: `flags.OverrideDefaultString(cfg.Service.Name, "proxy")`.


-transport-tls-cert |  $PROXY_TRANSPORT_TLS_CERT
: Certificate file for transport encryption. Default: `flags.OverrideDefaultString(cfg.HTTP.TLSCert, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.crt"))`.


-transport-tls-key |  $PROXY_TRANSPORT_TLS_KEY
: Secret file for transport encryption. Default: `flags.OverrideDefaultString(cfg.HTTP.TLSKey, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.key"))`.


-tls |  $PROXY_TLS
: Use TLS (disable only if proxy is behind a TLS-terminating reverse-proxy).. Default: `flags.OverrideDefaultBool(cfg.HTTP.TLS, true)`.


-jwt-secret |  $PROXY_JWT_SECRET , $OCIS_JWT_SECRET
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `flags.OverrideDefaultString(cfg.TokenManager.JWTSecret, "Pive-Fumkiu4")`.


-reva-gateway-addr |  $PROXY_REVA_GATEWAY_ADDR
: REVA Gateway Endpoint. Default: `flags.OverrideDefaultString(cfg.Reva.Address, "127.0.0.1:9142")`.


-insecure |  $PROXY_INSECURE_BACKENDS
: allow insecure communication to upstream servers. Default: `flags.OverrideDefaultBool(cfg.InsecureBackends, false)`.


-oidc-issuer |  $PROXY_OIDC_ISSUER , $OCIS_URL
: OIDC issuer. Default: `flags.OverrideDefaultString(cfg.OIDC.Issuer, "https://localhost:9200")`.


-oidc-insecure |  $PROXY_OIDC_INSECURE
: OIDC allow insecure communication. Default: `flags.OverrideDefaultBool(cfg.OIDC.Insecure, true)`.


-autoprovision-accounts |  $PROXY_AUTOPROVISION_ACCOUNTS
: create accounts from OIDC access tokens to learn new users. Default: `flags.OverrideDefaultBool(cfg.AutoprovisionAccounts, false)`.


-enable-presignedurls |  $PROXY_ENABLE_PRESIGNEDURLS
: Enable or disable handling the presigned urls in the proxy. Default: `flags.OverrideDefaultBool(cfg.PreSignedURL.Enabled, true)`.


-enable-basic-auth |  $PROXY_ENABLE_BASIC_AUTH
: enable basic authentication. Default: `flags.OverrideDefaultBool(cfg.EnableBasicAuth, false)`.


-account-backend-type |  $PROXY_ACCOUNT_BACKEND_TYPE
: account-backend-type. Default: `flags.OverrideDefaultString(cfg.AccountBackend, "accounts")`.



### proxy version

Print the versions of the running instances

Usage: `proxy version [command options] [arguments...]`



































-service-namespace |  $PROXY_SERVICE_NAMESPACE
: Set the base namespace for the service namespace. Default: `flags.OverrideDefaultString(cfg.OIDC.Issuer, "com.owncloud.web")`.


-service-name |  $PROXY_SERVICE_NAME
: Service name. Default: `flags.OverrideDefaultString(cfg.Service.Name, "proxy")`.

