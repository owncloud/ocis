Enhancement: Improve thumbnails API

Changed the thumbnails API to no longer transfer images via GRPC.
GRPC has a limited message size and isn't very efficient with large binary data.
The new API transports the images over HTTP.

https://github.com/owncloud/ocis/pull/3272
