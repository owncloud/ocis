Bugfix: Preserve pending activities until eventhistory expires them

Previously, listing activities for a resource deleted every stored
activity whose event ID was not yet returned by `eventhistory`. Because
`activitylog` and `eventhistory` consume NATS events independently, listing a
resource before `eventhistory` stored a newly emitted event permanently removed
the activity even after `eventhistory` caught up.

Missing events are now kept until `ACTIVITYLOG_STORE_TTL` (defaulting to 336h)
has elapsed since `activitylog` received the activity, while legacy entries
without a recorded receive time continue to be pruned on the first miss.
`ACTIVITYLOG_STORE_TTL` now also governs this pruning and should not be set
lower than `EVENTHISTORY_STORE_TTL`.

https://github.com/owncloud/ocis/issues/13490
