Enhancement: Allow to use libvips for generating thumbnails

To improve performance (and to be able to support a wider range of images formats in the future)
the thumbnails service is now able to utilize libvips (https://www.libvips.org/) for generating thumbnails.
Enabling the use of libvips is implemented as a build-time option which is currently disabled for the
"bare-metal" build of the ocis binary and enabled for the docker image builds.

https://github.com/owncloud/ocis/pull/10310
