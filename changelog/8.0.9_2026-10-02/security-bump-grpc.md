Security: Bump grpc-go dependency

We've updated the google.golang.org/grpc dependency to fix a reported
denial-of-service vulnerability related to gRPC server handling of requests
missing authority or Host headers. As no tagged release containing the fix was
available yet, we pulled in a pre-release snapshot of grpc-go that includes it,
which also required bumping the minimum Go version to 1.26.8.

https://github.com/owncloud/ocis/pull/13020
