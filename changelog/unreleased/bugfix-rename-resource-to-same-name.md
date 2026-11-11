Bugfix: Renaming a resource to its current name is a silent no-op

Renaming a file or folder to the name it already has failed with a
confusing `403 Forbidden` or `404 Not Found`, and `409 Conflict` when the
request addressed the resource by its id. Every client had to work around
it on its own. A `MOVE` whose source and destination resolve to the same
resource is now detected server side and answered with `204 No Content`
without touching the resource, so the behaviour is consistent for all
clients and for both dav routes. The no-op still requires the move
permission: a user who may not rename the resource keeps getting
`403 Forbidden`.

https://github.com/owncloud/ocis/pull/12784
https://github.com/owncloud/reva/pull/708
https://github.com/owncloud/reva/pull/713
