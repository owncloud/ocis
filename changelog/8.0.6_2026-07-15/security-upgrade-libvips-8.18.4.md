Security: Upgrade libvips to 8.18.4

Bumped libvips to 8.18.4 in all Docker images. The previous pin (8.18.3-r0)
was dropped from the Alpine edge/community repository, which broke the image
build.

https://github.com/owncloud/ocis/pull/12596
