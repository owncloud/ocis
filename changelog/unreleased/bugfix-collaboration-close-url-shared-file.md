Bugfix: Fix closing a shared document from the office editor

Closing a document or using "Open file location" in the office editor failed
with a private link error for users who had a single file shared with them.
The editor linked to the folder containing the file, which belongs to the
sharer and is not accessible to the sharee. The folder link is now only
provided if the user can access the folder, otherwise closing the editor
leads back to oCIS.

https://github.com/owncloud/ocis/pull/13198
