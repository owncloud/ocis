Bugfix: Update LDAP filters

With the separation of use and find filters we can now use a filter that taken into account a users uuid as well as his username. This is necessary to make sharing work with the new account service which assigns accounts an immutable account id that is different from the username. Furthermore, the separate find filters now allows searching users by their displayname or email as well.


```
userfilter = "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))"
findfilter = "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))"
```

https://github.com/owncloud/ocis-reva/pull/399
https://github.com/cs3org/reva/pull/996