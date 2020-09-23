Change: add middleware for x-access-token distmantling

We added a middleware that dismantles the `x-access-token` from the request header and makes
it available in the context.

https://github.com/owncloud/ocis-pkg/pull/46
