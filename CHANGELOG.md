# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-reva unreleased.

[unreleased]: https://github.com/owncloud/ocis-reva/compare/v0.1.1...master

## Summary

* Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)
* Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

## Details

* Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   We have added missing config options for the user sharing manager and added a dedicated eos
   storage command with pre configured settings for the eos-docker container. It configures a
   `Shares` folder in a users home when using eos as the storage driver.

   https://github.com/owncloud/ocis-reva/pull/127


* Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   Adds initial public sharing and ocm implementation.

   https://github.com/owncloud/ocis-reva/pull/127

# Changelog for [0.1.1] (2020-03-31)

The following sections list the changes in ocis-reva 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis-reva/compare/v0.1.0...v0.1.1

## Summary

* Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

## Details

* Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

   We have added missing config options for the home layout to the config struct that is passed to
   eos.

   https://github.com/owncloud/ocis-reva/pull/125

# Changelog for [0.1.0] (2020-03-23)

The following sections list the changes in ocis-reva 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-reva/compare/6702be7f9045a382d40691a9bcd04f572203e9ed...v0.1.0

## Summary

* Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)
* Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)
* Bugfix - Allow different namespaces for /webdav and /dav/files: [#68](https://github.com/owncloud/ocis-reva/pull/68)
* Change - Use /home as default namespace: [#68](https://github.com/owncloud/ocis-reva/pull/68)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)
* Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)
* Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)
* Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)
* Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)
* Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)
* Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

## Details

* Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)

   While upgrading to the micro/cli version 2 there where two instances of `StringFlag` which had
   not been changed to `StringSliceFlag`. This caused `ocis-reva users` and `ocis-reva
   storage-root` to fail on startup.

   https://github.com/owncloud/ocis-reva/issues/75
   https://github.com/owncloud/ocis-reva/pull/76


* Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   It was misspelled as `REVA_LDAP_SCHEMA_Mail`.

   https://github.com/owncloud/ocis-reva/pull/113


* Bugfix - Allow different namespaces for /webdav and /dav/files: [#68](https://github.com/owncloud/ocis-reva/pull/68)

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


* Change - Use /home as default namespace: [#68](https://github.com/owncloud/ocis-reva/pull/68)

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


* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)

   Just prepared an initial basic version to start a reva server and start integrating with the
   go-micro base dextension framework of ownCloud Infinite Scale.

   https://github.com/owncloud/ocis-reva/issues/1


* Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)

   The initial version would only allow us to use a set of reva configurations to start multiple
   services. We use a more opinionated set of commands to start dedicated services that allows us
   to configure them individually. It allows us to switch eg. the user backend to LDAP and fully use
   it on the cli.

   https://github.com/owncloud/ocis-reva/issues/6


* Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)

   The flags that let reva storage providers announce that they expose a data server now defaults
   to true:

   `REVA_STORAGE_HOME_EXPOSE_DATA_SERVER=1` `REVA_STORAGE_OC_EXPOSE_DATA_SERVER=1`

   https://github.com/owncloud/ocis-reva/issues/89


* Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   We changed the default configuration to integrate better with ocis.

   - We use ocis-glauth as the default ldap server on port 9125 with base `dc=example,dc=org`. - We
   use a dedicated technical `reva` user to make ldap binds - Clients are supposed to use the
   ocis-proxy endpoint `https://localhost:9200` - We removed unneeded ocis configuration
   from the frontend which no longer serves an oidc provider. - We changed the default user
   OpaqueID attribute from `sub` to `preferred_username`. The latter is a claim populated by
   konnectd that can also be used by the reva ldap user manager to look up users by their OpaqueId

   https://github.com/owncloud/ocis-reva/pull/113


* Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)

   Three new flags are now available:

   - scan files on startup to generate missing fileids default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_SCAN` cli option: `--storage-owncloud-scan`

   - autocreate home path for new users default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_AUTOCREATE` cli option: `--storage-owncloud-autocreate`

   - the address of the redis server default: `:6379` env var:
   `REVA_STORAGE_OWNCLOUD_REDIS_ADDR` cli option: `--storage-owncloud-redis`

   https://github.com/owncloud/ocis-reva/issues/87


* Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)

   Reva has seen a lot of changes that allow us to - reduce the configuration overhead - use the
   autocreato home folder option - use the home folder path layout option - no longer start the root
   storage

   https://github.com/owncloud/ocis-reva/pull/91
   related:


* Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

   We now default to `json` which persists shares in the sharing manager in a json file instead of an
   in memory db.

   https://github.com/owncloud/ocis-reva/pull/115

