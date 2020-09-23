Change: Add root path to static middleware

Currently the `Static` middleware always serves from the root path, but all our
HTTP handlers accept a custom root path which also got to be applied to the
static file handling.

<https://github.com/owncloud/ocis-pkg/issues/9>
