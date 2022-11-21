Bugfix: Make storage users mount ids unique by default

The mount ID of the storage users provider needs to be unique by default. We made this value configurable and added it to ocis init to be sure that we have a random uuid v4. This is important for federated instances.

https://github.com/owncloud/ocis/pull/5091
