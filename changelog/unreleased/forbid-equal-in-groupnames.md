Bugfix: Forbid `=` in group names

The underlying ldap library expects the name containing key-value pairs such as `uid=122`. It panics if you send just a `=`. 
It should not. However since we cannot rely on it we forbid using `=` in group names. A `BadRequest` will now be returned instead

https://github.com/owncloud/ocis/pull/5972
