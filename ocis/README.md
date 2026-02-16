# ocis

The ocis package includes the Infinite Scale runtime and commands for the Infinite Scale command-line interface (CLI), which are not bound to a service.

Table of Contents
=================

   * [Service Registry](README.md#service-registry)
   * [Memory limits](README.md#memory-limits)
   * [CLI Commands](README.md#cli-commands)

<!-- Created by https://github.com/ekalinin/github-markdown-toc -->

## Service Registry

This package also configures the service registry which will be used to look up the service addresses.

Available registries are:

-   nats-js-kv (default)
-   memory

To configure which registry to use, you have to set the environment variable `MICRO_REGISTRY`, and for all except `memory` you also have to set the registry address via `MICRO_REGISTRY_ADDRESS`.

## Memory limits

oCIS will automatically set the go native `GOMEMLIMIT` to `0.9`. To disable the limit set `AUTOMEMEMLIMIT=off`. For more information take a look at the official [Guide to the Go Garbage Collector](https://go.dev/doc/gc-guide).

## CLI Commands

See the [dev docs, Service Independent CLI](https://owncloud.dev/cli-commands/service_independent_cli/) for more details.
