Enhancement: Allow configuring grpc max connection age

We added a GRPC_MAX_CONNECTION_AGE env var that allows limiting the lifespan of connections. A closed connection triggers grpc clients to do a new DNS lookup to pick up new IPs.

https://github.com/owncloud/ocis/pull/9657
