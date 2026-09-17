Bugfix: Log unmapped thumbnail errors and always report a sabredav exception

The webdav service logged failures of the thumbnails service at debug level only.
Errors which are not a property of the requested file, such as the thumbnails
service being unreachable, were therefore invisible at production log levels: a
complete preview outage produced HTTP 500 responses without a single log line
explaining them. Unmapped errors are now logged at error level, while expected
per-file outcomes such as an unsupported file type or a file still being
processed stay at debug level. Two of the four thumbnail handlers also logged
without the request context, so their messages carried no request id.

In addition, `codesEnum` only mapped four status codes, so error responses for
all other codes were rendered with an empty `<s:exception></s:exception>`
element. The missing entries for 403, 425 and 429 have been added and any
remaining unmapped code now falls back to a generic exception name, so clients
always receive a usable exception.

https://github.com/owncloud/ocis/pull/12663
