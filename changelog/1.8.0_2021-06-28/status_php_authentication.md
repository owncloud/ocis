Bugfix: Remove authentication from /status.php completely

Despite requests without Authentication header being successful, requests with an
invalid bearer token in the Authentication header were rejected in the proxy with
an 401 unauthenticated. Now the Authentication header is completely ignored for the
/status.php route.

https://github.com/owncloud/ocis/pull/2188
https://github.com/owncloud/client/issues/8538
