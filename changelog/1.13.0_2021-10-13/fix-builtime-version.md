Bugfix: Fix version information for extensions

We've fixed the behavior for `ocis version` which previously always showed `0.0.0` as
version for extensions. Now the real version of the extensions are shown.

https://github.com/owncloud/ocis/pull/2575
