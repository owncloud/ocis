Bugfix: Fix the data race on the web theme defaults

The `web` service kept the theme defaults in a package-level map and merged the
base and branding themes into that shared map on every `GET /themes/*` request.
Concurrent requests therefore read and wrote the same map at the same time,
which is a data race: it could corrupt the rendered theme, leak values between
themes, or crash the service with a concurrent map read and map write.

The theme defaults are now built fresh for each request, so no state is shared
between concurrent requests.

https://github.com/owncloud/ocis/pull/13251
