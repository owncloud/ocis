Bugfix: Search for usernames if memberUid is used for group membership

Ldap group membership usually use "member" or "uniqueMember" for groups,
which contains a DN as a value. For the case of "memberUid", the value isn't
a DN but just a username / uid. As such, getting the members of a group
didn't work because we expected the values to be DNs.

Using "memberUid" as group membership attribute is now supported, and it
will show the group members as expected.

https://github.com/owncloud/ocis/pull/12814
