Bugfix: Fix DN parsing issues and sizelimit handling in libregraph/idm

We fixed a couple on issues in libregraph/idm related to correctly parsing
LDAP DNs for usernames contain characters that require escaping.

Also libregraph/idm was not properly returning "Size limit exceeded" errors
when the result set exceeded the requested size.

https://github.com/owncloud/ocis/issues/3631
https://github.com/owncloud/ocis/issues/4039
https://github.com/owncloud/ocis/issues/4078
