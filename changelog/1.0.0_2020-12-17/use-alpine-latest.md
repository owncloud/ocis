Bugfix: build docker images with alpine:latest instead of alpine:edge

Tags: docker

ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

https://github.com/owncloud/ocis/pull/416
