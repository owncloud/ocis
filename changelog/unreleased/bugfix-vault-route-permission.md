Bugfix: Deny Safe/Vault access to users without permission

Users whose role did not grant Safe access could still reach it by editing the URL.
Such users are now denied and shown the access-denied page.

https://github.com/owncloud/ocis/pull/12884
https://github.com/owncloud/ocis/pull/12925
https://github.com/owncloud/ocis/pull/12928
