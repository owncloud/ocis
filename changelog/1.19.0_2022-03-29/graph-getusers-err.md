Bugfix: Fix error handling in GraphAPI GetUsers call

A missing return statement caused GetUsers to return misleading results when
the identity backend returned an error.

https://github.com/owncloud/ocis/pull/3357
