Enhancement: Add a --force flag to the search index command

`ocis search index` skipped every file that was already indexed and had not
changed since. There was no way to index such files again, for example after
enabling OCR in Tika, other than deleting the whole search index, which makes
search unavailable until it has been rebuilt.

The command now has a `--force` flag. With it, every file of the given space,
or of all spaces, is extracted and indexed again, while the existing index
entries stay searchable. Without the flag the behaviour is unchanged.

https://github.com/owncloud/ocis/issues/9106
