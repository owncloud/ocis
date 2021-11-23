Change: Restructure Configuration Parsing

Tags: ocis

- Adds a new dependency that does roughly what I started to do on a side project. Not reinventing the wheel. MIT license.
- Remove cli flags on subcommands, not on ocis root command.
- Sane default propagation.
- Lays the foundation for easy config dump and service restart on config watch.
- Support for environment variables.
- Support for merging config values.
- Support for defaults.
- And technically speaking, the most important aspect is that it is thread safe, as I can create as many config.Config structs as desired.

https://github.com/owncloud/ocis/pull/2708
