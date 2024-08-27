Enhancement: We now set the configured protocol transport for service metadata

This allows configuring services to listan on `tcp` or `unix` sockets and clients to use the `dns`, `kubernetes` or `unix` protocol URIs instead of service names.

https://github.com/owncloud/ocis/pull/9490
https://github.com/cs3org/reva/pull/4744