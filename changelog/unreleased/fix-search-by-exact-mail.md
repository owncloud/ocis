Bugfix: Fix search by exact email

Users can be searched by exact email by using double quotes on the search
parameter. Note that double quotes are required because the "@" char is
being interpreted by the parser.

https://github.com/owncloud/ocis/pull/8035
