Bugfix: Fix the oidc role assigner

The update role method did not allow to set a role when the user already has two roles.
This makes no sense as the user is supposed to have only one and the update will fix that.
We still log an error level log to make the admin aware of that.

https://github.com/owncloud/ocis/pull/6605
https://github.com/owncloud/ocis/pull/6618
