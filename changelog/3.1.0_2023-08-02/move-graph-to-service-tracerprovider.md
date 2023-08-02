Enhancement: Move graph to service tracerprovider

This moves the graph to initialise a service tracer provider at service initialisation time,
instead of using a package global tracer provider.

https://github.com/owncloud/ocis/pull/6695
