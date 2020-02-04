Bugfix: Allow different namespaces for /webdav and /dav/files

After fbf131c the path for the "new" webdav path does not contain a username `/remote.php/dav/files/textfile0.txt`. It used to be `/remote.php/dav/files/oc/einstein/textfile0.txt` So it lost `oc/einstein`.

This PR allows setting up different namespaces for `/webav` and `/dav/files`:

`/webdav` is jailed into `/home` - which uses the home storage driver and uses the logged in user to construct the path
`/dav/files` is jailed into `/oc` - which uses the owncloud storage driver and expects a username as the first path segment

This mimics oc10

The `WEBDAV_NAMESPACE_JAIL` environment variable is split into
- `WEBDAV_NAMESPACE` and
- `DAV_FILES_NAMESPACE` accordingly.

Related: https://github.com/owncloud/ocis-reva/pull/68
