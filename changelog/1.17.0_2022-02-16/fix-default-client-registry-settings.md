Bugfix: Make the default grpc client use the registry settings

We've fixed the default grpc client to use the registry settings. Previously it always
used mdns.

https://github.com/owncloud/ocis/pull/3041
