Security: Bump grpc dependency

We've bumped the grpc dependency to fix a vulnerability that could cause a
server panic when handling requests missing authority or Host headers.

https://github.com/owncloud/ocis/pull/13066
