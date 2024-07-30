Enhancement: Log user agent and remote addr on auth errors

The proxy will now log `user_agent`, `client.address`, `network.peer.address` and `network.peer.port` to help operations debug authentication errors. The latter three follow the [Semantic Conventions 1.26.0 / General / Attributes](https://opentelemetry.io/docs/specs/semconv/general/attributes/) naming to better integrate with log aggregation tools.

https://github.com/owncloud/ocis/pull/9475
