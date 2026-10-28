Bugfix: Return 400 for search queries that can not be compiled

A metadata search query that the KQL to bleve compiler can not compile, for
example one using an unsupported operator, returned an internal server error.
It is now reported as a bad request.

https://github.com/owncloud/ocis/issues/12230
