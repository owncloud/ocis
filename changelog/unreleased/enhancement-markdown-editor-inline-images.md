Enhancement: Inline images in the markdown editor

Uploading an image in the markdown editor did nothing at all - neither
drag-and-drop nor the "Upload Images" toolbar button produced a result, a
progress indicator, or an error. The feature had been deliberately suppressed
because there was no answer to where an uploaded image should be stored, and
because a recipient of a public link to a single markdown file has no
permission to read an image stored next to it.

Images dropped or picked in the markdown editor are now encoded into the
document itself as base64 data URIs, so a shared document carries its own
images and needs no further permissions. Because the image bytes are now part
of the document, two limits keep documents small enough to stay fully
searchable: a maximum size per image, and a budget for the total size of all
images in one document. An image that exceeds either limit is rejected with a
notification naming the image and the limit, instead of failing silently.

Both limits are configurable through the text editor's app configuration as
`maxImageSize` and `maxDocumentImageSize`, given in bytes and defaulting to
2 MB and 10 MB respectively.

The editor's "Crop And Upload" entry is no longer offered. Cropping was not part
of the reported problem, so it stays out until users ask for it, and images are
inlined exactly as they were picked.

https://github.com/owncloud/ocis/pull/12967
https://github.com/owncloud/web/issues/12407
