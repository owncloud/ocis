Enhancement: OCM related adjustments in graph

The /users enpdoint of the graph service was changed with respect to how
it handles OCM federeated users:
- The 'userType' property is now alway returned. As new usertype 'Federated'
  was introduced. To indicate that the user is a federated user.
- Supported for filtering users by 'userType' as added. Queries like
  "$filter=userType eq 'Federated'" are now possible.
- Federated users are only returned when explicitly requested via filter.
  When no filter is provider only 'Member' users are returned.

https://github.com/owncloud/ocis/pull/9788
https://github.com/owncloud/ocis/pull/9757
https://github.com/owncloud/ocis/issues/9702
