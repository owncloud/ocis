Bugfix: Use UUID attribute for computing "sub" claim in lico idp

By default the LDAP backend for lico uses the User DN for computing the "sub"
claim of a user. This caused the "sub" claim to stay the same even if a user
was deleted and recreated (and go a new UUID assgined with that). We now
use the user's unique id (`owncloudUUID` by default) for computing the `sub`
claim. So that user's recreated with the same name will be treated as different
users by the IDP.

https://github.com/owncloud/ocis/issues/904
https://github.com/owncloud/ocis/pull/6326
https://github.com/owncloud/ocis/pull/6338
https://github.com/owncloud/ocis/pull/6420
