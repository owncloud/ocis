Bugfix: Reject space images which are too large to be rendered

Setting a space image larger than the thumbnails service accepts appeared to
succeed but never changed the space image. Uploading the file and assigning it
as the space image are two separate requests, and both completed successfully,
so clients reported success. Only a later preview request, which the assigning
client never observes, was rejected because the image exceeded
`THUMBNAILS_MAX_INPUT_IMAGE_FILE_SIZE`. As space images are rendered
exclusively as thumbnails, the image was not merely degraded, it was invisible
with no error anywhere in the UI.

The graph service now checks the file size when a space image is assigned and
rejects the request with `413 Request Entity Too Large` instead of storing a
reference which cannot be displayed. The previous space image is left
untouched. The limit is read from the new `GRAPH_MAX_IMAGE_FILE_SIZE` setting,
which defaults to `50MB` and must be less than or equal to
`THUMBNAILS_MAX_INPUT_IMAGE_FILE_SIZE`. Space descriptions (`readme`) are not
affected.

https://github.com/owncloud/ocis/pull/12701
