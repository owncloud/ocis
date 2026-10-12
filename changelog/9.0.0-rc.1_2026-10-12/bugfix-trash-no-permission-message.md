Bugfix: Show an error instead of "no deleted files" for inaccessible trash

Members of a space who are not allowed to list its trash, for example members
with the "Can edit" role, were told that the space has no deleted files, even
though files had been deleted. The trash view now tells the user that the
deleted files could not be loaded and that they might lack the permission to
view them. A trash that is really empty still shows the "no deleted files"
message.

https://github.com/owncloud/ocis/pull/13105
