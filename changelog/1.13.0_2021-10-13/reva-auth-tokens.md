Enhancement: Use reva's Authenticate method instead of spawning token managers

When using the CS3 proxy backend, we previously obtained the user from reva's
userprovider service and minted the token ourselves. This required maintaining
a shared JWT secret between ocis and reva, as well duplication of logic. This
PR delegates this logic by using the `Authenticate` method provided by the reva
gateway service to obtain this token, making it an arbitrary, indestructible
entry. Currently, the changes have been made to the proxy service but will be
extended to others as well.

https://github.com/owncloud/ocis/pull/2528
