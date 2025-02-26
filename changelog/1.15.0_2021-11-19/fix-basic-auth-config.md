Bugfix: Fix basic auth config

Users could authenticate using basic auth even though `PROXY_ENABLE_BASIC_AUTH` was set to false.

https://github.com/owncloud/ocis/pull/2719
https://github.com/owncloud/ocis/issues/2466
