Bugfix: Cache compiled rego policies to avoid recompiling on every request

The policies service was reading and compiling .rego files from disk on
every request, causing unnecessary memory pressure and per-request
latency. The compiled PreparedEvalQuery is now cached per query string
so compilation happens at most once per query string over the lifetime
of the service.

https://github.com/owncloud/ocis/pull/12345