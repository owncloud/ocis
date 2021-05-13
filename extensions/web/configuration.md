---
title: "Configuration"
date: "2021-05-13T06:31:57+0000"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/web/templates
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

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-web reads `web.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/web/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

If multiple variables are listed for one option, they are in order of precedence. This means the leftmost variable will always win if given.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

Serve ownCloud Web for oCIS

Usage: `web [global options] command [command options] [arguments...]`

## Sub Commands

### web health

Check health status

Usage: `web health [command options] [arguments...]`

-debug-addr |  $WEB_DEBUG_ADDR
: Address to debug endpoint. Default: `flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9104")`.

### web server

Start integrated server

Usage: `web server [command options] [arguments...]`

