Bugfix: Fix execution when passing program flags

The program flags where not correctly recognized because we didn't pass them to
the micro framework when initializing a grpc service.

https://github.com/owncloud/ocis/thumbnails/issues/15
