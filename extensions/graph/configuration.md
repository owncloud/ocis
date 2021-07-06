---
title: "Configuration"
date: "2021-07-06T12:26:35+0000"
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

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### graph health

Check health status

Usage: `graph health [command options] [arguments...]`






-debug-addr |  $GRAPH_DEBUG_ADDR
: Address to debug endpoint. Default: `"0.0.0.0:9124"`.


























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


-tracing-enabled |  $GRAPH_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $GRAPH_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $GRAPH_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $GRAPH_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $GRAPH_TRACING_SERVICE
: Service name for tracing. Default: `"graph"`.


-debug-addr |  $GRAPH_DEBUG_ADDR
: Address to bind debug server. Default: `"0.0.0.0:9124"`.


-debug-token |  $GRAPH_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $GRAPH_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $GRAPH_DEBUG_ZPAGES
: Enable zpages debugging.


-http-addr |  $GRAPH_HTTP_ADDR
: Address to bind http server. Default: `"0.0.0.0:9120"`.


-http-root |  $GRAPH_HTTP_ROOT
: Root path of http server. Default: `"/graph"`.


-http-namespace |  $GRAPH_HTTP_NAMESPACE
: Set the base namespace for the http service for service discovery. Default: `"com.owncloud.web"`.


-ldap-network |  $GRAPH_LDAP_NETWORK
: Network protocol to use to connect to the Ldap server. Default: `"tcp"`.


-ldap-address |  $GRAPH_LDAP_ADDRESS
: Address to connect to the Ldap server. Default: `"0.0.0.0:9125"`.


-ldap-username |  $GRAPH_LDAP_USERNAME
: User to bind to the Ldap server. Default: `"cn=admin,dc=example,dc=org"`.


-ldap-password |  $GRAPH_LDAP_PASSWORD
: Password to bind to the Ldap server. Default: `"admin"`.


-ldap-basedn-users |  $GRAPH_LDAP_BASEDN_USERS
: BaseDN to look for users. Default: `"ou=users,dc=example,dc=org"`.


-ldap-basedn-groups |  $GRAPH_LDAP_BASEDN_GROUPS
: BaseDN to look for users. Default: `"ou=groups,dc=example,dc=org"`.


-oidc-endpoint |  $GRAPH_OIDC_ENDPOINT , $OCIS_URL
: OpenIDConnect endpoint. Default: `"https://localhost:9200"`.


-oidc-insecure |  $GRAPH_OIDC_INSECURE
: OpenIDConnect endpoint.


-oidc-realm |  $GRAPH_OIDC_REALM
: OpenIDConnect realm. Default: `""`.


-reva-gateway-addr |  $REVA_GATEWAY_ADDR
: REVA Gateway Endpoint. Default: `"127.0.0.1:9142"`.


-webdav-namespace |  $STORAGE_WEBDAV_NAMESPACE
: Namespace prefix for the webdav endpoint. Default: `"/home"`.


-extensions | 
: Run specific extensions during supervised mode. This flag is set by the runtime.

