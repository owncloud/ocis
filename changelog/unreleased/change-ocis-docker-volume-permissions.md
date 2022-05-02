Change: Reduce permissions on docker image predeclared volumes

We've lowered the permissions on the predeclared volumes of the oCIS
docker image from 777 to 750.

This change doesn't affect you, unless you use the docker image with
the non default uid/guid to start oCIS (default is 1000:1000).

https://github.com/owncloud/ocis/pull/3641
