Enhancement: Create a root span on proxy that propagates down to consumers

In order to propagate and correctly associate a span with a request we need a root span that gets sent to other services.

https://github.com/owncloud/ocis/proxy/pull/64
