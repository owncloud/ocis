Bugfix: Don't write empty externalID to LDAP

When creating new users in the graph service, the externalID attribute was being written to LDAP even when it was empty.
Now, the externalID attribute is only written when it has a non-empty value.

https://github.com/owncloud/ocis/pull/12085
