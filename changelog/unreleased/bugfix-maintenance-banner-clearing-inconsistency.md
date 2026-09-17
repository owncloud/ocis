Bugfix: Unify when the maintenance banner clears

The Web frontend showed the maintenance banner inconsistently depending on which
HTTP client happened to make the last request: the webdav client cleared the banner
on any error that was not a maintenance response, including a plain 404 or 500,
while the other client only cleared it on a genuine successful response. Whether the
banner disappeared during real maintenance could therefore depend on which request
happened to run last.

Both clients now clear the maintenance banner only on an actual successful response,
so an unrelated error no longer hides the banner while the server is still in
maintenance.

https://github.com/owncloud/ocis/pull/12966
