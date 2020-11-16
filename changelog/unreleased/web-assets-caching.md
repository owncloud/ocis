Change: Caching for static web assets

Tags: accounts, settings, web

We now set http caching headers for static web assets, so that they don't get force-reloaded on each request. The max-age for the caching is configurable and defaults to 7 days. The last modified date of the assets is set to the service start date, so that a service restart results in cache invalidation.

https://github.com/owncloud/ocis/pull/866
