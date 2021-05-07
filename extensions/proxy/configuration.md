---
title: "Configuration"
date: "2021-05-07T10:50:10+0000"
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

--debug-addr |  $PROXY_DEBUG_ADDR
: Address to debug endpoint. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9109")`.

### proxy ocis-proxy

proxy for oCIS

Usage: `proxy ocis-proxy [command options] [arguments...]`

--log-level |  $PROXY_LOG_LEVEL , $OCIS_LOG_LEVEL
: Set logging level.

--log-pretty |  $PROXY_LOG_PRETTY , $OCIS_LOG_PRETTY
: Enable pretty logging.

--log-color |  $PROXY_LOG_COLOR , $OCIS_LOG_COLOR
: Enable colored logging.

### proxy server

Start integrated server

Usage: `proxy server [command options] [arguments...]`

### proxy version

Print the versions of the running instances

Usage: `proxy version [command options] [arguments...]`

--service-namespace |  $PROXY_SERVICE_NAMESPACE
: Set the base namespace for the service namespace. Default: `flags.OverrideDefaultString(cfg.OIDC.Issuer, "com.owncloud.web")`.

--service-name |  $PROXY_SERVICE_NAME
: Service name. Default: `flags.OverrideDefaultString(cfg.Service.Name, "proxy")`.

