Enhancement: Make the LDAP base DN for new groups configurable

The LDAP backend for the Graph service introduced a new config option for setting the
Parent DN for new groups created via the `/groups/` endpoint. (`GRAPH_LDAP_GROUP_CREATE_BASE_DN`)

It defaults to the value of `GRAPH_LDAP_GROUP_BASE_DN`. If set to a different value the
`GRAPH_LDAP_GROUP_CREATE_BASE_DN` needs to be a subordinate DN of `GRAPH_LDAP_GROUP_BASE_DN`.

All existing groups with a DN outside the `GRAPH_LDAP_GROUP_CREATE_BASE_DN` tree will be treated as
read-only groups. So it is not possible to edit these groups.

https://github.com/owncloud/ocis/pull/5974

