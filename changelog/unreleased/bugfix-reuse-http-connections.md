Bugfix: Reuse HTTP connections in the proxy and in the metadata client

Under high load, oCIS opened and closed many TCP connections each second.
The host ran out of free ports. This caused failed logins and failed
requests.

The proxy now keeps up to 100 idle connections for each backend. Before, it
kept 2. The proxy has no global limit for idle connections. The OIDC client of
the proxy now keeps its connections alive. Before, it opened a new connection
with a new TLS handshake for each call.

The metadata client and the HTTP clients of reva use the same settings.

The proxy and the reva clients now close an idle connection after 30 seconds.
Before, the proxy used 90 seconds. A load balancer or an ingress controller
often closes an idle connection after 60 seconds. If the client closes first,
it never sends a request on a connection that the other side already closed.
The transport does not repeat such a request for methods like PUT, POST or
PROPFIND.

https://github.com/owncloud/ocis/pull/13541
https://github.com/owncloud/reva/pull/825
