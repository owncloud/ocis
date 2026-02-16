---
title: Service Independent Envvars
date: 2025-11-13T00:00:00+00:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/service-independent/
geekdocFilePath: service_independent_envvars.md
---

The following environment variables are service independent. You will find the respective code in the `ocis-pkg/` directory. See the [Admin Documentation - Environment Variables with Special Scopes](https://doc.owncloud.com/ocis/7.3/deployment/services/env-vars-special-scope.html) for a comprehensive list and explanation.

{{< toc >}}

{{< hint info >}}
See the [Environment Variables]({{< ref "../services/general-info/envvars/" >}}) documentation for common and important details on envvars.
{{< /hint >}}

## Service Registry

This package configures the service registry which will be used to look up for example the service addresses.

Available registries are:

-   nats-js-kv (default)
-   memory

To configure which registry to use, you have to set the environment variable `MICRO_REGISTRY`, and for all except `memory` you also have to set the registry address via `MICRO_REGISTRY_ADDRESS` and other envvars.

## Startup Related Envvars

These envvars define the startup of ocis and can for example add or remove services from the startup process such as `OCIS_ADD_RUN_SERVICES`.

## Memory Limits

{{< hint info >}}
Note that this envvar is for development purposes only and not described in the admin docs.
{{< /hint >}}

oCIS will automatically set the go native `GOMEMLIMIT` to `0.9`. To disable the limit set `AUTOMEMLIMIT=off`. For more information take a look at the official [Guide to the Go Garbage Collector](https://go.dev/doc/gc-guide).
