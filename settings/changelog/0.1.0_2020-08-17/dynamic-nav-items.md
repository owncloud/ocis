Change: Dynamically add navItems for extensions with settings bundles

We now make use of a new feature in ocis-web-core, allowing us to add
navItems not only through configuration, but also after app initialization.
With this we now have navItems available for all extensions within the
settings ui, that have at least one settings bundle registered.

<https://github.com/owncloud/ocis/settings/pull/25>
