Enhancement: Separate user and auth providers, add config for rest user

Previously, the auth and user provider services used to have the same driver,
which restricted using separate drivers and configs for both. This PR separates
the two and adds the config for the rest user driver and the gatewaysvc
parameter to EOS fs.

https://github.com/owncloud/ocis-reva/pull/412
https://github.com/cs3org/reva/pull/995
