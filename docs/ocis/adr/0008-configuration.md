---
title: "8. Configuration"
weight: 8
date: 2021-05-03T15:00:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0008-configuration.md
---

* Status: proposed
* Deciders: @refs, @butonic, @micbar, @dragotin, @pmaier1
* Date: 2021-05-03

## Context and Problem Statement

As per urfave/cli's doc:

>The precedence for flag value sources is as follows (highest to lowest):
>
>0. Command line flag value from user
>1. Environment variable (if specified)
>2. Configuration file (if specified)
>3. Default defined on the flag

An issue arises in point 2, in the sense that configuration file refers to a single file containing the value for the env variable. The CLI framework we use for flag parsing does not support merging config structs with CLI flags. This introduces an inconsistency with the framework: config structs are not supported, and we cannot hook to the lifecycle of the flags parsing to use a file as source and conform to these rules.

Because we solely rely on [structured configuration](https://github.com/owncloud/ocis/blob/master/ocis-pkg/config/config.go) we need a way to modify values in this struct using the provided means urfave/cli gives us (flags, env variables, config files and default value), but since we have different modes of operation (supervised Vs. unsupervised) we have to define a clear line.

### Decision Drivers
- Improve experience for the end user.
- Improve experience for developers.
- Sane defaults.
- Sane overrides.

### Considered Options

- Extend [FlagInputSourceExtension interface](https://github.com/urfave/cli/blob/master/altsrc/flag.go#L12-L17)
- Feature request: support for structured configuration (urfave/cli).
- Clearly defined boundaries of what can and cannot be done.
- Expose structured field values as CLI flags
- Drop support for structure configuration
- Adapt the "structured config files have the highest priority" within oCIS

### Decision Outcome

[STILL UNDECIDED]

#### Positive Consequences

[TBD, depends on Decision Outcome]

### Pros and Cons of the Options

#### Extend FlagInputSourceExtension interface
- Good, because we could still use Viper to load from config files here and apply values to the flags in the context.
- Bad, because urfave/cli team are [actively working on v3 of altsrc](https://github.com/urfave/cli/issues/1051#issuecomment-606311923) and we don't want to maintain yet another slice of the codebase.

notes: source is  [FlagInputSourceExtension interface](https://github.com/urfave/cli/blob/master/altsrc/flag.go#L12-L17)

#### Feature request: support for structured configuration (urfave/cli).
- Good, because we could remove Viper off the codebase and solely rely on urfave/cli's native code.
- Bad, because there are no plans to support this upstream.

#### Clearly defined boundaries of what can and cannot be done.

- Good, because no changes to the codebase required (not drastic changes.)
- Bad, because we're limited by the framework

#### Expose structured field values as CLI flags

- Good, because it has been already taken into account on large projects (kubernetes) [here.](https://docs.google.com/document/d/1Dvct469xfjkgy3tjWMAKvRAJo4CmGH4cgSVGTDpay6A) in point 5.
- Bad, because it requires quite a bit<sup>1</sup> of custom logic.
- Bad, because how should these flags be present in the `-h` menu of a subcommand? Probably some code generation needed.

*[1] this is a big uncertainty.

#### Drop support for structure configuration

- Good, because it makes the integration with the cli framework easier to grasp.
- Good, because it is not encouraged by the 12factor app spec.
- Bad, because we already support if and users make active use of it. At least for development.

#### Adapt the "structured config files have the highest priority" within oCIS

- Good, because that would mean little structural changes to the codebase since the Viper config parsing logic already uses the `Before` hook to parse prior to the command's action executes.

### Notes

#### Use Cases and Expected Behaviors

##### Supervised (`ocis server` or `ocis run extension`)

![grafik](https://user-images.githubusercontent.com/6905948/116872568-62b1a780-ac16-11eb-9f29-030a651ee39b.png)

- Use a global config file (ocis.yaml) to configure an entire set of services: `> ocis --config-file /etc/ocis.yaml service`
- Use a global config file (ocis.yaml) to configure a single extension: `> ocis --config-file /etc/ocis/yaml proxy`
- When running in supervised mode, config files from extensions are NOT evaluated (only when running `ocis server`, runs with `ocis run extension` do parse individual config files)
  - i.e: present config files: `ocis.yaml` and `proxy.yaml`; only the contents of `ocis.yaml` are loaded<sup>1</sup>.
- Flag parsing for subcommands are not allowed in this mode, since the runtime is in control. Configuration has to be done solely using config files.

*[1] see the development section for more on this topic.

###### Known Gotchas
- `> ocis --config-file /etc/ocis/ocis.yaml server` does not work. It currently only supports reading global config values from the predefined locations.

##### Unsupervised (`ocis proxy`)

![grafik](https://user-images.githubusercontent.com/6905948/116872534-54fc2200-ac16-11eb-8267-ffe7b03177b3.png)

- `ocis.yaml` is parsed first (since `proxy` is a subcommand of `ocis`)
- `proxy.yaml` is parsed if present, overriding values from `ocis.yaml` and any cli flag or env variable present.

#### Other known use cases

- Configure via env + some configuration files like WEB_UI_CONFIG or proxy routes
- Configure via flags + some configuration files like WEB_UI_CONFIG or proxy routes
- Configure via global (single file for all extensions) config file + some configuration files like WEB_UI_CONFIG or proxy routes
- configure via per extension config file + some configuration files like WEB_UI_CONFIG or proxy routes

Each individual use case DOES NOT mix sources (i.e: when using cli flags, do not use environment variables nor cli flags).

_Limitations on urfave/cli prevent us from providing structured configuration and framework support for cli flags + env variables._

#### Use Cases for Development

#### Config Loading

Sometimes is desired to decouple the main series of services from an individual instance. We want to use the runtime to startup all services, then do work only on a single service. To achieve that one could use `ocis server && ocis kill proxy && ocis run proxy`. This series of commands will 1. load all config from `ocis.yaml`, 2. kill the supervised proxy service and 3. start the same service with the contents from `proxy.yaml`.

#### Start an extension multiple times with different configs (in Supervised mode)

Flag parsing on subcommands in supervised mode is not yet allowed. The runtime will first parse the global `ocis.yaml` (if any) and run with the loaded configuration. This use case should provide support for having 2 different proxy config files and making use of the runtime start 2 proxy services, with different values.

For this to work, services started via `Service.Start` need to forward any args as flags:

```go
if err := client.Call("Service.Start", os.Args[2], &reply); err != nil {
  log.Fatal(err)
}
```

This should provide with enough flexibility for interpreting different config sources as: `> bin/ocis run proxy --config-file /etc/ocis/unexpected/proxy.yaml`

#### Developing Considered Alternatives Further

Let's develop further the following concept: Adapt the "structured config files have the highest priority" within oCIS.

Of course it directly contradicts urfave/cli priorities. When a command finished parsing its cli args and env variables, only after that `Before` is called. This mean by the time we reach a command `Before` hook, flags have already been parsed and its values loaded to their respective destinations within the `Config` struct.

This should still not prevent a developer from using different config files for a single service. Let's analyze the following use case:

1. global config file present (ocis.yaml)
2. single proxy.yaml config file
3. another proxy.yaml config file
4. running under supervision mode

The outcome of the following set of commands should be having all bootstrapped services running + 2 proxies on different addresses:

```console
> ocis server
> ocis kill proxy
> ocis run proxy --config-file proxy.yaml
> ocis run proxy --config-file proxy2.yaml
```

This is a desired use case that is yet not supported due to lacking of flags forwarding.

#### Follow up PR's

- Variadic runtime extensions to run (development mostly)
- Arg forwarding to command (when running in supervised mode, forward any --config-file flag to supervised subcommands)
- Ability to set `OCIS_URL` from a config file (this would require to extend the ocis-pkg/config/config.go file).

#### The case for `OCIS_URL`

`OCIS_URL` is a jack-of-all trades configuration. It is meant to ease up providing defaults and ensuring dependant services are well configured. It is an override to the following env vars:

```
OCIS_IDM_ADDRESS
PROXY_OIDC_ISSUER
STORAGE_OIDC_ISSUER
STORAGE_FRONTEND_PUBLIC_URL
STORAGE_LDAP_IDP
WEB_UI_CONFIG_SERVER
WEB_OIDC_AUTHORITY
OCIS_PUBLIC_URL
```

Because this functionality is only available as an env var, there is no current way to "normalize" its usage with a config file. That is, there is no way to individually set `OCIS_URL` via config file. This is clear technical debt, and should be added functionality.

#### State of the Art
- [Kubernetes proposal on this very same topic](https://docs.google.com/document/d/1Dvct469xfjkgy3tjWMAKvRAJo4CmGH4cgSVGTDpay6A)
- [Configuration \| Pulumi](https://www.pulumi.com/docs/intro/concepts/config/)
  - Configuration can be altered via setters through the CLI.
