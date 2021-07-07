Bugfix: Dont use port 80 as debug for GroupsProvider

A copy/paste error where the configuration for the groupsprovider's debug address was not present leaves go-micro to start the debug service in port 80 by default.

https://github.com/owncloud/ocis/pull/2271
