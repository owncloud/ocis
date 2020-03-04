# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Enh #33: Allow http services to register handlers

## Details

 * Enhancement #33: Allow http services to register handlers

   Added a handler option on http services

   https://github.com/owncloud/ocis-pkg/pull/33


# Changelog for 2.0.0

The following sections list the changes for 2.0.0.

## Summary

 * Fix #25: Fix Module Path
 * Fix #27: Change import paths to ocis-pkg/v2
 * Chg #22: Upgrade the micro libraries

## Details

 * Bugfix #25: Fix Module Path

   The module version must be in the path. See
   https://github.com/golang/go/wiki/Modules#semantic-import-versioning for more
   information. > If the module is version v2 or higher, the major version of the module must be
   included as a /vN at the end of the module paths used in go.mod files (e.g., module
   github.com/my/mod/v2, require github.com/my/mod/v2 v2.0.1) and in the package import path
   (e.g., import "github.com/my/mod/v2/mypkg"). This includes the paths used in go get
   commands (e.g., go get github.com/my/mod/v2@v2.0.1. Note there is both a /v2 and a @v2.0.1 in
   that example. One way to think about it is that the module name now includes the /v2, so include
   /v2 whenever you are using the module name).

   https://github.com/owncloud/ocis-pkg/pull/25

 * Bugfix #27: Change import paths to ocis-pkg/v2

   Changed the import paths to the current version

   https://github.com/owncloud/ocis-pkg/pull/27

 * Change #22: Upgrade the micro libraries

   Upgraded the go-micro libraries to v2.

   https://github.com/owncloud/ocis-pkg/pull/22


# Changelog for 2.0.0

The following sections list the changes for 2.0.0.

## Summary

 * Fix #25: Fix Module Path
 * Fix #27: Change import paths to ocis-pkg/v2
 * Chg #22: Upgrade the micro libraries

## Details

 * Bugfix #25: Fix Module Path

   The module version must be in the path. See
   https://github.com/golang/go/wiki/Modules#semantic-import-versioning for more
   information. > If the module is version v2 or higher, the major version of the module must be
   included as a /vN at the end of the module paths used in go.mod files (e.g., module
   github.com/my/mod/v2, require github.com/my/mod/v2 v2.0.1) and in the package import path
   (e.g., import "github.com/my/mod/v2/mypkg"). This includes the paths used in go get
   commands (e.g., go get github.com/my/mod/v2@v2.0.1. Note there is both a /v2 and a @v2.0.1 in
   that example. One way to think about it is that the module name now includes the /v2, so include
   /v2 whenever you are using the module name).

   https://github.com/owncloud/ocis-pkg/pull/25

 * Bugfix #27: Change import paths to ocis-pkg/v2

   Changed the import paths to the current version

   https://github.com/owncloud/ocis-pkg/pull/27

 * Change #22: Upgrade the micro libraries

   Upgraded the go-micro libraries to v2.

   https://github.com/owncloud/ocis-pkg/pull/22


# Changelog for 1.3.0

The following sections list the changes for 1.3.0.

## Summary

 * Fix #14: Fix serving static assets
 * Chg #19: Add TLS support for http services
 * Enh #8: Introduce OpenID Connect middleware

## Details

 * Bugfix #14: Fix serving static assets

   Ocis-hello used "/" as root. adding another / caused the static middleware to always fail
   stripping that prefix. All requests will return 404. Setting root to `""` in the `ocis-hello`
   flag does not work because Chi reports that routes need to start with `/`.
   `path.Clean(root+"/")` would yield `""` for `root="/"`.

   https://github.com/owncloud/ocis-pkg/pull/14

 * Change #19: Add TLS support for http services

   `ocis-pkg` http services support TLS. The idea behind is setting the issuer on phoenix's
   `config.json` to `https`. Or in other words, use https to access the Kopano extension, and
   authenticate using an SSL certificate.

   https://github.com/owncloud/ocis-pkg/issues/19

 * Enhancement #8: Introduce OpenID Connect middleware

   Added an openid connect middleware that will try to authenticate users using OpenID Connect.
   The claims will be added to the context of the request.

   https://github.com/owncloud/ocis-pkg/issues/8


# Changelog for 1.2.0

The following sections list the changes for 1.2.0.

## Summary

 * Chg #9: Add root path to static middleware

## Details

 * Change #9: Add root path to static middleware

   Currently the `Static` middleware always serves from the root path, but all our HTTP handlers
   accept a custom root path which also got to be applied to the static file handling.

   https://github.com/owncloud/ocis-pkg/issues/9


# Changelog for 1.1.0

The following sections list the changes for 1.1.0.

## Summary

 * Chg #2: Better log level handling within micro

## Details

 * Change #2: Better log level handling within micro

   Currently every log message from the micro internals are logged with the info level, we really
   need to respect the proper defined log level within our log wrapper package.

   https://github.com/owncloud/ocis-pkg/issues/2


# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #1: Initial release of basic version

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version to have some shared functionality published which can
   be used by all other ownCloud Infinite Scale extensions.

   https://github.com/owncloud/ocis-pkg/issues/1


