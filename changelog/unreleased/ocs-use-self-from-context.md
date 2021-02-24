Enhancement: OCS use user from context when rendering self

Instead of querying the accounts service we can render the result by using the user in the context, which has been populated by the proxy. This is possible because the proxy can handle basic auth.

https://github.com/owncloud/ocis/pull/855
