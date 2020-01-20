# Changelog for 1.3.0

The following sections list the changes for 1.3.0.

## Summary

 * Fix #14: Fix serving static assets
 * Chg #19: Add TLS support for http services
 * Enh #8: Introduce OpenID Connect middleware

## Details

 * Bugfix #14: Fix serving static assets

   Ocis-hello used "/" as root. adding another / caused the static middleware to always fail
   stripping that prefix. All requests will return 404. setting root to "" in the `ocis-hello`
   flag does not work because chi reports that routes need to start with a /.
   `path.Clean(root+"/")` would yield "" for root="/"

   https://github.com/owncloud/ocis-pkg/pull/14

 * Change #19: Add TLS support for http services

   `ocis-pkg` http services support TLS.

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


