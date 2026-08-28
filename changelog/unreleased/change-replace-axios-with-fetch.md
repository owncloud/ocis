Change: Replace axios with the native fetch API in Web

The Web frontend no longer depends on axios. All HTTP traffic now goes through a
single fetch-based client, and the libre-graph client is generated from the
`typescript-fetch` template instead of `typescript-axios`.

Two changes are visible to consumers of the `@ownclouders/web-client` package.
The `graph` and `ocs` factories take a `FetchClient` where they previously took
an axios instance, and graph fields carrying an OData annotation are now exposed
under their generated camelCase names — `atLibreGraphPermissionsActions` for
`@libre.graph.permissions.actions`, and likewise for the other annotated fields.
The names sent over the wire are unchanged.

https://github.com/owncloud/ocis/pull/TBD
