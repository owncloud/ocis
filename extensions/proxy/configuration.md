---
title: "Configuration"
date: "2021-07-14T11:12:54+0000"
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
: Address to debug endpoint. Default: `"0.0.0.0:9109"`.



































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
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $PROXY_TRACING_ENDPOINT
: Endpoint for the agent.


-tracing-collector |  $PROXY_TRACING_COLLECTOR
: Endpoint for the collector.


-tracing-service |  $PROXY_TRACING_SERVICE
: Service name for tracing. Default: `"proxy"`.


-debug-addr |  $PROXY_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9205"`.


-debug-token |  $PROXY_DEBUG_TOKEN
: Token to grant metrics access.


-debug-pprof |  $PROXY_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $PROXY_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $PROXY_HTTP_ADDR
: Address to bind http server. Default: `"0.0.0.0:9200"`.


-http-root |  $PROXY_HTTP_ROOT
: Root path of http server. Default: `"/"`.


-asset-path |  $PROXY_ASSET_PATH
: Path to custom assets. Default: `""`.


-service-namespace |  $PROXY_SERVICE_NAMESPACE
: Set the base namespace for the service namespace. Default: `"com.owncloud.web"`.


-service-name |  $PROXY_SERVICE_NAME
: Service name. Default: `"proxy"`.


-transport-tls-cert |  $PROXY_TRANSPORT_TLS_CERT
: Certificate file for transport encryption. Default: `flags.OverrideDefaultString(cfg.HTTP.TLSCert, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.crt"))`.


-transport-tls-key |  $PROXY_TRANSPORT_TLS_KEY
: Secret file for transport encryption. Default: `flags.OverrideDefaultString(cfg.HTTP.TLSKey, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.key"))`.


-tls |  $PROXY_TLS
: Use TLS (disable only if proxy is behind a TLS-terminating reverse-proxy).. Default: `true`.


-jwt-secret |  $PROXY_JWT_SECRET , $OCIS_JWT_SECRET
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `"Pive-Fumkiu4"`.


-reva-gateway-addr |  $PROXY_REVA_GATEWAY_ADDR
: REVA Gateway Endpoint. Default: `"127.0.0.1:9142"`.


-insecure |  $PROXY_INSECURE_BACKENDS
: allow insecure communication to upstream servers. Default: `false`.


-oidc-issuer |  $PROXY_OIDC_ISSUER , $OCIS_URL
: OIDC issuer. Default: `"https://localhost:9200"`.


-oidc-insecure |  $PROXY_OIDC_INSECURE
: OIDC allow insecure communication. Default: `true`.


-oidc-userinfo-cache-tll |  $PROXY_OIDC_USERINFO_CACHE_TTL
: Fallback TTL in seconds for caching userinfo, when no token lifetime can be identified. Default: `10`.


-oidc-userinfo-cache-size |  $PROXY_OIDC_USERINFO_CACHE_SIZE
: Max entries for caching userinfo. Default: `1024`.


-autoprovision-accounts |  $PROXY_AUTOPROVISION_ACCOUNTS
: create accounts from OIDC access tokens to learn new users. Default: `false`.


-presignedurl-allow-method |  $PRESIGNEDURL_ALLOWED_METHODS
: --presignedurl-allow-method GET [--presignedurl-allow-method POST]. Default: `cli.NewStringSlice("GET")`.


-enable-presignedurls |  $PROXY_ENABLE_PRESIGNEDURLS
: Enable or disable handling the presigned urls in the proxy. Default: `true`.


-enable-basic-auth |  $PROXY_ENABLE_BASIC_AUTH
: enable basic authentication. Default: `false`.


-account-backend-type |  $PROXY_ACCOUNT_BACKEND_TYPE
: account-backend-type. Default: `"accounts"`.


-proxy-user-agent-lock-in |  $PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT
: --user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of [UserAgent:challenge] it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the proxy it configures the authentication middlewares..



### proxy version

Print the versions of the running instances

Usage: `proxy version [command options] [arguments...]`







































-service-namespace |  $PROXY_SERVICE_NAMESPACE
: Set the base namespace for the service namespace. Default: `"com.owncloud.web"`.


-service-name |  $PROXY_SERVICE_NAME
: Service name. Default: `"proxy"`.

