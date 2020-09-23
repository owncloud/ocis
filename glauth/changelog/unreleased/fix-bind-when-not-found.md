Bugfix: return invalid credentials when user was not found

We were relying on an error code of the ListAccounts call when the username and password was wrong. But the list will be empty if no user with the given login was found. So we also need to check if the list of accounts is empty.

<https://github.com/owncloud/ocis/glauth/pull/30>
