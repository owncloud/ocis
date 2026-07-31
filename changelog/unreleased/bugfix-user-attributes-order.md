Bugfix: Stable order for user search attributes

The `attributes` field returned from the user search endpoint came back in a
random order because `getUsersAttributes` ranged over a Go map. The function
now iterates over the configured `UserSearchDisplayedAttributes` slice, so
the returned attribute values follow the configured order.

https://github.com/owncloud/ocis/pull/12337
