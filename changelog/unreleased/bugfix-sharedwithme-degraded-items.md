Bugfix: Keep shares visible in sharedWithMe when a resource cannot be statted

The graph `sharedWithMe` handler resolves each received share by fanning out a per-resource
`Stat` call with a concurrency limit, all sharing the request context. When a `Stat` failed
for any reason (slow or stuck downstream, deleted space, gateway error, deadline exceeded)
the worker logged at debug and returned without emitting a drive item, so the share was
silently omitted from the response and the handler still returned `200 OK` with a partial,
non-deterministic list. A single chronically-slow share could therefore make other,
recently-accepted shares intermittently invisible, and repeated calls returned different
subsets of the user's shares.

The handler now:

- bounds each per-share `Stat` with its own timeout derived from the request context, so one
  slow resource can no longer consume the deadline shared by all the other shares;
- returns a degraded drive item built from the data already present in the share record
  (ids, permissions, grantees, timestamps, mountpoint name) when the resource cannot be statted
  due to a transient or indeterminate failure (timeout, slow or unavailable downstream), so the
  share stays visible instead of intermittently disappearing. A genuinely missing resource or
  revoked access (for example after the sharer was deleted) still drops the share, as before;
- logs the dropped/degraded shares at warning level with a per-request count for operator
  visibility.

https://github.com/owncloud/ocis/pull/12430
