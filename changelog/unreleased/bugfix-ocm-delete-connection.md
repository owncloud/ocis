Bugfix: Remove federated connections

Removing a federated connection in the connections panel sent a malformed
request to the server, so the connection was not actually removed and could
reappear after a page reload.

The request is now sent correctly and the connection is removed for good.

https://github.com/owncloud/ocis/pull/13160
