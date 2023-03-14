Enhancement: Introduce policies-service

Introduces policies service. The policies-service provides a new grpc api which can be used to return whether a requested operation is allowed or not.
Open Policy Agent is used to determine the set of rules of what is permitted and what is not.

2 further levels of authorization build on this:

* Proxy Authorization
* Event Authorization (needs async post-processing enabled)

The simplest authorization layer is in the proxy, since every request is processed here, only simple decisions that can be processed quickly are made here, more complex queries such as file evaluation are explicitly excluded in this layer.

The next layer is event-based as a pipeline step in asynchronous post-processing, since processing at this point is asynchronous, the operations there can also take longer and be more expensive,
the bytes of a file can be examined here as an example.

Since the base block is a grpc api, it is also possible to use it directly.
The policies are written in the [rego query language](https://www.openpolicyagent.org/docs/latest/policy-language/).

https://github.com/owncloud/ocis/pull/5714
https://github.com/owncloud/ocis/issues/5580
