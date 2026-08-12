Bugfix: Fix uploads stuck in processing with the s3ng storage driver

Uploads to a `s3ng` backed storage-users service never left the processing
state, so the file stayed unreadable and unmodifiable forever. Every client was
affected equally, because the defect sits below the protocol layer: WebDAV, the
desktop sync client and the Web UI alike.

Decomposedfs ignores any event whose storage id does not match its configured
mount id. The `s3ng` driver never received a `mount_id`, so that comparison ran
against an empty string and every `PostprocessingFinished` event was discarded.
The `s3ng` driver now passes `mount_id` the same way the `ocis` driver already
does.

The same key is also added to `OcisNoEvents` and `S3NGNoEvents`. Those build the
storage provider rather than the data provider and configure no events, so today
the value is unused there. Setting it keeps both constructions of the same
storage consistent about their own identity, rather than leaving one of them to
report an empty mount id to any future reader of that option.

https://github.com/owncloud/ocis/pull/12756
