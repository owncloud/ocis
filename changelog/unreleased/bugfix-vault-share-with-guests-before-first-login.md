Bugfix: Share vault resources with guests before their first login

A user of another instance only gets a role on an instance at their first
login there. Until then, the vault share search did not show them and sharing
a vault resource with them was refused, even if that login would give them a
role with the "Vault mode" permission.

The graph service now derives that role from the user's `userClass` LDAP
attribute with `OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING`, like the proxy does
with the roles claim at login. If the role grants vault mode, the vault share
search shows the user, and sharing a vault resource with them assigns the role
first. Users who already have a role are not affected. The graph service needs
`OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING` set as well.

https://github.com/owncloud/ocis/pull/13273
