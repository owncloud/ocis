Bugfix: Omit invalid upload mtime instead of sending NaN

Web computed an upload's mtime by dividing a `File`'s `lastModified` by 1000,
unconditionally casting whatever `data` the caller supplied. A `Blob`-backed
upload - for example extracted entries from the Unzip web extension, which
hands Uppy `Blob`s rather than `File`s - has no `lastModified`, so the
division produced `NaN`. That `NaN` was sent to the server as the tus
`Upload-Metadata` `mtime` or the literal string `"NaN"` in the `x-oc-mtime`
header, which the upload coordinator now rejects with a 400, failing every
file in the batch. A shared helper now derives the mtime only from a finite,
strictly-positive `lastModified` and otherwise omits the field entirely, so
the server falls back to stamping the upload with its own current time.
Uploads of ordinary `File` objects with a valid `lastModified` are
unaffected.

https://github.com/owncloud/ocis/pull/13450
