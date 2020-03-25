# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Chg #1: Initial release of basic version
 * Enh #14: Configuration
 * Enh #13: Implement preview API

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version to serve webdav for the ownCloud Infinite Scale
   project. It just provides a minimal viable product to demonstrate the microservice pattern.

   https://github.com/owncloud/ocis-webdav/issues/1

 * Enhancement #14: Configuration

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-webdav/pull/14

 * Enhancement #13: Implement preview API

   Added the API endpoint for file previews.

   https://github.com/owncloud/ocis-webdav/pull/13


