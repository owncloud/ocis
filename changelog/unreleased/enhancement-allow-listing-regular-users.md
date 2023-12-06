Enhancement: Allow regular users to list other users

Regular users can search for other users and groups. The following limitations
apply:

* Only search queries are allowed (using the `$search=term` query parameter)
* The search term needs to have at least 3 characters
* for user searches the result set only contains the attributes `displayName`,
  `userType`, `mail` and `id`
* for group searches the result set only contains the attributes `displayName`,
  `groupTypes` and `id`

https://github.com/owncloud/ocis/pull/7887
https://github.com/owncloud/ocis/issues/7782
