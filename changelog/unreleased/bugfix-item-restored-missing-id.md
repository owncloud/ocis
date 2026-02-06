Bugfix: Populate ID field in ItemRestored events

The ItemRestored event was missing the restored file's resource ID, making it impossible for event consumers to identify the restored file without additional API calls. This fix populates the ID field by extracting the file's opaque_id from the restore request key, mirroring how the ItemTrashed event correctly includes the trashed file's ID.

https://github.com/owncloud/ocis/pull/11991
