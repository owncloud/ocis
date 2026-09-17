Bugfix: Restoring an already-restored trash item now returns 404

Issuing a second restore request for a trash item that had already been
restored or purged returned a 500 (internal server error) instead of a 404
(not found). This made benign duplicate restores, such as those caused by a
sync client retrying a request, look like server errors. The restore
operation now reports 404 for this case, matching the existing behavior of
purging an already-purged item.

https://github.com/owncloud/ocis/pull/TODO
