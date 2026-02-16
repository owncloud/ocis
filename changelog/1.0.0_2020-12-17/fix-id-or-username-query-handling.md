Bugfix: Fix id or username query handling

Tags: accounts

The code was stopping execution when encountering an error while loading an account by id. But for or queries we can continue execution.

https://github.com/owncloud/ocis/pull/745