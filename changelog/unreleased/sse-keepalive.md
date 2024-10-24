Bugfix: make SSE keepalive interval configurable

To prevent intermediate proxies from closing the SSE connection admins can now configure a `SSE_KEEPALIVE_INTERVAL`.

https://github.com/owncloud/ocis/pull/10411
