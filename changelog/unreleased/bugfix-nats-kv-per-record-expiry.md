Bugfix: Honor per-record expiry in the nats-js-kv store

The nats-js-kv store ignored the per-record expiry written by callers and kept
every entry for the bucket-wide MaxAge, which is fixed at bucket creation. An
entry could therefore outlive the deadline stored in it and be served as stale
data — for the proxy userinfo cache this signed still-authenticated clients out.
Writes now stamp an absolute deadline on each record and reads treat a record
past its deadline as a cache miss, matching the memory store's behavior.

https://github.com/owncloud/ocis/pull/12863
