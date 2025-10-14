Bugfix: Fix the kql-bleve search

We fixed the issue when 500 on searches that contain ":". Added the characters escaping according to https://blevesearch.com/docs/Query-String-Query/


https://github.com/owncloud/ocis/pull/7290
https://github.com/owncloud/ocis/issues/7282
