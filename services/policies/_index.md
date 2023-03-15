---
title: Policies
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/policies
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The policies service provides a new grpc api which can be used to return whether a requested operation is allowed or not. To do so, Open Policy Agent (OPA) is used to determine the set of rules of what is permitted and what is not.

## Table of Contents

{{< toc-tree >}}

## Rego

Policies are written in the [rego query language](https://www.openpolicyagent.org/docs/latest/policy-language/). The location of the rego files can be configured via yaml, a configuration via environment variables is not possible.

The Policies Service consists of the following modules:

*   Proxy Authorization (middleware)
*   Event Authorization (async post-processing)
*   GRPC API (can be used from other services)

To configure the Policies Service, three environment variables need to be defined:

*   `POLICIES_ENGINE_TIMEOUT`
*   `POLICIES_POSTPROCESSING_QUERY`
*   `PROXY_POLICIES_QUERY`

Note that each query setting defines the [Complete Rules](https://www.openpolicyagent.org/docs/latest/#complete-rules) variable defined in the rego rule set the corresponding step uses for the evaluation. If the variable is mistyped or not found, the evaluation defaults to deny. Individual query definitions can be defined for each module.

To activate a the policies service for a module, it must be started with a yaml configuration that points to one or more rego files. Note that if the service is scaled horizontally, each instance should have access to the same rego files to avoid unpredictable results. If a file path has been configured but the file it is not present or accessible, the evaluation defaults to deny.

When using async post-processing which is done via the postprocessing service, the value `policies` must be added to the `POSTPROCESSING_STEPS` configuration in postprocessing service in the order where the evaluation should take place.

## Modules

### GRPC Service

This service can be used from any other internal service. It can also be used for example by third parties to find out if an action is allowed or not. This layer is already used by the proxy middleware.

### Event Service

This layer is event-based and part of the postprocessing service. Since processing at this point is asynchronous, the operations can also take longer and be more expensive, like evaluating the bytes of a file.

### Proxy Middleware

The [ocis proxy](../proxy) already includes such a middleware which uses the [GRPC service](#grpc-service) to evaluate the policies by using a configurable query. Since the Proxy is in heavy use and every request is processed here, only simple and quick decisions should be evaluated. More complex queries such as file evaluation are strongly discouraged.

## Example Policies

The policies service contains a set of pre-configured example policies. Those policies can be found in the [examples directory](https://github.com/owncloud/ocis/tree/master/deployments/examples/service_policies/policies). The contained policies disallows ocis to create certain filetypes, both for the proxy middleware and the events service.

To use the example policies, it's required to configure ocis to use these files which can be done by adding:

```yaml
policies:
  engine:
    policies:
      - YOUR_PATH/examples/policies/proxy.rego
      - YOUR_PATH/examples/policies/postprocessing.rego
      - YOUR_PATH/examples/policies/utils.rego
```

Once the policies are configured correctly, the _QUERY configuration needs to be defined for the proxy middleware and for the events service.

### Proxy

```yaml
proxy:
  policies_middleware:
    query: data.proxy.granted
```

The same can be achieved by setting the `PROXY_POLICIES_QUERY=data.proxy.granted` environment variable.

### ASYNC Postprocessing

```yaml
policies:
  postprocessing:
    query: data.postprocessing.granted
```

The same can be achieved by setting the `POLICIES_POSTPROCESSING_QUERY=data.postprocessing.granted` environment variable. As soon as that query is configured correctly, postprocessing must be informed to use the policies step by setting the environment variable `POSTPROCESSING_STEPS=policies`. Note that additional steps can be configured and their appearance defines the order of processing. For details see the postprocessing service documentation.
