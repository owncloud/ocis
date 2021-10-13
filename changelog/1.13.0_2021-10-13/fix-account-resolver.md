Bugfix: Fix the account resolver middleware

The accounts resolver middleware put an empty token into the request when the user was already present.
Added a step to get the token for the user.

https://github.com/owncloud/ocis/pull/2557
