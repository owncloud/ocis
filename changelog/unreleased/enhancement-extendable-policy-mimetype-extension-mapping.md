Enhancement: Extendable policy mimetype extension mapping

The extension mimetype mappings known from rego can now be extended.
To do this, ocis must be informed where the mimetype file (apache mime.types file format) is located.

`export OCIS_MACHINE_AUTH_API_KEY=$OCIS_HOME/mime.types`

https://github.com/owncloud/ocis/pull/6869
