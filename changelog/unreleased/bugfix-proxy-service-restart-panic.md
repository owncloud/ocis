Bugfix: Restarting the proxy service on error was causing panics

Starting oCIS with the `ocis server` command causes oCIS to monitor the
services being started. When a service fails to start, this monitor tries
to restart the services. For the cause of the proxy service, restarting it
was a panic.

Now, the proxy service can be restarted as expected, without that panic.
Note that the service might still fail to start if the actual cause isn't
fixed, such as the service's port being taken by a different application.

https://github.com/owncloud/ocis/pull/12769
