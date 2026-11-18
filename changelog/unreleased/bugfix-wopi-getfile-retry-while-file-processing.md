Bugfix: Retry WOPI GetFile while the file is still processing

Opening a file for editing right after it was uploaded could fail with an
error instead of showing the new content. GetFile now retries until the file
is ready.

https://github.com/owncloud/ocis/pull/13527
