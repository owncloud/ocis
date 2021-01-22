---
title: "Configuration"
date: "2021-01-22T08:28:12+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/konnectd/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-konnectd reads `konnectd.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### konnectd health

Check health status

Usage: `konnectd health [command options] [arguments...]`

-debug-addr |  $KONNECTD_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9134`.

### konnectd ocis-konnectd

Serve Konnectd API for oCIS

Usage: `konnectd ocis-konnectd [command options] [arguments...]`

-log-level |  $KONNECTD_LOG_LEVEL
: Set logging level. Default: `info`.

-log-pretty |  $KONNECTD_LOG_PRETTY
: Enable pretty logging. Default: `true`.

-log-color |  $KONNECTD_LOG_COLOR
: Enable colored logging. Default: `true`.

### konnectd server

Start integrated server

Usage: `konnectd server [command options] [arguments...]`

-config-file |  $KONNECTD_CONFIG_FILE
: Path to config file.

-tracing-enabled |  $KONNECTD_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $KONNECTD_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.

-tracing-endpoint |  $KONNECTD_TRACING_ENDPOINT
: Endpoint for the agent.

-tracing-collector |  $KONNECTD_TRACING_COLLECTOR
: Endpoint for the collector.

-tracing-service |  $KONNECTD_TRACING_SERVICE
: Service name for tracing. Default: `konnectd`.

-debug-addr |  $KONNECTD_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9134`.

-debug-token |  $KONNECTD_DEBUG_TOKEN
: Token to grant metrics access.

-debug-pprof |  $KONNECTD_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $KONNECTD_DEBUG_ZPAGES
: Enable zpages debugging.

-http-addr |  $KONNECTD_HTTP_ADDR
: Address to bind http server. Default: `0.0.0.0:9130`.

-http-root |  $KONNECTD_HTTP_ROOT
: Root path of http server. Default: `/`.

-http-namespace |  $KONNECTD_HTTP_NAMESPACE
: Set the base namespace for service discovery. Default: `com.owncloud.web`.

-name |  $KONNECTD_NAME
: Service name. Default: `konnectd`.

-identity-manager |  $KONNECTD_IDENTITY_MANAGER
: Identity manager (one of ldap,kc,cookie,dummy). Default: `ldap`.

-transport-tls-cert |  $KONNECTD_TRANSPORT_TLS_CERT
: Certificate file for transport encryption.

-transport-tls-key |  $KONNECTD_TRANSPORT_TLS_KEY
: Secret file for transport encryption.

-iss |  $KONNECTD_ISS , $OCIS_URL
: OIDC issuer URL. Default: `https://localhost:9200`.

-signing-kid |  $KONNECTD_SIGNING_KID
: Value of kid field to use in created tokens (uniquely identifying the signing-private-key).

-validation-keys-path |  $KONNECTD_VALIDATION_KEYS_PATH
: Full path to a folder containg PEM encoded private or public key files used for token validaton (file name without extension is used as kid).

-encryption-secret |  $KONNECTD_ENCRYPTION_SECRET
: Full path to a file containing a %d bytes secret key.

-signing-method |  $KONNECTD_SIGNING_METHOD
: JWT default signing method. Default: `PS256`.

-uri-base-path |  $KONNECTD_URI_BASE_PATH
: Custom base path for URI endpoints.

-sign-in-uri |  $KONNECTD_SIGN_IN_URI
: Custom redirection URI to sign-in form.

-signed-out-uri |  $KONNECTD_SIGN_OUT_URI
: Custom redirection URI to signed-out goodbye page.

-authorization-endpoint-uri |  $KONNECTD_ENDPOINT_URI
: Custom authorization endpoint URI.

-endsession-endpoint-uri |  $KONNECTD_ENDSESSION_ENDPOINT_URI
: Custom endsession endpoint URI.

-asset-path |  $KONNECTD_ASSET_PATH
: Path to custom assets.

-identifier-client-path |  $KONNECTD_IDENTIFIER_CLIENT_PATH
: Path to the identifier web client base folder. Default: `/var/tmp/ocis/konnectd`.

-identifier-registration-conf |  $KONNECTD_IDENTIFIER_REGISTRATION_CONF
: Path to a identifier-registration.yaml configuration file. Default: `./config/identifier-registration.yaml`.

-identifier-scopes-conf |  $KONNECTD_IDENTIFIER_SCOPES_CONF
: Path to a scopes.yaml configuration file.

-insecure |  $KONNECTD_INSECURE
: Disable TLS certificate and hostname validation.

-tls |  $KONNECTD_TLS
: Use TLS (disable only if konnectd is behind a TLS-terminating reverse-proxy).. Default: `false`.

-allow-client-guests |  $KONNECTD_ALLOW_CLIENT_GUESTS
: Allow sign in of client controlled guest users.

-allow-dynamic-client-registration |  $KONNECTD_ALLOW_DYNAMIC_CLIENT_REGISTRATION
: Allow dynamic OAuth2 client registration. Default: `true`.

-disable-identifier-webapp |  $KONNECTD_DISABLE_IDENTIFIER_WEBAPP
: Disable built-in identifier-webapp to use a frontend hosted elsewhere.. Default: `true`.

### konnectd version

Print the versions of the running instances

Usage: `konnectd version [command options] [arguments...]`

-http-namespace |  $KONNECTD_HTTP_NAMESPACE
: Set the base namespace for service discovery. Default: `com.owncloud.web`.

-name |  $KONNECTD_NAME
: Service name. Default: `konnectd`.

