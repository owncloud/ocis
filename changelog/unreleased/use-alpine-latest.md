Bugfix: build docker images with alpine:latest instead of alpine:edge

ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

https://github.com/owncloud/ocis-graph-explorer/pull/7
