Bugfix: Reduce activitylog load to prevent dropped events under bursty traffic

Under bursty event traffic (bulk uploads, POSIX inotify churn) the single
activitylog consumer could not keep up. For every event it walked the resource's
parent chain, issuing a gateway stat per level and a full JSON read-modify-write
of the (up to 6000 entry) activity list at every level. The consumer fell behind,
the NATS push-subscription buffer overflowed and silently dropped messages
("nats: slow consumer ... main-queue"), and the unacknowledged messages were
redelivered, pinning CPU and growing the JetStream store.

The per-event cost is now reduced: activity writes are coalesced per resource
over a configurable window (`ACTIVITYLOG_WRITE_BUFFER_DURATION`, default 10s;
set to 0 to write synchronously) and flushed in a single read-modify-write,
which removes most of the marshal/unmarshal cycles on hot parent nodes that were
the main cause of the high CPU load; and activity lists are stored with msgpack
instead of JSON, with a JSON read fallback so existing records stay readable.

https://github.com/owncloud/ocis/issues/10825
https://github.com/owncloud/ocis/pull/12417
