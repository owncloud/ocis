Bugfix: Fix the proxy readiness check for NATS

The proxy readiness check passed the events cluster ID instead of the events
endpoint to the NATS reachability check. NATS then tried to resolve the cluster
ID (e.g. `ocis-cluster`) as a host name, which always failed, so the proxy
`/readyz` endpoint reported the service as not ready even when NATS was
reachable. The check now uses the events endpoint, consistent with every other
service.

https://github.com/owncloud/ocis/issues/10661
https://github.com/owncloud/ocis/pull/12421
