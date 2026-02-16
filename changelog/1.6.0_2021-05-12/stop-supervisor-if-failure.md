Bugfix: Stop the supervisor if a service fails to start

Steps to make the supervisor fail:

`PROXY_HTTP_ADDR=0.0.0.0:9144 bin/ocis server`

https://github.com/owncloud/ocis/pull/1963
