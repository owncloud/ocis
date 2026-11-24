Bugfix: Skip non-finite photo facet values

Listing spaces via `/me/drives` could fail with an HTTP 500 when a file, such as
a space image, had a photo facet holding a non-finite value like
`exposureDenominator="+Inf"`. Tika produced it for files whose EXIF exposure
time is `0`, and the value could not be encoded as JSON.

Non-finite values are no longer produced from the exposure time, no longer
stored as file metadata, and are ignored when reading metadata that is already
stored.

https://github.com/owncloud/ocis/pull/13571
