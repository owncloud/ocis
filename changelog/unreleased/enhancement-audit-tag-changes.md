Enhancement: Record tag changes in the audit log

Adding or removing tags on a file or folder left no trace in the audit log.
The graph service already publishes `TagsAdded` and `TagsRemoved` events for
every change, but only the search service listened to them.

The audit service now records both: `file_tags_add` and `file_tags_remove`
entries carry the acting user, the file id and path, the space owner and the
tags concerned. The graph service also sets the timestamp on the two events,
which it left empty before.

https://github.com/owncloud/ocis/issues/13054
