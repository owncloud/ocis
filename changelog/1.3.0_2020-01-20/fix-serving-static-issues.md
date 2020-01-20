Bugfix: Fix serving static assets

ocis-hello used "/" as root. adding another / caused the static middleware to always fail stripping that prefix. All requests will return 404.
setting root to "" in the `ocis-hello` flag does not work because chi reports that routes need to start with a /. `path.Clean(root+"/")` would yield "" for root="/"

https://github.com/owncloud/ocis-pkg/pull/14