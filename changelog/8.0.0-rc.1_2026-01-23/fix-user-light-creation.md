Bugfix: Fix user light creation

When trying to switch a user to user light before they logged in for the first time, an error would occur.
The server now correctly handles this case and allows switching to user light even before the first login.

https://github.com/owncloud/ocis/pull/11765
