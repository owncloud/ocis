Bugfix: trace proxy middlewares

We moved trace initialization to an early middleware to also trace requests made by other proxy middlewares.

https://github.com/owncloud/ocis/pull/6313
