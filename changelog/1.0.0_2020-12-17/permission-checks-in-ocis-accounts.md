Change: Account management permissions for Admin role

Tags: accounts, settings

We created an `AccountManagement` permission and added it to the default admin role. There are permission
checks in place to protected http endpoints in ocis-accounts against requests without the permission.
All existing default users (einstein, marie, richard) have the default user role now (doesn't have the
`AccountManagement` permission). Additionally, there is a new default Admin user with credentials `moss:vista`.

Known issue: for users without the `AccountManagement` permission, the accounts UI extension is still available
in the ocis-web app switcher, but the requests for loading the users will fail (as expected). We are working
on a way to hide the accounts UI extension if the user doesn't have the `AccountManagement` permission.

https://github.com/owncloud/product/issues/124
https://github.com/owncloud/ocis-settings/pull/59
https://github.com/owncloud/ocis-settings/pull/66
https://github.com/owncloud/ocis-settings/pull/67
https://github.com/owncloud/ocis-settings/pull/69
https://github.com/owncloud/ocis-proxy/pull/95
https://github.com/owncloud/ocis-pkg/pull/59
https://github.com/owncloud/ocis-accounts/pull/95
https://github.com/owncloud/ocis-accounts/pull/100
https://github.com/owncloud/ocis-accounts/pull/102
