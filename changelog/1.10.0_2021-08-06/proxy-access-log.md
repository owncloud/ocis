Bugfix: log all requests in the proxy access log

We now use a dedicated middleware to log all requests, regardless of routing selector outcome.
While the log now includes the remote address, the selected routing policy is only logged when log level
is set to debug because the request context cannot be changed in the `directorSelectionDirector`, as per
the `ReverseProxy.Director` documentation. 

https://github.com/owncloud/ocis/pull/2301
