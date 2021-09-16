Bugfix: Remove notifications placeholder

Since Reva was communicating its notification capabilities incorrectly, oCIS relied on a hardcoding string to overwrite them.
This has been fixed in [reva#1819](https://github.com/cs3org/reva/pull/1819) so we now removed the hardcoded string 
and don't modify Reva's notification capabilities anymore in order to fix clients having to poll a (non-existent) notifications endpoint.

https://github.com/owncloud/ocis/pull/2514
