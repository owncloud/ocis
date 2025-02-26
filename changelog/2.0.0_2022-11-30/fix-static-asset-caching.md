Bugfix: Disable cache for selected static web assets

We've disabled caching for some static web assets.
Files like the web index.html, oidc-callback.html or similar contain paths to timestamped resources and should not be cached.

https://github.com/owncloud/ocis/pull/4809
