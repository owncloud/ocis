Bugfix: Show an error when a copy gets no answer

When copying, moving or duplicating files and folders, a request that ended
without any answer from the server, for example because the connection was
closed after a timeout, showed no message at all and left the file list
unchanged, so it looked as if nothing had happened. Such a request is now
reported with an error message. A folder duplication that got no answer is
reported as incomplete, because the part copied until then stays behind.

https://github.com/owncloud/ocis/pull/13354
https://github.com/owncloud/ocis/pull/13407
