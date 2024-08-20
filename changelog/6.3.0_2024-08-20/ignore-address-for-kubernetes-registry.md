Bugfix: Ignore address for kubernetes registry

We no longer pass an address to the go micro kubernetes registry implementation. This causes the implementation to autodetect the namespace and not hardcode it to `default`.

https://github.com/owncloud/ocis/pull/9490
