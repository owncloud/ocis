---
title: "Configuration"
date: "2021-10-26T05:42:09+0000"
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

### Command-line flags

If you prefer to configure the service with command-line flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### glauth health

Check health status

Usage: `glauth health [command options] [arguments...]`





-debug-addr |  $GLAUTH_DEBUG_ADDR
: Address to debug endpoint. Default: `"127.0.0.1:9129"`.




































### glauth ocis-glauth

Serve GLAuth API for oCIS

Usage: `glauth ocis-glauth [command options] [arguments...]`


-log-level |  $GLAUTH_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.


-log-pretty |  $GLAUTH_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.


-log-color |  $GLAUTH_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.





































### glauth server

Start integrated server

Usage: `glauth server [command options] [arguments...]`






-log-file |  $GLAUTH_LOG_FILE , $OCIS_LOG_FILE
: Enable log to file.


-config-file |  $GLAUTH_CONFIG_FILE
: Path to config file. Default: `""`.


-tracing-enabled |  $GLAUTH_TRACING_ENABLED , $OCIS_TRACING_ENABLED
: Enable sending traces.


-tracing-type |  $GLAUTH_TRACING_TYPE , $OCIS_TRACING_TYPE
: Tracing backend type. Default: `"jaeger"`.


-tracing-endpoint |  $GLAUTH_TRACING_ENDPOINT , $OCIS_TRACING_ENDPOINT
: Endpoint for the agent. Default: `""`.


-tracing-collector |  $GLAUTH_TRACING_COLLECTOR , $OCIS_TRACING_COLLECTOR
: Endpoint for the collector. Default: `""`.


-tracing-service |  $GLAUTH_TRACING_SERVICE
: Service name for tracing. Default: `"glauth"`.


-debug-addr |  $GLAUTH_DEBUG_ADDR
: Address to bind debug server. Default: `"127.0.0.1:9129"`.


-debug-token |  $GLAUTH_DEBUG_TOKEN
: Token to grant metrics access. Default: `""`.


-debug-pprof |  $GLAUTH_DEBUG_PPROF
: Enable pprof debugging.


-debug-zpages |  $GLAUTH_DEBUG_ZPAGES
: Enable zpages debugging.


-role-bundle-id |  $GLAUTH_ROLE_BUNDLE_ID
: roleid used to make internal grpc requests. Default: `"71881883-1768-46bd-a24d-a356a2afdf7f"`.


-ldap-addr |  $GLAUTH_LDAP_ADDR
: Address to bind ldap server. Default: `"127.0.0.1:9125"`.


-ldap-enabled |  $GLAUTH_LDAP_ENABLED
: Enable ldap server. Default: `true`.


-ldaps-addr |  $GLAUTH_LDAPS_ADDR
: Address to bind ldaps server. Default: `"127.0.0.1:9126"`.


-ldaps-enabled |  $GLAUTH_LDAPS_ENABLED
: Enable ldaps server. Default: `true`.


-ldaps-cert |  $GLAUTH_LDAPS_CERT
: path to ldaps certificate in PEM format. Default: `flags.OverrideDefaultString(cfg.Ldaps.Cert, path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"))`.


-ldaps-key |  $GLAUTH_LDAPS_KEY
: path to ldaps key in PEM format. Default: `flags.OverrideDefaultString(cfg.Ldaps.Key, path.Join(defaults.BaseDataPath(), "ldap", "ldap.key"))`.


-backend-basedn |  $GLAUTH_BACKEND_BASEDN
: base distinguished name to expose. Default: `"dc=ocis,dc=test"`.


-backend-name-format |  $GLAUTH_BACKEND_NAME_FORMAT
: name attribute for entries to expose. typically cn or uid. Default: `"cn"`.


-backend-group-format |  $GLAUTH_BACKEND_GROUP_FORMAT
: name attribute for entries to expose. typically ou, cn or dc. Default: `"ou"`.


-backend-ssh-key-attr |  $GLAUTH_BACKEND_SSH_KEY_ATTR
: ssh key attribute for entries to expose. Default: `"sshPublicKey"`.


-backend-datastore |  $GLAUTH_BACKEND_DATASTORE
: datastore to use as the backend. one of accounts, ldap or owncloud. Default: `"accounts"`.


-backend-insecure |  $GLAUTH_BACKEND_INSECURE
: Allow insecure requests to the datastore. Default: `false`.


-backend-server |  $GLAUTH_BACKEND_SERVERS
: `--backend-server https://demo.owncloud.com/apps/graphapi/v1.0 [--backend-server "https://demo2.owncloud.com/apps/graphapi/v1.0"]`. Default: `cli.NewStringSlice()`.


-backend-use-graphapi |  $GLAUTH_BACKEND_USE_GRAPHAPI
: use Graph API, only for owncloud datastore. Default: `true`.


-fallback-basedn |  $GLAUTH_FALLBACK_BASEDN
: base distinguished name to expose. Default: `"dc=ocis,dc=test"`.


-fallback-name-format |  $GLAUTH_FALLBACK_NAME_FORMAT
: name attribute for entries to expose. typically cn or uid. Default: `"cn"`.


-fallback-group-format |  $GLAUTH_FALLBACK_GROUP_FORMAT
: name attribute for entries to expose. typically ou, cn or dc. Default: `"ou"`.


-fallback-ssh-key-attr |  $GLAUTH_FALLBACK_SSH_KEY_ATTR
: ssh key attribute for entries to expose. Default: `"sshPublicKey"`.


-fallback-datastore |  $GLAUTH_FALLBACK_DATASTORE
: datastore to use as the fallback. one of accounts, ldap or owncloud. Default: `""`.


-fallback-insecure |  $GLAUTH_FALLBACK_INSECURE
: Allow insecure requests to the datastore. Default: `false`.


-fallback-server |  $GLAUTH_FALLBACK_SERVERS
: `--fallback-server http://internal1.example.com [--fallback-server http://internal2.example.com]`. Default: `cli.NewStringSlice("https://demo.owncloud.com/apps/graphapi/v1.0")`.


-fallback-use-graphapi |  $GLAUTH_FALLBACK_USE_GRAPHAPI
: use Graph API, only for owncloud datastore. Default: `true`.


-extensions | 
: Run specific extensions during supervised mode. This flag is set by the runtime.

