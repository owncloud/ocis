Bugfix: Cleanup separated indices in memory

The accounts service was creating a bleve index instance in the service handler, thus creating separate in memory indices for the http and grpc servers. We moved the service handler creation out of the server creation so that the service handler, thus also the bleve index, is a shared instance of the servers.

This fixes a bug that accounts created through the web ui were not able to sign in until a service restart.

https://github.com/owncloud/product/issues/224
https://github.com/owncloud/ocis-accounts/pull/117
https://github.com/owncloud/ocis-accounts/pull/118
https://github.com/owncloud/ocis/pull/555
