Enhancement: Create a Space using the Graph API

Spaces can now be created on `POST /drive/{drive-name}`. Only users with the `create-space` permissions can perform this operation.

Allowed body form values are:

- `quota` (bytes) maximum amount of bytes stored in the space.
- `maxQuotaFiles` (integer) maximum amount of files supported by the space.

https://github.com/owncloud/ocis/pull/2471
