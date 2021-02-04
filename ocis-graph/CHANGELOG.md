# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Fix #32: Build docker images with alpine:latest instead of alpine:edge
 * Chg #7: Initial release of basic version
 * Enh #27: Configuration

## Details

 * Bugfix #32: Build docker images with alpine:latest instead of alpine:edge

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-graph/pull/32

 * Change #7: Initial release of basic version

   Just prepare an initial basic version to serve a graph world API that can be used by Phoenix or
   other extensions.

   https://github.com/owncloud/ocis-graph/issues/7

 * Enhancement #27: Configuration

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-graph/pull/27


