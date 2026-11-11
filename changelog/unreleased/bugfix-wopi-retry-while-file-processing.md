Bugfix: Retry WOPI rename and upload while the file is processing

Renaming a file or uploading its contents through the collaboration service
failed with a 500 when the file was still being processed, for example right
after it had been created by a PutRelativeFile request. The storage answers
with a "too early" status in that case, which was only handled when deleting
files.

The collaboration service now retries the rename and upload requests with an
increasing waiting time until the file is ready. The delete request uses the
same retry logic, which is now limited to 10 retries instead of waiting
indefinitely.

https://github.com/owncloud/ocis/pull/13141
