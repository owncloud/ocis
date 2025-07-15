Change: Split spaces webdav url and graph url in base and path

We've fixed the behavior for the spaces webdav url and graph explorer graph url settings, so that they respect the environment variable `OCIS_URL`. Previously oCIS admins needed to set these URLs manually to make spaces and the graph explorer work.

https://github.com/owncloud/ocis/pull/2660
https://github.com/owncloud/ocis/issues/2659
