# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-thumbnails unreleased.

[unreleased]: https://github.com/owncloud/ocis-thumbnails/compare/v0.1.0...master

## Summary

* Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)

## Details

* Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)

   The program flags where not correctly recognized because we didn't pass them to the micro
   framework when initializing a grpc service.

   https://github.com/owncloud/ocis-thumbnails/issues/15

# Changelog for [0.1.0] (2020-03-31)

The following sections list the changes in ocis-thumbnails 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-thumbnails/compare/c43f3a33cb0b57d7e25ebc88c138d22e95f88cfe...v0.1.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)
* Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)
* Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)

   Just prepare an initial basic version to embed a thumbnailer into our microservice
   infrastructure in the scope of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis-thumbnails/issues/1


* Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)

   We implemented predefined resolutions to prevent attacker from flooding the service with a
   large number of thumbnails. The requested resolution gets mapped to the closest matching
   predefined resolution.

   https://github.com/owncloud/ocis-thumbnails/issues/7


* Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

   We implemented the first simple version. It can load images via webdav and store them locally in
   the filesystem.

   https://github.com/owncloud/ocis-thumbnails/pull/3

