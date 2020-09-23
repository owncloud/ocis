Change: Use UUIDs instead of alphanumeric identifiers

`Bundles`, `Settings` and `Values` were identified by a set of alphanumeric identifiers so far. We switched to UUIDs
in order to achieve a flat file hierarchy on disk. Referencing the respective entities by their alphanumeric
identifiers (as used in UI code) is still supported.

<https://github.com/owncloud/ocis/settings/pull/46>
