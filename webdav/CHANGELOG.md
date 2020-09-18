# Changelog for 0.1.1

The following sections list the changes for 0.1.1.

## Summary

 * Fix #22: Build docker images with alpine:latest instead of alpine:edge
 * Chg #20: Change status not found on missing thumbnail

## Details

 * Bugfix #22: Build docker images with alpine:latest instead of alpine:edge

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-webdav/pull/22

 * Change #20: Change status not found on missing thumbnail

   The service returned a bad request when no thumbnail was generated. It is now changed to not
   found.

   https://github.com/owncloud/ocis-webdav/issues/20
   https://github.com/owncloud/ocis-webdav/pull/21


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #1: Initial release of basic version
 * Chg #16: Update ocis-pkg to version 2.2.0
 * Enh #14: Configuration
 * Enh #13: Implement preview API

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version to serve webdav for the ownCloud Infinite Scale
   project. It just provides a minimal viable product to demonstrate the microservice pattern.

   https://github.com/owncloud/ocis-webdav/issues/1

 * Change #16: Update ocis-pkg to version 2.2.0

   Updated ocis-pkg to include the cors header changes.

   https://github.com/owncloud/ocis-webdav/issues/16

 * Enhancement #14: Configuration

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-webdav/pull/14

 * Enhancement #13: Implement preview API

   Added the API endpoint for file previews.

   https://github.com/owncloud/ocis-webdav/pull/13


