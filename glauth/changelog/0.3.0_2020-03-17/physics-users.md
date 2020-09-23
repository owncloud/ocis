Change: use physicist demo users

Demo users like admin, demo and test don't allow you to tell a story. Which is why we changed the set of hard coded demo users to `einstein`, `marie` and `feynman`. You should know who they are. This also changes the ldap domain from `dc=owncloud,dc=com` to `dc=example,dc=org` because that is what these users use as their email domain. There are also `konnectd` and `reva` for technical purposes, eg. to allow konnectd and reva to bind to glauth.

<https://github.com/owncloud/ocis/glauth/issues/5>
