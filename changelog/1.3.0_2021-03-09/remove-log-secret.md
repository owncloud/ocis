Enhancement: Remove the JWT from the log 

We were logging the JWT in some places. Secrets should not be exposed in logs so it got removed.

https://github.com/owncloud/ocis/pull/1758
