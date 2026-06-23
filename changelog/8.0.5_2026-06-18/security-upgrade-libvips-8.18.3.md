Security: Upgrade libvips to 8.18.3

Bumped libvips to 8.18.3 in all Docker images. The previous pin (8.18.2-r0)
was dropped from the Alpine edge/community repository, which broke the image
build.

https://github.com/owncloud/ocis/pull/12446
