Change: Use account uuid from x-access-token

We are now using an ocis-pkg middleware for extracting the account uuid of the
authenticated user from the `x-access-token` of the http request header and inject
it into the Identifier protobuf messages wherever possible. This allows us to use
`me` instead of an actual account uuid, when the request comes through the proxy.

<https://github.com/owncloud/ocis/settings/pull/14>
