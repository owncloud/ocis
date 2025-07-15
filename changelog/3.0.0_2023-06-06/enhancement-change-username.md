Enhancement: allow username to be changed

When OnPremisesSamAccountName is present in a PATCH on `{apiRoot}/users/{userID}` it will change the
username of the user. This also changes the references to this user in the groups.

https://github.com/owncloud/ocis/pull/5509
https://github.com/owncloud/ocis/issues/4988
