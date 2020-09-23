Change: use /home as default namespace

Currently, cross storage etag propagation is not yet implemented, which prevents the desktop client from detecting changes via the PROPFIND to /. / is managed by the root storage provider which is independend of the home and oc storage providers. If a file changes in /home/foo, the etag change will only be propagated to the root of the home storage provider.

This change jails users into the `/home` namespace, and allows configuring the namespace to use for the two webdav endpoints using the new environment variable `WEBDAV_NAMESPACE_JAIL` which affects both endpoints `/dav/files` and `/webdav`.

This will allow us to focus on getting a single storage driver like eos or owncloud tested and better resembles what owncloud 10 does.

To get back the global namespace, which ultimately is the goal, just set the above environment variable to `/`.

<https://github.com/owncloud/ocis/ocis-revapull/68>
