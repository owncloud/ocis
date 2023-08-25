Enhancement: Extendable policy mimetype extension mapping

The extension mimetype mappings known from rego can now be extended.
To do this, ocis must be informed where the mimetype file (apache mime.types file format) is located.

`export POLICIES_ENGINE_MIMES=OCIS_CONFIG_DIR/mime.types`

https://github.com/owncloud/ocis/pull/6869
