Bugfix: Deduplicate CreateHome calls in proxy middleware

The CreateHome proxy middleware previously fired a CreateHome gRPC
request on every authenticated HTTP request with no deduplication.
On first login, the browser sends many parallel requests, each
triggering a redundant CreateHome call. This change uses singleflight
to collapse concurrent calls for the same user and caches successful
results in-process so subsequent requests skip the gRPC call entirely.

https://github.com/owncloud/ocis/pull/12115
https://github.com/owncloud/ocis/issues/12068
