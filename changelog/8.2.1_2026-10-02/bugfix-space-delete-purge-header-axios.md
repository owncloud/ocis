Bugfix: Accept an X-Purge alias for the space-delete Purge header

axios >=1.20.0 reserves "purge" as an internal per-HTTP-method header-bucket
name and strips any outgoing request header matching it case-insensitively.
The web client's permanent space-delete request set a header literally named
Purge, so it was silently dropped, and the server treated the request as a
disable instead of a permanent delete.

The graph service now also accepts an X-Purge header, and the web client
sends that instead. The original Purge header is still accepted for other
API clients.

https://github.com/owncloud/ocis/pull/13049
