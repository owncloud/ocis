Bugfix: Return to the Shares space after closing a shared document

Closing a document in the office editor took users who had a single file
shared with them to their Personal space instead of the Shares space. The
folder containing the file is not accessible to them, so the editor fell back
to the plain oCIS URL. In that case the editor now links to the shared file
itself, which leads back to the Shares space.

https://github.com/owncloud/ocis/pull/13293
