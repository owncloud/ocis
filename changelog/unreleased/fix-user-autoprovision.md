Bugfix: Fix user autoprovisioning

We've fixed the autoprovsioning feature that was introduced in beta2. Due to a bug
the role assignment of the privileged user that is used to create accounts wasn't
propagated correctly to the `graph` service.

https://github.com/owncloud/ocis/issues/3893
