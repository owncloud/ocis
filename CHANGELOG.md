# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Fix #7: Build docker images with alpine:latest instead of alpine:edge
 * Chg #3: Initial release of basic version

## Details

 * Bugfix #7: Build docker images with alpine:latest instead of alpine:edge

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-graph-explorer/pull/7

 * Change #3: Initial release of basic version

   Just prepared an initial basic version to serve Graph-Explorer for the ownCloud Infinite
   Scale project. It just provides a minimal viable product to demonstrate the microservice
   pattern.

   https://github.com/owncloud/ocis-graph-explorer/issues/3


