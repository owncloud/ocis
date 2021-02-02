---
title: "Configuration"
date: "2021-02-02T12:12:27+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/glauth/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-glauth reads `glauth.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/glauth/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### glauth health

Check health status

Usage: `glauth health [command options] [arguments...]`

-debug-addr |  $GLAUTH_DEBUG_ADDR
: Address to debug endpoint. Default: `0.0.0.0:9129`.

### glauth ocis-glauth

Serve GLAuth API for oCIS

Usage: `glauth ocis-glauth [command options] [arguments...]`

-log-level |  $GLAUTH_LOG_LEVEL
: Set logging level. Default: `info`.

-log-pretty |  $GLAUTH_LOG_PRETTY
: Enable pretty logging. Default: `true`.

-log-color |  $GLAUTH_LOG_COLOR
: Enable colored logging. Default: `true`.

### glauth server

Start integrated server

Usage: `glauth server [command options] [arguments...]`

-config-file |  $GLAUTH_CONFIG_FILE
: Path to config file.

-tracing-enabled |  $GLAUTH_TRACING_ENABLED
: Enable sending traces.

-tracing-type |  $GLAUTH_TRACING_TYPE
: Tracing backend type. Default: `jaeger`.

-tracing-endpoint |  $GLAUTH_TRACING_ENDPOINT
: Endpoint for the agent.

-tracing-collector |  $GLAUTH_TRACING_COLLECTOR
: Endpoint for the collector.

-tracing-service |  $GLAUTH_TRACING_SERVICE
: Service name for tracing. Default: `glauth`.

-debug-addr |  $GLAUTH_DEBUG_ADDR
: Address to bind debug server. Default: `0.0.0.0:9129`.

-debug-token |  $GLAUTH_DEBUG_TOKEN
: Token to grant metrics access.

-debug-pprof |  $GLAUTH_DEBUG_PPROF
: Enable pprof debugging.

-debug-zpages |  $GLAUTH_DEBUG_ZPAGES
: Enable zpages debugging.

-role-bundle-id |  $GLAUTH_ROLE_BUNDLE_ID
: roleid used to make internal grpc requests. Default: `71881883-1768-46bd-a24d-a356a2afdf7f`.

-ldap-addr |  $GLAUTH_LDAP_ADDR
: Address to bind ldap server. Default: `0.0.0.0:9125`.

-ldap-enabled |  $GLAUTH_LDAP_ENABLED
: Enable ldap server. Default: `true`.

-ldaps-addr |  $GLAUTH_LDAPS_ADDR
: Address to bind ldap server. Default: `0.0.0.0:9126`.

-ldaps-enabled |  $GLAUTH_LDAPS_ENABLED
: Enable ldap server. Default: `true`.

-ldaps-cert |  $GLAUTH_LDAPS_CERT
: path to ldaps certificate in PEM format. Default: `./ldap.crt`.

-ldaps-key |  $GLAUTH_LDAPS_KEY
: path to ldaps key in PEM format. Default: `./ldap.key`.

-backend-basedn |  $GLAUTH_BACKEND_BASEDN
: base distinguished name to expose. Default: `dc=example,dc=org`.

-backend-name-format |  $GLAUTH_BACKEND_NAME_FORMAT
: name attribute for entries to expose. typically cn or uid. Default: `cn`.

-backend-group-format |  $GLAUTH_BACKEND_GROUP_FORMAT
: name attribute for entries to expose. typically ou, cn or dc. Default: `ou`.

-backend-ssh-key-attr |  $GLAUTH_BACKEND_SSH_KEY_ATTR
: ssh key attribute for entries to expose. Default: `sshPublicKey`.

-backend-datastore |  $GLAUTH_BACKEND_DATASTORE
: datastore to use as the backend. one of accounts, ldap or owncloud. Default: `accounts`.

-backend-insecure |  $GLAUTH_BACKEND_INSECURE
: Allow insecure requests to the datastore. Default: `false`.

-backend-use-graphapi |  $GLAUTH_BACKEND_USE_GRAPHAPI
: use Graph API, only for owncloud datastore. Default: `true`.

-fallback-basedn |  $GLAUTH_FALLBACK_BASEDN
: base distinguished name to expose. Default: `dc=example,dc=org`.

-fallback-name-format |  $GLAUTH_FALLBACK_NAME_FORMAT
: name attribute for entries to expose. typically cn or uid. Default: `cn`.

-fallback-group-format |  $GLAUTH_FALLBACK_GROUP_FORMAT
: name attribute for entries to expose. typically ou, cn or dc. Default: `ou`.

-fallback-ssh-key-attr |  $GLAUTH_FALLBACK_SSH_KEY_ATTR
: ssh key attribute for entries to expose. Default: `sshPublicKey`.

-fallback-datastore |  $GLAUTH_FALLBACK_DATASTORE
: datastore to use as the fallback. one of accounts, ldap or owncloud.

-fallback-insecure |  $GLAUTH_FALLBACK_INSECURE
: Allow insecure requests to the datastore. Default: `false`.

-fallback-use-graphapi |  $GLAUTH_FALLBACK_USE_GRAPHAPI
: use Graph API, only for owncloud datastore. Default: `true`.

