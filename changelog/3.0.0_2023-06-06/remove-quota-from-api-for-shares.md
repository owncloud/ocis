Enhancement: Remove quota from share jails api responses

We have removed the quota object from api responses for share jails, 
which would permanently show exceeded due to restrictions in the permission system.

https://github.com/owncloud/ocis/pull/6309
https://github.com/owncloud/ocis/issues/4472