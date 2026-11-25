Bugfix: Preserve pending activities until eventhistory expires them

Previously, listing activities for a resource deleted every stored
activity whose event ID was not yet returned by `eventhistory`. Because
`activitylog` and `eventhistory` consume NATS events independently, listing a
resource before `eventhistory` stored a newly emitted event permanently removed
the activity even after `eventhistory` caught up.

Missing events are now only pruned once their activity timestamp is older
than the configured store TTL (defaulting to `eventhistory`'s 336h retention),
leaving recently created activities and entries without a timestamp intact.

https://github.com/owncloud/ocis/issues/13490
