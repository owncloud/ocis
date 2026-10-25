Bugfix: Clarify the message for a trash the user may not access

When a member of a space was not allowed to list the space's deleted files, the
trash view said the user "may not" have the permission, even though the server
had clearly denied access. The view now states that the user does not have
permission to view the deleted files. Other failures to load the trash show a
generic error message instead of a guess about permissions.

https://github.com/owncloud/ocis/pull/13217
