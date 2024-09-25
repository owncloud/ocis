Bugfix: Fix activity limit

When requesting a limit on activities, ocis would limit first, then filter and sort. Now it filters and sorts first, then limits.

https://github.com/owncloud/ocis/pull/10165
