# Changes in unreleased

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

