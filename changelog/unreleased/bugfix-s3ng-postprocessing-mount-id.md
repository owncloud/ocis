Bugfix: Fix uploads stuck in processing with the s3ng storage driver

Uploads to a `s3ng` backed storage-users service never left the processing
state, so the file stayed unreadable and unmodifiable forever. Restoring a file
version was affected in the same way.

Decomposedfs ignores any event whose storage id does not match its configured
mount id. The `s3ng` driver never received a `mount_id`, so that comparison was
made against an empty string and every `PostprocessingFinished`,
`RestartPostprocessing` and `RevertRevision` event was discarded. The `s3ng`
driver now passes `mount_id` the same way the `ocis` driver already does.

https://github.com/owncloud/ocis/pull/12756
