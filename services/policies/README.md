# Policies Service

The Policies Service provides a new grpc api which can be used to return whether a requested operation is allowed or not.
Open Policy Agent is used to determine the set of rules of what is permitted and what is not.

Policies are written in the [rego query language](https://www.openpolicyagent.org/docs/latest/policy-language/).

The Policies Service consists of the following modules

* Proxy Authorization (middleware)
* Event Authorization (async post-processing)
* GRPC API (can be used from other services too)

### GRPC service

This service can be used from any other internal service, it can also be used for example by third parties to find out if an action is allowed or not.
This layer is already used by the proxy middleware.

### Event service

This layer is event-based and part of the asynchronous post-processing, since processing at this point is asynchronous, the operations there can also take longer and be more expensive,
the bytes of a file can be examined here as an example.

### Proxy Middleware

The [ocis proxy](../proxy) already includes such a middleware which uses the [GRPC service](#grpc-service) to evaluate the policies by using a configurable query.
Since the Proxy is in heavy use and every request is processed here, only simple decisions that can be processed quickly are recommended here, more complex queries such as file evaluation should be avoided urgently.

## Example
The Policies Service contains a set of pre-configured example policies, thos policies can be found in the [examples directory](./examples).
The contained policies disallows ocis to create certain filetypes, both for the proxy middleware and the events service.

To use the example policies, it's required to configure ocis to use these files which can be done by adding:

```yaml
policies:
  engine:
    policies:
      - YOUR_PATH/examples/policies/proxy.rego
      - YOUR_PATH/examples/policies/postprocessing.rego
      - YOUR_PATH/examples/policies/utils.rego
```
Once the policies are configured correctly we need to set up the correct queries for the proxy middleware and for the events service.

### Proxy

```yaml
proxy:
  policies_middleware:
    query: data.proxy.granted
```

the same can be achieved by setting the `PROXY_POLICIES_QUERY=data.proxy.granted` environment variable.

### ASYNC postprocessing

```yaml
policies:
  postprocessing:
    query: data.postprocessing.granted
```

the same can be achieved by setting the `POLICIES_POSTPROCESSING_QUERY=data.postprocessing.granted` environment variable.
As soon as that query is configured correctly, postprocessing must be informed to use the policies step too by setting the environment variable `POSTPROCESSING_STEPS=policies`.






