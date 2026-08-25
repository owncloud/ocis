Bugfix: Prevent deleting your own account in the user management

In the admin settings user management, the account you are currently logged in
as could be selected and included in a delete action, even though the backend
always rejects self-deletion. Deleting a selection that included your own
account produced a "Failed to delete 1 user" error and made your own row
temporarily disappear from the list until the page was refreshed.

The delete action now excludes the current user from the request, so all other
selected users are deleted while your own account is left untouched. When your
own account is part of the selection, the delete confirmation dialog shows a
hint that it will not be deleted.

https://github.com/owncloud/ocis/pull/12661
https://github.com/owncloud/ocis/issues/12582
