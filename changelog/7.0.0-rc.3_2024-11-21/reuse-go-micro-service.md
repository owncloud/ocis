Bugfix: Reuse go-micro service clients

go micro clients must not be reinitialized. The internal selector will spawn a new go routine to watch for registry changes.

https://github.com/owncloud/ocis/pull/10582
