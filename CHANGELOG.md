# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Fix #75: Set correct flag type in the flagsets
 * Fix #68: Allow different namespaces for /webdav and /dav/files
 * Chg #68: Use /home as default namespace
 * Chg #1: Initial release of basic version
 * Chg #6: Start multiple services with dedicated commands
 * Enh #87: Expose owncloud storage driver config in flagset

## Details

 * Bugfix #75: Set correct flag type in the flagsets

   While upgrading to the micro/cli version 2 there where two instances of `StringFlag` which had
   not been changed to `StringSliceFlag`. This caused `ocis-reva users` and `ocis-reva
   storage-root` to fail on startup.

   https://github.com/owncloud/ocis-reva/issues/75
   https://github.com/owncloud/ocis-reva/pull/76

 * Bugfix #68: Allow different namespaces for /webdav and /dav/files

   After fbf131c the path for the "new" webdav path does not contain a username
   `/remote.php/dav/files/textfile0.txt`. It used to be
   `/remote.php/dav/files/oc/einstein/textfile0.txt` So it lost `oc/einstein`.

   This PR allows setting up different namespaces for `/webav` and `/dav/files`:

   `/webdav` is jailed into `/home` - which uses the home storage driver and uses the logged in user
   to construct the path `/dav/files` is jailed into `/oc` - which uses the owncloud storage
   driver and expects a username as the first path segment

   This mimics oc10

   The `WEBDAV_NAMESPACE_JAIL` environment variable is split into - `WEBDAV_NAMESPACE` and -
   `DAV_FILES_NAMESPACE` accordingly.

   https://github.com/owncloud/ocis-reva/pull/68
   related:

 * Change #68: Use /home as default namespace

   Currently, cross storage etag propagation is not yet implemented, which prevents the desktop
   client from detecting changes via the PROPFIND to /. / is managed by the root storage provider
   which is independend of the home and oc storage providers. If a file changes in /home/foo, the
   etag change will only be propagated to the root of the home storage provider.

   This change jails users into the `/home` namespace, and allows configuring the namespace to
   use for the two webdav endpoints using the new environment variable `WEBDAV_NAMESPACE_JAIL`
   which affects both endpoints `/dav/files` and `/webdav`.

   This will allow us to focus on getting a single storage driver like eos or owncloud tested and
   better resembles what owncloud 10 does.

   To get back the global namespace, which ultimately is the goal, just set the above environment
   variable to `/`.

   https://github.com/owncloud/ocis-reva/pull/68

 * Change #1: Initial release of basic version

   Just prepared an initial basic version to start a reva server and start integrating with the
   go-micro base dextension framework of ownCloud Infinite Scale.

   https://github.com/owncloud/ocis-reva/issues/1

 * Change #6: Start multiple services with dedicated commands

   The initial version would only allow us to use a set of reva configurations to start multiple
   services. We use a more opinionated set of commands to start dedicated services that allows us
   to configure them individually. It allowcs us to switch eg. the user backend to LDAP and fully it
   on the cli.

   https://github.com/owncloud/ocis-reva/issues/6

 * Enhancement #87: Expose owncloud storage driver config in flagset

   Three new flags are now available:

   - scan files on startup to generate missing fileids default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_SCAN` cli option: `--storage-owncloud-scan`

   - autocreate home path for new users default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_AUTOCREATE` cli option: `--storage-owncloud-autocreate`

   - the address of the redis server default: `:6379` env var:
   `REVA_STORAGE_OWNCLOUD_REDIS_ADDR` cli option: `--storage-owncloud-redis`

   https://github.com/owncloud/ocis-reva/issues/87


