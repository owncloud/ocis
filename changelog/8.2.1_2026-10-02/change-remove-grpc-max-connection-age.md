Change: Replace GRPC_MAX_CONNECTION_AGE with client keepalive

The grpc clients now send a keepalive ping while a request is in flight and fail
the requests on a connection whose peer stops answering, instead of waiting for
as long as the caller allows. This covers both the reva CS3 clients (gateway,
storage-users, storage-shares, ...) and the go-micro based clients used for
inter-service calls between the other oCIS services. Set GRPC_CLIENT_KEEPALIVE_TIME
and GRPC_CLIENT_KEEPALIVE_TIMEOUT to enable and tune this; leave them unset to keep
grpc's own default (no pings).

GRPC_MAX_CONNECTION_AGE has been removed. It only closed healthy connections on
a timer, never ended a request that was already in flight, and silently did
nothing when its value had no unit suffix.

https://github.com/owncloud/ocis/pull/13020
