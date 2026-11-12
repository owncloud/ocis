Security: Do not leak internal error details in webdav error responses

The webdav service returned raw internal error messages to the client when a server-side error occurred. For example, a thumbnail request with a null byte in the filename and `?preview=1` could trigger a `500 Internal Server Error` whose body exposed the internal storage filesystem path (information disclosure).

Server-side error handlers now return a generic, user-relevant message to the client, while the detailed error is only logged server-side.

https://github.com/owncloud/ocis/pull/12398
