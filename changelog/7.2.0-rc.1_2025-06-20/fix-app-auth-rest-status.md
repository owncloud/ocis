Bugfix: Fix app-auth, REST status code

Now app-auth REST returns status code 404 when creating token for non-existent user (Impersonation)

https://github.com/owncloud/ocis/pull/11190
https://github.com/owncloud/ocis/issues/10815
