Enhancement: Faster listing of storage spaces

Listing the storage spaces of a user, as done by `/graph/v1.0/me/drives`, no
longer reads the type of every space root from disk, and grants are parsed
without setting up a csv reader for each grant.

https://github.com/owncloud/ocis/pull/13595
https://github.com/owncloud/reva/pull/839
