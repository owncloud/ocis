Bugfix: Allow selected updates on graph users

We are now allowing a couple of update request to complete even if GRAPH_LDAP_SERVER_WRITE_ENABLED=false:

*   When using a group to disable users (OCIS_LDAP_DISABLE_USER_MECHANISM=group) updates to the accountEnabled property of a user will be allowed
*   When a distinct base dn for new groups is configured ( GRAPH_LDAP_GROUP_CREATE_BASE_DN is set to a different value than GRAPH_LDAP_GROUP_BASE_DN), allow the creation/update of local groups.

https://github.com/owncloud/ocis/pull/6233
