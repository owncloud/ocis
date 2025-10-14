Bugfix: Let clients cache web and theme assets

We needed to remove "must-revalidate" from the cache-control header to allow clients to cache the web and theme assets.

https://github.com/owncloud/ocis/pull/6914
