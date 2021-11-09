---
title: "Configuration"
date: "2021-11-09T00:03:16+0100"
weight: 2
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/ocis/templates
geekdocFilePath: config.md
---

{{< toc >}}

## Configuration Framework

In order to simplify deployments and development the configuration model from oCIS aims to be simple yet flexible.

## Overview of the approach

{{< svg src="ocis/static/ocis-config-redesign.drawio.svg" >}}

## In-depth configuration

Since we include a set of predefined extensions within the single binary, configuring an extension can be done in a variety of ways. Since we work with complex types, having as many cli per config value scales poorly, so we limited the options to config files and environment variables, leaving cli flags for shared values, such as config file sources (`--config-file`) or logging (`--log-level`, `--log-pretty`, `--log-file` or `--log-color`).

The hierarchy is clear enough, leaving us with:

_(each element above overwrites its precedent)_

1. env variables
2. extension config
3. ocis config

This is manifested in the previous diagram. We can then speak about "configuration file arithmetics", where resulting config transformations happen through a series of steps. An administrator must be aware of these sources, since mis-managing them can be a source of confusion, having undesired transformations on config files believed not to be applied.

## Flows

Let's explore the various flows with examples and workflows.

### Examples

Let's explore with examples this approach.

####  Only config files

The following config files are present in the default loading locations:

_ocis.yaml_
```yaml
proxy:
  http:
    addr: localhost:1111
  log:
    pretty: false
    color: false
    level: info
accounts:
  http:
    addr: localhost:2222
  log:
    level: debug
    color: false
    pretty: false
log:
  pretty: true
  color: true
  level: info
```

_proxy.yaml_
```yaml
http:
  addr: localhost:3333
```

_accounts.yaml_
```yaml
http:
  addr: localhost:4444
```

Note that the extension files will overwrite values from the main `ocis.yaml`, causing `ocis server` to run with the following configuration:

```yaml
proxy:
  http:
    addr: localhost:3333
accounts:
  http:
    addr: localhost:4444
log:
  pretty: true
  color: true
  level: info
```

#### Using ENV variables

The logging configuration if defined in the main ocis.yaml is inherited by all extensions. It can be, however, overwritten by a single extension file if desired. The same example can be used to demonstrate environment values overwrites. With the same set of config files now we have the following command `PROXY_HTTP_ADDR=localhost:5555 ocis server`, now the resulting config looks like:

```yaml
proxy:
  http:
    addr: localhost:5555
accounts:
  http:
    addr: localhost:4444
log:
  pretty: true
  color: true
  level: info
```

### Workflows

Since one can run an extension using the runtime (supervised) or not (unsupervised), we ensure correct behavior in both modes, expecting the same outputs.

#### Supervised

You are using the supervised mode whenever you issue the `ocis server` command. We start the runtime on port `9250` (by default) that listens for commands regarding the lifecycle of the supervised extensions. When an extension runs supervised and is killed, the only way to provide / overwrite configuration values will be through an extension config file. This is due to the parent process has already started, and it already has its own environment.

#### Unsupervised

All the points from the priority section hold true. An unsupervised extension can be started with the format: `ocis [extension]` i.e: `ocis proxy`. First, `ocis.yaml` is parsed, then `proxy.yaml` followed by environment variables.

## Shared Values

When running in supervised mode (`ocis server`) it is beneficial to have common values for logging, so that the log output is correctly formatted, or everything is piped to the same file without duplicating config keys and values all over the place. This is possible using the global `log` config key:

_ocis.yaml_
```yaml
log:
  level: error
  color: true
  pretty: true
  file: /var/tmp/ocis_output.log
```

There is, however, the option for extensions to overwrite this global values by declaring their own logging directives:

_ocis.yaml_
```yaml
log:
  level: info
  color: false
  pretty: false
```

One can go as far as to make the case of an extension overwriting its shared logging config that received from the main `ocis.yaml` file. Because things can get out of hands pretty fast we recommend not mixing logging configuration values and either use the same global logging values for all extensions.

{{< hint warning >}}
When overwriting a globally shared logging values, one *MUST* specify all values.
{{< /hint >}}

### Log config keys

```yaml
log:
  level: [ error | warning | info | debug ]
  color: [ true | false ]
  pretty: [ true | false ]
  file: [ path/to/log/file ] # MUST not be used with pretty = true
```

## Default config values (in yaml)

TBD. Needs to be generated and merged with the env mappings.
